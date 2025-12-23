package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	myerrors "github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/repository/mocks"
	services "github.com/m13ha/asiko/services"
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
		cancellationCount int64
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
			cancellationCount: 3,
			setupMock: func(mockRepo *mocks.AnalyticsRepository, userID uuid.UUID) {
				start, _ := time.Parse("2006-01-02", "2025-01-01")
				end, _ := time.Parse("2006-01-02", "2025-01-31")
				end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

				mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(5), nil)
				mockRepo.On("GetUserBookingCount", userID, start, end).Return(int64(12), nil)
				mockRepo.On("GetBookingsPerDay", userID, start, end).Return([]repository.DateCount{}, nil)
				mockRepo.On("GetUserCancellationCount", userID, start, end).Return(int64(3), nil)
				mockRepo.On("GetCancellationsPerDay", userID, start, end).Return([]repository.DateCount{}, nil)
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

				repoErr := myerrors.NewAppError(myerrors.CodeInternalError, "internal", 500, "db error", nil)
				mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(0), repoErr)
			},
			expectedError: "db error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.AnalyticsRepository)
			tc.setupMock(mockRepo, tc.userID)

			service := services.NewAnalyticsService(mockRepo)

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
				assert.Equal(t, int(tc.cancellationCount), result.TotalCancellations)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserAnalytics_DetailedSuccess(t *testing.T) {
	userID := uuid.New()
	startDate := "2025-01-01"
	endDate := "2025-01-31"

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	mockRepo := new(mocks.AnalyticsRepository)

	mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(5), nil)
	mockRepo.On("GetUserBookingCount", userID, start, end).Return(int64(12), nil)

	bookingsPerDay := []repository.DateCount{{Date: "2025-01-01", Count: 2}, {Date: "2025-01-02", Count: 3}}
	mockRepo.On("GetBookingsPerDay", userID, start, end).Return(bookingsPerDay, nil)
	mockRepo.On("GetUserCancellationCount", userID, start, end).Return(int64(4), nil)
	cancellationsPerDay := []repository.DateCount{{Date: "2025-01-05", Count: 1}, {Date: "2025-01-06", Count: 2}}
	mockRepo.On("GetCancellationsPerDay", userID, start, end).Return(cancellationsPerDay, nil)

	svc := services.NewAnalyticsService(mockRepo)
	resp, err := svc.GetUserAnalytics(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 5, resp.TotalAppointments)
	assert.Equal(t, 12, resp.TotalBookings)
	assert.Equal(t, 4, resp.TotalCancellations)
	assert.Len(t, resp.BookingsPerDay, 2)
	assert.Len(t, resp.CancellationsPerDay, 2)

	mockRepo.AssertExpectations(t)
}

func TestGetUserAnalytics_FailsOnBookingsPerDayError(t *testing.T) {
	userID := uuid.New()
	startDate := "2025-01-01"
	endDate := "2025-01-31"

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	mockRepo := new(mocks.AnalyticsRepository)

	mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(5), nil)
	mockRepo.On("GetUserBookingCount", userID, start, end).Return(int64(12), nil)

	repoErr := myerrors.NewAppError(myerrors.CodeInternalError, "internal", 500, "boom", nil)
	mockRepo.On("GetBookingsPerDay", userID, start, end).Return([]repository.DateCount{}, repoErr)

	svc := services.NewAnalyticsService(mockRepo)
	resp, err := svc.GetUserAnalytics(userID, startDate, endDate)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "boom")
	mockRepo.AssertExpectations(t)
}
