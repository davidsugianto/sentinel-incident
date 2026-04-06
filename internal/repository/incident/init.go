package incident

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, incident *incidentModel.Incident) error
	GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error)
	List(ctx context.Context, params ListParams) ([]incidentModel.Incident, int64, error)
	Update(ctx context.Context, incident *incidentModel.Incident) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

type Dependencies struct {
	Database *gorm.DB
}

func New(deps Dependencies) Repository {
	return &repository{db: deps.Database}
}

type ListParams struct {
	TeamID   string
	Status   incidentModel.Status
	Severity incidentModel.Severity
	Page     int
	PageSize int
}
