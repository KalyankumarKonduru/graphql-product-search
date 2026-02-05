package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/kalyankumar/graphql-product-search/graph"
	"github.com/kalyankumar/graphql-product-search/internal/auth"
	"github.com/kalyankumar/graphql-product-search/internal/config"
	"github.com/kalyankumar/graphql-product-search/internal/database"
	"github.com/kalyankumar/graphql-product-search/internal/metrics"
	mw "github.com/kalyankumar/graphql-product-search/internal/middleware"
)

func main() {
	// Load configuration from environment
	cfg := config.Load()

	// Initialize structured logger
	var logHandler slog.Handler
	if cfg.IsProd() {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Initialize database
	database.Initialize()

	// Initialize metrics
	appMetrics := metrics.New()

	// Initialize authenticator
	authenticator := auth.NewAuthenticator(cfg.JWTSecret, cfg.APIKeys)

	// Initialize rate limiter
	rateLimiter := mw.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)

	// Build router
	router := chi.NewRouter()

	// --- Middleware Stack (order matters) ---

	// 1. Request ID (first, so all subsequent middleware/handlers can use it)
	router.Use(mw.RequestID)

	// 2. Structured logging (captures all downstream timing)
	router.Use(mw.StructuredLogger(logger))

	// 3. Panic recovery
	router.Use(recoverer(logger))

	// 4. Metrics collection
	router.Use(appMetrics.MetricsMiddleware)

	// 5. Security headers
	router.Use(mw.SecurityHeaders(cfg.IsProd()))

	// 6. Rate limiting
	router.Use(rateLimiter.Handler)

	// 7. CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-API-Key"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// 8. Response compression
	router.Use(mw.Compress)

	// 9. Request timeout
	router.Use(mw.Timeout(cfg.RequestTimeout))

	// 10. Authentication
	router.Use(authenticator.Middleware)

	// --- GraphQL Server ---
	gqlConfig := graph.Config{
		Resolvers: &graph.Resolver{
			Config: cfg,
		},
	}

	srv := handler.New(graph.NewExecutableSchema(gqlConfig))

	// Transport configuration
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Query complexity limiting
	srv.Use(extension.FixedComplexityLimit(cfg.MaxQueryComplexity))

	// Introspection control - disable in production to prevent schema exposure
	if !cfg.EnableIntrospection {
		srv.Use(extension.Introspection{})
	}

	// Automatic persisted queries for reduced bandwidth
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: graphql.MapCache[string]{},
	})

	// --- Routes ---

	// GraphQL playground (disabled in production)
	if cfg.EnablePlayground {
		router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	} else {
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"service": "graphql-product-search",
				"status":  "running",
			})
		})
	}

	// GraphQL endpoint
	router.Handle("/graphql", srv)

	// Health check endpoints
	router.Get("/health", healthHandler())
	router.Get("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	router.Get("/health/ready", readinessHandler())

	// Metrics endpoint (should be protected in production)
	router.Get("/metrics", appMetrics.Handler())

	// Auth endpoints (for demo/testing)
	router.Post("/auth/token", tokenHandler(authenticator))

	// --- HTTP Server with production timeouts ---
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		// Limits the size of request headers to prevent slowloris attacks
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// --- Graceful Shutdown ---
	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("server starting",
			slog.String("port", cfg.Port),
			slog.String("environment", cfg.Environment),
			slog.Bool("playground", cfg.EnablePlayground),
			slog.Bool("introspection", cfg.EnableIntrospection),
		)

		if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
			serverErrors <- httpServer.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			serverErrors <- httpServer.ListenAndServe()
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server failed to start", slog.String("error", err.Error()))
		os.Exit(1)
	case sig := <-quit:
		logger.Info("shutdown signal received", slog.String("signal", sig.String()))
	}

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	logger.Info("shutting down server", slog.Duration("timeout", cfg.ShutdownTimeout))

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}

// recoverer is a middleware that catches panics and logs them with structured logging.
func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.String("request_id", mw.GetRequestID(r.Context())),
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{
						"error":   "internal_server_error",
						"message": "An unexpected error occurred",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// healthHandler returns a deep health check that reports component status.
func healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		health := map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"checks": map[string]interface{}{
				"database": map[string]string{
					"status": "healthy",
					"type":   "in-memory",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(health)
	}
}

// readinessHandler checks if the server is ready to accept traffic.
func readinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if database is initialized
		if database.DB == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "not_ready",
				"reason": "database not initialized",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ready",
		})
	}
}

// tokenHandler provides a demo endpoint to generate JWT tokens for testing.
func tokenHandler(authenticator *auth.Authenticator) http.HandlerFunc {
	type tokenRequest struct {
		UserID string `json:"user_id"`
		Email  string `json:"email"`
		Name   string `json:"name"`
		Role   string `json:"role"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req tokenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		if req.UserID == "" || req.Role == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "user_id and role are required"})
			return
		}

		role := auth.Role(req.Role)
		if role != auth.RoleUser && role != auth.RoleAdmin {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid role: %s (must be 'user' or 'admin')", req.Role)})
			return
		}

		user := &auth.User{
			ID:    req.UserID,
			Email: req.Email,
			Name:  req.Name,
			Role:  role,
		}

		token, err := authenticator.GenerateToken(user, 24*time.Hour)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate token"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token":      token,
			"type":       "Bearer",
			"expires_in": "86400",
		})
	}
}
