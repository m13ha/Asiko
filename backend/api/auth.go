package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/asiko/errors"
	"github.com/m13ha/asiko/middleware"
	"github.com/m13ha/asiko/models/requests"
	"github.com/m13ha/asiko/models/responses"
	"github.com/m13ha/asiko/utils"
)

// @Summary User Login
// @Description Authenticate a user and receive a JWT token.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   login  body   requests.LoginRequest  true  "Login Credentials"
// @Success 200 {object} responses.LoginResponse
// @Success 202 {object} errors.APIErrorResponse "Registration pending verification"
// @Failure 400 {object} errors.APIErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} errors.APIErrorResponse "Invalid email or password"
// @Failure 500 {object} errors.APIErrorResponse "Could not generate token"
// @Router /login [post]
// @ID loginUser
func (h *Handler) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request body"))
		return
	}

	if err := utils.Validate(req); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	userEntity, err := h.userService.AuthenticateUser(utils.NormalizeEmail(req.Email), req.Password)
	if err != nil {
		c.Error(errors.FromError(err))
		return
	}

	accessToken, err := middleware.GenerateToken(userEntity.ID.String())
	if err != nil {
		c.Error(errors.New(errors.CodeInternalError).WithKind(errors.KindInternal).WithHTTP(500).WithMessage("Could not generate token"))
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(userEntity.ID.String())
	if err != nil {
		c.Error(errors.New(errors.CodeInternalError).WithKind(errors.KindInternal).WithHTTP(500).WithMessage("Could not generate refresh token"))
		return
	}

	c.JSON(http.StatusOK, responses.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    middleware.AccessTokenTTLSeconds(),
		User: responses.UserResponse{
			ID:    userEntity.ID,
			Name:  userEntity.Name,
			Email: userEntity.Email,
		},
	})
}

// @Summary User Logout
// @Description Invalidate the user's session.
// @Tags Authentication
// @Produce  application/json
// @Success 200 {object} responses.SimpleMessage
// @Router /logout [post]
// @ID logoutUser
func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "Logged out successfully"})
}

// @Summary Refresh access token
// @Description Exchange a refresh token for a new access token
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   refresh body requests.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} responses.TokenResponse
// @Failure 400 {object} errors.APIErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} errors.APIErrorResponse "Invalid refresh token"
// @Failure 500 {object} errors.APIErrorResponse "Could not generate token"
// @Router /auth/refresh [post]
// @ID refreshToken
func (h *Handler) Refresh(c *gin.Context) {
	var req requests.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request body"))
		return
	}

	if err := utils.Validate(req); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	userID, err := middleware.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.Error(errors.New(errors.CodeLoginInvalidCredentials).WithKind(errors.KindUnauthorized).WithHTTP(401).WithMessage("Invalid refresh token"))
		return
	}

	accessToken, err := middleware.GenerateToken(userID)
	if err != nil {
		c.Error(errors.New(errors.CodeInternalError).WithKind(errors.KindInternal).WithHTTP(500).WithMessage("Could not generate token"))
		return
	}

	newRefreshToken, err := middleware.GenerateRefreshToken(userID)
	if err != nil {
		c.Error(errors.New(errors.CodeInternalError).WithKind(errors.KindInternal).WithHTTP(500).WithMessage("Could not generate refresh token"))
		return
	}

	c.JSON(http.StatusOK, responses.TokenResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    middleware.AccessTokenTTLSeconds(),
	})
}

// @Summary Generate Device Token
// @Description Generate a short-lived token for a given device ID to be used in booking requests.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   device   body   requests.DeviceTokenRequest  true  "Device ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errors.APIErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} errors.APIErrorResponse "Could not generate token"
// @Router /auth/device-token [post]
// @ID generateDeviceToken
func (h *Handler) GenerateDeviceTokenHandler(c *gin.Context) {
	var req requests.DeviceTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.New(errors.CodeValidationFailed).WithKind(errors.KindValidation).WithHTTP(400).WithMessage("Invalid request body"))
		return
	}

	if err := utils.Validate(req); err != nil {
		c.Error(errors.FromValidation(err, "Validation failed"))
		return
	}

	token, err := middleware.GenerateDeviceToken(req.DeviceID)
	if err != nil {
		c.Error(errors.New(errors.CodeInternalError).WithKind(errors.KindInternal).WithHTTP(500).WithMessage("Could not generate device token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"device_token": token})
}
