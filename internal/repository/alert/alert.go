package alert

import (
	"context"
	"fmt"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/logger"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/slack"
)

func (r *repository) SendAlert(ctx context.Context, data *incidentModel.Incident) error {
	if !r.slackConfig.Enabled || r.slackConfig.WebhookURL == "" {
		logger.Info(ctx, "Slack alerts disabled or webhook not configured", nil)
		return nil
	}

	if r.slackClient == nil {
		logger.Info(ctx, "Slack client not initialized", nil)
		return nil
	}

	msg := r.buildSlackMessage(data)

	if err := r.slackClient.SendMessage(ctx, msg); err != nil {
		logger.Error(ctx, "Failed to send Slack alert", map[string]interface{}{
			"error":       err.Error(),
			"incident_id": data.ID,
		})
		return err
	}

	logger.Info(ctx, "Slack alert sent successfully", map[string]interface{}{
		"incident_id": data.ID,
	})
	return nil
}

func (r *repository) buildSlackMessage(incident *incidentModel.Incident) *slack.Message {
	color := r.getSeverityColor(incident.Severity)
	statusEmoji := r.getStatusEmoji(incident.Status)

	return &slack.Message{
		Attachments: []slack.Attachment{
			{
				Color: color,
				Title: fmt.Sprintf("Incident: %s", incident.Title),
				Fields: []slack.Field{
					{Title: "Status", Value: fmt.Sprintf("%s %s", statusEmoji, incident.Status), Short: true},
					{Title: "Severity", Value: string(incident.Severity), Short: true},
					{Title: "Team", Value: incident.TeamID, Short: true},
					{Title: "ID", Value: incident.ID.String(), Short: true},
				},
			},
		},
	}
}

func (r *repository) getSeverityColor(severity incidentModel.Severity) string {
	switch severity {
	case incidentModel.SeverityCritical:
		return "#FF0000" // Red
	case incidentModel.SeverityHigh:
		return "#FFA500" // Orange
	case incidentModel.SeverityMedium:
		return "#FFFF00" // Yellow
	case incidentModel.SeverityLow:
		return "#00FF00" // Green
	default:
		return "#808080" // Gray
	}
}

func (r *repository) getStatusEmoji(status incidentModel.Status) string {
	switch status {
	case incidentModel.StatusOpen:
		return ":red_circle:"
	case incidentModel.StatusAcknowledged:
		return ":large_yellow_circle:"
	case incidentModel.StatusResolved:
		return ":large_green_circle:"
	default:
		return ":white_circle:"
	}
}
