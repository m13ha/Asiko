package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/m13ha/appointment_master/utils"
)

// @Summary Create a new user (initiate registration)
// @Description Register a new user in the system. This will trigger an email verification.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   user  body   requests.UserRequest  true  "User Registration Details"
// @Success 202 {object} responses.SimpleMessage
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /users [post]
// @ID createUser
func (h *Handler) CreateUser(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errors.BadRequest(c.Writer, "Failed to read request body")
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.FormatValidationErrors(c.Writer, err)
		return
	}

	_, err = h.userService.CreateUser(req)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusAccepted, responses.SimpleMessage{Message: "Registration pending. Please check your email for a verification code."})
}

// @Summary Verify user registration
// @Description Verify a user's email address with a code to complete registration.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   verification  body   requests.VerificationRequest  true  "Email and Verification Code"
// @Success 201 {object} responses.LoginResponse
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request payload or verification error"
// @Failure 500 {object} errors.ApiErrorResponse "Internal server error"
// @Router /auth/verify-registration [post]
// @ID verifyRegistration
func (h *Handler) VerifyRegistrationHandler(c *gin.Context) {
	var req requests.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	token, err := h.userService.VerifyRegistration(req.Email, req.Code)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, responses.LoginResponse{Token: token})
}
