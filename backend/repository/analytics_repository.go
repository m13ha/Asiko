package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
	GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error)
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
	return count, err
}

func (r *gormAnalyticsRepository) GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("bookings").
		Joins("JOIN appointments ON bookings.appointment_id = appointments.id").
		Where("appointments.owner_id = ? AND bookings.created_at BETWEEN ? AND ? AND bookings.status = ?", 
			userID, startDate, endDate, "active").
		Count(&count).Error
	return count, err
}