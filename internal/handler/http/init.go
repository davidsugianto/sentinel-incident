package http

import (
	incident "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
)

type Handler struct {
	incidentUseCase incident.Usecase
}

type Dependencies struct {
	IncidentUseCase incident.Usecase
}

func New(deps Dependencies) *Handler {
	return &Handler{
		incidentUseCase: deps.IncidentUseCase,
	}
}
