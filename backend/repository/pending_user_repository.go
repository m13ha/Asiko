package repository

import (
    repoerrors "github.com/m13ha/asiko/errors/repoerrors"
    "github.com/m13ha/asiko/models/entities"
    "gorm.io/gorm"
)

type PendingUserRepository interface {
	Create(user *entities.PendingUser) error
	FindByEmail(email string) (*entities.PendingUser, error)
	Update(user *entities.PendingUser) error
	Delete(email string) error
}

type gormPendingUserRepository struct {
	db *gorm.DB
}

func NewGormPendingUserRepository(db *gorm.DB) PendingUserRepository {
	return &gormPendingUserRepository{db: db}
}

func (r *gormPendingUserRepository) Create(user *entities.PendingUser) error {
    if err := r.db.Create(user).Error; err != nil {
        return repoerrors.InternalError("failed to create pending user: " + err.Error())
    }
    return nil
}

func (r *gormPendingUserRepository) FindByEmail(email string) (*entities.PendingUser, error) {
	var user entities.PendingUser
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, repoerrors.NotFoundError("pending user not found with email: " + email)
        }
        return nil, repoerrors.InternalError("failed to find pending user by email: " + err.Error())
    }
    return &user, nil
}

func (r *gormPendingUserRepository) Update(user *entities.PendingUser) error {
    if err := r.db.Save(user).Error; err != nil {
        return repoerrors.InternalError("failed to update pending user: " + err.Error())
    }
    return nil
}

func (r *gormPendingUserRepository) Delete(email string) error {
    if err := r.db.Where("email = ?", email).Delete(&entities.PendingUser{}).Error; err != nil {
        return repoerrors.InternalError("failed to delete pending user: " + err.Error())
    }
    return nil
}
