package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS middleware to allow cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed origins from environment variable or use wildcard as fallback
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		origin := c.Request.Header.Get("Origin")

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
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// If not in the allowed list, use wildcard (less secure but maintains compatibility)
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			// No specific origins configured, use wildcard
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}
