package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/m13ha/asiko/errors/apierrors" // New import
)

// ErrorHandler captures panics and writes standardized JSON errors for any
// accumulated errors in the Gin context.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// Pick the last error
		ge := c.Errors.Last()
		apierrors.HandleAppError(c, ge.Err) // Use the centralized error handler
		// c.Abort() is handled within HandleAppError implicitly through c.JSON which writes the header and body,
		// and the current Gin context stops further processing implicitly.
	}
}
