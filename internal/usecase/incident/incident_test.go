package incident

import (
	"context"
	"testing"

	"github.com/davidsugianto/sentinel-incident/internal/mocks"
	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type IncidentUsecaseTestSuite struct {
	suite.Suite
	mockIncidentRepo *mocks.MockIncidentRepository
	mockAlertRepo    *mocks.MockAlertRepository
	usecase          Usecase
	ctx              context.Context
}

func (s *IncidentUsecaseTestSuite) SetupTest() {
	s.mockIncidentRepo = new(mocks.MockIncidentRepository)
	s.mockAlertRepo = new(mocks.MockAlertRepository)
	s.usecase = New(Dependencies{
		IncidentRepo: s.mockIncidentRepo,
		AlertRepo:    s.mockAlertRepo,
	})
	s.ctx = context.Background()
}

func (s *IncidentUsecaseTestSuite) TestCreate() {
	incident := &incidentModel.Incident{
		TeamID:      "team-1",
		Title:       "Test Incident",
		Description: "Test Description",
		Severity:    incidentModel.SeverityHigh,
	}

	s.mockIncidentRepo.On("Create", mock.Anything, mock.AnythingOfType("*incident.Incident")).Return(nil)
	s.mockAlertRepo.On("SendAlert", mock.Anything, mock.AnythingOfType("*incident.Incident")).Return(nil)

	result, err := s.usecase.Create(s.ctx, incident)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	s.mockIncidentRepo.AssertExpectations(s.T())
}

func (s *IncidentUsecaseTestSuite) TestGetByID() {
	id := uuid.New()
	incident := &incidentModel.Incident{
		ID:          id,
		TeamID:      "team-1",
		Title:       "Test Incident",
		Description: "Test Description",
	}

	s.mockIncidentRepo.On("GetByID", mock.Anything, id).Return(incident, nil)

	result, err := s.usecase.GetByID(s.ctx, id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), id, result.ID)
	s.mockIncidentRepo.AssertExpectations(s.T())
}

func (s *IncidentUsecaseTestSuite) TestList() {
	incidents := []incidentModel.Incident{
		{ID: uuid.New(), Title: "Incident 1"},
		{ID: uuid.New(), Title: "Incident 2"},
	}
	params := ListParams{Page: 1, PageSize: 10}

	s.mockIncidentRepo.On("List", mock.Anything, mock.AnythingOfType("incident.ListParams")).Return(incidents, int64(2), nil)

	result, err := s.usecase.List(s.ctx, params)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Incidents, 2)
	assert.Equal(s.T(), int64(2), result.Total)
	s.mockIncidentRepo.AssertExpectations(s.T())
}

func (s *IncidentUsecaseTestSuite) TestUpdate() {
	id := uuid.New()
	existingIncident := &incidentModel.Incident{ID: id, TeamID: "team-1"}
	updates := &UpdateIncidentRequest{
		Title:       strPtr("Updated Title"),
		Description: strPtr("Updated Description"),
	}

	s.mockIncidentRepo.On("GetByID", mock.Anything, id).Return(existingIncident, nil)
	s.mockIncidentRepo.On("Update", mock.Anything, mock.AnythingOfType("*incident.Incident")).Return(nil)

	result, err := s.usecase.Update(s.ctx, id, updates)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	s.mockIncidentRepo.AssertExpectations(s.T())
}

func (s *IncidentUsecaseTestSuite) TestDelete() {
	id := uuid.New()

	s.mockIncidentRepo.On("Delete", mock.Anything, id).Return(nil)

	err := s.usecase.Delete(s.ctx, id)

	assert.NoError(s.T(), err)
	s.mockIncidentRepo.AssertExpectations(s.T())
}

func TestIncidentUsecaseSuite(t *testing.T) {
	suite.Run(t, new(IncidentUsecaseTestSuite))
}

func strPtr(s string) *string {
	return &s
}
