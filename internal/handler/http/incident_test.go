package http

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	incidentUsecase "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockIncidentUsecase is a local mock to avoid import cycle
type MockIncidentUsecase struct {
	mock.Mock
}

func (m *MockIncidentUsecase) Create(ctx context.Context, incident *incidentModel.Incident) (*incidentModel.Incident, error) {
	args := m.Called(ctx, incident)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incidentModel.Incident), args.Error(1)
}

func (m *MockIncidentUsecase) GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incidentModel.Incident), args.Error(1)
}

func (m *MockIncidentUsecase) List(ctx context.Context, params incidentUsecase.ListParams) (*incidentUsecase.ListResult, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incidentUsecase.ListResult), args.Error(1)
}

func (m *MockIncidentUsecase) Update(ctx context.Context, id uuid.UUID, updates *incidentUsecase.UpdateIncidentRequest) (*incidentModel.Incident, error) {
	args := m.Called(ctx, id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*incidentModel.Incident), args.Error(1)
}

func (m *MockIncidentUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type IncidentHandlerTestSuite struct {
	suite.Suite
	router      *gin.Engine
	mockUsecase *MockIncidentUsecase
	handler     *Handler
}

func (s *IncidentHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockUsecase = new(MockIncidentUsecase)
	s.handler = New(Dependencies{
		IncidentUseCase: s.mockUsecase,
	})

	s.router = gin.New()
	v1 := s.router.Group("/v1")
	incidents := v1.Group("/incidents")
	incidents.GET("", s.handler.ListIncidents)
	incidents.GET("/:id", s.handler.GetIncident)
	incidents.POST("", s.handler.CreateIncident)
	incidents.PUT("/:id", s.handler.UpdateIncident)
	incidents.DELETE("/:id", s.handler.DeleteIncident)
}

func (s *IncidentHandlerTestSuite) TestCreateIncident() {
	incident := &incidentModel.Incident{
		ID:          uuid.New(),
		TeamID:      "team-1",
		Title:       "Test Incident",
		Description: "Test Description",
		Severity:    incidentModel.SeverityHigh,
	}

	s.mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*incident.Incident")).Return(incident, nil)

	body := `{"team_id": "team-1", "title": "Test Incident", "description": "Test Description", "severity": "high"}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/incidents", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
}

func (s *IncidentHandlerTestSuite) TestGetIncident() {
	id := uuid.New()
	incident := &incidentModel.Incident{
		ID:          id,
		TeamID:      "team-1",
		Title:       "Test Incident",
		Description: "Test Description",
	}

	s.mockUsecase.On("GetByID", mock.Anything, id).Return(incident, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/incidents/"+id.String(), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *IncidentHandlerTestSuite) TestListIncidents() {
	incidents := []incidentModel.Incident{
		{ID: uuid.New(), Title: "Incident 1"},
		{ID: uuid.New(), Title: "Incident 2"},
	}
	result := &incidentUsecase.ListResult{
		Incidents: incidents,
		Total:     2,
	}

	s.mockUsecase.On("List", mock.Anything, mock.AnythingOfType("incident.ListParams")).Return(result, nil)

	req, _ := http.NewRequest(http.MethodGet, "/v1/incidents", nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *IncidentHandlerTestSuite) TestDeleteIncident() {
	id := uuid.New()

	s.mockUsecase.On("Delete", mock.Anything, id).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/v1/incidents/"+id.String(), nil)
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func TestIncidentHandlerSuite(t *testing.T) {
	suite.Run(t, new(IncidentHandlerTestSuite))
}
