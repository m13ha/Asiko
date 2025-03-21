package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

func GetUserRegisteredBookings(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	bookings, err := services.GetUserBookings(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve bookings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

// BookAppointment handles both registered user and guest bookings
func BookAppointment(w http.ResponseWriter, r *http.Request) {
	userIDStr := ExtractUserIDfromContext(r)

	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("USER IS STRING: %s", userIDStr)
	if userIDStr == "" {
		if req.Name == "" || (req.Email == "" && req.Phone == "") {
			http.Error(w, "Name and either email or phone are required for guest bookings", http.StatusBadRequest)
			return
		}
	}

	if err := utils.Validate(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(formatValidationErrors(err))
		return
	}

	booking, err := services.BookAppointment(req, userIDStr)
	if err != nil {
		switch {
		case err.Error() == "no available slot found":
			http.Error(w, "No available slots", http.StatusNotFound)
		case err.Error() == "attendee count exceeds maximum allowed":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Failed to book appointment: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := models.BookingResponse{
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
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	appcode := chi.URLParam(r, "id")

	slots, err := services.GetAvailableSlots(appcode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get available slots: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}
