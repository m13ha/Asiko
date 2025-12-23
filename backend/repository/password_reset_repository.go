package repository

import (
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"github.com/m13ha/asiko/models/entities"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	Create(token *entities.PasswordResetToken) error
	FindByToken(token string) (*entities.PasswordResetToken, error)
	DeleteAllForUser(userID string) error
}

type gormPasswordResetRepository struct {
	db *gorm.DB
}

func NewGormPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &gormPasswordResetRepository{db: db}
}

func (r *gormPasswordResetRepository) Create(token *entities.PasswordResetToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return repoerrors.InternalError("failed to create password reset token: " + err.Error())
	}
	return nil
}

func (r *gormPasswordResetRepository) FindByToken(token string) (*entities.PasswordResetToken, error) {
	var resetToken entities.PasswordResetToken
	if err := r.db.Where("token = ?", token).First(&resetToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("invalid or expired token")
		}
		return nil, repoerrors.InternalError("failed to find token: " + err.Error())
	}
	return &resetToken, nil
}

func (r *gormPasswordResetRepository) DeleteAllForUser(userID string) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&entities.PasswordResetToken{}).Error; err != nil {
		return repoerrors.InternalError("failed to delete user tokens: " + err.Error())
	}
	return nil
}
