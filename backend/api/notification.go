package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/responses"
)

// @Summary Get user notifications
// @Description Retrieves a paginated list of notifications for the currently authenticated user.
// @Tags Notifications
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Notification}
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /notifications [get]
// @ID getNotifications
func (h *Handler) GetNotificationsHandler(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	notifications, err := h.eventNotificationService.GetUserNotifications(c.Request.Context(), userIDStr)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// @Summary Mark all notifications as read
// @Description Marks all notifications for the currently authenticated user as read.
// @Tags Notifications
// @Produce  json
// @Security BearerAuth
// @Success 200 {object} responses.ResponsesSimpleMessageResponse
// @Failure 401 {object} errors.ApiErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /notifications/read-all [put]
// @ID markAllNotificationsAsRead
func (h *Handler) MarkAllNotificationsAsReadHandler(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "Unauthorized")
		return
	}

	if err := h.eventNotificationService.MarkAllNotificationsAsRead(userIDStr); err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, responses.ResponsesSimpleMessageResponse{Message: "All notifications marked as read."})
}
