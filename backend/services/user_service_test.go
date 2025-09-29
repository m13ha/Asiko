package services

import (
	"testing"
	"time"

	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/notifications/mocks"
	repoMocks "github.com/m13ha/appointment_master/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		request       requests.UserRequest
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository, *mocks.NotificationService)
		expectedError string
	}{
		{
			name: "Failure - Email already registered",
			request: requests.UserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) {
				userRepo.On("FindByEmail", "test@example.com").Return(&entities.User{}, nil).Once()
			},
			expectedError: "Email already registered.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			mockNotificationSvc := new(mocks.NotificationService)
			tc.setupMocks(mockUserRepo, mockPendingRepo, mockNotificationSvc)

			userService := NewUserService(mockUserRepo, mockPendingRepo, mockNotificationSvc)
			_, err := userService.CreateUser(tc.request)

			assert.Error(t, err)
			assert.Equal(t, tc.expectedError, err.Error())

			mockUserRepo.AssertExpectations(t)
			mockPendingRepo.AssertExpectations(t)
			mockNotificationSvc.AssertExpectations(t)
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &entities.User{
		Email:          "test@example.com",
		HashedPassword: string(hashedPassword),
	}

	testCases := []struct {
		name          string
		email         string
		password      string
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository)
		expectedError string
	}{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password123",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				userRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil).Once()
			},
			expectedError: "",
		},
		{
			name:     "Failure - Pending Verification",
			email:    "pending@example.com",
			password: "password123",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				userRepo.On("FindByEmail", "pending@example.com").Return(nil, gorm.ErrRecordNotFound).Once()
				pendingRepo.On("FindByEmail", "pending@example.com").Return(&entities.PendingUser{}, nil).Once()
			},
			expectedError: "Registration is pending verification. Please check your email for a verification code.",
		},
		{
			name:     "Failure - User not found",
			email:    "notfound@example.com",
			password: "password123",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				userRepo.On("FindByEmail", "notfound@example.com").Return(nil, gorm.ErrRecordNotFound).Once()
				pendingRepo.On("FindByEmail", "notfound@example.com").Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedError: "Invalid email or password.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			tc.setupMocks(mockUserRepo, mockPendingRepo)

			userService := NewUserService(mockUserRepo, mockPendingRepo, nil)
			_, err := userService.AuthenticateUser(tc.email, tc.password)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVerifyRegistration(t *testing.T) {
	mockPendingUser := &entities.PendingUser{
		Email:                     "test@example.com",
		VerificationCode:          "123456",
		VerificationCodeExpiresAt: time.Now().Add(15 * time.Minute),
	}

	testCases := []struct {
		name          string
		email         string
		code          string
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository)
		expectedError bool
	}{
		{
			name:  "Success",
			email: "test@example.com",
			code:  "123456",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				pendingRepo.On("FindByEmail", "test@example.com").Return(mockPendingUser, nil).Once()
				userRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(nil).Once()
				pendingRepo.On("Delete", "test@example.com").Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name:  "Failure - Invalid Code",
			email: "test@example.com",
			code:  "654321",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				pendingRepo.On("FindByEmail", "test@example.com").Return(mockPendingUser, nil).Once()
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			tc.setupMocks(mockUserRepo, mockPendingRepo)

			userService := NewUserService(mockUserRepo, mockPendingRepo, nil)
			_, err := userService.VerifyRegistration(tc.email, tc.code)

			assert.Equal(t, tc.expectedError, err != nil)
		})
	}
}
