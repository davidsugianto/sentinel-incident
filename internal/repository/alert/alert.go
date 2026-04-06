package alert

import (
	"context"
	"fmt"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/lark"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/logger"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/slack"
)

func (r *repository) SendAlert(ctx context.Context, data *incidentModel.Incident) error {
	// Send to Slack
	if err := r.sendSlackAlert(ctx, data); err != nil {
		logger.Error(ctx, "Failed to send Slack alert", map[string]interface{}{
			"error":       err.Error(),
			"incident_id": data.ID,
		})
		// Continue to other channels even if Slack fails
	}

	// Send to Lark
	if err := r.sendLarkAlert(ctx, data); err != nil {
		logger.Error(ctx, "Failed to send Lark alert", map[string]interface{}{
			"error":       err.Error(),
			"incident_id": data.ID,
		})
		// Continue even if Lark fails
	}

	return nil
}

func (r *repository) sendSlackAlert(ctx context.Context, data *incidentModel.Incident) error {
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
		return err
	}

	logger.Info(ctx, "Slack alert sent successfully", map[string]interface{}{
		"incident_id": data.ID,
	})
	return nil
}

func (r *repository) sendLarkAlert(ctx context.Context, data *incidentModel.Incident) error {
	if !r.larkConfig.Enabled || r.larkConfig.WebhookURL == "" {
		logger.Info(ctx, "Lark alerts disabled or webhook not configured", nil)
		return nil
	}

	if r.larkClient == nil {
		logger.Info(ctx, "Lark client not initialized", nil)
		return nil
	}

	msg := r.buildLarkMessage(data)

	if err := r.larkClient.SendMessage(ctx, msg); err != nil {
		return err
	}

	logger.Info(ctx, "Lark alert sent successfully", map[string]interface{}{
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

func (r *repository) buildLarkMessage(incident *incidentModel.Incident) *lark.CardMessage {
	template := r.getLarkTemplate(incident.Severity)

	return &lark.CardMessage{
		MsgType: "interactive",
		Card: &lark.Card{
			Config: &lark.CardConfig{
				WideScreenMode: true,
				EnableForward:  true,
			},
			Header: &lark.CardHeader{
				Title: &lark.CardTitle{
					Tag:     "plain_text",
					Content: fmt.Sprintf("Incident: %s", incident.Title),
				},
				Template: template,
			},
			Elements: []lark.Element{
				lark.FieldElement{
					Tag: "div",
					Fields: []lark.Field{
						{
							IsShort: true,
							Text: lark.TextItem{
								Tag:     "lark_md",
								Content: fmt.Sprintf("**Status:** %s", incident.Status),
							},
						},
						{
							IsShort: true,
							Text: lark.TextItem{
								Tag:     "lark_md",
								Content: fmt.Sprintf("**Severity:** %s", incident.Severity),
							},
						},
						{
							IsShort: true,
							Text: lark.TextItem{
								Tag:     "lark_md",
								Content: fmt.Sprintf("**Team:** %s", incident.TeamID),
							},
						},
						{
							IsShort: true,
							Text: lark.TextItem{
								Tag:     "lark_md",
								Content: fmt.Sprintf("**ID:** %s", incident.ID.String()),
							},
						},
					},
				},
				lark.DivElement{
					Tag: "div",
					Text: &lark.TextItem{
						Tag:     "lark_md",
						Content: fmt.Sprintf("**Description:** %s", incident.Description),
					},
				},
			},
		},
	}
}

func (r *repository) getLarkTemplate(severity incidentModel.Severity) string {
	switch severity {
	case incidentModel.SeverityCritical:
		return "red"
	case incidentModel.SeverityHigh:
		return "orange"
	case incidentModel.SeverityMedium:
		return "yellow"
	case incidentModel.SeverityLow:
		return "green"
	default:
		return "blue"
	}
}
