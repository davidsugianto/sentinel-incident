package team

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug      string    `gorm:"type:varchar(255);uniqueIndex" json:"slug"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Team) TableName() string {
	return "teams"
}

type ChannelType string

const (
	ChannelTypeSlack   ChannelType = "slack"
	ChannelTypeLark    ChannelType = "lark"
	ChannelTypeWebhook ChannelType = "webhook"
)

type TeamChannel struct {
	ID          uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TeamID      uuid.UUID              `gorm:"type:uuid;not null;index" json:"team_id"`
	ChannelType ChannelType            `gorm:"type:varchar(50);not null" json:"channel_type"`
	Config      map[string]interface{} `gorm:"type:jsonb" json:"config"`
	IsActive    bool                   `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
}

func (TeamChannel) TableName() string {
	return "team_channels"
}
