package alert

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"gorm.io/gorm"
)

type Repository interface {
	SendAlert(ctx context.Context, data *incidentModel.Incident) error
}

type repository struct {
	db *gorm.DB
}

type Dependencies struct {
	Database *gorm.DB
}

func New(deps Dependencies) Repository {
	return &repository{
		db: deps.Database,
	}
}
