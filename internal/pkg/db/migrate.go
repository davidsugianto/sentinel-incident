package db

import (
	"log"

	alertModel "github.com/davidsugianto/sentinel-incident/internal/model/alert"
	incidentModel "github.com/davidsugianto/sentinel-incident/internal/model/incident"
	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	teamModel "github.com/davidsugianto/sentinel-incident/internal/model/team"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&incidentModel.Incident{},
		&alertModel.Alert{},
		&teamModel.Team{},
		&teamModel.TeamChannel{},
		&oncallModel.OnCallSchedule{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}
