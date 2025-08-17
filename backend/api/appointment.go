package api

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/errors"
	"github.com/m13ha/appointment_master/middleware"
	"github.com/m13ha/appointment_master/models/requests"
	"github.com/m13ha/appointment_master/services"
	"github.com/m13ha/appointment_master/utils"
)

func CreateAppointment(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errors.BadRequest(c.Writer, "Invalid user ID")
		return
	}

	// Read and log raw body for debugging
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errors.BadRequest(c.Writer, "Failed to read request body")
		return
	}
	log.Printf("Raw request body: %s", body)

	// Reset body for decoding
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	var req requests.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(c.Writer, "Invalid request payload: "+err.Error())
		return
	}

	if err := utils.Validate(req); err != nil {
		errors.BadRequest(c.Writer, "Validation failed: "+err.Error())
		return
	}

	appointmentResponse, err := services.CreateAppointment(req, userID)
	if err != nil {
		switch {
		case err.Error() == "end time cannot be before start time" ||
			err.Error() == "end date cannot be before start date" ||
			err.Error() == "booking duration exceeds available time window":
			errors.BadRequest(c.Writer, err.Error())
		default:
			errors.InternalServerError(c.Writer, "Failed to create appointment: "+err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, appointmentResponse)
}

func GetAppointmentsCreatedByUser(c *gin.Context) {
	userIDStr := middleware.GetUserIDFromContext(c)
	if userIDStr == "" {
		errors.Unauthorized(c.Writer, "")
		return
	}

	appointments, err := services.GetAllAppointmentsCreatedByUser(userIDStr, nil)
	if err != nil {
		errors.InternalServerError(c.Writer, "Failed to retrieve appointments: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, appointments)
}
