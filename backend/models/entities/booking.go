package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID                  uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AppointmentID       uuid.UUID      `json:"appointment_id" gorm:"type:uuid;not null"`
	Appointment         Appointment    `json:"-" gorm:"foreignKey:AppointmentID"`
	AppCode             string         `json:"app_code" gorm:"not null"`
	UserID              *uuid.UUID     `json:"user_id" gorm:"type:uuid"`
	User                User           `json:"-" gorm:"foreignKey:UserID"`
	Name                string         `json:"name" gorm:""`
	Email               string         `json:"email" gorm:""`
	Phone               string         `json:"phone" gorm:""`
	Date                time.Time      `json:"date" gorm:"not null"`
	StartTime           time.Time      `json:"start_time" gorm:"not null"`
	EndTime             time.Time      `json:"end_time" gorm:"not null"`
	Available           bool           `json:"available" gorm:"not null;default:true"`
	AttendeeCount       int            `json:"attendee_count" gorm:"default:1"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time"`
	BookingCode         string         `json:"booking_code" gorm:"uniqueIndex;not null"` // Permanent booking code for all bookings
	NotificationStatus  string         `json:"notification_status" gorm:"default:''"`
	NotificationChannel string         `json:"notification_channel" gorm:"default:''"`
	Status              string         `json:"status" gorm:"default:'active'"` // Booking status: active, cancelled, etc.
	Description         string         `json:"description" gorm:"type:text"`   // Additional info from the booker
}

func (a *Booking) BeforeCreate(tx *gorm.DB) error {
	// Generate unique appointment code
	code, err := generateUniqueCode(tx, "bookings", "booking_code = ?", Booking{}, "BK")
	if err != nil {
		return err
	}
	a.BookingCode = code

	return nil
}
