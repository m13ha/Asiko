package services

import (
	"fmt"
	"testing"

	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name          string
		request       requests.UserRequest
		setupMock     func(mockRepo *mocks.UserRepository)
		expectedError string
	}{
		{
			name: "Success",
			request: requests.UserRequest{
				Name:        "Test User",
				Email:       "test@example.com",
				Password:    "password123",
				PhoneNumber: "1234567890",
			},
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, fmt.Errorf("not found")).Once()
				mockRepo.On("FindByPhone", "1234567890").Return(nil, fmt.Errorf("not found")).Once()
				mockRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name: "Failure - Email already exists",
			request: requests.UserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(&entities.User{}, nil).Once()
			},
			expectedError: "Email already registered.",
		},
		{
			name: "Failure - Phone already exists",
			request: requests.UserRequest{
				Name:        "Test User",
				Email:       "test@example.com",
				Password:    "password123",
				PhoneNumber: "1234567890",
			},
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, fmt.Errorf("not found")).Once()
				mockRepo.On("FindByPhone", "1234567890").Return(&entities.User{}, nil).Once()
			},
			expectedError: "Phone number already registered.",
		},
		{
			name: "Failure - Database error on create",
			request: requests.UserRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(nil, fmt.Errorf("not found")).Once()
				mockRepo.On("FindByPhone", "").Return(nil, fmt.Errorf("not found")).Maybe()
				mockRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(fmt.Errorf("db error")).Once()
			},
			expectedError: "internal error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mocks.UserRepository)
			tc.setupMock(mockUserRepo)
			userService := NewUserService(mockUserRepo)

			// Act
			userResponse, err := userService.CreateUser(tc.request)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, userResponse)
				assert.Equal(t, tc.request.Name, userResponse.Name)
			} else {
				assert.Error(t, err)
				assert.Nil(t, userResponse)
				assert.Equal(t, tc.expectedError, err.Error())
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &entities.User{
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: string(hashedPassword),
	}

	testCases := []struct {
		name          string
		email         string
		password      string
		setupMock     func(mockRepo *mocks.UserRepository)
		expectedError string
	}{
		{
			name:     "Success",
			email:    "test@example.com",
			password: "password123",
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil).Once()
			},
			expectedError: "",
		},
		{
			name:     "Failure - User not found",
			email:    "notfound@example.com",
			password: "password123",
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, errors.NewUserError("not found")).Once()
			},
			expectedError: "Invalid email or password.",
		},
		{
			name:     "Failure - Invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			setupMock: func(mockRepo *mocks.UserRepository) {
				mockRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil).Once()
			},
			expectedError: "Invalid email or password.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockUserRepo := new(mocks.UserRepository)
			tc.setupMock(mockUserRepo)
			userService := NewUserService(mockUserRepo)

			// Act
			authenticatedUser, err := userService.AuthenticateUser(tc.email, tc.password)

			// Assert
			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.NotNil(t, authenticatedUser)
			} else {
				assert.Error(t, err)
				assert.Nil(t, authenticatedUser)
				assert.Equal(t, tc.expectedError, err.Error())
			}
			mockUserRepo.AssertExpectations(t)
		})
	}
}
