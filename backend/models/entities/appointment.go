package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/asiko/utils"
	"gorm.io/gorm"
)

type AppointmentType string

const (
	Single AppointmentType = "single"
	Group  AppointmentType = "group"
	Party  AppointmentType = "party"
)

type AntiScalpingLevel string

const (
	ScalpingNone     AntiScalpingLevel = "none"
	ScalpingStandard AntiScalpingLevel = "standard"
	ScalpingStrict   AntiScalpingLevel = "strict"
)

type Appointment struct {
	ID                uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title             string            `json:"title" gorm:"not null"`
	StartTime         time.Time         `json:"start_time" gorm:"not null"`
	EndTime           time.Time         `json:"end_time" gorm:"not null"`
	StartDate         time.Time         `json:"start_date" gorm:"not null"`
	EndDate           time.Time         `json:"end_date" gorm:"not null"`
	BookingDuration   int               `json:"booking_duration" gorm:"not null"` // in minutes
	MaxAttendees      int               `json:"max_attendees" gorm:"default:1"`   // for group appointments
	Type              AppointmentType   `json:"type" gorm:"not null;default:'single'"`
	AntiScalpingLevel AntiScalpingLevel `json:"anti_scalping_level" gorm:"type:anti_scalping_level;not null;default:'none'"`
	OwnerID           uuid.UUID         `json:"owner_id" gorm:"type:uuid;not null"`
	User              User              `json:"-" gorm:"foreignKey:OwnerID"`
	AppCode           string            `json:"app_code" gorm:"unique;not null"`
	Bookings          []Booking         `json:"bookings" gorm:"foreignKey:AppointmentID"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         gorm.DeletedAt    `json:"deleted_at,omitempty" gorm:"index" swaggertype:"string" format:"date-time"`
	Description       string            `json:"description" gorm:"type:text"` // Additional info for the appointment
	AttendeesBooked   int               `json:"attendees_booked" gorm:"default:0"`
}

func (a *Appointment) BeforeCreate(tx *gorm.DB) error {
	// Generate unique appointment code
	code, err := generateUniqueCode(tx, "appointments", "app_code = ?", Appointment{}, "AP")
	if err != nil {
		return err
	}
	a.AppCode = code

	a.Type = AppointmentType(utils.NormalizeString(string(a.Type)))

	return nil
}

func (a *Appointment) AfterCreate(tx *gorm.DB) error {
	if a.Type == Party {
		return nil
	}
	slots := a.generateBookings()
	return tx.Create(&slots).Error
}

func (a *Appointment) generateBookings() []Booking {
	var slots []Booking
	duration := time.Duration(a.BookingDuration) * time.Minute
	defaultCapacity := 1
	if a.Type == Group && a.MaxAttendees > 0 {
		defaultCapacity = a.MaxAttendees
	}

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
				Available:     true,
				IsSlot:        true,
				Capacity:      defaultCapacity,
				SeatsBooked:   0,
			})
			currentSlotStart = slotEnd
		}
	}
	return slots
}

func isAppCodeAvailable(tx *gorm.DB, code string, table string, query string, entity interface{}) (bool, error) {
	err := tx.Table(table).Where(query, code).Order("end_date desc").First(&entity).Error
	if err != nil {
		return true, nil
	}

	return false, nil
}

func generateUniqueCode(tx *gorm.DB, table string, query string, entity interface{}, codeType string) (string, error) {
	for i := 0; i < 10; i++ {
		code := utils.GenerateCode(codeType)
		available, err := isAppCodeAvailable(tx, code, table, query, entity)
		if err != nil {
			return "", err
		}
		if available {
			return code, nil
		}
	}
	return "", fmt.Errorf("could not generate unique AppCode after 10 attempts")
}
