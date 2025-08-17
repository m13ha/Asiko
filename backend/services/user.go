package services

import (
	"fmt"

	"github.com/m13ha/appointment_master/db"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// ToUserResponse converts an entities.User to a responses.UserResponse
func ToUserResponse(user *entities.User) *responses.UserResponse {
	return &responses.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func CreateUser(userReq requests.UserRequest) (*responses.UserResponse, error) {
	if err := utils.Validate(userReq); err != nil {
		log.Error().Err(err).Msg("User validation failed")
		return nil, myerrors.NewUserError("Invalid user data. Please check your input.")
	}

	// Log normalized email for debugging
	normalizedEmail := utils.NormalizeEmail(userReq.Email)
	log.Debug().
		Str("original_email", userReq.Email).
		Str("normalized_email", normalizedEmail).
		Msg("Email normalization")

	// Check if user with email already exists
	var existingUser entities.User
	if result := db.DB.Where("email = ?", normalizedEmail).First(&existingUser); result.Error == nil {
		log.Warn().
			Str("email", normalizedEmail).
			Msg("Attempted registration with existing email")
		return nil, myerrors.NewUserError("Email already registered.")
	}

	// Check if user with phone number already exists
	if userReq.PhoneNumber != "" {
		if result := db.DB.Where("phone_number = ?", userReq.PhoneNumber).First(&existingUser); result.Error == nil {
			log.Warn().
				Str("phone", userReq.PhoneNumber).
				Msg("Attempted registration with existing phone number")
			return nil, myerrors.NewUserError("Phone number already registered.")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		return nil, fmt.Errorf("internal error")
	}

	user := &entities.User{
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
		return nil, fmt.Errorf("internal error")
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("User created successfully")

	return ToUserResponse(user), nil
}

func AuthenticateUser(email, password string) (*entities.User, error) {
	var user entities.User
	if err := db.DB.Where("email = ?", utils.NormalizeEmail(email)).First(&user).Error; err != nil {
		return nil, myerrors.NewUserError("Invalid email or password.")
	}

	if !user.CheckPassword(password) {
		return nil, myerrors.NewUserError("Invalid email or password.")
	}

	return &user, nil
}
