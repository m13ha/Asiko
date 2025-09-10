package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
)

// @Summary Get user analytics
// @Description Get analytics data for the authenticated user including total appointments created and total bookings received
// @Tags Analytics
// @Produce json
// @Security BearerAuth
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} responses.AnalyticsResponse
// @Failure 400 {object} errors.ApiErrorResponse "Invalid date format or missing parameters"
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /analytics [get]
// @ID getUserAnalytics
func (h *Handler) GetUserAnalytics(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.BadRequest(c.Writer, "Invalid user ID")
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		errors.BadRequest(c.Writer, "start_date and end_date query parameters are required")
		return
	}

	analytics, err := h.analyticsService.GetUserAnalytics(userID, startDate, endDate)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, analytics)
}