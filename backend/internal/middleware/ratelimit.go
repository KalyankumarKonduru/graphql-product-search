package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

// tokenBucket implements a simple token bucket rate limiter per client IP.
type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

// RateLimiter provides per-IP rate limiting using the token bucket algorithm.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*tokenBucket
	rps      float64
	burst    int
	cleanup  time.Duration
	lastClean time.Time
}

// NewRateLimiter creates a new rate limiter.
// rps: requests per second allowed per IP.
// burst: maximum burst size (bucket capacity).
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		buckets:   make(map[string]*tokenBucket),
		rps:       rps,
		burst:     burst,
		cleanup:   5 * time.Minute,
		lastClean: time.Now(),
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Periodic cleanup of stale buckets to prevent memory leaks
	if now.Sub(rl.lastClean) > rl.cleanup {
		for key, bucket := range rl.buckets {
			if now.Sub(bucket.lastRefill) > rl.cleanup {
				delete(rl.buckets, key)
			}
		}
		rl.lastClean = now
	}

	bucket, exists := rl.buckets[ip]
	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(rl.burst),
			maxTokens:  float64(rl.burst),
			refillRate: rl.rps,
			lastRefill: now,
		}
		rl.buckets[ip] = bucket
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	if bucket.tokens < 1 {
		return false
	}

	bucket.tokens--
	return true
}

// extractIP gets the client IP address, checking X-Forwarded-For and
// X-Real-IP headers for proxied requests.
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header (first IP is the client)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := net.ParseIP(xff)
		if parts != nil {
			return parts.String()
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}

	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// Handler returns an HTTP middleware that enforces rate limiting.
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)

		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests. Please try again later.",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
