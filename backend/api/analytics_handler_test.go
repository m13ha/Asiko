package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUserAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name               string
		queryParams        string
		userID             string
		setupMock          func(mockService *mocks.AnalyticsService, userID string)
		expectedStatusCode int
		expectedContains   string
		expectedError      *apiErrorPayload
	}{
		{
			name:        "Success",
			queryParams: "?start_date=2025-01-01&end_date=2025-01-31",
			userID:      uuid.New().String(),
			setupMock: func(mockService *mocks.AnalyticsService, userID string) {
				userUUID, _ := uuid.Parse(userID)
				mockService.On("GetUserAnalytics", userUUID, "2025-01-01", "2025-01-31").Return(&responses.AnalyticsResponse{
					TotalAppointments: 5,
					TotalBookings:     12,
					StartDate:         time.Now(),
					EndDate:           time.Now(),
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedContains:   `"total_appointments":5`,
		},
		{
			name:               "Missing query parameters",
			queryParams:        "?start_date=2025-01-01",
			userID:             uuid.New().String(),
			setupMock:          func(mockService *mocks.AnalyticsService, userID string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedError: &apiErrorPayload{
				Status:  http.StatusBadRequest,
				Code:    apperrors.CodeValidationFailed,
				Message: "start_date and end_date query parameters are required",
			},
		},
		{
			name:        "Service error",
			queryParams: "?start_date=2025-01-01&end_date=2025-01-31",
			userID:      uuid.New().String(),
			setupMock: func(mockService *mocks.AnalyticsService, userID string) {
				userUUID, _ := uuid.Parse(userID)
				err := apperrors.New(apperrors.CodeInternalError).
					WithKind(apperrors.KindInternal).
					WithHTTP(http.StatusInternalServerError).
					WithMessage("Internal server error")
				mockService.On("GetUserAnalytics", userUUID, "2025-01-01", "2025-01-31").Return((*responses.AnalyticsResponse)(nil), err)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedError: &apiErrorPayload{
				Status:  http.StatusInternalServerError,
				Code:    apperrors.CodeInternalError,
				Message: "Internal server error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAnalyticsService := new(mocks.AnalyticsService)
			tc.setupMock(mockAnalyticsService, tc.userID)

			handler := &Handler{
				analyticsService: mockAnalyticsService,
			}

			router := gin.New()
			router.Use(middleware.RequestID())
			router.Use(gin.Recovery())
			router.Use(middleware.ErrorHandler())
			router.GET("/analytics", func(c *gin.Context) {
				if tc.userID != "" {
					c.Set("userID", tc.userID)
					if uid, err := uuid.Parse(tc.userID); err == nil {
						c.Set("userUUID", uid)
					}
				}
				handler.GetUserAnalytics(c)
			})

			req, _ := http.NewRequest("GET", "/analytics"+tc.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				assert.Equal(t, tc.expectedError.Message, resp.Message)
			} else if tc.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tc.expectedContains)
			}
			mockAnalyticsService.AssertExpectations(t)
		})
	}
}
