package repository

import (
	"time"

	"github.com/google/uuid"
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetBookingsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error)
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
		return count, repoerrors.InternalError("failed to get user appointment count: " + err.Error())
	}
	return count, nil
}

func (r *gormAnalyticsRepository) GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	type result struct {
		Total int64
	}
	var res result

	// Count bookings across all appointment types:
	// - Party: count booking records (is_slot=false)
	// - Single: count booked slots (is_slot=true AND available=false)
	// - Group: sum seats_booked from slots (is_slot=true)
	err := r.db.Raw(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN appointments.type = 'party' AND bookings.is_slot = FALSE THEN 1
				WHEN appointments.type = 'single' AND bookings.is_slot = TRUE AND bookings.available = FALSE THEN 1
				WHEN appointments.type = 'group' AND bookings.is_slot = TRUE THEN bookings.seats_booked
				ELSE 0
			END
		), 0) as total
		FROM bookings
		JOIN appointments ON bookings.appointment_id = appointments.id
		WHERE appointments.owner_id = ? 
			AND (bookings.created_at BETWEEN ? AND ? OR bookings.date BETWEEN ? AND ?)
			AND bookings.status = 'active'
	`, userID, startDate, endDate, startDate, endDate).Scan(&res).Error

	if err != nil {
		return 0, repoerrors.InternalError("failed to get user booking count: " + err.Error())
	}

	return res.Total, nil
}

type DateCount struct {
	Date  string
	Count int
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
		return nil, repoerrors.InternalError("failed to get bookings per day: " + err.Error())
	}
	return rows, nil
}
