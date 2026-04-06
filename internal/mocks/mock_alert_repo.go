package mocks

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/stretchr/testify/mock"
)

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) SendAlert(ctx context.Context, data *incidentModel.Incident) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
