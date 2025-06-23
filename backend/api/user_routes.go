package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
	"github.com/rs/zerolog/log"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read and log raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the raw request for debugging
	log.Debug().RawJSON("raw_request", body).Msg("User registration request received")

	// Reset body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode user registration request")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Log the parsed request (without password)
	log.Info().
		Str("name", req.Name).
		Str("email", req.Email).
		Str("phone", req.PhoneNumber).
		Msg("Processing user registration")

	if err := utils.Validate(req); err != nil {
		validationErrors := formatValidationErrors(err)
		log.Error().
			Interface("validation_errors", validationErrors).
			Msg("Validation failed for user registration")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrors)
		return
	}

	// Check for field name discrepancy between frontend and backend
	if req.PhoneNumber == "" && req.Phone != "" {
		log.Info().Msg("Converting 'phone' field to 'phoneNumber' for compatibility")
		req.PhoneNumber = req.Phone
	}

	user, err := services.CreateUser(req)
	if err != nil {
		// Check for specific database errors
		if strings.Contains(err.Error(), "duplicate key") {
			if strings.Contains(err.Error(), "email") {
				log.Error().Err(err).Str("email", req.Email).Msg("Email already exists")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(models.NewDatabaseErrorResponse("Email already registered", "duplicate_email"))
				return
			} else if strings.Contains(err.Error(), "phone") {
				log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Phone number already exists")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(models.NewDatabaseErrorResponse("Phone number already registered", "duplicate_phone"))
				return
			}
		}

		log.Error().Err(err).Msg("Failed to create user")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewDatabaseErrorResponse("Failed to create user", "database_error"))
		return
	}

	response := models.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("email", user.Email).
		Msg("User registered successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
