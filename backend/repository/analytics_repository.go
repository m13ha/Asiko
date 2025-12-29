package repository

import (
	"time"

	"github.com/google/uuid"
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"github.com/m13ha/asiko/models/entities"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetBookingsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error)
	GetUserCancellationCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetCancellationsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error)
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
		Where("owner_id = ? AND created_at BETWEEN ? AND ? AND status <> ?", userID, startDate, endDate, entities.AppointmentStatusCanceled).
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
	// - Party: sum attendee_count from reservation rows (is_slot=false)
	// - Single: count booked slots (is_slot=true AND available=false)
	// - Group: sum attendee_count from reservation rows (is_slot=false)
	err := r.db.Raw(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN appointments.type = 'party' AND bookings.is_slot = FALSE THEN bookings.attendee_count
				WHEN appointments.type = 'single' AND bookings.is_slot = TRUE AND bookings.available = FALSE THEN 1
				WHEN appointments.type = 'group' AND bookings.is_slot = FALSE THEN bookings.attendee_count
				ELSE 0
			END
		), 0) as total
		FROM bookings
		JOIN appointments ON bookings.appointment_id = appointments.id
		WHERE appointments.owner_id = ? 
			AND (bookings.created_at BETWEEN ? AND ? OR bookings.date BETWEEN ? AND ?)
			AND LOWER(bookings.status) IN (?, ?, ?, ?)
	`, userID, startDate, endDate, startDate, endDate,
		entities.BookingStatusPending,
		entities.BookingStatusOngoing,
		entities.BookingStatusConfirmed,
		entities.BookingStatusExpired,
	).Scan(&res).Error

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
		SELECT
			TO_CHAR(
				CASE
					WHEN bookings.is_slot = TRUE THEN bookings.updated_at::date
					ELSE bookings.created_at::date
				END,
				'YYYY-MM-DD'
			) as date,
			COALESCE(SUM(
				CASE 
					WHEN appointments.type = 'party' AND bookings.is_slot = FALSE THEN bookings.attendee_count
					WHEN appointments.type = 'single' AND bookings.is_slot = TRUE AND bookings.available = FALSE THEN 1
					WHEN appointments.type = 'group' AND bookings.is_slot = FALSE THEN bookings.attendee_count
					ELSE 0
				END
			), 0) as count
		FROM bookings
		JOIN appointments ON bookings.appointment_id = appointments.id
		WHERE appointments.owner_id = ?
			AND (
				CASE
					WHEN bookings.is_slot = TRUE THEN bookings.updated_at
					ELSE bookings.created_at
				END
			) BETWEEN ? AND ?
			AND LOWER(bookings.status) IN (?, ?, ?, ?)
		GROUP BY 1
		ORDER BY 1
	`, userID, startDate, endDate,
		entities.BookingStatusPending,
		entities.BookingStatusOngoing,
		entities.BookingStatusConfirmed,
		entities.BookingStatusExpired,
	).Scan(&rows).Error
	if err != nil {
		return nil, repoerrors.InternalError("failed to get bookings per day: " + err.Error())
	}
	return rows, nil
}

func (r *gormAnalyticsRepository) GetUserCancellationCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	type result struct {
		Total int64
	}
	var res result

	err := r.db.Raw(`
		SELECT COALESCE(SUM(
			CASE 
				WHEN appointments.type = 'party' AND bookings.is_slot = FALSE THEN 1
				WHEN appointments.type = 'single' AND bookings.is_slot = TRUE THEN 1
				WHEN appointments.type = 'group' AND bookings.is_slot = FALSE THEN bookings.attendee_count
				ELSE 0
			END
		), 0) as total
		FROM bookings
		JOIN appointments ON bookings.appointment_id = appointments.id
		WHERE appointments.owner_id = ?
			AND bookings.updated_at BETWEEN ? AND ?
			AND bookings.status = 'cancelled'
	`, userID, startDate, endDate).Scan(&res).Error

	if err != nil {
		return 0, repoerrors.InternalError("failed to get user cancellation count: " + err.Error())
	}

	return res.Total, nil
}

func (r *gormAnalyticsRepository) GetCancellationsPerDay(userID uuid.UUID, startDate, endDate time.Time) ([]DateCount, error) {
	var rows []DateCount
	err := r.db.Raw(`
		SELECT
			TO_CHAR(bookings.updated_at::date, 'YYYY-MM-DD') as date,
			COALESCE(SUM(
				CASE 
					WHEN appointments.type = 'party' AND bookings.is_slot = FALSE THEN 1
					WHEN appointments.type = 'single' AND bookings.is_slot = TRUE THEN 1
					WHEN appointments.type = 'group' AND bookings.is_slot = FALSE THEN bookings.attendee_count
					ELSE 0
				END
			), 0) as count
		FROM bookings
		JOIN appointments ON bookings.appointment_id = appointments.id
		WHERE appointments.owner_id = ?
			AND bookings.updated_at BETWEEN ? AND ?
			AND bookings.status = 'cancelled'
		GROUP BY bookings.updated_at::date
		ORDER BY bookings.updated_at::date
	`, userID, startDate, endDate).Scan(&rows).Error
	if err != nil {
		return nil, repoerrors.InternalError("failed to get cancellations per day: " + err.Error())
	}
	return rows, nil
}
