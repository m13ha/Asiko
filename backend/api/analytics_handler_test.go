package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/services/mocks"
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
		},
		{
			name:               "Missing query parameters",
			queryParams:        "?start_date=2025-01-01",
			userID:             uuid.New().String(),
			setupMock:          func(mockService *mocks.AnalyticsService, userID string) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:        "Service error",
			queryParams: "?start_date=2025-01-01&end_date=2025-01-31",
			userID:      uuid.New().String(),
			setupMock: func(mockService *mocks.AnalyticsService, userID string) {
				userUUID, _ := uuid.Parse(userID)
				mockService.On("GetUserAnalytics", userUUID, "2025-01-01", "2025-01-31").Return((*responses.AnalyticsResponse)(nil), fmt.Errorf("service error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
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
			router.GET("/analytics", func(c *gin.Context) {
				c.Set("userID", tc.userID)
				handler.GetUserAnalytics(c)
			})

			req, _ := http.NewRequest("GET", "/analytics"+tc.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockAnalyticsService.AssertExpectations(t)
		})
	}
}
