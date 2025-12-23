package repository

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	repoerrors "github.com/m13ha/asiko/errors/repoerrors"
	"github.com/m13ha/asiko/models/entities"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	Create(booking *entities.Booking) error
	FindAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error)
	FindAndLockAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error)
	FindAndLockSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error)
	Update(booking *entities.Booking) error
	GetBookingsByAppCode(ctx context.Context, req *http.Request, appCode string, available bool) paginate.Page
	GetBookingsByUserID(ctx context.Context, req *http.Request, userID uuid.UUID, statuses []string) paginate.Page
	GetAvailableSlots(ctx context.Context, req *http.Request, appCode string) paginate.Page
	GetAvailableSlotsByDay(ctx context.Context, req *http.Request, appCode string, date time.Time) paginate.Page
	GetBookingByCode(bookingCode string) (*entities.Booking, error)
	FindActiveBookingByEmail(appointmentID uuid.UUID, email string) (*entities.Booking, error)
	FindActiveBookingByDevice(appointmentID uuid.UUID, deviceID string) (*entities.Booking, error)
	MarkBookingsOngoing(ctx context.Context, now time.Time) (int64, error)
	MarkBookingsExpired(ctx context.Context, now time.Time) (int64, error)
	UpdateNotificationStatus(id uuid.UUID, status string, channel string) error
	GetAvailableDates(ctx context.Context, appCode string) ([]time.Time, error)
	WithTx(tx *gorm.DB) BookingRepository
}

type gormBookingRepository struct {
	db *gorm.DB
}

func NewGormBookingRepository(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{db: db}
}

func (r *gormBookingRepository) WithTx(tx *gorm.DB) BookingRepository {
	return &gormBookingRepository{db: tx}
}

func (r *gormBookingRepository) Create(booking *entities.Booking) error {
	if err := r.db.Select("*").Create(booking).Error; err != nil {
		return repoerrors.InternalError("failed to create booking: " + err.Error())
	}
	return nil
}

func (r *gormBookingRepository) FindAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error) {
	var slot entities.Booking
	if err := r.db.Where("app_code = ? AND date = ? AND start_time = ? AND available = true AND is_slot = true AND seats_booked < capacity", appCode, date, startTime).
		First(&slot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("no available slot found")
		}
		return nil, repoerrors.InternalError("failed to find available slot: " + err.Error())
	}
	return &slot, nil
}

// FindAndLockAvailableSlot locks the row to prevent concurrent updates when reserving a slot
func (r *gormBookingRepository) FindAndLockAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error) {
	var slot entities.Booking
	if err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("app_code = ? AND date = ? AND start_time = ? AND available = true AND is_slot = true AND seats_booked < capacity", appCode, date, startTime).
		First(&slot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("no available slot found for locking")
		}
		return nil, repoerrors.InternalError("failed to find and lock available slot: " + err.Error())
	}
	return &slot, nil
}

func (r *gormBookingRepository) FindAndLockSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error) {
	var slot entities.Booking
	if err := r.db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("app_code = ? AND date = ? AND start_time = ? AND is_slot = true", appCode, date, startTime).
		First(&slot).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("slot not found for locking")
		}
		return nil, repoerrors.InternalError("failed to find and lock slot: " + err.Error())
	}
	return &slot, nil
}

func (r *gormBookingRepository) Update(booking *entities.Booking) error {
	if err := r.db.Save(booking).Error; err != nil {
		return repoerrors.InternalError("failed to update booking: " + err.Error())
	}
	return nil
}

func (r *gormBookingRepository) GetBookingsByAppCode(ctx context.Context, req *http.Request, appCode string, available bool) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("app_code = ? AND available = ?", appCode, available).
		Order("created_at DESC")
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetBookingsByUserID(ctx context.Context, req *http.Request, userID uuid.UUID, statuses []string) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if len(statuses) > 0 {
		db = db.Where("status IN ?", statuses)
	}
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetAvailableSlots(ctx context.Context, req *http.Request, appCode string) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("bookings.app_code = ? AND bookings.available = true AND bookings.is_slot = true", appCode).
		Where("(appointments.type != 'party' AND bookings.seats_booked < bookings.capacity) OR (appointments.type = 'party' AND appointments.attendees_booked < appointments.max_attendees)").
		Order("bookings.date ASC, bookings.start_time ASC").
		Select("bookings.*")
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetAvailableSlotsByDay(ctx context.Context, req *http.Request, appCode string, date time.Time) paginate.Page {
	pg := paginate.New()
	now := time.Now()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("bookings.app_code = ? AND bookings.date >= ? AND bookings.date < ? AND bookings.start_time >= ? AND bookings.available = true AND bookings.is_slot = true", appCode, date, date.AddDate(0, 0, 1), now).
		Where("(appointments.type != 'party' AND bookings.seats_booked < bookings.capacity) OR (appointments.type = 'party' AND appointments.attendees_booked < appointments.max_attendees)").
		Order("bookings.start_time ASC").
		Select("bookings.*")
	var request interface{}
	if req != nil {
		request = req
	} else {
		request = &paginate.Request{}
	}
	return pg.With(db).Request(request).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetBookingByCode(bookingCode string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("booking_code = ?", bookingCode).First(&booking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("booking not found with code: " + bookingCode)
		}
		return nil, repoerrors.InternalError("failed to get booking by code: " + err.Error())
	}
	return &booking, nil
}

func (r *gormBookingRepository) FindActiveBookingByEmail(appointmentID uuid.UUID, email string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("appointment_id = ? AND email = ? AND status IN ?", appointmentID, email, []string{"active", "ongoing", "pending", "confirmed"}).
		First(&booking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("active booking not found")
		}
		return nil, repoerrors.InternalError("failed to find active booking by email: " + err.Error())
	}
	return &booking, nil
}

func (r *gormBookingRepository) FindActiveBookingByDevice(appointmentID uuid.UUID, deviceID string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("appointment_id = ? AND device_id = ? AND status IN ?", appointmentID, deviceID, []string{"active", "ongoing", "pending", "confirmed"}).
		First(&booking).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repoerrors.NotFoundError("active booking not found for device")
		}
		return nil, repoerrors.InternalError("failed to find active booking by device: " + err.Error())
	}
	return &booking, nil
}

func (r *gormBookingRepository) MarkBookingsOngoing(ctx context.Context, now time.Time) (int64, error) {
	startExpr := "date_trunc('day', date) + (start_time - date_trunc('day', start_time))"
	endExpr := "date_trunc('day', date) + (end_time - date_trunc('day', end_time))"
	res := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("available = ?", false).
		Where("status IN ?", []string{"active", "confirmed", "pending"}).
		Where(startExpr+" <= ?", now).
		Where(endExpr+" > ?", now).
		Updates(map[string]interface{}{
			"status":     "ongoing",
			"updated_at": now,
		})
	if res.Error != nil {
		return 0, repoerrors.InternalError("failed to mark bookings ongoing: " + res.Error.Error())
	}
	return res.RowsAffected, nil
}

func (r *gormBookingRepository) MarkBookingsExpired(ctx context.Context, now time.Time) (int64, error) {
	endExpr := "date_trunc('day', date) + (end_time - date_trunc('day', end_time))"
	res := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("available = ?", false).
		Where("status IN ?", []string{"active", "ongoing", "confirmed", "pending"}).
		Where(endExpr+" < ?", now).
		Updates(map[string]interface{}{
			"status":     "expired",
			"updated_at": now,
		})
	if res.Error != nil {
		return 0, repoerrors.InternalError("failed to mark bookings expired: " + res.Error.Error())
	}
	return res.RowsAffected, nil
}

func (r *gormBookingRepository) UpdateNotificationStatus(id uuid.UUID, status string, channel string) error {
	if err := r.db.Model(&entities.Booking{}).Where("id = ?", id).Updates(map[string]interface{}{
		"notification_status":  status,
		"notification_channel": channel,
	}).Error; err != nil {
		return repoerrors.InternalError("failed to update notification status: " + err.Error())
	}
	return nil
}

func (r *gormBookingRepository) GetAvailableDates(ctx context.Context, appCode string) ([]time.Time, error) {
	var dates []time.Time
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Distinct("bookings.date").
		Where("bookings.app_code = ? AND bookings.start_time >= ? AND bookings.available = true AND bookings.is_slot = true", appCode, now).
		Where("(appointments.type != 'party' AND bookings.seats_booked < bookings.capacity) OR (appointments.type = 'party' AND appointments.attendees_booked < appointments.max_attendees)").
		Order("bookings.date ASC").
		Pluck("bookings.date", &dates).Error
	if err != nil {
		return nil, repoerrors.InternalError("failed to get available dates: " + err.Error())
	}
	return dates, nil
}
