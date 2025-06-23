package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/db"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/utils"
)

func CreateAppointment(req models.AppointmentRequest, userId uuid.UUID) (*models.Appointment, error) {
	// Validate request
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, fmt.Errorf("end time cannot be before start time")
	}

	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Validate booking duration fits within time window
	duration := req.EndTime.Sub(req.StartTime)
	if duration.Minutes() < float64(req.BookingDuration) {
		return nil, fmt.Errorf("booking duration exceeds available time window")
	}

	appointment := &models.Appointment{
		Title:           req.Title,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		BookingDuration: req.BookingDuration,
		Type:            models.AppointmentType(utils.NormalizeString(string(req.Type))),
		MaxAttendees:    req.MaxAttendees,
		OwnerID:         userId,
		Description:     req.Description,
	}

	if err := db.DB.Create(appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func BookRegisteredUserAppointment(req models.BookingRequest, userIDStr string) (*models.Booking, error) {
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if appointment exists
	var appointment models.Appointment
	if err := db.DB.Where("app_code = ?", req.AppCode).First(&appointment).Error; err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	// Find matching slot
	var slot models.Booking
	if err := db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error; err != nil {
		return nil, fmt.Errorf("no available slot found: %w", err)
	}

	// Handle capacity for group appointments
	if appointment.Type == models.Group {
		if req.AttendeeCount > appointment.MaxAttendees {
			return nil, fmt.Errorf("attendee count exceeds maximum allowed")
		}
		slot.AttendeeCount = req.AttendeeCount
	} else {
		slot.AttendeeCount = 1
	}

	// Use user details unless overridden by request
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	slot.UserID = &userID
	slot.Name = user.Name
	slot.Email = user.Email
	slot.Phone = user.PhoneNumber

	// Set description from request
	slot.Description = req.Description

	slot.Available = false
	// Generate and assign a permanent booking code if not already set
	if slot.BookingCode == "" {
		slot.BookingCode = utils.GenerateBookingCode()
	}
	if err := db.DB.Save(&slot).Error; err != nil {
		return nil, fmt.Errorf("failed to book appointment: %w", err)
	}

	return &slot, nil
}

// BookGuestAppointment books an appointment for a guest
func BookGuestAppointment(req models.BookingRequest) (*models.Booking, error) {
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if appointment exists
	var appointment models.Appointment
	if err := db.DB.Where("app_code = ?", req.AppCode).First(&appointment).Error; err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	// Find matching slot
	var slot models.Booking
	if err := db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error; err != nil {
		return nil, fmt.Errorf("no available slot found: %w", err)
	}

	// Handle capacity for group appointments
	if appointment.Type == models.Group {
		if req.AttendeeCount > appointment.MaxAttendees {
			return nil, fmt.Errorf("attendee count exceeds maximum allowed")
		}
		slot.AttendeeCount = req.AttendeeCount
	} else {
		slot.AttendeeCount = 1
	}

	// Guest bookings require name and either email or phone
	if req.Name == "" || (req.Email == "" && req.Phone == "") {
		return nil, fmt.Errorf("name and either email or phone are required for guest bookings")
	}
	slot.Name = req.Name
	slot.Email = utils.NormalizeEmail(req.Email)
	slot.Phone = req.Phone

	// Set description from request
	slot.Description = req.Description

	slot.Available = false
	// Generate and assign a permanent booking code if not already set
	if slot.BookingCode == "" {
		slot.BookingCode = utils.GenerateBookingCode()
	}
	if err := db.DB.Save(&slot).Error; err != nil {
		return nil, fmt.Errorf("failed to book appointment: %w", err)
	}

	return &slot, nil
}

func GetAllBookingsForAppointment(appcode string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := db.DB.Where("app_code = ? AND available = false", appcode).Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	return bookings, nil
}

func GetAvailableSlots(appcode string) ([]models.Booking, error) {
	var slots []models.Booking
	if err := db.DB.Where("app_code = ? AND available = true", appcode).
		Find(&slots).Error; err != nil {
		return nil, fmt.Errorf("failed to get available slots: %w", err)
	}
	return slots, nil
}

func GetAllAppointmentsCreatedByUser(userID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	if err := db.DB.Where("owner_id = ?", userID).Find(&appointments).Error; err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}
	return appointments, nil
}

// GetBookingByCode retrieves a booking by its permanent booking_code
func GetBookingByCode(bookingCode string) (*models.Booking, error) {
	var booking models.Booking
	if err := db.DB.Where("booking_code = ?", bookingCode).First(&booking).Error; err != nil {
		return nil, fmt.Errorf("booking not found: %w", err)
	}
	return &booking, nil
}

// UpdateBookingByCode allows rescheduling a booking if the new slot is available
func UpdateBookingByCode(bookingCode string, req models.BookingRequest) (*models.Booking, error) {
	booking, err := GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}
	// Check if new slot is available
	var slot models.Booking
	err = db.DB.Where("app_code = ? AND date = ? AND start_time = ? AND available = true",
		req.AppCode, req.Date, req.StartTime).First(&slot).Error
	if err != nil {
		return nil, fmt.Errorf("requested slot is not available: %w", err)
	}
	// Mark old slot as available
	booking.Available = true
	db.DB.Save(booking)
	// Update booking to new slot
	booking.Date = req.Date
	booking.StartTime = req.StartTime
	booking.EndTime = req.EndTime
	booking.AttendeeCount = req.AttendeeCount
	booking.Description = req.Description
	booking.Available = false
	booking.Status = "active"
	if err := db.DB.Save(booking).Error; err != nil {
		return nil, fmt.Errorf("failed to reschedule booking: %w", err)
	}
	return booking, nil
}

// CancelBookingByCode cancels a booking by booking_code
func CancelBookingByCode(bookingCode string) (*models.Booking, error) {
	booking, err := GetBookingByCode(bookingCode)
	if err != nil {
		return nil, err
	}
	booking.Available = true
	booking.Status = "cancelled"
	if err := db.DB.Save(booking).Error; err != nil {
		return nil, fmt.Errorf("failed to cancel booking: %w", err)
	}
	return booking, nil
}
