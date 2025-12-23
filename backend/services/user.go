package services

import (
	"errors"
	"strings"
	"time"

	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/notifications"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepo          repository.UserRepository
	pendingUserRepo   repository.PendingUserRepository
	passwordResetRepo repository.PasswordResetRepository
	notificationSvc   notifications.NotificationService
}

func NewUserService(userRepo repository.UserRepository, pendingUserRepo repository.PendingUserRepository, passwordResetRepo repository.PasswordResetRepository, notificationSvc notifications.NotificationService) UserService {
	return &userServiceImpl{userRepo: userRepo, pendingUserRepo: pendingUserRepo, passwordResetRepo: passwordResetRepo, notificationSvc: notificationSvc}
}

func sanitizePhone(phone *string) *string {
	if phone == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*phone)
	if trimmed == "" {
		return nil
	}

	value := trimmed
	return &value
}

func (s *userServiceImpl) CreateUser(userReq requests.UserRequest) (*responses.UserResponse, error) {
	if err := utils.Validate(userReq); err != nil {
		return nil, serviceerrors.ValidationError("Invalid user data.")
	}

	normalizedEmail := utils.NormalizeEmail(userReq.Email)

	// Check if user already exists in main table
	if _, err := s.userRepo.FindByEmail(normalizedEmail); err == nil {
		return nil, serviceerrors.EmailAlreadyRegisteredError("Email already registered.")
	} else if err != nil && !isNotFoundError(err) {
		return nil, serviceerrors.FromError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, serviceerrors.InternalError("Internal error")
	}

	verificationCode := utils.GenerateRandomCode(6)
	expiresAt := time.Now().Add(15 * time.Minute)

	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if isNotFoundError(err) {
			pendingUser = nil
		} else {
			return nil, serviceerrors.FromError(err)
		}
	}

	if pendingUser != nil {
		pendingUser.PhoneNumber = sanitizePhone(pendingUser.PhoneNumber)
		// Update existing pending user
		pendingUser.Name = userReq.Name
		pendingUser.HashedPassword = string(hashedPassword)
		pendingUser.VerificationCode = verificationCode
		pendingUser.VerificationCodeExpiresAt = expiresAt
		if err := s.pendingUserRepo.Update(pendingUser); err != nil {
			return nil, serviceerrors.FromError(err)
		}
	} else {
		// Create new pending user
		pendingUser = &entities.PendingUser{
			Name:                      userReq.Name,
			Email:                     normalizedEmail,
			HashedPassword:            string(hashedPassword),
			VerificationCode:          verificationCode,
			VerificationCodeExpiresAt: expiresAt,
		}
		if err := s.pendingUserRepo.Create(pendingUser); err != nil {
			return nil, serviceerrors.FromError(err)
		}
	}

	// Send verification email asynchronously with error logging
	go func(email, code string) {
		if err := s.notificationSvc.SendVerificationCode(email, code); err != nil {
			log.Error().Err(err).Str("email", email).Msg("notifications: failed to send verification email")
		} else {
			log.Info().Str("email", email).Msg("notifications: verification email queued/sent")
		}
	}(normalizedEmail, verificationCode)

	return &responses.UserResponse{
		ID:          pendingUser.ID,
		Name:        pendingUser.Name,
		Email:       pendingUser.Email,
		PhoneNumber: pendingUser.PhoneNumber,
		CreatedAt:   pendingUser.CreatedAt,
	}, nil
}

func (s *userServiceImpl) VerifyRegistration(email, code string) (string, error) {
	normalizedEmail := utils.NormalizeEmail(email)
	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		return "", serviceerrors.FromError(err)
	}

	if time.Now().After(pendingUser.VerificationCodeExpiresAt) {
		return "", serviceerrors.VerificationExpiredError("Verification code expired. Request a new code.")
	}

	if pendingUser.VerificationCode != code {
		return "", serviceerrors.InvalidVerificationCodeError("Invalid verification code.")
	}

	cleanPhone := sanitizePhone(pendingUser.PhoneNumber)
	pendingUser.PhoneNumber = cleanPhone

	user := &entities.User{
		Name:           pendingUser.Name,
		Email:          pendingUser.Email,
		PhoneNumber:    cleanPhone,
		HashedPassword: pendingUser.HashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", serviceerrors.FromError(err)
	}

	// Delete from pending users table
	_ = s.pendingUserRepo.Delete(normalizedEmail)

	// Generate JWT token for immediate login
	token, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		return "", serviceerrors.InternalError("Internal error")
	}

	return token, nil
}

func (s *userServiceImpl) ResendVerificationCode(email string) error {
	normalizedEmail := utils.NormalizeEmail(email)

	if _, err := s.userRepo.FindByEmail(normalizedEmail); err == nil {
		return serviceerrors.EmailAlreadyRegisteredError("Account already verified. Please login.")
	} else if err != nil && !isNotFoundError(err) {
		return serviceerrors.FromError(err)
	}

	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if isNotFoundError(err) {
			return serviceerrors.NotFoundError("No pending verification found for this email.")
		}
		return serviceerrors.FromError(err)
	}

	pendingUser.PhoneNumber = sanitizePhone(pendingUser.PhoneNumber)
	pendingUser.VerificationCode = utils.GenerateRandomCode(6)
	pendingUser.VerificationCodeExpiresAt = time.Now().Add(15 * time.Minute)
	if err := s.pendingUserRepo.Update(pendingUser); err != nil {
		return serviceerrors.FromError(err)
	}

	go func(email, code string) {
		if err := s.notificationSvc.SendVerificationCode(email, code); err != nil {
			log.Error().Err(err).Str("email", email).Msg("notifications: failed to resend verification email")
		} else {
			log.Info().Str("email", email).Msg("notifications: verification email resent")
		}
	}(normalizedEmail, pendingUser.VerificationCode)

	return nil
}

func (s *userServiceImpl) AuthenticateUser(email, password string) (*entities.User, error) {
	normalizedEmail := utils.NormalizeEmail(email)
	user, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if !isNotFoundError(err) {
			return nil, serviceerrors.FromError(err)
		}

		pendingUser, pendingErr := s.pendingUserRepo.FindByEmail(normalizedEmail)
		if pendingErr != nil {
			if isNotFoundError(pendingErr) {
				return nil, serviceerrors.LoginInvalidCredentialsError("Invalid email or password.")
			}
			return nil, serviceerrors.FromError(pendingErr)
		}

		if pendingUser != nil {
			return nil, serviceerrors.UserPendingVerificationError("Registration is pending verification. Please check your email for a verification code.")
		}
		return nil, serviceerrors.LoginInvalidCredentialsError("Invalid email or password.")
	}

	if !user.CheckPassword(password) {
		return nil, serviceerrors.LoginInvalidCredentialsError("Invalid email or password.")
	}

	return user, nil
}

func (s *userServiceImpl) ForgotPassword(email string) error {
	normalizedEmail := utils.NormalizeEmail(email)
	user, err := s.userRepo.FindByEmail(normalizedEmail)
	if err != nil {
		// If user not found, we should not reveal it for security reasons,
		// but for now we'll just return nil to simulate success
		if isNotFoundError(err) {
			return nil
		}
		return serviceerrors.FromError(err)
	}

	// Generate reset token
	token := utils.GenerateRandomCode(6) // Using 6 digit code for simplicity, could be UUID
	expiresAt := time.Now().Add(15 * time.Minute)

	resetToken := &entities.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}

	// Invalidate previous tokens
	if err := s.passwordResetRepo.DeleteAllForUser(user.ID.String()); err != nil {
		return serviceerrors.FromError(err)
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return serviceerrors.FromError(err)
	}

	// Send email asynchronously
	go func(email, code string) {
		if err := s.notificationSvc.SendPasswordResetEmail(email, code); err != nil {
			log.Error().Err(err).Str("email", email).Msg("notifications: failed to send password reset email")
		} else {
			log.Info().Str("email", email).Msg("notifications: password reset email sent")
		}
	}(normalizedEmail, token)

	return nil
}

func (s *userServiceImpl) ResetPassword(token, newPassword string) error {
	resetToken, err := s.passwordResetRepo.FindByToken(token)
	if err != nil {
		if isNotFoundError(err) {
			return serviceerrors.ValidationError("Invalid or expired reset token.")
		}
		return serviceerrors.FromError(err)
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return serviceerrors.ValidationError("Reset token has expired.")
	}

	user, err := s.userRepo.FindByID(resetToken.UserID.String())
	if err != nil {
		return serviceerrors.FromError(err)
	}

	if err := user.SetPassword(newPassword); err != nil {
		return serviceerrors.InternalError("Failed to set new password")
	}

	if err := s.userRepo.Update(user); err != nil {
		return serviceerrors.FromError(err)
	}

	// Invalidate used token
	_ = s.passwordResetRepo.DeleteAllForUser(user.ID.String())

	return nil
}

func (s *userServiceImpl) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return serviceerrors.FromError(err)
	}

	if !user.CheckPassword(oldPassword) {
		return serviceerrors.ValidationError("Incorrect old password.")
	}

	if err := user.SetPassword(newPassword); err != nil {
		return serviceerrors.InternalError("Failed to set new password")
	}

	if err := s.userRepo.Update(user); err != nil {
		return serviceerrors.FromError(err)
	}

	return nil
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

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if serviceErr := serviceerrors.FromError(err); serviceErr != nil {
		// Check for new error system
		return serviceErr.Kind == "not_found"
	}
	return false
}
