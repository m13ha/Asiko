package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apierrors "github.com/m13ha/asiko/errors/apierrors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/responses"
)

// @Summary Get user notifications
// @Description Retrieves a paginated list of notifications for the currently authenticated user.
// @Tags Notifications
// @Produce  application/json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Page size (default: 10)"
// @Success 200 {object} responses.PaginatedResponse{items=[]entities.Notification}
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /notifications [get]
// @ID getNotifications
func (h *Handler) GetNotificationsHandler(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	notifications, err := h.eventNotificationService.GetUserNotifications(c.Request.Context(), c.Request, userID.String())
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// @Summary Mark all notifications as read
// @Description Marks all notifications for the currently authenticated user as read.
// @Tags Notifications
// @Produce  application/json
// @Security BearerAuth
// @Success 200 {object} responses.SimpleMessage
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /notifications/read-all [put]
// @ID markAllNotificationsAsRead
func (h *Handler) MarkAllNotificationsAsReadHandler(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	if err := h.eventNotificationService.MarkAllNotificationsAsRead(userID.String()); err != nil {
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "All notifications marked as read."})
}

// @Summary Get unread notifications count
// @Description Retrieves the number of unread notifications for the currently authenticated user.
// @Tags Notifications
// @Produce  application/json
// @Security BearerAuth
// @Success 200 {object} map[string]int64 "count"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /notifications/unread-count [get]
// @ID getUnreadNotificationsCount
func (h *Handler) GetUnreadNotificationsCountHandler(c *gin.Context) {
	userID, ok := middleware.GetUUIDFromContext(c)
	if !ok {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	count, err := h.eventNotificationService.GetUserUnreadNotificationsCount(userID.String())
	if err != nil {
		apierrors.InternalServerError(c, "Internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
