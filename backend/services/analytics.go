package services

import (
	"github.com/google/uuid"
	serviceerrors "github.com/m13ha/asiko/errors/serviceerrors"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/repository"
	"github.com/m13ha/asiko/utils"
)

type AnalyticsService interface {
	GetUserAnalytics(userID uuid.UUID, startDate, endDate string) (*responses.AnalyticsResponse, error)
}

type analyticsServiceImpl struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsServiceImpl{analyticsRepo: analyticsRepo}
}

func (s *analyticsServiceImpl) GetUserAnalytics(userID uuid.UUID, startDate, endDate string) (*responses.AnalyticsResponse, error) {
	timeRange, err := utils.ParseTimeRange(startDate, endDate)
	if err != nil {
		return nil, serviceerrors.ValidationError(err.Error())
	}

	appointmentCount, err := s.analyticsRepo.GetUserAppointmentCount(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	bookingCount, err := s.analyticsRepo.GetUserBookingCount(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	// Time series and insights
	bookingsPerDay, err := s.analyticsRepo.GetBookingsPerDay(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	// Build response
	resp := &responses.AnalyticsResponse{
		TotalAppointments: int(appointmentCount),
		TotalBookings:     int(bookingCount),
		StartDate:         timeRange.Start,
		EndDate:           timeRange.End,
		BookingsPerDay:    make([]responses.TimeSeriesPoint, len(bookingsPerDay)),
	}

	return resp, nil
}
