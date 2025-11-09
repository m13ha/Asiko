package repository

import (
	"time"

	"github.com/google/uuid"
	apperr "github.com/m13ha/asiko/errors"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)

	// Breakdowns
	GetAppointmentsByTypeCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error)
	GetBookingsByStatusCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error)
	GetGuestVsRegisteredCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error)
	GetDistinctAndRepeatCustomers(userID uuid.UUID, startDate, endDate time.Time) (distinct int, repeat int, err error)

	// Utilization & capacity
	GetSlotUtilization(userID uuid.UUID, startDate, endDate time.Time) (totalSlots int64, bookedSlots int64, err error)
	GetAvgAttendeesPerBooking(userID uuid.UUID, startDate, endDate time.Time) (float64, error)
	GetPartyCapacity(userID uuid.UUID, startDate, endDate time.Time) (total int64, used int64, err error)

	// Timing
	GetLeadTimeStatsHours(userID uuid.UUID, startDate, endDate time.Time) (avg float64, median float64, err error)

	// Time series and insights
	GetBookingsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error)
	GetStatusPerDay(userID uuid.UUID, status string, startDate, endDate time.Time) ([]DateCount, error)
	GetPeakHours(userID uuid.UUID, startDate, endDate time.Time) ([]KeyCount, error)
	GetPeakDays(userID uuid.UUID, startDate, endDate time.Time) ([]KeyCount, error)
	GetTopAppointments(userID uuid.UUID, startDate, endDate time.Time, limit int) ([]TopAppointmentRow, error)
}

type gormAnalyticsRepository struct {
	db *gorm.DB
}

func NewGormAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &gormAnalyticsRepository{db: db}
}

func (r *gormAnalyticsRepository) GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("appointments").
		Where("owner_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Count(&count).Error
	if err != nil {
		return count, apperr.TranslateRepoError("repository.analytics.UserAppointmentCount", err)
	}
	return count, nil
}

func (r *gormAnalyticsRepository) GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("bookings").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = ? AND bookings.is_slot = FALSE",
			userID, startDate, endDate, "active").
		Count(&count).Error
	if err != nil {
		return count, apperr.TranslateRepoError("repository.analytics.UserBookingCount", err)
	}
	return count, nil
}

func (r *gormAnalyticsRepository) GetAppointmentsByTypeCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error) {
	type row struct {
		Type string
		Cnt  int
	}
	var rows []row
	err := r.db.Table("appointments").
		Select("type, COUNT(*) as cnt").
		Where("owner_id = ? AND created_at BETWEEN ? AND ?", userID, startDate, endDate).
		Group("type").
		Find(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.AppointmentsByType", err)
	}
	out := map[string]int{}
	for _, v := range rows {
		out[v.Type] = v.Cnt
	}
	return out, nil
}

func (r *gormAnalyticsRepository) GetBookingsByStatusCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error) {
	type row struct {
		Status string
		Cnt    int
	}
	var rows []row
	err := r.db.Table("bookings").
		Select("bookings.status as status, COUNT(*) as cnt").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.is_slot = FALSE", userID, startDate, endDate).
		Group("bookings.status").
		Find(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.BookingsByStatus", err)
	}
	out := map[string]int{}
	for _, v := range rows {
		out[v.Status] = v.Cnt
	}
	return out, nil
}

func (r *gormAnalyticsRepository) GetGuestVsRegisteredCounts(userID uuid.UUID, startDate, endDate time.Time) (map[string]int, error) {
	type row struct {
		Kind string
		Cnt  int
	}
	var rows []row
	err := r.db.Table("bookings").
		Select("CASE WHEN bookings.user_id IS NULL THEN 'guest' ELSE 'registered' END as kind, COUNT(*) as cnt").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.is_slot = FALSE", userID, startDate, endDate).
		Group("kind").
		Find(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.GuestVsRegistered", err)
	}
	out := map[string]int{}
	for _, v := range rows {
		out[v.Kind] = v.Cnt
	}
	return out, nil
}

func (r *gormAnalyticsRepository) GetDistinctAndRepeatCustomers(userID uuid.UUID, startDate, endDate time.Time) (int, int, error) {
	// Distinct users by (coalesce(user_id::text, email))
	type rec struct {
		Key string
		Cnt int
	}
	var rows []rec
	err := r.db.Raw(`
        SELECT COALESCE(CAST(bookings.user_id AS TEXT), LOWER(bookings.email)) AS key, COUNT(*) as cnt
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.is_slot = FALSE
        GROUP BY key
    `, userID, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return 0, 0, apperr.TranslateRepoError("repository.analytics.DistinctRepeat", err)
	}
	distinct := len(rows)
	repeat := 0
	for _, r2 := range rows {
		if r2.Cnt > 1 {
			repeat++
		}
	}
	return distinct, repeat, nil
}

func (r *gormAnalyticsRepository) GetSlotUtilization(userID uuid.UUID, startDate, endDate time.Time) (int64, int64, error) {
	var total int64
	var booked int64
	// total slots = all slot rows for single/group appointments
	err := r.db.Table("bookings").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND appointments.type IN ('single','group') AND bookings.date BETWEEN ? AND ? AND bookings.is_slot = TRUE", userID, startDate, endDate).
		Count(&total).Error
	if err != nil {
		return 0, 0, apperr.TranslateRepoError("repository.analytics.TotalSlots", err)
	}
	err = r.db.Table("bookings").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND appointments.type IN ('single','group') AND bookings.date BETWEEN ? AND ? AND bookings.is_slot = TRUE AND bookings.available = FALSE", userID, startDate, endDate).
		Count(&booked).Error
	if err != nil {
		return 0, 0, apperr.TranslateRepoError("repository.analytics.BookedSlots", err)
	}
	return total, booked, nil
}

func (r *gormAnalyticsRepository) GetAvgAttendeesPerBooking(userID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	type row struct{ Avg float64 }
	var o row
	err := r.db.Raw(`
        SELECT COALESCE(AVG(attendee_count), 0) as avg
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = 'active' AND bookings.is_slot = FALSE
    `, userID, startDate, endDate).Scan(&o).Error
	if err != nil {
		return 0, apperr.TranslateRepoError("repository.analytics.AvgAttendees", err)
	}
	return o.Avg, nil
}

func (r *gormAnalyticsRepository) GetPartyCapacity(userID uuid.UUID, startDate, endDate time.Time) (int64, int64, error) {
	type row struct {
		Total int64
		Used  int64
	}
	var o row
	err := r.db.Raw(`
        SELECT COALESCE(SUM(max_attendees),0) as total, COALESCE(SUM(attendees_booked),0) as used
        FROM appointments
        WHERE owner_id = ? AND type = 'party' AND created_at BETWEEN ? AND ?
    `, userID, startDate, endDate).Scan(&o).Error
	if err != nil {
		return 0, 0, apperr.TranslateRepoError("repository.analytics.PartyCapacity", err)
	}
	return o.Total, o.Used, nil
}

func (r *gormAnalyticsRepository) GetLeadTimeStatsHours(userID uuid.UUID, startDate, endDate time.Time) (float64, float64, error) {
	type row struct {
		Avg    float64
		Median float64
	}
	var o row
	err := r.db.Raw(`
        SELECT 
            COALESCE(AVG(EXTRACT(EPOCH FROM (bookings.start_time - bookings.created_at)))/3600, 0) AS avg,
            COALESCE(percentile_cont(0.5) WITHIN GROUP (ORDER BY EXTRACT(EPOCH FROM (bookings.start_time - bookings.created_at))/3600), 0) AS median
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = 'active' AND bookings.is_slot = FALSE
    `, userID, startDate, endDate).Scan(&o).Error
	if err != nil {
		return 0, 0, apperr.TranslateRepoError("repository.analytics.LeadTimeStats", err)
	}
	return o.Avg, o.Median, nil
}

type DateCount struct {
	Date  string
	Count int
}
type KeyCount struct {
	Key   string
	Count int
}
type TopAppointmentRow struct {
	AppCode              string
	Title                string
	Bookings             int
	CapacityUsagePercent float64
}

func (r *gormAnalyticsRepository) GetBookingsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error) {
	var rows []DateCount
	err := r.db.Raw(`
        SELECT TO_CHAR(bookings.created_at::date, 'YYYY-MM-DD') as date, COUNT(*) as count
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = 'active' AND bookings.is_slot = FALSE
        GROUP BY bookings.created_at::date
        ORDER BY bookings.created_at::date
    `, userID, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.BookingsPerDay", err)
	}
	return rows, nil
}

func (r *gormAnalyticsRepository) GetStatusPerDay(userID uuid.UUID, status string, startDate, endDate time.Time) ([]DateCount, error) {
	var rows []DateCount
	err := r.db.Raw(`
        SELECT TO_CHAR(bookings.created_at::date, 'YYYY-MM-DD') as date, COUNT(*) as count
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = ? AND bookings.is_slot = FALSE
        GROUP BY bookings.created_at::date
        ORDER BY bookings.created_at::date
    `, userID, startDate, endDate, status).Scan(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.StatusPerDay", err)
	}
	return rows, nil
}

func (r *gormAnalyticsRepository) GetPeakHours(userID uuid.UUID, startDate, endDate time.Time) ([]KeyCount, error) {
	var rows []KeyCount
	err := r.db.Raw(`
        SELECT LPAD(CAST(EXTRACT(HOUR FROM bookings.start_time) AS TEXT),2,'0') as key, COUNT(*) as count
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = 'active' AND bookings.is_slot = FALSE
        GROUP BY key
        ORDER BY key
    `, userID, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.PeakHours", err)
	}
	return rows, nil
}

func (r *gormAnalyticsRepository) GetPeakDays(userID uuid.UUID, startDate, endDate time.Time) ([]KeyCount, error) {
	var rows []KeyCount
	err := r.db.Raw(`
        SELECT CAST(EXTRACT(DOW FROM bookings.start_time) AS TEXT) as key, COUNT(*) as count
        FROM bookings
        JOIN appointments ON bookings.appointment_id = appointments.id
        WHERE appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = 'active' AND bookings.is_slot = FALSE
        GROUP BY key
        ORDER BY key
    `, userID, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.PeakDays", err)
	}
	return rows, nil
}

func (r *gormAnalyticsRepository) GetTopAppointments(userID uuid.UUID, startDate, endDate time.Time, limit int) ([]TopAppointmentRow, error) {
	// Compute bookings count and capacity usage percent for party appointments
	var rows []TopAppointmentRow
	err := r.db.Raw(`
        WITH base AS (
            SELECT appointments.id, appointments.app_code, appointments.title,
                   COUNT(CASE WHEN bookings.status = 'active' AND bookings.is_slot = FALSE THEN 1 END) AS bookings,
                   CASE 
                     WHEN appointments.type = 'party' AND appointments.max_attendees > 0 
                       THEN (appointments.attendees_booked::float / appointments.max_attendees::float) * 100.0
                     ELSE 0
                   END AS capacity_usage_percent
            FROM appointments
            LEFT JOIN bookings ON bookings.appointment_id = appointments.id 
              AND bookings.created_at BETWEEN ? AND ?
              AND bookings.is_slot = FALSE
            WHERE appointments.owner_id = ? AND appointments.created_at <= ?
            GROUP BY appointments.id
        )
        SELECT app_code as app_code, title as title, bookings as bookings, capacity_usage_percent
        FROM base
        ORDER BY bookings DESC
        LIMIT ?
    `, startDate, endDate, userID, endDate, limit).Scan(&rows).Error
	if err != nil {
		return nil, apperr.TranslateRepoError("repository.analytics.TopAppointments", err)
	}
	return rows, nil
}
