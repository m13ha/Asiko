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

    // Breakdowns
    AppointmentsByType map[string]int `json:"appointments_by_type,omitempty"`
    BookingsByStatus   map[string]int `json:"bookings_by_status,omitempty"`
    GuestVsRegistered  map[string]int `json:"guest_vs_registered,omitempty"`
    DistinctCustomers  int            `json:"distinct_customers,omitempty"`
    RepeatCustomers    int            `json:"repeat_customers,omitempty"`

    // Utilization & Capacity
    SlotUtilizationPercent float64 `json:"slot_utilization_percent,omitempty"`
    AvgAttendeesPerBooking float64 `json:"avg_attendees_per_booking,omitempty"`
    PartyCapacity          struct {
        Total   int `json:"total,omitempty"`
        Used    int `json:"used,omitempty"`
        Percent float64 `json:"percent,omitempty"`
    } `json:"party_capacity,omitempty"`

    // Timing
    AvgLeadTimeHours    float64 `json:"avg_lead_time_hours,omitempty"`
    MedianLeadTimeHours float64 `json:"median_lead_time_hours,omitempty"`

    // Time series
    BookingsPerDay     []TimeSeriesPoint `json:"bookings_per_day,omitempty"`
    CancellationsPerDay []TimeSeriesPoint `json:"cancellations_per_day,omitempty"`
    RejectionsPerDay   []TimeSeriesPoint `json:"rejections_per_day,omitempty"`

    // Insights
    PeakHours []BucketCount `json:"peak_hours,omitempty"`
    PeakDays  []BucketCount `json:"peak_days,omitempty"`
    TopAppointments []TopAppointment `json:"top_appointments,omitempty"`
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
    AppCode   string  `json:"app_code"`
    Title     string  `json:"title"`
    Bookings  int     `json:"bookings"`
    CapacityUsagePercent float64 `json:"capacity_usage_percent,omitempty"`
}
