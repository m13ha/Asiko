package repository

import (
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

	mock.ExpectQuery(`(?s)SELECT count\(\*\) FROM "appointments".*status (<>|!=)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

	count, err := repo.GetUserAppointmentCount(userID, start, end)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookingsPerDay(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(`(?s)SELECT.*TO_CHAR.*FROM bookings.*JOIN appointments.*LOWER\(bookings\.status\) IN.*GROUP BY.*ORDER BY`).
		WillReturnRows(sqlmock.NewRows([]string{"date", "count"}).AddRow("2025-01-01", 2).AddRow("2025-01-02", 1))

	rows, err := repo.GetBookingsPerDay(userID, start, end)
	assert.NoError(t, err)
	if assert.Len(t, rows, 2) {
		assert.Equal(t, "2025-01-01", rows[0].Date)
		assert.Equal(t, 2, rows[0].Count)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserCancellationCount(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(`(?s)SELECT.*FROM bookings.*JOIN appointments.*bookings\.status = 'cancelled'`).
		WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(int64(4)))

	count, err := repo.GetUserCancellationCount(userID, start, end)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCancellationsPerDay(t *testing.T) {
	gdb, mock := setupMockDB(t)
	repo := NewGormAnalyticsRepository(gdb)

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	mock.ExpectQuery(`(?s)SELECT.*TO_CHAR\(bookings\.updated_at::date, 'YYYY-MM-DD'\).*FROM bookings.*JOIN appointments.*bookings\.status = 'cancelled'.*GROUP BY.*ORDER BY`).
		WillReturnRows(sqlmock.NewRows([]string{"date", "count"}).AddRow("2025-01-01", 1).AddRow("2025-01-02", 2))

	rows, err := repo.GetCancellationsPerDay(userID, start, end)
	assert.NoError(t, err)
	assert.Len(t, rows, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
