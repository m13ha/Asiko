package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
	}{
		{
			name: "Success",
			body: `{"email": "test@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockUser := &entities.User{Name: "Test User", Email: "test@example.com"}
				mockService.On("AuthenticateUser", "test@example.com", "password123").Return(mockUser, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Failure - Bad Request (Invalid JSON)",
			body:               `{"email": "test@example.com"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Unauthorized",
			body: `{"email": "wrong@example.com", "password": "wrongpassword"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("AuthenticateUser", "wrong@example.com", "wrongpassword").Return(nil, errors.NewUserError("Invalid email or password.")).Once()
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			router, mockUserService, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
	}{
		{
			name: "Success",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(&responses.UserResponse{}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Failure - Bad Request",
			body:               `{"name": "New User"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Failure - Service Error",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(nil, fmt.Errorf("service error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError, // Service errors return 500
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			router, mockUserService, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}
