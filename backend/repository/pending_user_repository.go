package repository

import (
	"github.com/m13ha/appointment_master/models/entities"
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
	return r.db.Create(user).Error
}

func (r *gormPendingUserRepository) FindByEmail(email string) (*entities.PendingUser, error) {
	var user entities.PendingUser
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormPendingUserRepository) Update(user *entities.PendingUser) error {
	return r.db.Save(user).Error
}

func (r *gormPendingUserRepository) Delete(email string) error {
	return r.db.Where("email = ?", email).Delete(&entities.PendingUser{}).Error
}
