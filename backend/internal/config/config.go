package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	Port             string
	Environment      string // development, staging, production
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
	ShutdownTimeout  time.Duration
	RequestTimeout   time.Duration

	// Security
	JWTSecret          string
	APIKeys            []string
	AllowedOrigins     []string
	EnablePlayground   bool
	EnableIntrospection bool

	// Rate Limiting
	RateLimitRPS   float64
	RateLimitBurst int

	// GraphQL Limits
	MaxQueryDepth      int
	MaxQueryComplexity int
	MaxPageSize        int

	// TLS
	TLSCertFile string
	TLSKeyFile  string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	env := getEnv("APP_ENV", "development")
	isProd := env == "production"

	return &Config{
		// Server
		Port:            getEnv("PORT", "4000"),
		Environment:     env,
		ReadTimeout:     getDuration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getDuration("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:     getDuration("IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		RequestTimeout:  getDuration("REQUEST_TIMEOUT", 10*time.Second),

		// Security
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production"),
		APIKeys:             getStringSlice("API_KEYS", nil),
		AllowedOrigins:      getStringSlice("ALLOWED_ORIGINS", defaultOrigins(isProd)),
		EnablePlayground:    getBool("ENABLE_PLAYGROUND", !isProd),
		EnableIntrospection: getBool("ENABLE_INTROSPECTION", !isProd),

		// Rate Limiting
		RateLimitRPS:   getFloat("RATE_LIMIT_RPS", 100),
		RateLimitBurst: getInt("RATE_LIMIT_BURST", 200),

		// GraphQL Limits
		MaxQueryDepth:      getInt("MAX_QUERY_DEPTH", 10),
		MaxQueryComplexity: getInt("MAX_QUERY_COMPLEXITY", 200),
		MaxPageSize:        getInt("MAX_PAGE_SIZE", 100),

		// TLS
		TLSCertFile: getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:  getEnv("TLS_KEY_FILE", ""),
	}
}

func (c *Config) IsProd() bool {
	return c.Environment == "production"
}

func (c *Config) IsDev() bool {
	return c.Environment == "development"
}

func defaultOrigins(isProd bool) []string {
	if isProd {
		return []string{} // Must be explicitly configured in production
	}
	return []string{
		"http://localhost:3000",
		"http://localhost:5173",
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getFloat(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return fallback
}

func getBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getStringSlice(key string, fallback []string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return fallback
}
