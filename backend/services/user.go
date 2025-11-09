package services

import (
	"errors"
	"strings"
	"time"

	myerrors "github.com/m13ha/asiko/errors"
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
	userRepo        repository.UserRepository
	pendingUserRepo repository.PendingUserRepository
	notificationSvc notifications.NotificationService
}

func NewUserService(userRepo repository.UserRepository, pendingUserRepo repository.PendingUserRepository, notificationSvc notifications.NotificationService) UserService {
	return &userServiceImpl{userRepo: userRepo, pendingUserRepo: pendingUserRepo, notificationSvc: notificationSvc}
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
		return nil, myerrors.New(myerrors.CodeValidationFailed).WithKind(myerrors.KindValidation).WithHTTP(400).WithMessage("Invalid user data.")
	}

	normalizedEmail := utils.NormalizeEmail(userReq.Email)

	// Check if user already exists in main table
	if _, err := s.userRepo.FindByEmail(normalizedEmail); err == nil {
		return nil, myerrors.New(myerrors.CodeEmailAlreadyRegistered).WithKind(myerrors.KindConflict).WithHTTP(409).WithMessage("Email already registered.")
	} else if err != nil && !isNotFoundError(err) {
		return nil, myerrors.FromError(err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, myerrors.New(myerrors.CodeInternalError).WithKind(myerrors.KindInternal).WithHTTP(500).WithMessage("Internal error").WithCause(err)
	}

	verificationCode := utils.GenerateRandomCode(6)
	expiresAt := time.Now().Add(15 * time.Minute)

	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if isNotFoundError(err) {
			pendingUser = nil
		} else {
			return nil, myerrors.FromError(err)
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
			return nil, myerrors.FromError(err)
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
			return nil, myerrors.FromError(err)
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
		return "", myerrors.FromError(err)
	}

	if time.Now().After(pendingUser.VerificationCodeExpiresAt) {
		return "", myerrors.New(myerrors.CodeVerificationExpired).WithKind(myerrors.KindValidation).WithHTTP(400).WithMessage("Verification code expired. Request a new code.")
	}

	if pendingUser.VerificationCode != code {
		return "", myerrors.New(myerrors.CodeInvalidVerificationCode).WithKind(myerrors.KindValidation).WithHTTP(400).WithMessage("Invalid verification code.")
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
		return "", myerrors.FromError(err)
	}

	// Delete from pending users table
	_ = s.pendingUserRepo.Delete(normalizedEmail)

	// Generate JWT token for immediate login
	token, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		return "", myerrors.New(myerrors.CodeInternalError).WithKind(myerrors.KindInternal).WithHTTP(500).WithMessage("Internal error")
	}

	return token, nil
}

func (s *userServiceImpl) ResendVerificationCode(email string) error {
	normalizedEmail := utils.NormalizeEmail(email)

	if _, err := s.userRepo.FindByEmail(normalizedEmail); err == nil {
		return myerrors.New(myerrors.CodeEmailAlreadyRegistered).WithKind(myerrors.KindConflict).WithHTTP(409).WithMessage("Account already verified. Please login.")
	} else if err != nil && !isNotFoundError(err) {
		return myerrors.FromError(err)
	}

	pendingUser, err := s.pendingUserRepo.FindByEmail(normalizedEmail)
	if err != nil {
		if isNotFoundError(err) {
			return myerrors.New(myerrors.CodeResourceNotFound).WithKind(myerrors.KindNotFound).WithHTTP(404).WithMessage("No pending verification found for this email.")
		}
		return myerrors.FromError(err)
	}

	pendingUser.PhoneNumber = sanitizePhone(pendingUser.PhoneNumber)
	pendingUser.VerificationCode = utils.GenerateRandomCode(6)
	pendingUser.VerificationCodeExpiresAt = time.Now().Add(15 * time.Minute)
	if err := s.pendingUserRepo.Update(pendingUser); err != nil {
		return myerrors.FromError(err)
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
			return nil, myerrors.FromError(err)
		}

		pendingUser, pendingErr := s.pendingUserRepo.FindByEmail(normalizedEmail)
		if pendingErr != nil {
			if isNotFoundError(pendingErr) {
				return nil, myerrors.New(myerrors.CodeLoginInvalidCredentials).WithKind(myerrors.KindUnauthorized).WithHTTP(401).WithMessage("Invalid email or password.")
			}
			return nil, myerrors.FromError(pendingErr)
		}

		if pendingUser != nil {
			return nil, myerrors.New(myerrors.CodeUserPendingVerification).WithKind(myerrors.KindPrecondition).WithHTTP(202).WithMessage("Registration is pending verification. Please check your email for a verification code.")
		}
		return nil, myerrors.New(myerrors.CodeLoginInvalidCredentials).WithKind(myerrors.KindUnauthorized).WithHTTP(401).WithMessage("Invalid email or password.")
	}

	if !user.CheckPassword(password) {
		return nil, myerrors.New(myerrors.CodeLoginInvalidCredentials).WithKind(myerrors.KindUnauthorized).WithHTTP(401).WithMessage("Invalid email or password.")
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

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	if appErr := myerrors.FromError(err); appErr != nil {
		return appErr.Kind == myerrors.KindNotFound
	}
	return false
}
