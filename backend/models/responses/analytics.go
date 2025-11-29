package responses

import "time"

// AnalyticsResponse now includes richer, backward-compatible analytics.
// Existing fields are preserved and new sections add detail.
type AnalyticsResponse struct {
	// Summary
	TotalAppointments int       `json:"total_appointments"`
	TotalBookings     int       `json:"total_bookings"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`

	// Time series
	BookingsPerDay []TimeSeriesPoint `json:"bookings_per_day,omitempty"`
}

type TimeSeriesPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type BucketCount struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type TopAppointment struct {
	AppCode              string  `json:"app_code"`
	Title                string  `json:"title"`
	Bookings             int     `json:"bookings"`
	CapacityUsagePercent float64 `json:"capacity_usage_percent,omitempty"`
}
