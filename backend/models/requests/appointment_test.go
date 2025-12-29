package requests

import (
	"testing"
	"time"

	"github.com/m13ha/asiko/models/entities"
	"github.com/stretchr/testify/assert"
)

func TestAppointmentRequestValidateRejectsPastStartDate(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startTime := time.Date(2000, 1, 1, 9, 0, 0, 0, now.Location())
	endTime := time.Date(2000, 1, 1, 10, 0, 0, 0, now.Location())

	req := &AppointmentRequest{
		Title:           "Past Start Date",
		StartDate:       yesterday,
		EndDate:         yesterday,
		StartTime:       startTime,
		EndTime:         endTime,
		BookingDuration: 30,
		Type:            entities.Single,
		MaxAttendees:    1,
	}

	err := req.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Start date cannot be in the past")
}

func TestAppointmentRequestValidateRejectsPastStartTime(t *testing.T) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	pastStart := now.Add(-30 * time.Minute)
	startTime := time.Date(2000, 1, 1, pastStart.Hour(), pastStart.Minute(), 0, 0, now.Location())
	endTime := time.Date(2000, 1, 1, pastStart.Add(30*time.Minute).Hour(), pastStart.Add(30*time.Minute).Minute(), 0, 0, now.Location())

	req := &AppointmentRequest{
		Title:           "Past Start Time",
		StartDate:       startDate,
		EndDate:         startDate,
		StartTime:       startTime,
		EndTime:         endTime,
		BookingDuration: 30,
		Type:            entities.Single,
		MaxAttendees:    1,
	}

	err := req.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Start time cannot be in the past")
}
