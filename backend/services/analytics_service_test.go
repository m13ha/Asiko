package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/repository/mocks"
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

				// Return zero/empty defaults for remaining calls used by the service
				mockRepo.On("GetAppointmentsByTypeCounts", userID, start, end).Return(map[string]int{}, nil)
				mockRepo.On("GetBookingsByStatusCounts", userID, start, end).Return(map[string]int{}, nil)
				mockRepo.On("GetGuestVsRegisteredCounts", userID, start, end).Return(map[string]int{}, nil)
				mockRepo.On("GetDistinctAndRepeatCustomers", userID, start, end).Return(0, 0, nil)
				mockRepo.On("GetSlotUtilization", userID, start, end).Return(int64(0), int64(0), nil)
				mockRepo.On("GetAvgAttendeesPerBooking", userID, start, end).Return(0.0, nil)
				mockRepo.On("GetPartyCapacity", userID, start, end).Return(int64(0), int64(0), nil)
				mockRepo.On("GetLeadTimeStatsHours", userID, start, end).Return(0.0, 0.0, nil)
				mockRepo.On("GetBookingsPerDay", userID, start, end).Return([]repository.DateCount{}, nil)
				mockRepo.On("GetStatusPerDay", userID, "cancelled", start, end).Return([]repository.DateCount{}, nil)
				mockRepo.On("GetStatusPerDay", userID, "rejected", start, end).Return([]repository.DateCount{}, nil)
				mockRepo.On("GetPeakHours", userID, start, end).Return([]repository.KeyCount{}, nil)
				mockRepo.On("GetPeakDays", userID, start, end).Return([]repository.KeyCount{}, nil)
				mockRepo.On("GetTopAppointments", userID, start, end, 5).Return([]repository.TopAppointmentRow{}, nil)
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

	appsByType := map[string]int{"single": 3, "group": 2}
	mockRepo.On("GetAppointmentsByTypeCounts", userID, start, end).Return(appsByType, nil)

	bookingsByStatus := map[string]int{"active": 10, "cancelled": 1, "rejected": 1}
	mockRepo.On("GetBookingsByStatusCounts", userID, start, end).Return(bookingsByStatus, nil)

	guestVsReg := map[string]int{"guest": 4, "registered": 8}
	mockRepo.On("GetGuestVsRegisteredCounts", userID, start, end).Return(guestVsReg, nil)

	mockRepo.On("GetDistinctAndRepeatCustomers", userID, start, end).Return(10, 2, nil)

	mockRepo.On("GetSlotUtilization", userID, start, end).Return(int64(100), int64(40), nil)
	mockRepo.On("GetAvgAttendeesPerBooking", userID, start, end).Return(1.5, nil)
	mockRepo.On("GetPartyCapacity", userID, start, end).Return(int64(50), int64(20), nil)

	mockRepo.On("GetLeadTimeStatsHours", userID, start, end).Return(24.0, 12.0, nil)

	bookingsPerDay := []repository.DateCount{{Date: "2025-01-01", Count: 2}, {Date: "2025-01-02", Count: 3}}
	cancelsPerDay := []repository.DateCount{{Date: "2025-01-02", Count: 1}}
	rejectsPerDay := []repository.DateCount{{Date: "2025-01-03", Count: 1}}
	mockRepo.On("GetBookingsPerDay", userID, start, end).Return(bookingsPerDay, nil)
	mockRepo.On("GetStatusPerDay", userID, "cancelled", start, end).Return(cancelsPerDay, nil)
	mockRepo.On("GetStatusPerDay", userID, "rejected", start, end).Return(rejectsPerDay, nil)

	peakHours := []repository.KeyCount{{Key: "09", Count: 4}, {Key: "10", Count: 5}}
	peakDays := []repository.KeyCount{{Key: "1", Count: 3}, {Key: "2", Count: 6}}
	mockRepo.On("GetPeakHours", userID, start, end).Return(peakHours, nil)
	mockRepo.On("GetPeakDays", userID, start, end).Return(peakDays, nil)

	topRows := []repository.TopAppointmentRow{{AppCode: "AP1", Title: "A", Bookings: 5, CapacityUsagePercent: 40.0}}
	mockRepo.On("GetTopAppointments", userID, start, end, 5).Return(topRows, nil)

	svc := NewAnalyticsService(mockRepo)
	resp, err := svc.GetUserAnalytics(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 5, resp.TotalAppointments)
	assert.Equal(t, 12, resp.TotalBookings)
	assert.Equal(t, appsByType, resp.AppointmentsByType)
	assert.Equal(t, bookingsByStatus, resp.BookingsByStatus)
	assert.Equal(t, guestVsReg, resp.GuestVsRegistered)
	assert.Equal(t, 10, resp.DistinctCustomers)
	assert.Equal(t, 2, resp.RepeatCustomers)
	assert.InDelta(t, 40.0, resp.SlotUtilizationPercent, 0.001)
	assert.InDelta(t, 1.5, resp.AvgAttendeesPerBooking, 0.001)
	assert.Equal(t, 50, resp.PartyCapacity.Total)
	assert.Equal(t, 20, resp.PartyCapacity.Used)
	assert.InDelta(t, 40.0, resp.PartyCapacity.Percent, 0.001)
	assert.InDelta(t, 24.0, resp.AvgLeadTimeHours, 0.001)
	assert.InDelta(t, 12.0, resp.MedianLeadTimeHours, 0.001)

	// series mapping
	if assert.Len(t, resp.BookingsPerDay, len(bookingsPerDay)) {
		assert.Equal(t, bookingsPerDay[0].Date, resp.BookingsPerDay[0].Date)
		assert.Equal(t, bookingsPerDay[0].Count, resp.BookingsPerDay[0].Count)
	}
	if assert.Len(t, resp.CancellationsPerDay, len(cancelsPerDay)) {
		assert.Equal(t, cancelsPerDay[0].Date, resp.CancellationsPerDay[0].Date)
		assert.Equal(t, cancelsPerDay[0].Count, resp.CancellationsPerDay[0].Count)
	}
	if assert.Len(t, resp.RejectionsPerDay, len(rejectsPerDay)) {
		assert.Equal(t, rejectsPerDay[0].Date, resp.RejectionsPerDay[0].Date)
		assert.Equal(t, rejectsPerDay[0].Count, resp.RejectionsPerDay[0].Count)
	}
	if assert.Len(t, resp.PeakHours, len(peakHours)) {
		assert.Equal(t, peakHours[0].Key, resp.PeakHours[0].Key)
		assert.Equal(t, peakHours[0].Count, resp.PeakHours[0].Count)
	}
	if assert.Len(t, resp.PeakDays, len(peakDays)) {
		assert.Equal(t, peakDays[0].Key, resp.PeakDays[0].Key)
		assert.Equal(t, peakDays[0].Count, resp.PeakDays[0].Count)
	}
	if assert.Len(t, resp.TopAppointments, len(topRows)) {
		assert.Equal(t, topRows[0].AppCode, resp.TopAppointments[0].AppCode)
		assert.Equal(t, topRows[0].Title, resp.TopAppointments[0].Title)
		assert.Equal(t, topRows[0].Bookings, resp.TopAppointments[0].Bookings)
		assert.InDelta(t, topRows[0].CapacityUsagePercent, resp.TopAppointments[0].CapacityUsagePercent, 0.001)
	}

	mockRepo.AssertExpectations(t)
}

func TestGetUserAnalytics_FailsOnBreakdownError(t *testing.T) {
	userID := uuid.New()
	startDate := "2025-01-01"
	endDate := "2025-01-31"

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	mockRepo := new(mocks.AnalyticsRepository)

	mockRepo.On("GetUserAppointmentCount", userID, start, end).Return(int64(5), nil)
	mockRepo.On("GetUserBookingCount", userID, start, end).Return(int64(12), nil)
	mockRepo.On("GetAppointmentsByTypeCounts", userID, start, end).Return(map[string]int{"single": 3}, nil)
	mockRepo.On("GetBookingsByStatusCounts", userID, start, end).Return((map[string]int)(nil), fmt.Errorf("boom"))

	svc := NewAnalyticsService(mockRepo)
	resp, err := svc.GetUserAnalytics(userID, startDate, endDate)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "boom")
	mockRepo.AssertExpectations(t)
}
