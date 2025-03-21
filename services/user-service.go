package services

import (
	"fmt"

	"github.com/m13ha/appointment_master/db"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(userReq models.UserRequest) (*models.User, error) {
	if err := utils.Validate(userReq); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:           userReq.Name,
		Email:          userReq.Email,
		PhoneNumber:    userReq.PhoneNumber,
		HashedPassword: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func GetUserBookings(userID string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := db.DB.Where("user_id = ?", userID).Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	return bookings, nil
}

func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
