package repository

import (
    "context"
    "net/http"

    "github.com/google/uuid"
    repoerrors "github.com/m13ha/asiko/errors/repoerrors"
    "github.com/m13ha/asiko/models/entities"
    "github.com/morkid/paginate"
    "gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entities.Notification) error
	GetByUserID(ctx context.Context, req *http.Request, userID uuid.UUID) paginate.Page
	MarkAllAsRead(userID uuid.UUID) error
}

type gormNotificationRepository struct {
	db *gorm.DB
}

func NewGormNotificationRepository(db *gorm.DB) NotificationRepository {
	return &gormNotificationRepository{db: db}
}

func (r *gormNotificationRepository) Create(notification *entities.Notification) error {
    if err := r.db.Create(notification).Error; err != nil {
        return repoerrors.InternalError("failed to create notification: " + err.Error())
    }
    return nil
}

func (r *gormNotificationRepository) GetByUserID(ctx context.Context, req *http.Request, userID uuid.UUID) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Notification{}).Where("user_id = ?", userID).Order("created_at DESC")
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]entities.Notification{})
}

func (r *gormNotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
    if err := r.db.Model(&entities.Notification{}).Where("user_id = ? AND is_read = false", userID).Update("is_read", true).Error; err != nil {
        return repoerrors.InternalError("failed to mark all notifications as read: " + err.Error())
    }
    return nil
}
