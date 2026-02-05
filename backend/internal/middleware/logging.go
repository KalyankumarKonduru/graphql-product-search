package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

// StructuredLogger is a middleware that logs HTTP requests using Go's
// structured logging (slog). It captures method, path, status code,
// duration, response size, client IP, and request ID for each request.
func StructuredLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			recorder := &responseRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			duration := time.Since(start)
			requestID := GetRequestID(r.Context())

			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", recorder.statusCode),
				slog.Duration("duration", duration),
				slog.Int("bytes", recorder.bytes),
				slog.String("ip", extractIP(r)),
				slog.String("user_agent", r.UserAgent()),
			}

			if requestID != "" {
				attrs = append(attrs, slog.String("request_id", requestID))
			}

			// Log at appropriate level based on status code
			msg := "HTTP request"
			switch {
			case recorder.statusCode >= 500:
				logger.LogAttrs(r.Context(), slog.LevelError, msg, attrs...)
			case recorder.statusCode >= 400:
				logger.LogAttrs(r.Context(), slog.LevelWarn, msg, attrs...)
			default:
				logger.LogAttrs(r.Context(), slog.LevelInfo, msg, attrs...)
			}
		})
	}
}
