package incident

import (
	"time"

	"github.com/google/uuid"
)

type Status string
type Severity string

const (
	StatusOpen         Status = "open"
	StatusAcknowledged Status = "acknowledged"
	StatusResolved     Status = "resolved"

	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

type Incident struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TeamID      string                 `gorm:"type:varchar(255);not null;index" json:"team_id"`
	Title       string                 `gorm:"type:varchar(255)" json:"title"`
	Description string                 `gorm:"type:text" json:"description"`
	Content     map[string]interface{} `gorm:"type:jsonb" json:"content"`
	Status      Status                 `gorm:"type:varchar(50);default:'open';index" json:"status"`
	Severity    Severity               `gorm:"type:varchar(50)" json:"severity"`
	Resolved    bool                   `gorm:"default:false" json:"resolved"`
	CreatedAt   time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	ResolvedAt  *time.Time             `json:"resolved_at,omitempty"`
}

func (Incident) TableName() string {
	return "incidents"
}
