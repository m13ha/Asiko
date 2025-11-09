package repository

import (
    "context"

    "github.com/google/uuid"
    apperr "github.com/m13ha/asiko/errors"
    "github.com/m13ha/asiko/models/entities"
    "github.com/morkid/paginate"
    "gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *entities.Notification) error
	GetByUserID(ctx context.Context, userID uuid.UUID) paginate.Page
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
        return apperr.TranslateRepoError("repository.notification.Create", err)
    }
    return nil
}

func (r *gormNotificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Notification{}).Where("user_id = ?", userID).Order("created_at DESC")
	return pg.With(db).Request(ctx).Response(&[]entities.Notification{})
}

func (r *gormNotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
    if err := r.db.Model(&entities.Notification{}).Where("user_id = ? AND is_read = false", userID).Update("is_read", true).Error; err != nil {
        return apperr.TranslateRepoError("repository.notification.MarkAllAsRead", err)
    }
    return nil
}
