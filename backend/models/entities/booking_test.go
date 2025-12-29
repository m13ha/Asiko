package entities

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBookingBeforeCreateSlotDefaults(t *testing.T) {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	date := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	start := time.Date(2025, 12, 31, 8, 30, 0, 0, time.UTC)

	booking := &Booking{
		ID:        id,
		IsSlot:    true,
		Status:    "",
		Available: true,
		Date:      date,
		StartTime: start,
	}

	err := booking.BeforeCreate(nil)

	assert.NoError(t, err)
	assert.Equal(t, BookingStatusActive, booking.Status)
	assert.True(t, booking.Available)
	assert.Equal(t, "BK25123108301111", booking.BookingCode)
}

func TestBookingBeforeCreateNonSlotDefaults(t *testing.T) {
	id := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	date := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	start := time.Date(2026, 1, 2, 14, 15, 0, 0, time.UTC)

	booking := &Booking{
		ID:        id,
		IsSlot:    false,
		Status:    "",
		Available: true,
		Date:      date,
		StartTime: start,
	}

	err := booking.BeforeCreate(nil)

	assert.NoError(t, err)
	assert.False(t, booking.Available)
	assert.Equal(t, BookingStatusPending, booking.Status)
	assert.Equal(t, "BK26010214152222", booking.BookingCode)
}
