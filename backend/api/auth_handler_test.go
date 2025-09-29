package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/entities"
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
		expectedBody       string
	}{
		{
			name:        "Success",
			body:        `{"email": "test@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockUser := &entities.User{Name: "Test User", Email: "test@example.com"}
				mockService.On("AuthenticateUser", "test@example.com", "password123").Return(mockUser, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"token":"`,
		},
		{
			name:               "Failure - Bad Request (Invalid JSON)",
			body:               `{"email": "test@example.com"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":[`,
		},
		{
			name: "Failure - Unauthorized (Invalid Credentials)",
			body: `{"email": "wrong@example.com", "password": "wrongpassword"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("AuthenticateUser", "wrong@example.com", "wrongpassword").Return(nil, errors.NewUserError("Invalid email or password")).Once()
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid email or password"}`,
		},
		{
			name: "Failure - Pending Verification",
			body: `{"email": "pending@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("AuthenticateUser", "pending@example.com", "password123").Return(nil, errors.NewUserError("Invalid email or password")).Once()
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       `{"error":"Invalid email or password"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, mockUserService, _, _, _, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
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
		expectedBody       string
	}{
		{
			name: "Success - Pending Registration",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(nil, nil).Once()
			},
			expectedStatusCode: http.StatusAccepted,
			expectedBody:       `{"message":"Registration pending. Please check your email for a verification code."}`,
		},
		{
			name:               "Failure - Bad Request",
			body:               `{"name": "New User"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"errors":[`,
		},
		{
			name: "Failure - Service Error",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(nil, fmt.Errorf("service error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"error":"internal server error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, mockUserService, _, _, _, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestVerifyRegistrationAPI(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "Success",
			body: `{"email": "verify@example.com", "code": "123456"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("VerifyRegistration", "verify@example.com", "123456").Return("mock-jwt-token", nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `{"token":"mock-jwt-token"}`,
		},
		{
			name:               "Failure - Bad Request (Invalid JSON)",
			body:               `{invalid}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid request payload"}`,
		},
		{
			name: "Failure - Service Error",
			body: `{"email": "verify@example.com", "code": "wrong"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("VerifyRegistration", "verify@example.com", "wrong").Return("", errors.NewUserError("Invalid or expired verification code.")).Once()
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"Invalid or expired verification code."}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, mockUserService, _, _, _, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/auth/verify-registration", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)
			mockUserService.AssertExpectations(t)
		})
	}
}
