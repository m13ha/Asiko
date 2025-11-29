package services

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/repository"
	"github.com/morkid/paginate"
)

type EventNotificationService interface {
	CreateEventNotification(userID uuid.UUID, eventType string, message string, resourceID uuid.UUID) error
	GetUserNotifications(ctx context.Context, req *http.Request, userID string) (paginate.Page, error)
	MarkAllNotificationsAsRead(userID string) error
}

type eventNotificationServiceImpl struct {
	notificationRepo repository.NotificationRepository
}

func NewEventNotificationService(notificationRepo repository.NotificationRepository) EventNotificationService {
	return &eventNotificationServiceImpl{notificationRepo: notificationRepo}
}

func (s *eventNotificationServiceImpl) CreateEventNotification(userID uuid.UUID, eventType string, message string, resourceID uuid.UUID) error {
	notification := &entities.Notification{
		UserID:     userID,
		EventType:  eventType,
		Message:    message,
		ResourceID: resourceID,
	}
	return s.notificationRepo.Create(notification)
}

func (s *eventNotificationServiceImpl) GetUserNotifications(ctx context.Context, req *http.Request, userID string) (paginate.Page, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.Page{}, serviceerrors.ValidationError("Invalid user ID.")
	}
	return s.notificationRepo.GetByUserID(ctx, req, uid), nil
}

func (s *eventNotificationServiceImpl) MarkAllNotificationsAsRead(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return serviceerrors.ValidationError("Invalid user ID.")
	}
	return s.notificationRepo.MarkAllAsRead(uid)
}
