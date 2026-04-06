package oncall

import (
	"context"
	"errors"

	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/pagerduty"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("on-call schedule not found")

// Repository interface for on-call operations
type Repository interface {
	SyncSchedules(ctx context.Context, teamID uuid.UUID, provider oncallModel.Provider) error
	GetOnCallUser(ctx context.Context, teamID uuid.UUID) (*OnCallUserInfo, error)
	GetScheduleByID(ctx context.Context, id uuid.UUID) (*oncallModel.OnCallSchedule, error)
	CreateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error
	UpdateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	ListSchedules(ctx context.Context, teamID uuid.UUID) ([]oncallModel.OnCallSchedule, error)
}

// OnCallUserInfo represents the current on-call user information
type OnCallUserInfo struct {
	UserID     string
	UserName   string
	UserEmail  string
	PolicyID   string
	PolicyName string
	ScheduleID string
}

type repository struct {
	db              *gorm.DB
	pagerdutyClient *pagerduty.Client
	pagerdutyConfig *config.PagerDuty
}

type Dependencies struct {
	Database        *gorm.DB
	PagerDutyClient *pagerduty.Client
	PagerDutyConfig *config.PagerDuty
}

func New(deps Dependencies) Repository {
	return &repository{
		db:              deps.Database,
		pagerdutyClient: deps.PagerDutyClient,
		pagerdutyConfig: deps.PagerDutyConfig,
	}
}
