package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
	"github.com/rs/zerolog/log"
)

func CreateUser(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read request body")
		errors.BadRequest(c.Writer, "Failed to read request body")
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req requests.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Failed to decode user registration request")
		errors.BadRequest(c.Writer, "Invalid request payload")
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.FormatValidationErrors(c.Writer, err)
		log.Error().
			Interface("validation_errors", err).
			Msg("Validation failed for user registration")
		return
	}

	userResponse, err := services.CreateUser(req)
	if err != nil {
		errors.HandleServiceError(c.Writer, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, userResponse)
}
