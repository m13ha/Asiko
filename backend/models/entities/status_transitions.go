package entities

import "strings"

const (
	BookingStatusActive    = "active"
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusOngoing   = "ongoing"
	BookingStatusCancelled = "cancelled"
	BookingStatusCanceled  = "canceled"
	BookingStatusRejected  = "rejected"
	BookingStatusExpired   = "expired"
)

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusOngoing   AppointmentStatus = "ongoing"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCanceled  AppointmentStatus = "canceled"
)

var bookingStatusTransitions = map[string]map[string]struct{}{
	BookingStatusActive: {
		BookingStatusCancelled: {},
		BookingStatusRejected:  {},
	},
	BookingStatusPending: {
		BookingStatusConfirmed: {},
		BookingStatusCancelled: {},
		BookingStatusRejected:  {},
	},
	BookingStatusConfirmed: {
		BookingStatusCancelled: {},
		BookingStatusRejected:  {},
	},
	BookingStatusOngoing: {
		BookingStatusCancelled: {},
		BookingStatusRejected:  {},
	},
}

var appointmentStatusTransitions = map[AppointmentStatus]map[AppointmentStatus]struct{}{
	AppointmentStatusPending: {
		AppointmentStatusCanceled: {},
	},
	AppointmentStatusOngoing: {
		AppointmentStatusCanceled: {},
	},
}

func CanTransitionBookingStatus(from, to string) bool {
	from = strings.ToLower(strings.TrimSpace(from))
	to = strings.ToLower(strings.TrimSpace(to))
	if from == "" || to == "" {
		return false
	}
	if _, terminal := map[string]struct{}{
		BookingStatusCancelled: {},
		BookingStatusCanceled:  {},
		BookingStatusRejected:  {},
		BookingStatusExpired:   {},
	}[from]; terminal {
		return false
	}
	next, ok := bookingStatusTransitions[from]
	if !ok {
		return false
	}
	_, ok = next[to]
	return ok
}

func CanTransitionAppointmentStatus(from, to AppointmentStatus) bool {
	next, ok := appointmentStatusTransitions[from]
	if !ok {
		return false
	}
	_, ok = next[to]
	return ok
}
