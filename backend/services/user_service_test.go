package services_test

import (
	"testing"
	"time"

	myerrors "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/notifications/mocks"
	repoMocks "github.com/m13ha/asiko/repository/mocks"
	services "github.com/m13ha/asiko/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		request       requests.UserRequest
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository, *mocks.NotificationService) func(*testing.T)
		expectedError string
	}{
		{
			name: "Failure - Email already registered",
			request: requests.UserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) func(*testing.T) {
				userRepo.On("FindByEmail", "test@example.com").Return(&entities.User{}, nil).Once()
				return nil
			},
			expectedError: "EMAIL_ALREADY_REGISTERED: Email already registered.",
		},
		{
			name: "Success - Creates pending user when none exists",
			request: requests.UserRequest{
				Name:     "New User",
				Email:    "new@example.com",
				Password: "password123",
			},
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) func(*testing.T) {
				userRepo.On("FindByEmail", "new@example.com").
					Return((*entities.User)(nil), repoNotFoundError()).Once()
				pendingRepo.On("FindByEmail", "new@example.com").
					Return((*entities.PendingUser)(nil), repoNotFoundError()).Once()
				pendingRepo.On("Create", mock.AnythingOfType("*entities.PendingUser")).Return(nil).Once()

				done := make(chan struct{}, 1)
				notificationSvc.On("SendVerificationCode", "new@example.com", mock.AnythingOfType("string")).
					Return(nil).
					Run(func(args mock.Arguments) {
						select {
						case done <- struct{}{}:
						default:
						}
					}).Once()

				return func(t *testing.T) {
					select {
					case <-done:
					case <-time.After(100 * time.Millisecond):
						t.Fatal("expected SendVerificationCode to be called")
					}
				}
			},
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			mockNotificationSvc := new(mocks.NotificationService)
			waitFn := tc.setupMocks(mockUserRepo, mockPendingRepo, mockNotificationSvc)

			userService := services.NewUserService(mockUserRepo, mockPendingRepo, mockNotificationSvc)
			resp, err := userService.CreateUser(tc.request)

			if waitFn != nil {
				waitFn(t)
			}

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				if assert.NotNil(t, resp) {
					assert.Equal(t, tc.request.Name, resp.Name)
					assert.Equal(t, tc.request.Email, resp.Email)
					assert.Nil(t, resp.PhoneNumber)
				}
			}

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
				userRepo.On("FindByEmail", "pending@example.com").Return(nil, repoNotFoundError()).Once()
				pendingRepo.On("FindByEmail", "pending@example.com").Return(&entities.PendingUser{}, nil).Once()
			},
			expectedError: "USER_PENDING_VERIFICATION: Registration is pending verification. Please check your email for a verification code.",
		},
		{
			name:     "Failure - User not found",
			email:    "notfound@example.com",
			password: "password123",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				userRepo.On("FindByEmail", "notfound@example.com").Return(nil, repoNotFoundError()).Once()
				pendingRepo.On("FindByEmail", "notfound@example.com").Return(nil, repoNotFoundError()).Once()
			},
			expectedError: "LOGIN_INVALID_CREDENTIALS: Invalid email or password.",
		},
		{
			name:     "Failure - Repository error surfaces",
			email:    "error@example.com",
			password: "password123",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				repoErr := myerrors.New(myerrors.CodeInternalError).WithKind(myerrors.KindInternal).WithHTTP(500).WithMessage("repository failure")
				userRepo.On("FindByEmail", "error@example.com").Return(nil, repoErr).Once()
			},
			expectedError: "INTERNAL_ERROR: repository failure",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			tc.setupMocks(mockUserRepo, mockPendingRepo)

			userService := services.NewUserService(mockUserRepo, mockPendingRepo, nil)
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

func repoNotFoundError() error {
	return myerrors.New(myerrors.CodeResourceNotFound).
		WithKind(myerrors.KindNotFound).
		WithHTTP(404).
		WithMessage("Resource not found")
}

func TestVerifyRegistration(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		code          string
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository)
		expectedError string
	}{
		{
			name:  "Success",
			email: "test@example.com",
			code:  "123456",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				mockPendingUser := &entities.PendingUser{
					Email:                     "test@example.com",
					VerificationCode:          "123456",
					VerificationCodeExpiresAt: time.Now().Add(15 * time.Minute),
					Name:                      "Test User",
					HashedPassword:            "hashed",
				}
				pendingRepo.On("FindByEmail", "test@example.com").Return(mockPendingUser, nil).Once()
				userRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(nil).Once()
				pendingRepo.On("Delete", "test@example.com").Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name:  "Success - Empty phone sanitized",
			email: "phone@example.com",
			code:  "123456",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				empty := ""
				mockPendingUser := &entities.PendingUser{
					Email:                     "phone@example.com",
					VerificationCode:          "123456",
					VerificationCodeExpiresAt: time.Now().Add(15 * time.Minute),
					Name:                      "Phone User",
					HashedPassword:            "hashed",
					PhoneNumber:               &empty,
				}
				pendingRepo.On("FindByEmail", "phone@example.com").Return(mockPendingUser, nil).Once()
				userRepo.On("Create", mock.MatchedBy(func(u *entities.User) bool {
					return u.Email == "phone@example.com" && u.PhoneNumber == nil
				})).Return(nil).Once()
				pendingRepo.On("Delete", "phone@example.com").Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name:  "Failure - Invalid Code",
			email: "test@example.com",
			code:  "654321",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				mockPendingUser := &entities.PendingUser{
					Email:                     "test@example.com",
					VerificationCode:          "123456",
					VerificationCodeExpiresAt: time.Now().Add(15 * time.Minute),
				}
				pendingRepo.On("FindByEmail", "test@example.com").Return(mockPendingUser, nil).Once()
			},
			expectedError: "INVALID_VERIFICATION_CODE: Invalid verification code.",
		},
		{
			name:  "Failure - Code Expired",
			email: "test@example.com",
			code:  "123456",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository) {
				mockPendingUser := &entities.PendingUser{
					Email:                     "test@example.com",
					VerificationCode:          "123456",
					VerificationCodeExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				pendingRepo.On("FindByEmail", "test@example.com").Return(mockPendingUser, nil).Once()
			},
			expectedError: "VERIFICATION_EXPIRED: Verification code expired. Request a new code.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo := new(repoMocks.UserRepository)
			mockPendingRepo := new(repoMocks.PendingUserRepository)
			tc.setupMocks(mockUserRepo, mockPendingRepo)

			userService := services.NewUserService(mockUserRepo, mockPendingRepo, nil)
			_, err := userService.VerifyRegistration(tc.email, tc.code)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResendVerificationCode(t *testing.T) {
	testCases := []struct {
		name          string
		email         string
		setupMocks    func(*repoMocks.UserRepository, *repoMocks.PendingUserRepository, *mocks.NotificationService) func(*testing.T)
		expectedError string
	}{
		{
			name:  "Success - Pending user found",
			email: "pending@example.com",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) func(*testing.T) {
				userRepo.On("FindByEmail", "pending@example.com").Return((*entities.User)(nil), repoNotFoundError()).Once()
				pendingUser := &entities.PendingUser{
					Email:                     "pending@example.com",
					VerificationCode:          "000000",
					VerificationCodeExpiresAt: time.Now().Add(-1 * time.Minute),
				}
				pendingRepo.On("FindByEmail", "pending@example.com").Return(pendingUser, nil).Once()
				pendingRepo.On("Update", mock.AnythingOfType("*entities.PendingUser")).Return(nil).Once()
				done := make(chan struct{}, 1)
				notificationSvc.On("SendVerificationCode", "pending@example.com", mock.AnythingOfType("string")).Return(nil).Run(func(args mock.Arguments) {
					select {
					case done <- struct{}{}:
					default:
					}
				}).Once()
				return func(t *testing.T) {
					select {
					case <-done:
					case <-time.After(100 * time.Millisecond):
						t.Fatal("expected SendVerificationCode to be called")
					}
				}
			},
			expectedError: "",
		},
		{
			name:  "Failure - Already verified",
			email: "existing@example.com",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) func(*testing.T) {
				userRepo.On("FindByEmail", "existing@example.com").Return(&entities.User{}, nil).Once()
				return nil
			},
			expectedError: "EMAIL_ALREADY_REGISTERED: Account already verified. Please login.",
		},
		{
			name:  "Failure - Pending not found",
			email: "missing@example.com",
			setupMocks: func(userRepo *repoMocks.UserRepository, pendingRepo *repoMocks.PendingUserRepository, notificationSvc *mocks.NotificationService) func(*testing.T) {
				userRepo.On("FindByEmail", "missing@example.com").Return((*entities.User)(nil), repoNotFoundError()).Once()
				pendingRepo.On("FindByEmail", "missing@example.com").Return((*entities.PendingUser)(nil), repoNotFoundError()).Once()
				return nil
			},
			expectedError: "RESOURCE_NOT_FOUND: No pending verification found for this email.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := new(repoMocks.UserRepository)
			pendingRepo := new(repoMocks.PendingUserRepository)
			notificationSvc := new(mocks.NotificationService)
			waitFn := tc.setupMocks(userRepo, pendingRepo, notificationSvc)

			service := services.NewUserService(userRepo, pendingRepo, notificationSvc)
			err := service.ResendVerificationCode(tc.email)

			if waitFn != nil {
				waitFn(t)
			}

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			pendingRepo.AssertExpectations(t)
			notificationSvc.AssertExpectations(t)
		})
	}
}
