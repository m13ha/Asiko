package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/utils"
)

// @Summary Create a new user (initiate registration)
// @Description Register a new user in the system. This will trigger an email verification.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   user  body   requests.UserRequest  true  "User Registration Details"
// @Success 202 {object} responses.SimpleMessage
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /users [post]
// @ID createUser
func (h *Handler) CreateUser(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Failed to read request body"))
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	if err := utils.Validate(req); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusAccepted, responses.SimpleMessage{
		Message: "Registration pending. Please check your email for a verification code.",
		Data:    user,
	})
}

// @Summary Verify user registration
// @Description Verify a user's email address with a code to complete registration.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   verification  body   requests.VerificationRequest  true  "Email and Verification Code"
// @Success 201 {object} responses.LoginResponse
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload or verification error"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /auth/verify-registration [post]
// @ID verifyRegistration
func (h *Handler) VerifyRegistrationHandler(c *gin.Context) {
	var req requests.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	token, err := h.userService.VerifyRegistration(req.Email, req.Code)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusCreated, responses.LoginResponse{Token: token})
}

// @Summary Resend verification code
// @Description Resend a verification code for a pending user registration.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   resend  body   requests.ResendVerificationRequest  true  "Email to resend verification code to"
// @Success 202 {object} responses.SimpleMessage
// @Failure 400 {object} errors.APIErrorResponse "Invalid request payload"
// @Failure 404 {object} errors.APIErrorResponse "Pending registration not found"
// @Failure 409 {object} errors.APIErrorResponse "Account already verified"
// @Failure 500 {object} errors.APIErrorResponse "Internal server error"
// @Router /auth/resend-verification [post]
// @ID resendVerification
func (h *Handler) ResendVerificationHandler(c *gin.Context) {
	var req requests.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request payload"))
		return
	}

	if err := utils.Validate(req); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	if err := h.userService.ResendVerificationCode(req.Email); err != nil {
		c.Error(errors.FromError(err))
		return
	}

	c.JSON(http.StatusAccepted, responses.SimpleMessage{
		Message: "Verification code resent if a pending registration exists for this email.",
	})
}
