package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeRange(t *testing.T) {
	testCases := []struct {
		name        string
		startDate   string
		endDate     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid date range",
			startDate:   "2025-01-01",
			endDate:     "2025-01-31",
			expectError: false,
		},
		{
			name:        "Missing start date",
			startDate:   "",
			endDate:     "2025-01-31",
			expectError: true,
			errorMsg:    "start_date and end_date are required",
		},
		{
			name:        "Missing end date",
			startDate:   "2025-01-01",
			endDate:     "",
			expectError: true,
			errorMsg:    "start_date and end_date are required",
		},
		{
			name:        "Invalid start date format",
			startDate:   "2025/01/01",
			endDate:     "2025-01-31",
			expectError: true,
			errorMsg:    "invalid start_date format, use YYYY-MM-DD",
		},
		{
			name:        "Invalid end date format",
			startDate:   "2025-01-01",
			endDate:     "2025/01/31",
			expectError: true,
			errorMsg:    "invalid end_date format, use YYYY-MM-DD",
		},
		{
			name:        "End date before start date",
			startDate:   "2025-01-31",
			endDate:     "2025-01-01",
			expectError: true,
			errorMsg:    "end_date cannot be before start_date",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseTimeRange(tc.startDate, tc.endDate)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				
				expectedStart, _ := time.Parse("2006-01-02", tc.startDate)
				expectedEnd, _ := time.Parse("2006-01-02", tc.endDate)
				expectedEnd = expectedEnd.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				
				assert.Equal(t, expectedStart, result.Start)
				assert.Equal(t, expectedEnd, result.End)
			}
		})
	}
}