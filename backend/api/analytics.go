package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apierrors "github.com/m13ha/asiko/errors/apierrors"
	"github.com/m13ha/asiko/middleware"
)

// @Summary Get user analytics
// @Description Get analytics for the authenticated user over a date window.
// @Description Includes totals, breakdowns (by type/status, guest vs registered), utilization,
// @Description lead-time stats, daily series, peak hours/days, and top appointments.
// @Tags Analytics
// @Produce application/json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} responses.AnalyticsResponse
// @Failure 400 {object} responses.APIErrorResponse "Invalid date format or missing parameters"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /analytics [get]
// @ID getUserAnalytics
func (h *Handler) GetUserAnalytics(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		apierrors.ValidationError(c, "start_date and end_date query parameters are required")
		return
	}

	analytics, err := h.analyticsService.GetUserAnalytics(userID, startDate, endDate)
	if err != nil {
		// For service errors, we'll use internal server error unless we have specific error handling
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, analytics)
}
