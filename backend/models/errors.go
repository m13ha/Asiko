package models

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
