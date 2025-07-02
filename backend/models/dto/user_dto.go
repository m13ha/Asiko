package dto

import (
	"time"

	"github.com/google/uuid"
)

// UserRequest represents the request payload for creating or updating a user.
type UserRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=64"`
	PhoneNumber string `json:"phone_number"` // Primary field
	Phone       string `json:"phone"`        // Alternative field for frontend compatibility
}

// UserResponse represents the response payload for user-related requests.
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// LoginRequest represents the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
