package services

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/m13ha/appointment_master/db"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(userReq models.UserRequest) (*models.User, error) {
	if err := utils.Validate(userReq); err != nil {
		log.Error().Err(err).Msg("User validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Log normalized email for debugging
	normalizedEmail := utils.NormalizeEmail(userReq.Email)
	log.Debug().
		Str("original_email", userReq.Email).
		Str("normalized_email", normalizedEmail).
		Msg("Email normalization")

	// Check if user with email already exists
	var existingUser models.User
	if result := db.DB.Where("email = ?", normalizedEmail).First(&existingUser); result.Error == nil {
		log.Warn().
			Str("email", normalizedEmail).
			Msg("Attempted registration with existing email")
		return nil, fmt.Errorf("failed to create user: duplicate key value violates unique constraint on email")
	}

	// Check if user with phone number already exists
	if userReq.PhoneNumber != "" {
		if result := db.DB.Where("phone_number = ?", userReq.PhoneNumber).First(&existingUser); result.Error == nil {
			log.Warn().
				Str("phone", userReq.PhoneNumber).
				Msg("Attempted registration with existing phone number")
			return nil, fmt.Errorf("failed to create user: duplicate key value violates unique constraint on phone")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:           userReq.Name,
		Email:          normalizedEmail,
		PhoneNumber:    userReq.PhoneNumber,
		HashedPassword: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Error().Err(err).
			Str("name", user.Name).
			Str("email", user.Email).
			Msg("Database error when creating user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("User created successfully")

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
	if err := db.DB.Where("email = ?", utils.NormalizeEmail(email)).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.CheckPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}
