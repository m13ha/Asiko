package responses

import "time"

type AnalyticsResponse struct {
	TotalAppointments int       `json:"total_appointments"`
	TotalBookings     int       `json:"total_bookings"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
}