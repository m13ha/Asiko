package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

// Helper for consistent error responses
func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func GetUserRegisteredBookings(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	bookings, err := services.GetUserBookings(userIDStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve bookings: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// BookAppointment handles both registered user and guest bookings
func BookGuestAppointment(w http.ResponseWriter, r *http.Request) {
	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	booking, err := services.BookGuestAppointment(req)
	// Handle specific errors from the service layer
	if err != nil {
		switch {
		case err.Error() == "no available slot found":
			writeError(w, http.StatusNotFound, "No available slots")
		case err.Error() == "attendee count exceeds maximum allowed":
			writeError(w, http.StatusBadRequest, err.Error())
		case err.Error() == "name and either email or phone are required for guest bookings":
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to book appointment: "+err.Error())
		}
		return
	}

	response := map[string]interface{}{
		"booking_code": booking.BookingCode,
		"message":      "Booking created successfully",
		"booking": models.BookingResponse{
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
			AppCode:       booking.AppCode,
			Description:   booking.Description,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// BookRegisteredUserAppointment handles registered user bookings
func BookRegisteredUserAppointment(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	booking, err := services.BookRegisteredUserAppointment(req, userIDStr)
	if err != nil {
		switch {
		case err.Error() == "no available slot found":
			writeError(w, http.StatusNotFound, "No available slots")
		case err.Error() == "attendee count exceeds maximum allowed":
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to book appointment: "+err.Error())
		}
		return
	}

	response := map[string]interface{}{
		"booking_code": booking.BookingCode,
		"message":      "Booking created successfully",
		"booking": models.BookingResponse{
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
			AppCode:       booking.AppCode,
			Description:   booking.Description,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	appcode := chi.URLParam(r, "id")
	if appcode == "" {
		writeError(w, http.StatusBadRequest, "Missing appointment code parameter")
		return
	}

	slots, err := services.GetAvailableSlots(appcode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get available slots: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// GetBookingByCodeHandler returns booking details by booking_code
func GetBookingByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "booking_code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "Missing booking_code parameter")
		return
	}
	booking, err := services.GetBookingByCode(code)
	if err != nil {
		writeError(w, http.StatusNotFound, "Booking not found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

// UpdateBookingByCodeHandler reschedules a booking if slot is available
func UpdateBookingByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "booking_code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "Missing booking_code parameter")
		return
	}
	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	booking, err := services.UpdateBookingByCode(code, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := map[string]interface{}{
		"booking_code": booking.BookingCode,
		"message":      "Booking updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// CancelBookingByCodeHandler cancels a booking by booking_code
func CancelBookingByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "booking_code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "Missing booking_code parameter")
		return
	}
	booking, err := services.CancelBookingByCode(code)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	resp := map[string]interface{}{
		"booking_code": booking.BookingCode,
		"message":      "Booking cancelled successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
