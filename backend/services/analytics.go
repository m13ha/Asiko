package services

import (
	"math"

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

	cancellationCount, err := s.analyticsRepo.GetUserCancellationCount(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	cancellationsPerDay, err := s.analyticsRepo.GetCancellationsPerDay(userID, timeRange.Start, timeRange.End)
	if err != nil {
		return nil, serviceerrors.FromError(err)
	}

	days := math.Floor(timeRange.End.Sub(timeRange.Start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}
	avgBookingsPerDay := float64(bookingCount) / days

	denominator := float64(bookingCount + cancellationCount)
	cancellationRate := 0.0
	if denominator > 0 {
		cancellationRate = (float64(cancellationCount) / denominator) * 100
	}

	// Build response
	resp := &responses.AnalyticsResponse{
		TotalAppointments:   int(appointmentCount),
		TotalBookings:       int(bookingCount),
		TotalCancellations:  int(cancellationCount),
		CancellationRate:    cancellationRate,
		AvgBookingsPerDay:   avgBookingsPerDay,
		StartDate:           timeRange.Start,
		EndDate:             timeRange.End,
		BookingsPerDay:      make([]responses.TimeSeriesPoint, 0, len(bookingsPerDay)),
		CancellationsPerDay: make([]responses.TimeSeriesPoint, 0, len(cancellationsPerDay)),
	}

	for _, row := range bookingsPerDay {
		resp.BookingsPerDay = append(resp.BookingsPerDay, responses.TimeSeriesPoint{
			Date:  row.Date,
			Count: row.Count,
		})
	}

	for _, row := range cancellationsPerDay {
		resp.CancellationsPerDay = append(resp.CancellationsPerDay, responses.TimeSeriesPoint{
			Date:  row.Date,
			Count: row.Count,
		})
	}

	return resp, nil
}
