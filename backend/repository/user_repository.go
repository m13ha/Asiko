package repository

import (
	"github.com/m13ha/appointment_master/models/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	FindByPhone(phone string) (*entities.User, error)
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
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) FindByPhone(phone string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}
