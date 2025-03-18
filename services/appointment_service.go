package services

import (
	"fmt"

	"github.com/m13ha/appointment_master/db"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/utils"
)

func CreateAppointment(req models.AppointmentRequest) (*models.Appointment, error) {
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
		OwnerID:         req.UserID,
	}

	if err := db.DB.Create(appointment).Error; err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func BookAppointment(req models.BookingRequest) (*models.Booking, error) {
	if err := utils.Validate(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if appointment exists and get its details
	var appointment models.Appointment
	if err := db.DB.First(&appointment, req.AppointmentID).Error; err != nil {
		return nil, fmt.Errorf("appointment not found: %w", err)
	}

	// Find matching slot
	var slot models.Booking
	if err := db.DB.Where("appointment_id = ? AND date = ? AND start_time = ? AND available = true",
		req.AppointmentID, req.Date, req.StartTime).First(&slot).Error; err != nil {
		return nil, fmt.Errorf("no available slot found: %w", err)
	}

	// For group appointments, check capacity
	if appointment.Type == models.Group {
		if req.AttendeeCount > appointment.MaxAttendees {
			return nil, fmt.Errorf("attendee count exceeds maximum allowed")
		}
		slot.AttendeeCount = req.AttendeeCount
	} else {
		slot.AttendeeCount = 1 // Single appointments always have 1 attendee
	}

	// Update slot with booking details
	slot.Available = false
	if req.GuestName != "" { // Guest booking
		slot.GuestName = req.GuestName
		slot.GuestEmail = req.GuestEmail
		slot.GuestPhone = req.GuestPhone
	} else { // Registered user booking
		userID := req.UserID
		slot.UserID = &userID
	}

	if err := db.DB.Save(&slot).Error; err != nil {
		return nil, fmt.Errorf("failed to book appointment: %w", err)
	}

	return &slot, nil
}

func GetAllBookingsForAppointment(appointmentID string) ([]models.Booking, error) {
	var bookings []models.Booking
	if err := db.DB.Where("appointment_id = ?", appointmentID).Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}
	return bookings, nil
}

func GetAvailableSlots(appointmentID string) ([]models.Booking, error) {
	var slots []models.Booking
	if err := db.DB.Where("appointment_id = ? AND available = true", appointmentID).
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
