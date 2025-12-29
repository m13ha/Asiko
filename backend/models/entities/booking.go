package entities

import (
	"fmt"
	"strings"
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
	IsSlot              bool           `json:"is_slot" gorm:"not null;default:false"`
	Capacity            int            `json:"capacity" gorm:"not null;default:1"`
	SeatsBooked         int            `json:"seats_booked" gorm:"not null;default:0"`
	AttendeeCount       int            `json:"attendee_count" gorm:"default:1"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time"`
	BookingCode         string         `json:"booking_code" gorm:"uniqueIndex;not null"` // Permanent booking code for all bookings
	NotificationStatus  string         `json:"notification_status" gorm:"default:''"`
	NotificationChannel string         `json:"notification_channel" gorm:"default:''"`
	Status              string         `json:"status" gorm:"default:'active'"` // Booking status: active, cancelled, etc.
	Description         string         `json:"description" gorm:"type:text"`   // Additional info from the booker
	DeviceID            string         `json:"-"`
}

func (a *Booking) BeforeCreate(tx *gorm.DB) error {
	if a.Capacity < 1 {
		a.Capacity = 1
	}
	if a.SeatsBooked < 0 {
		a.SeatsBooked = 0
	}
	if !a.IsSlot {
		a.Available = false
	}
	if a.IsSlot && strings.TrimSpace(a.Status) == "" {
		a.Status = BookingStatusActive
	}
	if !a.IsSlot && strings.TrimSpace(a.Status) == "" {
		a.Status = BookingStatusPending
	}
	if a.BookingCode == "" {
		a.BookingCode = generateDeterministicBookingCode(a)
	}
	return nil
}

func generateDeterministicBookingCode(b *Booking) string {
	datePart := b.Date.Format("060102")
	timePart := b.StartTime.Format("1504")
	base := fmt.Sprintf("BK%s%s", datePart, timePart)
	if b.ID != uuid.Nil {
		return fmt.Sprintf("%s%s", base, shortUUID(b.ID))
	}
	return fmt.Sprintf("%s%s", base, uuid.New().String()[:4])
}

func shortUUID(id uuid.UUID) string {
	return id.String()[:4]
}

func (b *Booking) NormalizeState() {
	if b.Capacity < 1 {
		b.Capacity = 1
	}
	if b.SeatsBooked < 0 {
		b.SeatsBooked = 0
	}
	if b.SeatsBooked >= b.Capacity {
		b.Available = false
	} else {
		b.Available = true
	}
	remaining := b.Capacity - b.SeatsBooked
	if remaining < 0 {
		remaining = 0
	}
	b.AttendeeCount = remaining
}
