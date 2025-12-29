package entities

import (
	"errors"
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

type AppointmentStatus string

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
	Status            AppointmentStatus `json:"status" gorm:"type:appointment_status;not null;default:'pending'"`
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
	if a.Status == "" {
		a.Status = AppointmentStatusPending
	}

	return nil
}

func (a *Appointment) AfterCreate(tx *gorm.DB) error {
	slots := a.GenerateBookings()
	return tx.Create(&slots).Error
}

func (a *Appointment) GenerateBookings() []Booking {
	var slots []Booking
	duration := time.Duration(a.BookingDuration) * time.Minute
	defaultCapacity := 1
	if (a.Type == Group || a.Type == Party) && a.MaxAttendees > 0 {
		defaultCapacity = a.MaxAttendees
	}

	startDate := a.StartDate.UTC()
	endDate := a.EndDate.UTC()
	startTime := a.StartTime.UTC()
	endTime := a.EndTime.UTC()

	if a.Type == Party {
		startDateTime := time.Date(
			startDate.Year(), startDate.Month(), startDate.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0,
			time.UTC,
		)
		endDateTime := time.Date(
			endDate.Year(), endDate.Month(), endDate.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0,
			time.UTC,
		)
		if endDateTime.After(startDateTime) {
			slots = append(slots, Booking{
				AppointmentID: a.ID,
				AppCode:       a.AppCode,
				Date:          startDate,
				StartTime:     startDateTime,
				EndTime:       endDateTime,
				Available:     true,
				IsSlot:        true,
				Capacity:      defaultCapacity,
				SeatsBooked:   0,
			})
		}
		return slots
	}

	overnight := !endTime.After(startTime)
	appendSlots := func(date time.Time, windowStart time.Time, windowEnd time.Time) {
		currentSlotStart := windowStart
		for currentSlotStart.Before(windowEnd) {
			slotEnd := currentSlotStart.Add(duration)
			if slotEnd.After(windowEnd) {
				break
			}

			slots = append(slots, Booking{
				AppointmentID: a.ID,
				AppCode:       a.AppCode,
				Date:          date,
				StartTime:     currentSlotStart.UTC(),
				EndTime:       slotEnd.UTC(),
				Available:     true,
				IsSlot:        true,
				Capacity:      defaultCapacity,
				SeatsBooked:   0,
			})
			currentSlotStart = slotEnd
		}
	}

	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if !overnight {
			windowStart := time.Date(
				currentDate.Year(), currentDate.Month(), currentDate.Day(),
				startTime.Hour(), startTime.Minute(), 0, 0,
				time.UTC,
			)
			windowEnd := time.Date(
				currentDate.Year(), currentDate.Month(), currentDate.Day(),
				endTime.Hour(), endTime.Minute(), 0, 0,
				time.UTC,
			)
			appendSlots(currentDate, windowStart, windowEnd)
			continue
		}

		windowStart := time.Date(
			currentDate.Year(), currentDate.Month(), currentDate.Day(),
			startTime.Hour(), startTime.Minute(), 0, 0,
			time.UTC,
		)
		windowEnd := time.Date(
			currentDate.Year(), currentDate.Month(), currentDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		).AddDate(0, 0, 1)
		appendSlots(currentDate, windowStart, windowEnd)

		nextDate := currentDate.AddDate(0, 0, 1)
		if nextDate.After(endDate) {
			continue
		}
		nextWindowStart := time.Date(
			nextDate.Year(), nextDate.Month(), nextDate.Day(),
			0, 0, 0, 0,
			time.UTC,
		)
		nextWindowEnd := time.Date(
			nextDate.Year(), nextDate.Month(), nextDate.Day(),
			endTime.Hour(), endTime.Minute(), 0, 0,
			time.UTC,
		)
		appendSlots(nextDate, nextWindowStart, nextWindowEnd)
	}
	return slots
}

func isAppCodeAvailable(tx *gorm.DB, code string, table string, query string, entity interface{}) (bool, error) {
	err := tx.Table(table).Where(query, code).First(&entity).Error
	if err == nil {
		return false, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
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

func (a *Appointment) IsBookable() (bool, string) {
	switch a.Status {
	case AppointmentStatusCanceled:
		return false, "appointment has been canceled"
	case AppointmentStatusCompleted:
		return false, "appointment is no longer active"
	default:
		return true, ""
	}
}
