package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/m13ha/appointment_master/models/dto"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
	"github.com/rs/zerolog/log"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Read and log raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Log the raw request for debugging
	log.Debug().RawJSON("raw_request", body).Msg("User registration request received")

	// Reset body for decoding
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var req dto.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode user registration request")
		writeError(w, http.StatusBadRequest, "Invalid request payload")
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

	userResponse, err := services.CreateUser(req)
	if err != nil {
		HandleServiceError(w, err, http.StatusBadRequest, writeError)
		return
	}

	log.Info().
		Str("user_id", userResponse.ID.String()).
		Str("email", userResponse.Email).
		Msg("User registered successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}
