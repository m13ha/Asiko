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

    // Breakdowns
    appsByType, err := s.analyticsRepo.GetAppointmentsByTypeCounts(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    bookingsByStatus, err := s.analyticsRepo.GetBookingsByStatusCounts(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    guestVsRegistered, err := s.analyticsRepo.GetGuestVsRegisteredCounts(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    distinctCustomers, repeatCustomers, err := s.analyticsRepo.GetDistinctAndRepeatCustomers(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    // Utilization & capacity
    totalSlots, bookedSlots, err := s.analyticsRepo.GetSlotUtilization(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }
    var slotUtilization float64
    if totalSlots > 0 {
        slotUtilization = (float64(bookedSlots) / float64(totalSlots)) * 100.0
    }

    avgAttendees, err := s.analyticsRepo.GetAvgAttendeesPerBooking(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    partyTotal, partyUsed, err := s.analyticsRepo.GetPartyCapacity(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }
    var partyPercent float64
    if partyTotal > 0 {
        partyPercent = (float64(partyUsed) / float64(partyTotal)) * 100.0
    }

    // Timing
    avgLead, medLead, err := s.analyticsRepo.GetLeadTimeStatsHours(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    // Time series and insights
    bookingsPerDay, err := s.analyticsRepo.GetBookingsPerDay(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }
    cancelsPerDay, err := s.analyticsRepo.GetStatusPerDay(userID, "cancelled", timeRange.Start, timeRange.End)
    if err != nil { return nil, err }
    rejectsPerDay, err := s.analyticsRepo.GetStatusPerDay(userID, "rejected", timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    peakHours, err := s.analyticsRepo.GetPeakHours(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }
    peakDays, err := s.analyticsRepo.GetPeakDays(userID, timeRange.Start, timeRange.End)
    if err != nil { return nil, err }

    topRows, err := s.analyticsRepo.GetTopAppointments(userID, timeRange.Start, timeRange.End, 5)
    if err != nil { return nil, err }
    top := make([]responses.TopAppointment, 0, len(topRows))
    for _, r := range topRows {
        top = append(top, responses.TopAppointment{
            AppCode: r.AppCode,
            Title: r.Title,
            Bookings: r.Bookings,
            CapacityUsagePercent: r.CapacityUsagePercent,
        })
    }

    // Build response
    resp := &responses.AnalyticsResponse{
        TotalAppointments: int(appointmentCount),
        TotalBookings:     int(bookingCount),
        StartDate:         timeRange.Start,
        EndDate:           timeRange.End,
        AppointmentsByType: appsByType,
        BookingsByStatus:   bookingsByStatus,
        GuestVsRegistered:  guestVsRegistered,
        DistinctCustomers:  distinctCustomers,
        RepeatCustomers:    repeatCustomers,
        SlotUtilizationPercent: slotUtilization,
        AvgAttendeesPerBooking: avgAttendees,
        BookingsPerDay:         make([]responses.TimeSeriesPoint, 0, len(bookingsPerDay)),
        CancellationsPerDay:    make([]responses.TimeSeriesPoint, 0, len(cancelsPerDay)),
        RejectionsPerDay:       make([]responses.TimeSeriesPoint, 0, len(rejectsPerDay)),
        PeakHours:              make([]responses.BucketCount, 0, len(peakHours)),
        PeakDays:               make([]responses.BucketCount, 0, len(peakDays)),
        TopAppointments:        top,
    }
    resp.PartyCapacity.Total = int(partyTotal)
    resp.PartyCapacity.Used = int(partyUsed)
    resp.PartyCapacity.Percent = partyPercent
    resp.AvgLeadTimeHours = avgLead
    resp.MedianLeadTimeHours = medLead

    for _, p := range bookingsPerDay {
        resp.BookingsPerDay = append(resp.BookingsPerDay, responses.TimeSeriesPoint{Date: p.Date, Count: p.Count})
    }
    for _, p := range cancelsPerDay {
        resp.CancellationsPerDay = append(resp.CancellationsPerDay, responses.TimeSeriesPoint{Date: p.Date, Count: p.Count})
    }
    for _, p := range rejectsPerDay {
        resp.RejectionsPerDay = append(resp.RejectionsPerDay, responses.TimeSeriesPoint{Date: p.Date, Count: p.Count})
    }
    for _, b := range peakHours {
        resp.PeakHours = append(resp.PeakHours, responses.BucketCount{Key: b.Key, Count: b.Count})
    }
    for _, b := range peakDays {
        resp.PeakDays = append(resp.PeakDays, responses.BucketCount{Key: b.Key, Count: b.Count})
    }

    return resp, nil
}
