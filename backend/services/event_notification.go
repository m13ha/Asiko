package services

import (
	"context"

	"github.com/google/uuid"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/repository"
	"github.com/morkid/paginate"
)

type EventNotificationService interface {
	CreateEventNotification(userID uuid.UUID, eventType string, message string, resourceID uuid.UUID) error
	GetUserNotifications(ctx context.Context, userID string) (paginate.Page, error)
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

func (s *eventNotificationServiceImpl) GetUserNotifications(ctx context.Context, userID string) (paginate.Page, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return paginate.Page{}, myerrors.NewUserError("Invalid user ID.")
	}
	return s.notificationRepo.GetByUserID(ctx, uid), nil
}

func (s *eventNotificationServiceImpl) MarkAllNotificationsAsRead(userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return myerrors.NewUserError("Invalid user ID.")
	}
	return s.notificationRepo.MarkAllAsRead(uid)
}
