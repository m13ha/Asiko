package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *entities.Booking) error
	FindAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error)
	Update(booking *entities.Booking) error
	GetBookingsByAppCode(ctx context.Context, appCode string, available bool) paginate.Page
	GetBookingsByUserID(ctx context.Context, userID uuid.UUID) paginate.Page
	GetAvailableSlots(ctx context.Context, appCode string) paginate.Page
	GetAvailableSlotsByDay(ctx context.Context, appCode string, date time.Time) paginate.Page
	GetBookingByCode(bookingCode string) (*entities.Booking, error)
	FindActiveBookingByEmail(appointmentID uuid.UUID, email string) (*entities.Booking, error)
	FindActiveBookingByDevice(appointmentID uuid.UUID, deviceID string) (*entities.Booking, error)
	UpdateNotificationStatus(id uuid.UUID, status string, channel string) error
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
	return r.db.Create(booking).Error
}

func (r *gormBookingRepository) FindAvailableSlot(appCode string, date time.Time, startTime time.Time) (*entities.Booking, error) {
	var slot entities.Booking
	if err := r.db.Where("app_code = ? AND date = ? AND start_time = ? AND available = true", appCode, date, startTime).First(&slot).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *gormBookingRepository) Update(booking *entities.Booking) error {
	return r.db.Save(booking).Error
}

func (r *gormBookingRepository) GetBookingsByAppCode(ctx context.Context, appCode string, available bool) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("app_code = ? AND available = ?", appCode, available).
		Order("created_at DESC")
	return pg.With(db).Request(ctx).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetBookingsByUserID(ctx context.Context, userID uuid.UUID) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("user_id = ?", userID).
		Order("created_at DESC")
	return pg.With(db).Request(ctx).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetAvailableSlots(ctx context.Context, appCode string) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("app_code = ? AND available = true", appCode).
		Order("date ASC, start_time ASC")
	return pg.With(db).Request(ctx).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetAvailableSlotsByDay(ctx context.Context, appCode string, date time.Time) paginate.Page {
	pg := paginate.New()
	db := r.db.WithContext(ctx).Model(&entities.Booking{}).
		Where("app_code = ? AND date = ? AND available = true", appCode, date).
		Order("start_time ASC")
	return pg.With(db).Request(ctx).Response(&[]entities.Booking{})
}

func (r *gormBookingRepository) GetBookingByCode(bookingCode string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("booking_code = ?", bookingCode).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *gormBookingRepository) FindActiveBookingByEmail(appointmentID uuid.UUID, email string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("appointment_id = ? AND email = ? AND status = ?", appointmentID, email, "active").First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *gormBookingRepository) FindActiveBookingByDevice(appointmentID uuid.UUID, deviceID string) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.db.Where("appointment_id = ? AND device_id = ? AND status = ?", appointmentID, deviceID, "active").First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *gormBookingRepository) UpdateNotificationStatus(id uuid.UUID, status string, channel string) error {
	return r.db.Model(&entities.Booking{}).Where("id = ?", id).Updates(map[string]interface{}{
		"notification_status":  status,
		"notification_channel": channel,
	}).Error
}