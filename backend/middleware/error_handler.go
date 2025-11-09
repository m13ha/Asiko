package middleware

import (
	"github.com/gin-gonic/gin"
	appErr "github.com/m13ha/asiko/errors"
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
		ae := appErr.FromError(ge.Err)
		status := ae.HTTP
		if status == 0 {
			status = appErr.StatusFromKind(ae.Kind)
		}
		reqID := c.Writer.Header().Get("X-Request-ID")
		if reqID == "" {
			reqID = c.GetHeader("X-Request-ID")
		}
		resp := appErr.APIErrorResponse{
			Status:    status,
			Code:      ae.Code,
			Message:   ae.Message,
			Fields:    ae.Fields,
			RequestID: reqID,
			Meta:      ae.Meta,
		}
		c.JSON(status, resp)
		// Stop other handlers from writing
		c.Abort()
	}
}
