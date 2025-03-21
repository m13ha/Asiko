package api

import (
	"encoding/json"
	"net/http"

	"github.com/m13ha/appointment_master/models"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := utils.Validate(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(formatValidationErrors(err))
		return
	}

	user, err := services.CreateUser(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewDatabaseErrorResponse("Failed to create user", err.Error()))
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
