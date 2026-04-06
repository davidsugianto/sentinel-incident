package oncall

import (
	"context"

	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	oncallRepo "github.com/davidsugianto/sentinel-incident/internal/repository/oncall"
	"github.com/google/uuid"
)

type Usecase interface {
	SyncSchedules(ctx context.Context, teamID uuid.UUID, provider oncallModel.Provider) error
	GetOnCallUser(ctx context.Context, teamID uuid.UUID) (*oncallRepo.OnCallUserInfo, error)
	GetScheduleByID(ctx context.Context, id uuid.UUID) (*oncallModel.OnCallSchedule, error)
	CreateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error
	UpdateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error
	DeleteSchedule(ctx context.Context, id uuid.UUID) error
	ListSchedules(ctx context.Context, teamID uuid.UUID) ([]oncallModel.OnCallSchedule, error)
}

type usecase struct {
	oncallRepo oncallRepo.Repository
}

type Dependencies struct {
	OncallRepo oncallRepo.Repository
}

func New(deps Dependencies) Usecase {
	return &usecase{
		oncallRepo: deps.OncallRepo,
	}
}
