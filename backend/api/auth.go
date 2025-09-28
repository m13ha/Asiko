package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

// @Summary User Login
// @Description Authenticate a user and receive a JWT token.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   login  body   requests.LoginRequest  true  "Login Credentials"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} errors.ApiErrorResponse "Invalid email or password"
// @Failure 500 {object} errors.ApiErrorResponse "Could not generate token"
// @Router /login [post]
// @ID loginUser
func (h *Handler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.FormatValidationErrors(c.Writer, err)
		return
	}

	userEntity, err := h.userService.AuthenticateUser(utils.NormalizeEmail(req.Email), req.Password)
	if err != nil {
		errors.Unauthorized(c.Writer, "Invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(userEntity.ID.String())
	if err != nil {
		errors.InternalServerError(c.Writer, "Could not generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  services.ToUserResponse(userEntity),
	})
}

// @Summary User Logout
// @Description Invalidate the user's session.
// @Tags Authentication
// @Produce  json
// @Success 200 {object} responses.SimpleMessageResponse
// @Router /logout [post]
// @ID logoutUser
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// @Summary Generate Device Token
// @Description Generate a short-lived token for a given device ID to be used in booking requests.
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   device   body   requests.DeviceTokenRequest  true  "Device ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.ApiErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} errors.ApiErrorResponse "Could not generate token"
// @Router /auth/device-token [post]
// @ID generateDeviceToken
func (h *Handler) GenerateDeviceTokenHandler(c *gin.Context) {
	var req requests.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.FormatValidationErrors(c.Writer, err)
		return
	}

	token, err := middleware.GenerateDeviceToken(req.DeviceID)
	if err != nil {
		errors.InternalServerError(c.Writer, "Could not generate device token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"device_token": token})
}