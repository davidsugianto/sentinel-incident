package main

import (
	"fmt"
	"log"

	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/db"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/lark"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/pagerduty"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/slack"

	alertRepository "github.com/davidsugianto/sentinel-incident/internal/repository/alert"
	incidentRepository "github.com/davidsugianto/sentinel-incident/internal/repository/incident"
	oncallRepository "github.com/davidsugianto/sentinel-incident/internal/repository/oncall"
	incidentUsecase "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
	oncallUsecase "github.com/davidsugianto/sentinel-incident/internal/usecase/oncall"
)

// @title Sentinel Incident API
// @version 1.0
// @description Incident management and alerting service API with multi-channel support (Slack, Lark) and on-call integrations (PagerDuty).
// @contact.name API Support
// @contact.email support@sentinel-incident.local
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Config
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}
	cfg := config.GetConfig()

	// DB
	dbClient, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	// Run migrations
	if err := db.Migrate(dbClient); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Initialize Slack client
	var slackClient *slack.Client
	if cfg.Slack.Enabled && cfg.Slack.WebhookURL != "" {
		slackClient = slack.NewClient(cfg.Slack.WebhookURL)
	}

	// Initialize Lark client
	var larkClient *lark.Client
	if cfg.Lark.Enabled && cfg.Lark.WebhookURL != "" {
		larkClient = lark.NewClient(cfg.Lark.WebhookURL)
	}

	// Initialize PagerDuty client
	var pagerdutyClient *pagerduty.Client
	if cfg.PagerDuty.Enabled && cfg.PagerDuty.APIToken != "" {
		pagerdutyClient = pagerduty.NewClient(cfg.PagerDuty.APIToken)
	}

	// Repositories
	incidentRepo := incidentRepository.New(incidentRepository.Dependencies{
		Database: dbClient,
	})

	alertRepo := alertRepository.New(alertRepository.Dependencies{
		Database:    dbClient,
		SlackClient: slackClient,
		SlackConfig: &cfg.Slack,
		LarkClient:  larkClient,
		LarkConfig:  &cfg.Lark,
	})

	oncallRepo := oncallRepository.New(oncallRepository.Dependencies{
		Database:        dbClient,
		PagerDutyClient: pagerdutyClient,
		PagerDutyConfig: &cfg.PagerDuty,
	})

	// Usecases
	incidentUC := incidentUsecase.New(incidentUsecase.Dependencies{
		IncidentRepo: incidentRepo,
		AlertRepo:    alertRepo,
	})

	oncallUC := oncallUsecase.New(oncallUsecase.Dependencies{
		OncallRepo: oncallRepo,
	})

	server := New(Dependencies{
		IncidentUseCase: incidentUC,
		OncallUseCase:   oncallUC,
		Config:          cfg,
	})

	log.Printf("listening on :%d", cfg.Server.Port)
	if err := server.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
