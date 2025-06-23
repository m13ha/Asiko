package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS middleware to allow cross-origin requests
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from environment variable or use wildcard as fallback
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		origin := r.Header.Get("Origin")
		
		if allowedOrigins != "" {
			// Check if the request origin is in the allowed origins list
			origins := strings.Split(allowedOrigins, ",")
			allowed := false
			for _, o := range origins {
				if o == origin {
					allowed = true
					break
				}
			}
			
			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// If not in the allowed list, use wildcard (less secure but maintains compatibility)
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			// No specific origins configured, use wildcard
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
