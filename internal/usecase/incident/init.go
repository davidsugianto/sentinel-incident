package incident

import (
	"context"

	alertRepo "github.com/davidsugianto/sentinel-incident/internal/repository/alert"
)

type Usecase interface {
	CreateIncident(ctx context.Context, teamID string, content *map[string]interface{}, params ...*map[string]string) error
}

type usecase struct {
	alertRepo alertRepo.Repository
}

type Dependencies struct {
	AlertRepo alertRepo.Repository
}

func New(deps Dependencies) Usecase {
	return &usecase{
		alertRepo: deps.AlertRepo,
	}
}
