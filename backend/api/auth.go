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

func Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request body")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.FormatValidationErrors(c.Writer, err)
		return
	}

	user, err := services.AuthenticateUser(utils.NormalizeEmail(req.Email), req.Password)
	if err != nil {
		errors.Unauthorized(c.Writer, "Invalid email or password")
		return
	}

	token, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		errors.InternalServerError(c.Writer, "Could not generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  services.ToUserResponse(user),
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
