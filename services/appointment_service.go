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
		Type:            req.Type,
		MaxAttendees:    req.MaxAttendees,
		OwnerID:         userId,
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

	slot.Available = false
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
	slot.Email = req.Email
	slot.Phone = req.Phone

	slot.Available = false
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
