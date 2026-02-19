package alert

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
)

func (r *repository) SendAlert(ctx context.Context, data *incidentModel.Incident) error {
	return nil
}
