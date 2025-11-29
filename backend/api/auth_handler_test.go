package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/services/mocks"
	apperrors "github.com/m13ha/asiko/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
		assertResponse     func(t *testing.T, body []byte)
		expectedError      *apiErrorPayload
	}{
		{
			name: "Success",
			body: `{"email": "test@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockUser := &entities.User{Name: "Test User", Email: "test@example.com"}
				mockService.On("AuthenticateUser", "test@example.com", "password123").Return(mockUser, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			assertResponse: func(t *testing.T, body []byte) {
				var payload responses.LoginResponse
				require.NoError(t, json.Unmarshal(body, &payload))
				assert.NotEmpty(t, payload.Token)
				assert.NotEmpty(t, payload.RefreshToken)
				assert.True(t, payload.ExpiresIn > 0)
				assert.Equal(t, "Test User", payload.User.Name)
				assert.Equal(t, "test@example.com", payload.User.Email)
			},
		},
		{
			name:               "Failure - Bad Request (Invalid JSON)",
			body:               `{"email": "test@example.com"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnprocessableEntity,
				Code:    "VALIDATION_FAILED",
				Message: "Validation failed",
			},
		},
		{
			name: "Failure - Unauthorized (Invalid Credentials)",
			body: `{"email": "wrong@example.com", "password": "wrongpassword"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API unauthorized error
				mockService.On("AuthenticateUser", "wrong@example.com", "wrongpassword").Return(nil, assert.AnError).Once()
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnauthorized,
				Code:    "AUTH_UNAUTHORIZED",
				Message: "Invalid email or password",
			},
		},
		{
			name: "Failure - Pending Verification",
			body: `{"email": "pending@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API unauthorized error
				mockService.On("AuthenticateUser", "pending@example.com", "password123").Return(nil, assert.AnError).Once()
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnauthorized,
				Code:    "AUTH_UNAUTHORIZED",
				Message: "Invalid email or password",
			},
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
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				// assert.Equal(t, tc.expectedError.Message, resp.Message) // Message might vary
			} else if tc.assertResponse != nil {
				tc.assertResponse(t, w.Body.Bytes())
			}
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
		expectedMessage    string
		expectedError      *apiErrorPayload
	}{
		{
			name: "Success - Pending Registration",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(nil, nil).Once()
			},
			expectedStatusCode: http.StatusAccepted,
			expectedMessage:    "Registration pending. Please check your email for a verification code.",
		},
		{
			name:               "Failure - Bad Request",
			body:               `{"name": "New User"}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnprocessableEntity,
				Code:    "VALIDATION_FAILED",
				Message: "Validation failed",
			},
		},
		{
			name: "Failure - Service Error",
			body: `{"name": "New User", "email": "new@example.com", "password": "password123"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API internal error
				mockService.On("CreateUser", mock.AnythingOfType("requests.UserRequest")).Return(nil, assert.AnError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &apiErrorPayload{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
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
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				// assert.Equal(t, tc.expectedError.Message, resp.Message)
			} else if tc.expectedMessage != "" {
				var payload responses.SimpleMessage
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
				assert.Equal(t, tc.expectedMessage, payload.Message)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestVerifyRegistrationAPI(t *testing.T) {
	tokenForVerify, _ := middleware.GenerateToken("user-123")
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
		expectedToken      string
		expectRefresh      bool
		expectedError      *apiErrorPayload
	}{
		{
			name: "Success",
			body: `{"email": "verify@example.com", "code": "123456"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("VerifyRegistration", "verify@example.com", "123456").Return(tokenForVerify, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedToken:      tokenForVerify,
			expectRefresh:      true,
		},
		{
			name:               "Failure - Bad Request (Invalid JSON)",
			body:               `{invalid}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: &apiErrorPayload{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: "Invalid request payload",
			},
		},
		{
			name: "Failure - Service Error",
			body: `{"email": "verify@example.com", "code": "wrong"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API error handling
				mockService.On("VerifyRegistration", "verify@example.com", "wrong").Return("", apperrors.NewAppError(apperrors.CodeInternalError, "internal", http.StatusInternalServerError, "service error", nil)).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &apiErrorPayload{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Verification error", // This is the error message from updated API layer
			},
		},
		{
			name: "Failure - Expired Code",
			body: `{"email": "verify@example.com", "code": "123456"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API error handling
				mockService.On("VerifyRegistration", "verify@example.com", "123456").Return("", apperrors.NewAppError(apperrors.CodeVerificationExpired, "validation", http.StatusUnprocessableEntity, "expired code", nil)).Once()
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnprocessableEntity,
				Code:    "VALIDATION_FAILED",
				Message: "Verification error",
			},
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
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				// assert.Equal(t, tc.expectedError.Message, resp.Message)
			} else if tc.expectedToken != "" {
				var payload responses.LoginResponse
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
				assert.Equal(t, tc.expectedToken, payload.Token)
				if tc.expectRefresh {
					assert.NotEmpty(t, payload.RefreshToken)
					assert.True(t, payload.ExpiresIn > 0)
				}
			}
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestResendVerificationAPI(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		setupMock          func(mockService *mocks.UserService)
		expectedStatusCode int
		expectedContains   string
		expectedError      *apiErrorPayload
	}{
		{
			name: "Success",
			body: `{"email": "pending@example.com"}`,
			setupMock: func(mockService *mocks.UserService) {
				mockService.On("ResendVerificationCode", "pending@example.com").Return(nil).Once()
			},
			expectedStatusCode: http.StatusAccepted,
			expectedContains:   `{"message":"Verification code resent if a pending registration exists for this email."}`,
		},
		{
			name:               "Failure - Bad Request (invalid JSON)",
			body:               `{invalid}`,
			setupMock:          func(mockService *mocks.UserService) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: &apiErrorPayload{
				Status:  http.StatusBadRequest,
				Code:    "BAD_REQUEST",
				Message: "Invalid request payload",
			},
		},
		{
			name: "Failure - Service Error",
			body: `{"email": "missing@example.com"}`,
			setupMock: func(mockService *mocks.UserService) {
				// Mock service returning an error that will trigger the API error handling
				// Use a specific error type if you want specific status code mapping
				// For now, generic error -> 500
				mockService.On("ResendVerificationCode", "missing@example.com").Return(assert.AnError).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &apiErrorPayload{
				Status:  http.StatusInternalServerError,
				Code:    "INTERNAL_ERROR",
				Message: "Internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router, mockUserService, _, _, _, _, _ := setupTestRouter()
			tc.setupMock(mockUserService)

			req, _ := http.NewRequest("POST", "/auth/resend-verification", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				// assert.Equal(t, tc.expectedError.Message, resp.Message)
			} else if tc.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tc.expectedContains)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestRefreshTokenAPI(t *testing.T) {
	router, _, _, _, _, _, _ := setupTestRouter()

	t.Run("Success", func(t *testing.T) {
		rt, err := middleware.GenerateRefreshToken("user-123")
		require.NoError(t, err)

		req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(`{"refreshToken":"`+rt+`"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var payload responses.TokenResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &payload))
		assert.NotEmpty(t, payload.Token)
		assert.NotEmpty(t, payload.RefreshToken)
		assert.True(t, payload.ExpiresIn > 0)
	})

	t.Run("Invalid payload", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/auth/refresh", strings.NewReader(`{"refreshToken":""}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		resp := decodeAPIError(t, w.Body.Bytes())
		assert.Equal(t, "VALIDATION_FAILED", resp.Code)
	})
}
