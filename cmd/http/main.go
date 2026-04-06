package main

import (
	"fmt"
	"log"

	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/db"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/slack"

	alertRepository "github.com/davidsugianto/sentinel-incident/internal/repository/alert"
	incidentRepository "github.com/davidsugianto/sentinel-incident/internal/repository/incident"
	incidentUsecase "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
)

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

	// Repositories
	incidentRepo := incidentRepository.New(incidentRepository.Dependencies{
		Database: dbClient,
	})

	alertRepo := alertRepository.New(alertRepository.Dependencies{
		Database:    dbClient,
		SlackClient: slackClient,
		SlackConfig: &cfg.Slack,
	})

	// Usecases
	incidentUC := incidentUsecase.New(incidentUsecase.Dependencies{
		IncidentRepo: incidentRepo,
		AlertRepo:    alertRepo,
	})

	server := New(Dependencies{
		IncidentUseCase: incidentUC,
		Config:          cfg,
	})

	log.Printf("listening on :%d", cfg.Server.Port)
	if err := server.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal(err)
	}
}
