package utils

import (
	"fmt"
	"time"
)

// IsTimeInFuture checks if a given time is in the future
func IsTimeInFuture(t time.Time) bool {
	return t.After(time.Now())
}

// ValidateAppointmentTimes checks if appointment start/end times and dates are valid
func ValidateAppointmentTimes(startDateStr, endDateStr, startTimeStr, endTimeStr string) error {
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start date format. Use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return fmt.Errorf("invalid end date format. Use YYYY-MM-DD")
	}

	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		return fmt.Errorf("invalid start time format. Use HH:MM")
	}

	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil {
		return fmt.Errorf("invalid end time format. Use HH:MM")
	}

	// Check if start datetime is in the future
	startDateTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, startDate.Location())
	if !startDateTime.After(time.Now()) {
		return fmt.Errorf("appointment start time must be in the future")
	}

	// Check if end date is not before start date
	if endDate.Before(startDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	// Check if end time is not before start time
	if endTime.Before(startTime) {
		return fmt.Errorf("end time cannot be before start time")
	}

	return nil
}

// ParseAppointmentTimes converts string dates/times to time.Time objects
func ParseAppointmentTimes(startDateStr, endDateStr, startTimeStr, endTimeStr string) (startDate, endDate, startTime, endTime time.Time, err error) {
	startDate, err = time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return
	}

	endDate, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return
	}

	startTime, err = time.Parse("15:04", startTimeStr)
	if err != nil {
		return
	}

	endTime, err = time.Parse("15:04", endTimeStr)
	if err != nil {
		return
	}

	return
}
