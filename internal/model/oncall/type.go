package oncall

import (
	"time"

	"github.com/google/uuid"
)

type Provider string

const (
	ProviderPagerDuty Provider = "pagerduty"
	ProviderOpsGenie  Provider = "opsgenie"
)

type OnCallSchedule struct {
	ID         uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TeamID     uuid.UUID              `gorm:"type:uuid;not null;index" json:"team_id"`
	Provider   Provider               `gorm:"type:varchar(50);not null" json:"provider"`
	ScheduleID string                 `gorm:"type:varchar(255);not null" json:"schedule_id"`
	Config     map[string]interface{} `gorm:"type:jsonb" json:"config"`
	IsActive   bool                   `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
}

func (OnCallSchedule) TableName() string {
	return "oncall_schedules"
}
