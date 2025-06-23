package api

import (
	"bytes"
	"encoding/json"
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
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Read and log raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	log.Printf("Raw request body: %s", body)

	// Reset body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.AppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if err := utils.Validate(req); err != nil {
		writeError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	appointment, err := services.CreateAppointment(req, userID)
	if err != nil {
		switch {
		case err.Error() == "end time cannot be before start time" ||
			err.Error() == "end date cannot be before start date" ||
			err.Error() == "booking duration exceeds available time window":
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "Failed to create appointment: "+err.Error())
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
		Description:     appointment.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetUsersRegisteredForAppointment(w http.ResponseWriter, r *http.Request) {
	app_code := chi.URLParam(r, "id")
	if app_code == "" {
		writeError(w, http.StatusBadRequest, "Missing appointment code parameter")
		return
	}

	bookings, err := services.GetAllBookingsForAppointment(app_code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve bookings: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func GetAppointmentsCreatedByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := GetUserIDFromContext(r)
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	appointments, err := services.GetAllAppointmentsCreatedByUser(userIDStr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve appointments: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}
