package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetUserAnalytics(t *testing.T) {
	testCases := []struct {
		name             string
		userID           uuid.UUID
		startDate        string
		endDate          string
		appointmentCount int64
		bookingCount     int64
		setupMock        func(mockRepo *mocks.AnalyticsRepository, userID uuid.UUID)
		expectedError    string
	}{
		{
			name:             "Success",
			userID:           uuid.New(),
			startDate:        "2025-01-01",
			endDate:          "2025-01-31",
			appointmentCount: 5,
			bookingCount:     12,
			setupMock: func(mockRepo *mocks.AnalyticsRepository, userID uuid.UUID) {
				start, _ := time.Parse("2006-01-02", "2025-01-01")
				end, _ := time.Parse("2006-01-02", "2025-01-31")
				end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

				mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(5), nil)
				mockRepo.On("GetUserBookingCount", userID, start, end).Return(int64(12), nil)
			},
			expectedError: "",
		},
		{
			name:          "Invalid date format",
			userID:        uuid.New(),
			startDate:     "2025/01/01",
			endDate:       "2025-01-31",
			setupMock:     func(mockRepo *mocks.AnalyticsRepository, userID uuid.UUID) {},
			expectedError: "invalid start_date format, use YYYY-MM-DD",
		},
		{
			name:      "Repository error",
			userID:    uuid.New(),
			startDate: "2025-01-01",
			endDate:   "2025-01-31",
			setupMock: func(mockRepo *mocks.AnalyticsRepository, userID uuid.UUID) {
				start, _ := time.Parse("2006-01-02", "2025-01-01")
				end, _ := time.Parse("2006-01-02", "2025-01-31")
				end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

				mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(0), fmt.Errorf("db error"))
			},
			expectedError: "db error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.AnalyticsRepository)
			tc.setupMock(mockRepo, tc.userID)

			service := NewAnalyticsService(mockRepo)

			result, err := service.GetUserAnalytics(tc.userID, tc.startDate, tc.endDate)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, int(tc.appointmentCount), result.TotalAppointments)
				assert.Equal(t, int(tc.bookingCount), result.TotalBookings)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
