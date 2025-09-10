package utils

import (
	"fmt"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

// ParseTimeRange parses start and end date strings into a TimeRange
func ParseTimeRange(startDate, endDate string) (*TimeRange, error) {
	if startDate == "" || endDate == "" {
		return nil, fmt.Errorf("start_date and end_date are required")
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format, use YYYY-MM-DD")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end_date format, use YYYY-MM-DD")
	}

	if end.Before(start) {
		return nil, fmt.Errorf("end_date cannot be before start_date")
	}

	// Set end time to end of day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return &TimeRange{Start: start, End: end}, nil
}