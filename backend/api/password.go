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

// @Summary Forgot Password
// @Description Request a password reset email.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   request  body   requests.ForgotPasswordRequest  true  "Email"
// @Success 200 {object} responses.SimpleMessage
// @Failure 400 {object} responses.APIErrorResponse "Invalid request body or validation error"
// @Failure 500 {object} responses.APIErrorResponse "Could not initiate password reset"
// @Router /auth/forgot-password [post]
// @ID forgotPassword
func (h *Handler) ForgotPasswordHandler(c *gin.Context) {
	var req requests.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	if err := h.userService.ForgotPassword(req.Email); err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "If an account exists with this email, a reset code has been sent."})
}

// @Summary Reset Password
// @Description Reset password using a valid token.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   request  body   requests.ResetPasswordRequest  true  "Token and New Password"
// @Success 200 {object} responses.SimpleMessage
// @Failure 400 {object} responses.APIErrorResponse "Invalid request body or validation error"
// @Failure 422 {object} responses.APIErrorResponse "Invalid or expired reset token"
// @Failure 500 {object} responses.APIErrorResponse "Could not reset password"
// @Router /auth/reset-password [post]
// @ID resetPassword
func (h *Handler) ResetPasswordHandler(c *gin.Context) {
	var req requests.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	if err := h.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "Password has been reset successfully."})
}

// @Summary Change Password
// @Description Change password for authenticated user.
// @Tags Authentication
// @Accept  application/json
// @Produce  application/json
// @Param   request  body   requests.ChangePasswordRequest  true  "Old and New Password"
// @Success 200 {object} responses.SimpleMessage
// @Failure 400 {object} responses.APIErrorResponse "Invalid request body or validation error"
// @Failure 401 {object} responses.APIErrorResponse "Unauthorized"
// @Failure 422 {object} responses.APIErrorResponse "Incorrect old password"
// @Failure 500 {object} responses.APIErrorResponse "Could not change password"
// @Router /auth/change-password [post]
// @ID changePassword
func (h *Handler) ChangePasswordHandler(c *gin.Context) {
	userID := middleware.GetUserIDFromContext(c)
	if userID == "" {
		apierrors.UnauthorizedError(c, "Unauthorized")
		return
	}

	var req requests.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierrors.BadRequestError(c, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		apierrors.ValidationError(c, "Validation failed")
		return
	}

	if err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		apierrors.HandleAppError(c, err)
		return
	}

	c.JSON(http.StatusOK, responses.SimpleMessage{Message: "Password changed successfully."})
}
