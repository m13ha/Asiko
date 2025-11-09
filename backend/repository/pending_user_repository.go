package repository

import (
    apperr "github.com/m13ha/asiko/errors"
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
        return apperr.TranslateRepoError("repository.pending.Create", err)
    }
    return nil
}

func (r *gormPendingUserRepository) FindByEmail(email string) (*entities.PendingUser, error) {
	var user entities.PendingUser
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.pending.FindByEmail", err)
    }
    return &user, nil
}

func (r *gormPendingUserRepository) Update(user *entities.PendingUser) error {
    if err := r.db.Save(user).Error; err != nil {
        return apperr.TranslateRepoError("repository.pending.Update", err)
    }
    return nil
}

func (r *gormPendingUserRepository) Delete(email string) error {
    if err := r.db.Where("email = ?", email).Delete(&entities.PendingUser{}).Error; err != nil {
        return apperr.TranslateRepoError("repository.pending.Delete", err)
    }
    return nil
}
