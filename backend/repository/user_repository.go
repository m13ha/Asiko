package repository

import (
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"github.com/m13ha/asiko/models/entities"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
	FindByPhone(phone string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
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
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("user not found with email: " + email)
		}
		return nil, repoerrors.InternalError("failed to find user by email: " + err.Error())
	}
	return &user, nil
}

func (r *gormUserRepository) FindByPhone(phone string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("phone_number = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("user not found with phone: " + phone)
		}
		return nil, repoerrors.InternalError("failed to find user by phone: " + err.Error())
	}
	return &user, nil
}

func (r *gormUserRepository) FindByID(id string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("user not found with id: " + id)
		}
		return nil, repoerrors.InternalError("failed to find user by id: " + err.Error())
	}
	return &user, nil
}

func (r *gormUserRepository) Create(user *entities.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return repoerrors.InternalError("failed to create user: " + err.Error())
	}
	return nil
}

func (r *gormUserRepository) Update(user *entities.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return repoerrors.InternalError("failed to update user: " + err.Error())
	}
	return nil
}
