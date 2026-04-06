package incident

import (
	"context"
	"time"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	incidentRepo "github.com/davidsugianto/sentinel-incident/internal/repository/incident"
	"github.com/google/uuid"
)

var (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (u *usecase) Create(ctx context.Context, incident *incidentModel.Incident) (*incidentModel.Incident, error) {
	// Set defaults
	if incident.Status == "" {
		incident.Status = incidentModel.StatusOpen
	}

	if err := u.incidentRepo.Create(ctx, incident); err != nil {
		return nil, err
	}

	// Send alert asynchronously
	go func() {
		_ = u.alertRepo.SendAlert(context.Background(), incident)
	}()

	return incident, nil
}

func (u *usecase) GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error) {
	return u.incidentRepo.GetByID(ctx, id)
}

func (u *usecase) List(ctx context.Context, params ListParams) (*ListResult, error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = DefaultPageSize
	}
	if params.PageSize > MaxPageSize {
		params.PageSize = MaxPageSize
	}

	// Convert to repository params
	repoParams := incidentRepo.ListParams{
		TeamID:   params.TeamID,
		Status:   params.Status,
		Severity: params.Severity,
		Page:     params.Page,
		PageSize: params.PageSize,
	}

	incidents, total, err := u.incidentRepo.List(ctx, repoParams)
	if err != nil {
		return nil, err
	}

	return &ListResult{
		Incidents: incidents,
		Total:     total,
		Page:      params.Page,
		PageSize:  params.PageSize,
	}, nil
}

func (u *usecase) Update(ctx context.Context, id uuid.UUID, updates *UpdateIncidentRequest) (*incidentModel.Incident, error) {
	incident, err := u.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if updates.Title != nil {
		incident.Title = *updates.Title
	}
	if updates.Description != nil {
		incident.Description = *updates.Description
	}
	if updates.Content != nil {
		incident.Content = *updates.Content
	}
	if updates.Status != nil {
		incident.Status = *updates.Status
		// Handle status transitions
		if *updates.Status == incidentModel.StatusResolved && incident.ResolvedAt == nil {
			now := time.Now()
			incident.ResolvedAt = &now
			incident.Resolved = true
		}
	}
	if updates.Severity != nil {
		incident.Severity = *updates.Severity
	}

	if err := u.incidentRepo.Update(ctx, incident); err != nil {
		return nil, err
	}

	return incident, nil
}

func (u *usecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.incidentRepo.Delete(ctx, id)
}
