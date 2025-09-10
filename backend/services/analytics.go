package services

import (
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/repository"
	"github.com/m13ha/appointment_master/utils"
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
		return nil, err
	}

	appointmentCount, err := s.analyticsRepo.GetUserAppointmentCount(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, err
	}

	bookingCount, err := s.analyticsRepo.GetUserBookingCount(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, err
	}

	return &responses.AnalyticsResponse{
		TotalAppointments: int(appointmentCount),
		TotalBookings:     int(bookingCount),
		StartDate:         timeRange.Start,
		EndDate:           timeRange.End,
	}, nil
}