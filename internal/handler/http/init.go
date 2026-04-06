package http

import (
	incident "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
	oncall "github.com/davidsugianto/sentinel-incident/internal/usecase/oncall"
)

type Handler struct {
	incidentUseCase incident.Usecase
	oncallUseCase   oncall.Usecase
}

type Dependencies struct {
	IncidentUseCase incident.Usecase
	OncallUseCase   oncall.Usecase
}

func New(deps Dependencies) *Handler {
	return &Handler{
		incidentUseCase: deps.IncidentUseCase,
		oncallUseCase:   deps.OncallUseCase,
	}
}
