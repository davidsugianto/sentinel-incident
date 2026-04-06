package mocks

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	incidentRepo "github.com/davidsugianto/sentinel-incident/internal/repository/incident"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockIncidentRepository struct {
	mock.Mock
}

func (m *MockIncidentRepository) Create(ctx context.Context, incident *incidentModel.Incident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockIncidentRepository) GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incidentModel.Incident), args.Error(1)
}

func (m *MockIncidentRepository) List(ctx context.Context, params incidentRepo.ListParams) ([]incidentModel.Incident, int64, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]incidentModel.Incident), args.Get(1).(int64), args.Error(2)
}

func (m *MockIncidentRepository) Update(ctx context.Context, incident *incidentModel.Incident) error {
	args := m.Called(ctx, incident)
	return args.Error(0)
}

func (m *MockIncidentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
