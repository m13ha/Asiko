package services

import (
	"fmt"
	"time"

	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/notifications"
	"github.com/m13ha/appointment_master/repository"
	"github.com/m13ha/appointment_master/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepo        repository.UserRepository
	pendingUserRepo repository.PendingUserRepository
	notificationSvc notifications.NotificationService
}

func NewUserService(userRepo repository.UserRepository, pendingUserRepo repository.PendingUserRepository, notificationSvc notifications.NotificationService) UserService {
	return &userServiceImpl{userRepo: userRepo, pendingUserRepo: pendingUserRepo, notificationSvc: notificationSvc}
}

func (s *userServiceImpl) CreateUser(userReq requests.UserRequest) (*responses.UserResponse, error) {
	if err := utils.Validate(userReq); err != nil {
		return nil, myerrors.NewUserError("Invalid user data.")
	}

	normalizedEmail := utils.NormalizeEmail(userReq.Email)

	// Check if user already exists in main table
	_, err := s.userRepo.FindByEmail(normalizedEmail)
	if err == nil {
		return nil, myerrors.NewUserError("Email already registered.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	verificationCode := utils.GenerateRandomCode(6)
	expiresAt := time.Now().Add(15 * time.Minute)

	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("internal error")
	}

	if pendingUser != nil {
		// Update existing pending user
		pendingUser.Name = userReq.Name
		pendingUser.HashedPassword = string(hashedPassword)
		pendingUser.PhoneNumber = userReq.PhoneNumber
		pendingUser.VerificationCode = verificationCode
		pendingUser.VerificationCodeExpiresAt = expiresAt
		if err := s.pendingUserRepo.Update(pendingUser); err != nil {
			return nil, fmt.Errorf("internal error")
		}
	} else {
		// Create new pending user
		pendingUser = &entities.PendingUser{
			Name:                      userReq.Name,
			Email:                     normalizedEmail,
			HashedPassword:            string(hashedPassword),
			PhoneNumber:               userReq.PhoneNumber,
			VerificationCode:          verificationCode,
			VerificationCodeExpiresAt: expiresAt,
		}
		if err := s.pendingUserRepo.Create(pendingUser); err != nil {
			return nil, fmt.Errorf("internal error")
		}
	}

	// Send verification email
	go s.notificationSvc.SendVerificationCode(normalizedEmail, verificationCode)

	return nil, myerrors.NewUserError("Registration pending. Please check your email for a verification code.")
}

func (s *userServiceImpl) VerifyRegistration(email, code string) (string, error) {
	normalizedEmail := utils.NormalizeEmail(email)
	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		return "", myerrors.NewUserError("Invalid email or verification code.")
	}

	if pendingUser.VerificationCode != code || time.Now().After(pendingUser.VerificationCodeExpiresAt) {
		return "", myerrors.NewUserError("Invalid or expired verification code.")
	}

	user := &entities.User{
		Name:           pendingUser.Name,
		Email:          pendingUser.Email,
		PhoneNumber:    pendingUser.PhoneNumber,
		HashedPassword: pendingUser.HashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", fmt.Errorf("internal error")
	}

	// Delete from pending users table
	_ = s.pendingUserRepo.Delete(normalizedEmail)

	// Generate JWT token for immediate login
	token, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		return "", fmt.Errorf("internal error")
	}

	return token, nil
}

func (s *userServiceImpl) AuthenticateUser(email, password string) (*entities.User, error) {
	normalizedEmail := utils.NormalizeEmail(email)
	user, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil {
		// User not found in main table, check pending users
		_, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
		if err == nil {
			return nil, myerrors.NewUserError("Registration is pending verification. Please check your email for a verification code.")
		}
		return nil, myerrors.NewUserError("Invalid email or password.")
	}

	if !user.CheckPassword(password) {
		return nil, myerrors.NewUserError("Invalid email or password.")
	}

	return user, nil
}

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
