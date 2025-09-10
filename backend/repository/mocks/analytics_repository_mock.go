package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type AnalyticsRepository struct {
	mock.Mock
}

func (m *AnalyticsRepository) GetUserAppointmentCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AnalyticsRepository) GetUserBookingCount(userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(int64), args.Error(1)
}