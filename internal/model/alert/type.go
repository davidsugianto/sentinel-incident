package alert

import (
	"time"

	"github.com/google/uuid"
)

type Channel string
type AlertStatus string

const (
	ChannelSlack Channel = "slack"
	ChannelLark  Channel = "lark"
	ChannelEmail Channel = "email"

	AlertStatusPending AlertStatus = "pending"
	AlertStatusSent    AlertStatus = "sent"
	AlertStatusFailed  AlertStatus = "failed"
)

type Alert struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IncidentID   uuid.UUID   `gorm:"type:uuid;not null;index" json:"incident_id"`
	Channel      Channel     `gorm:"type:varchar(50);not null" json:"channel"`
	Status       AlertStatus `gorm:"type:varchar(50);default:'pending';index" json:"status"`
	SentAt       *time.Time  `json:"sent_at,omitempty"`
	ErrorMessage string      `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
}

func (Alert) TableName() string {
	return "alerts"
}
