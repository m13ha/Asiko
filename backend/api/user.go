package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apierrors "github.com/m13ha/asiko/errors/apierrors"
	"github.com/m13ha/asiko/middleware"
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
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload or validation error"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /users [post]
// @ID createUser
func (h *Handler) CreateUser(c *gin.Context) {
	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		apierrors.HandleAppError(c, err)
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
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload or verification error"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /auth/verify-registration [post]
// @ID verifyRegistration
func (h *Handler) VerifyRegistrationHandler(c *gin.Context) {
	var req requests.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	token, err := h.userService.VerifyRegistration(req.Email, req.Code)
	if err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	userID, err := middleware.ParseUserIDFromToken(token)
	if err != nil {
		apierrors.InternalServerError(c, "Could not parse token")
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(userID)
	if err != nil {
		apierrors.InternalServerError(c, "Could not generate refresh token")
		return
	}

	c.JSON(http.StatusCreated, responses.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    middleware.AccessTokenTTLSeconds(),
	})
}

// @Summary Resend verification code
// @Description Resend a verification code for a pending user registration.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   resend  body   requests.ResendVerificationRequest  true  "Email to resend verification code to"
// @Success 202 {object} responses.SimpleMessage
// @Failure 400 {object} responses.APIErrorResponse "Invalid request payload"
// @Failure 404 {object} responses.APIErrorResponse "Pending registration not found"
// @Failure 409 {object} responses.APIErrorResponse "Account already verified"
// @Failure 500 {object} responses.APIErrorResponse "Internal server error"
// @Router /auth/resend-verification [post]
// @ID resendVerification
func (h *Handler) ResendVerificationHandler(c *gin.Context) {
	var req requests.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	err := h.userService.ResendVerificationCode(req.Email)
	if err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusAccepted, responses.SimpleMessage{
		Message: "Verification code resent if a pending registration exists for this email.",
	})
}
