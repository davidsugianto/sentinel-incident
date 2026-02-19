package main

import (
	"fmt"
	"log"

	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/db"

	alertRepository "github.com/davidsugianto/sentinel-incident/internal/repository/alert"
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
		return
	}

	// Repository
	alertRepo := alertRepository.New(alertRepository.Dependencies{
		Database: dbClient,
	})

	// Usecases
	incidentUC := incidentUsecase.New(incidentUsecase.Dependencies{
		AlertRepo: alertRepo,
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
