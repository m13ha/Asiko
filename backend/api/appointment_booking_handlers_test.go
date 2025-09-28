package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/services/mocks"
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
	}{
		{
			name:        "Success",
			body:        `{"title": "Test App", "type": "single", "start_time": "2025-01-01T10:00:00Z", "end_time": "2025-01-01T11:00:00Z", "start_date": "2025-01-01T00:00:00Z", "end_date": "2025-01-01T00:00:00Z", "booking_duration": 30, "max_attendees": 1}`,
			tokenUserID: uuid.New().String(),
			setupMock: func(mockService *mocks.AppointmentService) {
				mockService.On("CreateAppointment", mock.Anything, mock.Anything).Return(&entities.Appointment{Title: "Test App"}, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Failure - Unauthorized (No Token)",
			body:               `{}`,
			tokenUserID:        "", // No user ID set in context
			setupMock:          func(mockService *mocks.AppointmentService) {},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Failure - Bad Request",
			body:               `{"title": ""}`,
			tokenUserID:        uuid.New().String(),
			setupMock:          func(mockService *mocks.AppointmentService) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			_, _, mockAppointmentService, _ := setupTestRouter()
			tc.setupMock(mockAppointmentService)

			req, _ := http.NewRequest("POST", "/appointments", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			// Simulate auth middleware by setting (or not setting) the userID
			if tc.tokenUserID != "" {
				ctx.Set("userID", tc.tokenUserID)
			}

			// Act
			h := NewHandler(nil, mockAppointmentService, nil, nil, nil)
			h.CreateAppointment(ctx)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
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
	}{
		{
			name:        "Success",
			bookingCode: "BK123XYZ",
			setupMock: func(mockService *mocks.BookingService) {
				mockService.On("CancelBookingByCode", "BK123XYZ").Return(&entities.Booking{Status: "cancelled"}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:        "Failure - Booking Not Found",
			bookingCode: "NOTFOUND",
			setupMock: func(mockService *mocks.BookingService) {
				mockService.On("CancelBookingByCode", "NOTFOUND").Return(nil, fmt.Errorf("not found")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError, // Service errors return 500
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			router, _, _, mockBookingService := setupTestRouter()
			tc.setupMock(mockBookingService)

			req, _ := http.NewRequest("DELETE", "/bookings/"+tc.bookingCode, nil)
			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			mockBookingService.AssertExpectations(t)
		})
	}
}
