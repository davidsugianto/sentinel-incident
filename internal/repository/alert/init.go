package alert

import (
	"context"

	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/slack"
	"gorm.io/gorm"
)

type Repository interface {
	SendAlert(ctx context.Context, data *incidentModel.Incident) error
}

type repository struct {
	db          *gorm.DB
	slackClient *slack.Client
	slackConfig *config.Slack
}

type Dependencies struct {
	Database    *gorm.DB
	SlackClient *slack.Client
	SlackConfig *config.Slack
}

func New(deps Dependencies) Repository {
	return &repository{
		db:          deps.Database,
		slackClient: deps.SlackClient,
		slackConfig: deps.SlackConfig,
	}
}
