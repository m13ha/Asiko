package responses

import "time"

// AnalyticsResponse now includes richer, backward-compatible analytics.
// Existing fields are preserved and new sections add detail.
type AnalyticsResponse struct {
	// Summary
	TotalAppointments int       `json:"total_appointments"`
	TotalBookings     int       `json:"total_bookings"`
	TotalCancellations int      `json:"total_cancellations,omitempty"`
	CancellationRate   float64  `json:"cancellation_rate,omitempty"`    // percent 0-100
	AvgBookingsPerDay  float64  `json:"avg_bookings_per_day,omitempty"` // derived from total bookings / days in range
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`

	// Time series
	BookingsPerDay []TimeSeriesPoint `json:"bookings_per_day,omitempty"`
	CancellationsPerDay []TimeSeriesPoint `json:"cancellations_per_day,omitempty"`
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
