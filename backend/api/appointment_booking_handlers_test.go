package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/entities"
	"github.com/m13ha/asiko/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAppointmentAPI(t *testing.T) {
	testCases := []struct {
		name               string
		body               string
		tokenUserID        string
		setupMock          func(mockService *mocks.AppointmentService)
		expectedStatusCode int
		expectedContains   string
		expectedError      *apiErrorPayload
	}{
		{
			name:        "Success",
			body:        `{"title": "Test App", "type": "single", "start_time": "2025-01-01T10:00:00Z", "end_time": "2025-01-01T11:00:00Z", "start_date": "2025-01-01T00:00:00Z", "end_date": "2025-01-01T00:00:00Z", "booking_duration": 30, "max_attendees": 1}`,
			tokenUserID: uuid.New().String(),
			setupMock: func(mockService *mocks.AppointmentService) {
				mockService.On("CreateAppointment", mock.Anything, mock.Anything).Return(&entities.Appointment{Title: "Test App"}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
			expectedContains:   "Test App",
		},
		{
			name:               "Failure - Unauthorized (No Token)",
			body:               `{}`,
			tokenUserID:        "", // No user ID set in context
			setupMock:          func(mockService *mocks.AppointmentService) {},
			expectedStatusCode: http.StatusUnauthorized,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnauthorized,
				Code:    "AUTH_UNAUTHORIZED",
				Message: "authentication required",
			},
		},
		{
			name:               "Failure - Bad Request",
			body:               `{"title": ""}`,
			tokenUserID:        uuid.New().String(),
			setupMock:          func(mockService *mocks.AppointmentService) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedError: &apiErrorPayload{
				Status:  http.StatusUnprocessableEntity,
				Code:    "VALIDATION_FAILED",
				Message: "validation failed", // note: this is the new error message from the updated API layer
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAppointmentService := new(mocks.AppointmentService)
			tc.setupMock(mockAppointmentService)

			handler := NewHandler(nil, mockAppointmentService, nil, nil, nil, nil)
			router := gin.New()
			router.Use(middleware.RequestID())
			router.Use(gin.Recovery())
			router.Use(middleware.ErrorHandler())
			router.POST("/appointments", func(c *gin.Context) {
				if tc.tokenUserID != "" {
					if uid, err := uuid.Parse(tc.tokenUserID); err == nil {
						c.Set("userID", tc.tokenUserID)
						c.Set("userUUID", uid)
					}
				}
				handler.CreateAppointment(c)
			})

			req, _ := http.NewRequest("POST", "/appointments", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
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
			mockAppointmentService.AssertExpectations(t)
		})
	}
}

func TestCancelBookingAPI(t *testing.T) {
	testCases := []struct {
		name               string
		bookingCode        string
		setupMock          func(mockService *mocks.BookingService)
		expectedStatusCode int
		expectedContains   string
		expectedError      *apiErrorPayload
	}{
		{
			name:        "Success",
			bookingCode: "BK123XYZ",
			setupMock: func(mockService *mocks.BookingService) {
				mockService.On("CancelBookingByCode", "BK123XYZ").Return(&entities.Booking{Status: entities.BookingStatusCancelled}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedContains:   fmt.Sprintf(`"status":"%s"`, entities.BookingStatusCancelled),
		},
		{
			name:        "Failure - Booking Not Found",
			bookingCode: "NOTFOUND",
			setupMock: func(mockService *mocks.BookingService) {
				// Mock service returning an error that will trigger the API error handling
				mockService.On("CancelBookingByCode", "NOTFOUND").Return((*entities.Booking)(nil), apperrors.NewAppError(apperrors.CodeBookingNotFound, "resource_not_found", http.StatusNotFound, "booking not found", nil)).Once()
			},
			expectedStatusCode: http.StatusNotFound,
			expectedError: &apiErrorPayload{
				Status:  http.StatusNotFound,
				Code:    "RESOURCE_NOT_FOUND",
				Message: "booking not found", // This message is from the updated API layer error handling
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			router, _, _, mockBookingService, _, _, _ := setupTestRouter()
			tc.setupMock(mockBookingService)

			req, _ := http.NewRequest("DELETE", "/bookings/"+tc.bookingCode, nil)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedError != nil {
				resp := decodeAPIError(t, w.Body.Bytes())
				assert.Equal(t, tc.expectedError.Code, resp.Code)
				assert.Equal(t, tc.expectedError.Message, resp.Message)
			} else if tc.expectedContains != "" {
				assert.Contains(t, w.Body.String(), tc.expectedContains)
			}
			mockBookingService.AssertExpectations(t)
		})
	}
}
