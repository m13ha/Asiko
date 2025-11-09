package repository

import (
    apperr "github.com/m13ha/asiko/errors"
    "github.com/m13ha/asiko/models/entities"
    "gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	FindByPhone(phone string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	Create(user *entities.User) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) FindByEmail(email string) (*entities.User, error) {
    var user entities.User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.user.FindByEmail", err)
    }
    return &user, nil
}

func (r *gormUserRepository) FindByPhone(phone string) (*entities.User, error) {
    var user entities.User
    if err := r.db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.user.FindByPhone", err)
    }
    return &user, nil
}

func (r *gormUserRepository) FindByID(id string) (*entities.User, error) {
    var user entities.User
    if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
        return nil, apperr.TranslateRepoError("repository.user.FindByID", err)
    }
    return &user, nil
}

func (r *gormUserRepository) Create(user *entities.User) error {
    if err := r.db.Create(user).Error; err != nil {
        return apperr.TranslateRepoError("repository.user.Create", err)
    }
    return nil
}
