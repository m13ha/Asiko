package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating sqlmock: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("error opening gorm: %v", err)
	}
	return gdb, mock
}

func TestGetUserAppointmentCount(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "appointments"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.GetUserAppointmentCount(userID, start, end)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsByStatusCounts(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	// GORM builds a grouped SELECT; match loosely on FROM bookings
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT bookings.status as status, COUNT(*) as cnt FROM "bookings" JOIN appointments ON bookings.appointment_id = appointments.id`)).
		WillReturnRows(sqlmock.NewRows([]string{"status", "cnt"}).AddRow("active", 3).AddRow("cancelled", 1))

	m, err := repo.GetBookingsByStatusCounts(userID, start, end)
	assert.NoError(t, err)
	assert.Equal(t, map[string]int{"active": 3, "cancelled": 1}, m)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsPerDay(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT TO_CHAR(bookings.created_at::date, 'YYYY-MM-DD') as date, COUNT(*) as count`)).
		WillReturnRows(sqlmock.NewRows([]string{"date", "count"}).AddRow("2025-01-01", 2).AddRow("2025-01-02", 1))

	rows, err := repo.GetBookingsPerDay(userID, start, end)
	assert.NoError(t, err)
	if assert.Len(t, rows, 2) {
		assert.Equal(t, "2025-01-01", rows[0].Date)
		assert.Equal(t, 2, rows[0].Count)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopAppointments(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	// Match loosely on WITH base ... FROM base
	mock.ExpectQuery(`WITH\s+base[\s\S]*FROM\s+base`).
		WillReturnRows(sqlmock.NewRows([]string{"app_code", "title", "bookings", "capacity_usage_percent"}).
			AddRow("AP1", "A", 5, 40.0))

	out, err := repo.GetTopAppointments(userID, start, end, 5)
	assert.NoError(t, err)
	if assert.Len(t, out, 1) {
		assert.Equal(t, "AP1", out[0].AppCode)
		assert.Equal(t, "A", out[0].Title)
		assert.Equal(t, 5, out[0].Bookings)
		assert.InDelta(t, 40.0, out[0].CapacityUsagePercent, 0.001)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}
