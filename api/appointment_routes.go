package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

var (
	jwtKey          = []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenExpiration = time.Hour * 24
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Context key for user ID
type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Read and log raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Raw request body: %s", body)

	// Reset body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.AppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decoding error: %v", err) // Log detailed error
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if err := utils.Validate(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(formatValidationErrors(err))
		return
	}

	appointment, err := services.CreateAppointment(req, userID)
	if err != nil {
		switch {
		case err.Error() == "end time cannot be before start time" ||
			err.Error() == "end date cannot be before start date" ||
			err.Error() == "booking duration exceeds available time window":
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Failed to create appointment: %v", err), http.StatusInternalServerError)
		}
		return
	}

	response := models.AppointmentResponse{
		ID:              appointment.ID,
		Title:           appointment.Title,
		StartTime:       appointment.StartTime,
		EndTime:         appointment.EndTime,
		StartDate:       appointment.StartDate,
		EndDate:         appointment.EndDate,
		BookingDuration: appointment.BookingDuration,
		Type:            appointment.Type,
		MaxAttendees:    appointment.MaxAttendees,
		AppCode:         appointment.AppCode,
		CreatedAt:       appointment.CreatedAt,
		UpdatedAt:       appointment.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetUsersRegisteredForAppointment(w http.ResponseWriter, r *http.Request) {
	app_code := chi.URLParam(r, "id")

	bookings, err := services.GetAllBookingsForAppointment(app_code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve bookings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func GetAppointmentsCreatedByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	appointments, err := services.GetAllAppointmentsCreatedByUser(userIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve appointments: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}
