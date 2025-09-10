package mocks

import (
	"github.com/google/uuid"
	"github.com/m13ha/appointment_master/models/responses"
	"github.com/stretchr/testify/mock"
)

type AnalyticsService struct {
	mock.Mock
}

func (m *AnalyticsService) GetUserAnalytics(userID uuid.UUID, startDate, endDate string) (*responses.AnalyticsResponse, error) {
	args := m.Called(userID, startDate, endDate)
	return args.Get(0).(*responses.AnalyticsResponse), args.Error(1)
}