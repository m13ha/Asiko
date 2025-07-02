package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/utils"
	"gorm.io/gorm"
)

type AppointmentType string

const (
	Single AppointmentType = "single"
	Group  AppointmentType = "group"
)

type Appointment struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title           string          `json:"title" gorm:"not null"`
	StartTime       time.Time       `json:"start_time" gorm:"not null"`
	EndTime         time.Time       `json:"end_time" gorm:"not null"`
	StartDate       time.Time       `json:"start_date" gorm:"not null"`
	EndDate         time.Time       `json:"end_date" gorm:"not null"`
	BookingDuration int             `json:"booking_duration" gorm:"not null"` // in minutes
	MaxAttendees    int             `json:"max_attendees" gorm:"default:1"`   // for group appointments
	Type            AppointmentType `json:"type" gorm:"not null;default:'single'"`
	OwnerID         uuid.UUID       `json:"owner_id" gorm:"type:uuid;not null"`
	User            User            `json:"-" gorm:"foreignKey:OwnerID"`
	AppCode         string          `json:"app_code" gorm:"unique;not null"`
	Bookings        []Booking       `json:"bookings" gorm:"foreignKey:AppointmentID"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
	Description     string          `json:"description" gorm:"type:text"` // Additional info for the appointment
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	// Generate unique appointment code
	code := utils.GenerateAppCode()
	a.AppCode = code

	// Normalize the appointment type to lowercase
	a.Type = AppointmentType(utils.NormalizeString(string(a.Type)))

	return nil
}

func (a *Appointment) AfterCreate(tx *gorm.DB) error {
	slots := a.generateBookings()
	return tx.Create(&slots).Error
}

func (a *Appointment) generateBookings() []Booking {
	var slots []Booking
	duration := time.Duration(a.BookingDuration) * time.Minute

	for currentDate := a.StartDate; !currentDate.After(a.EndDate); currentDate = currentDate.AddDate(0, 0, 1) {
		currentSlotStart := time.Date(
			currentDate.Year(), currentDate.Month(), currentDate.Day(),
			a.StartTime.Hour(), a.StartTime.Minute(), 0, 0,
			currentDate.Location(),
		)
		dailyEndTime := time.Date(
			currentDate.Year(), currentDate.Month(), currentDate.Day(),
			a.EndTime.Hour(), a.EndTime.Minute(), 0, 0,
			currentDate.Location(),
		)

		for currentSlotStart.Before(dailyEndTime) {
			slotEnd := currentSlotStart.Add(duration)
			if slotEnd.After(dailyEndTime) {
				break
			}

			slots = append(slots, Booking{
				AppointmentID: a.ID,
				AppCode:       a.AppCode,
				Date:          currentDate,
				StartTime:     currentSlotStart,
				EndTime:       slotEnd,
				Available:     true, // Changed to true by default
			})
			currentSlotStart = slotEnd
		}
	}
	return slots
}
