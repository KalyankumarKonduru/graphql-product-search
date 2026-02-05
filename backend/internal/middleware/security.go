package middleware

import (
	"net/http"
)

// SecurityHeaders adds security-related HTTP headers to all responses.
// These headers protect against common web vulnerabilities including
// XSS, clickjacking, MIME sniffing, and information disclosure.
func SecurityHeaders(isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")

			// XSS protection (legacy but still useful for older browsers)
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Control referrer information
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Restrict browser features
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

			// Prevent caching of sensitive data
			if r.URL.Path == "/graphql" {
				w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
				w.Header().Set("Pragma", "no-cache")
			}

			if isProd {
				// HSTS - enforce HTTPS (2 years, include subdomains, preload eligible)
				w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

				// Content Security Policy
				w.Header().Set("Content-Security-Policy",
					"default-src 'self'; "+
						"script-src 'self'; "+
						"style-src 'self' 'unsafe-inline'; "+
						"img-src 'self' https://picsum.photos data:; "+
						"font-src 'self'; "+
						"connect-src 'self'; "+
						"frame-ancestors 'none'; "+
						"base-uri 'self'; "+
						"form-action 'self'")
			}

			// Remove server identification header
			w.Header().Del("Server")
			w.Header().Del("X-Powered-By")

			next.ServeHTTP(w, r)
		})
	}
}
