package incident

import (
	"context"
	"errors"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("incident not found")

func (r *repository) Create(ctx context.Context, incident *incidentModel.Incident) error {
	return r.db.WithContext(ctx).Create(incident).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*incidentModel.Incident, error) {
	var incident incidentModel.Incident
	err := r.db.WithContext(ctx).First(&incident, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &incident, nil
}

func (r *repository) List(ctx context.Context, params ListParams) ([]incidentModel.Incident, int64, error) {
	var incidents []incidentModel.Incident
	var total int64

	query := r.db.WithContext(ctx).Model(&incidentModel.Incident{})

	if params.TeamID != "" {
		query = query.Where("team_id = ?", params.TeamID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Severity != "" {
		query = query.Where("severity = ?", params.Severity)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&incidents).Error; err != nil {
		return nil, 0, err
	}

	return incidents, total, nil
}

func (r *repository) Update(ctx context.Context, incident *incidentModel.Incident) error {
	return r.db.WithContext(ctx).Save(incident).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&incidentModel.Incident{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
