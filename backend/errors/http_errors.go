package errors

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// WriteError is a helper for consistent error responses
func WriteError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "Unauthorized"
	}
	WriteError(w, http.StatusUnauthorized, msg)
}

// BadRequest sends a 400 Bad Request error
func BadRequest(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusBadRequest, msg)
}

// NotFound sends a 404 Not Found error
func NotFound(w http.ResponseWriter, msg string) {
	WriteError(w, http.StatusNotFound, msg)
}

// InternalServerError sends a 500 Internal Server Error
func InternalServerError(w http.ResponseWriter, msg string) {
	if msg == "" {
		msg = "internal server error"
	}
	WriteError(w, http.StatusInternalServerError, msg)
}

// HandleServiceError checks for UserError and writes a user-friendly message, otherwise writes a generic error
func HandleServiceError(w http.ResponseWriter, err error, status int) {
	if err == nil {
		return
	}
	if ue, ok := err.(*UserError); ok {
		WriteError(w, status, ue.Message)
	} else {
		InternalServerError(w, "")
	}
}

// FormatValidationErrors writes validation errors as a structured HTTP response
func FormatValidationErrors(w http.ResponseWriter, err error) {
	var errors []ValidationError
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errors = append(errors, ValidationError{
				Field:   e.Field(),
				Message: e.Error(),
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(NewValidationErrorResponse(errors...))
}

// ValidationError represents an error response for validation issues
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse is the structure for validation error responses
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

// DatabaseErrorResponse represents an error response for database issues
type DatabaseErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// ApiErrorResponse represents a general API error response
type ApiErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// NewValidationErrorResponse creates a new ValidationErrorResponse
func NewValidationErrorResponse(errors ...ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{Errors: errors}
}

// NewDatabaseErrorResponse creates a new DatabaseErrorResponse
func NewDatabaseErrorResponse(message, code string) DatabaseErrorResponse {
	return DatabaseErrorResponse{Message: message, Code: code}
}

// NewApiErrorResponse creates a new ApiErrorResponse
func NewApiErrorResponse(message, code string, details string) ApiErrorResponse {
	return ApiErrorResponse{
		Message: message,
		Code:    code,
		Details: details,
	}
}
