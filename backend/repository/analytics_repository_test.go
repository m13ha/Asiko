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
