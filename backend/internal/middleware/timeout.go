package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// Timeout is a middleware that sets a deadline on request processing.
// If the handler does not complete within the specified duration,
// the request context is cancelled, allowing downstream handlers
// to detect the cancellation and return early.
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			done := make(chan struct{})
			tw := &timeoutWriter{ResponseWriter: w, header: make(http.Header)}

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				// Copy headers from the timeout writer
				for k, v := range tw.header {
					for _, val := range v {
						w.Header().Set(k, val)
					}
				}
				if tw.code != 0 {
					w.WriteHeader(tw.code)
				}
				w.Write(tw.body)
			case <-ctx.Done():
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusGatewayTimeout)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "request_timeout",
					"message": "Request timed out. Please try again.",
				})
			}
		})
	}
}

type timeoutWriter struct {
	http.ResponseWriter
	header http.Header
	body   []byte
	code   int
}

func (tw *timeoutWriter) Header() http.Header {
	return tw.header
}

func (tw *timeoutWriter) Write(data []byte) (int, error) {
	tw.body = append(tw.body, data...)
	return len(data), nil
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.code = code
}
