package db

import (
	"fmt"

	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
