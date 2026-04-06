package incident

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	alertRepo "github.com/davidsugianto/sentinel-incident/internal/repository/alert"
	incidentRepo "github.com/davidsugianto/sentinel-incident/internal/repository/incident"
	"github.com/google/uuid"
)

type Usecase interface {
	Create(ctx context.Context, incident *incidentModel.Incident) (*incidentModel.Incident, error)
	GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error)
	List(ctx context.Context, params ListParams) (*ListResult, error)
	Update(ctx context.Context, id uuid.UUID, updates *UpdateIncidentRequest) (*incidentModel.Incident, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type usecase struct {
	incidentRepo incidentRepo.Repository
	alertRepo    alertRepo.Repository
}

type Dependencies struct {
	IncidentRepo incidentRepo.Repository
	AlertRepo    alertRepo.Repository
}

func New(deps Dependencies) Usecase {
	return &usecase{
		incidentRepo: deps.IncidentRepo,
		alertRepo:    deps.AlertRepo,
	}
}

type ListParams struct {
	TeamID   string
	Status   incidentModel.Status
	Severity incidentModel.Severity
	Page     int
	PageSize int
}

type ListResult struct {
	Incidents []incidentModel.Incident `json:"incidents"`
	Total     int64                    `json:"total"`
	Page      int                      `json:"page"`
	PageSize  int                      `json:"page_size"`
}

type UpdateIncidentRequest struct {
	Title       *string                 `json:"title,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Content     *map[string]interface{} `json:"content,omitempty"`
	Status      *incidentModel.Status   `json:"status,omitempty"`
	Severity    *incidentModel.Severity `json:"severity,omitempty"`
}
