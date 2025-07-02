package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m13ha/appointment_master/models/dto"
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

	bookings, err := services.GetUserBookings(userIDStr, r)
	HandleServiceError(w, err, http.StatusInternalServerError, writeError)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// BookAppointment handles both registered user and guest bookings
func BookGuestAppointment(w http.ResponseWriter, r *http.Request) {
	var req dto.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	bookingResponse, err := services.BookGuestAppointment(req)
	if err != nil {
		HandleServiceError(w, err, http.StatusBadRequest, writeError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookingResponse)
}

// BookRegisteredUserAppointment handles registered user bookings
func BookRegisteredUserAppointment(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dto.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	bookingResponse, err := services.BookRegisteredUserAppointment(req, userIDStr)
	if err != nil {
		HandleServiceError(w, err, http.StatusBadRequest, writeError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bookingResponse)
}

func GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	appcode := chi.URLParam(r, "id")
	if appcode == "" {
		writeError(w, http.StatusBadRequest, "Missing appointment code parameter")
		return
	}

	slots, err := services.GetAvailableSlots(appcode, r)
	HandleServiceError(w, err, http.StatusInternalServerError, writeError)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

// GetAvailableSlotsByDay returns available slots for an appointment on a specific day
func GetAvailableSlotsByDay(w http.ResponseWriter, r *http.Request) {
	appcode := chi.URLParam(r, "id")
	if appcode == "" {
		writeError(w, http.StatusBadRequest, "Missing appointment code parameter")
		return
	}
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		writeError(w, http.StatusBadRequest, "Missing date parameter")
		return
	}
	slots, err := services.GetAvailableSlotsByDay(appcode, dateStr, r)
	HandleServiceError(w, err, http.StatusInternalServerError, writeError)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

func GetUsersRegisteredForAppointment(w http.ResponseWriter, r *http.Request) {
	app_code := chi.URLParam(r, "id")
	if app_code == "" {
		writeError(w, http.StatusBadRequest, "Missing appointment code parameter")
		return
	}
	bookings, err := services.GetAllBookingsForAppointment(app_code, r)
	HandleServiceError(w, err, http.StatusInternalServerError, writeError)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// GetBookingByCodeHandler returns booking details by booking_code
func GetBookingByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "booking_code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "Missing booking_code parameter")
		return
	}
	booking, err := services.GetBookingByCode(code)
	HandleServiceError(w, err, http.StatusNotFound, writeError)
	if err != nil {
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
	var req dto.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}
	bookingResponse, err := services.UpdateBookingByCode(code, req)
	if err != nil {
		HandleServiceError(w, err, http.StatusBadRequest, writeError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookingResponse)
}

// CancelBookingByCodeHandler cancels a booking by booking_code
func CancelBookingByCodeHandler(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "booking_code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "Missing booking_code parameter")
		return
	}
	bookingResponse, err := services.CancelBookingByCode(code)
	if err != nil {
		HandleServiceError(w, err, http.StatusBadRequest, writeError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookingResponse)
}
