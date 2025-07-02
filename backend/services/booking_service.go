package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/db"
	myerrors "github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/dto"
	"github.com/m13ha/appointment_master/models/entities"
	"github.com/m13ha/appointment_master/utils"
	"github.com/morkid/paginate"
)

// ToBookingResponse converts an entities.Booking to a dto.BookingResponse
func ToBookingResponse(booking *entities.Booking) *dto.BookingResponse {
	return &dto.BookingResponse{
		AppCode:       booking.AppCode,
		ID:            booking.ID,
		AppointmentID: booking.AppointmentID,
		UserID:        booking.UserID,
		Name:          booking.Name,
		Email:         booking.Email,
		Phone:         booking.Phone,
		Date:          booking.Date,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		AttendeeCount: booking.AttendeeCount,
		CreatedAt:     booking.CreatedAt,
		UpdatedAt:     booking.UpdatedAt,
		Description:   booking.Description,
	}
}

func bookSlot(req dto.BookingRequest, slot *entities.Booking, appointment *entities.Appointment) (*dto.BookingResponse, error) {
	// Handle capacity for group appointments
	if appointment.Type == entities.Group {
		if req.AttendeeCount > appointment.MaxAttendees {
			return nil, myerrors.NewUserError("Attendee count exceeds maximum allowed.")
		}
		slot.AttendeeCount = req.AttendeeCount
	} else {
		slot.AttendeeCount = 1
	}

	// Set description from request
	slot.Description = req.Description

	slot.Available = false
	// Generate and assign a globally unique booking code if not already set
	if slot.BookingCode == "" {
		code, err := generateUniqueBookingCode()
		if err != nil {
			log.Printf("[bookSlot] Internal error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		slot.BookingCode = code
	}
	if err := db.DB.Save(&slot).Error; err != nil {
		log.Printf("[bookSlot] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}

	return ToBookingResponse(slot), nil
}

func BookRegisteredUserAppointment(req dto.BookingRequest, userIDStr string) (*dto.BookingResponse, error) {
	if err := utils.Validate(req); err != nil {
		return nil, myerrors.NewUserError("Invalid booking data. Please check your input.")
	}

	// Check if appointment exists
	var appointment entities.Appointment
	if err := db.DB.Where("app_code = ?", req.AppCode).First(&appointment).Error; err != nil {
		return nil, myerrors.NewUserError("Appointment not found.")
	}

	// Find matching slot
	var slot entities.Booking
	if err := db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error; err != nil {
		return nil, myerrors.NewUserError("No available slot found.")
	}

	// Use user details unless overridden by request
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, myerrors.NewUserError("Invalid user ID.")
	}
	var user entities.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, myerrors.NewUserError("User not found.")
	}
	slot.UserID = &userID
	slot.Name = user.Name
	slot.Email = user.Email
	slot.Phone = user.PhoneNumber

	return bookSlot(req, &slot, &appointment)
}

// BookGuestAppointment books an appointment for a guest
func BookGuestAppointment(req dto.BookingRequest) (*dto.BookingResponse, error) {
	if err := utils.Validate(req); err != nil {
		return nil, myerrors.NewUserError("Invalid booking data. Please check your input.")
	}

	// Check if appointment exists
	var appointment entities.Appointment
	if err := db.DB.Where("app_code = ?", req.AppCode).First(&appointment).Error; err != nil {
		return nil, myerrors.NewUserError("Appointment not found.")
	}

	// Find matching slot
	var slot entities.Booking
	if err := db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error; err != nil {
		return nil, myerrors.NewUserError("No available slot found.")
	}

	// Guest bookings require name and either email or phone
	if req.Name == "" || (req.Email == "" && req.Phone == "") {
		return nil, myerrors.NewUserError("Name and either email or phone are required for guest bookings.")
	}
	slot.Name = req.Name
	slot.Email = utils.NormalizeEmail(req.Email)
	slot.Phone = req.Phone

	return bookSlot(req, &slot, &appointment)
}

func GetAllBookingsForAppointment(appcode string, r *http.Request) (any, error) {
	query := db.DB.Model(&entities.Booking{}).Where("app_code = ? AND available = false", appcode)
	if r == nil {
		var bookings []entities.Booking
		if err := query.Find(&bookings).Error; err != nil {
			log.Printf("[GetAllBookingsForAppointment] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		return bookings, nil
	}
	p := paginate.New()
	result := p.With(query).Request(r).Response(&[]dto.BookingResponse{})
	return &result, nil
}

// GetUserBookings returns all bookings for a user
func GetUserBookings(userID string, r *http.Request) (any, error) {
	query := db.DB.Model(&entities.Booking{}).Where("user_id = ?", userID)
	if r == nil {
		var bookings []entities.Booking
		if err := query.Find(&bookings).Error; err != nil {
			log.Printf("[GetUserBookings] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		return bookings, nil
	}
	p := paginate.New()
	result := p.With(query).Request(r).Response(&[]dto.BookingResponse{})
	return &result, nil
}

func GetAvailableSlots(appcode string, r *http.Request) (any, error) {
	query := db.DB.Model(&entities.Booking{}).Where("app_code = ? AND available = true", appcode)
	if r == nil {
		var slots []entities.Booking
		if err := query.Find(&slots).Error; err != nil {
			log.Printf("[GetAvailableSlots] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		return slots, nil
	}
	p := paginate.New()
	result := p.With(query).Request(r).Response(&[]dto.BookingResponse{})
	return &result, nil
}

// GetAvailableSlotsByDay returns available slots for an appointment on a specific day
func GetAvailableSlotsByDay(appcode string, dateStr string, r *http.Request) (any, error) {
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, myerrors.NewUserError("Invalid date format. Use YYYY-MM-DD.")
	}
	query := db.DB.Model(&entities.Booking{}).Where("app_code = ? AND date = ? AND available = true", appcode, parsedDate)
	if r == nil {
		var slots []entities.Booking
		if err := query.Find(&slots).Error; err != nil {
			log.Printf("[GetAvailableSlotsByDay] DB error: %v", err)
			return nil, fmt.Errorf("internal error")
		}
		return slots, nil
	}
	p := paginate.New()
	result := p.With(query).Request(r).Response(&[]dto.BookingResponse{})
	return &result, nil
}

// GetBookingByCode retrieves a booking by its permanent booking_code
func GetBookingByCode(bookingCode string) (*entities.Booking, error) {
	var booking entities.Booking
	if err := db.DB.Where("booking_code = ?", bookingCode).First(&booking).Error; err != nil {
		return nil, myerrors.NewUserError("Booking not found.")
	}
	return &booking, nil
}

// UpdateBookingByCode allows rescheduling a booking if the new slot is available
func UpdateBookingByCode(bookingCode string, req dto.BookingRequest) (*dto.BookingResponse, error) {
	booking, err := GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}
	var slot entities.Booking
	err = db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error
	if err != nil {
		return nil, myerrors.NewUserError("Requested slot is not available.")
	}
	booking.Available = true
	db.DB.Save(booking)
	booking.Date = req.Date
	booking.StartTime = req.StartTime
	booking.EndTime = req.EndTime
	booking.AttendeeCount = req.AttendeeCount
	booking.Description = req.Description
	booking.Available = false
	booking.Status = "active"
	if err := db.DB.Save(booking).Error; err != nil {
		log.Printf("[UpdateBookingByCode] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}
	return ToBookingResponse(booking), nil
}

// CancelBookingByCode cancels a booking by booking_code
func CancelBookingByCode(bookingCode string) (*dto.BookingResponse, error) {
	booking, err := GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}
	booking.Available = true
	booking.Status = "cancelled"
	if err := db.DB.Save(booking).Error; err != nil {
		log.Printf("[CancelBookingByCode] DB error: %v", err)
		return nil, fmt.Errorf("internal error")
	}
	return ToBookingResponse(booking), nil
}

// isBookingCodeAvailable checks if a BookingCode is available for use (not in use by an active or recently expired booking)
func isBookingCodeAvailable(bookingCode string) (bool, error) {
	var booking entities.Booking
	err := db.DB.Where("booking_code = ?", bookingCode).Order("end_time desc").First(&booking).Error
	if err != nil {
		// Not found, so available
		return true, nil
	}
	now := time.Now()
	holdUntil := booking.EndTime.Add(7 * 24 * time.Hour)
	if now.Before(holdUntil) {
		return false, nil
	}
	return true, nil
}

// generateUniqueBookingCode generates a globally unique BookingCode, respecting the 2-week hold after expiry
func generateUniqueBookingCode() (string, error) {
	for i := 0; i < 10; i++ {
		code := utils.GenerateBookingCode()
		available, err := isBookingCodeAvailable(code)
		if err != nil {
			return "", err
		}
		if available {
			return code, nil
		}
	}
	return "", fmt.Errorf("could not generate unique BookingCode after 10 attempts")
}
