package oncall

import (
	"context"

	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	oncallRepo "github.com/davidsugianto/sentinel-incident/internal/repository/oncall"
	"github.com/google/uuid"
)

func (u *usecase) SyncSchedules(ctx context.Context, teamID uuid.UUID, provider oncallModel.Provider) error {
	return u.oncallRepo.SyncSchedules(ctx, teamID, provider)
}

func (u *usecase) GetOnCallUser(ctx context.Context, teamID uuid.UUID) (*oncallRepo.OnCallUserInfo, error) {
	return u.oncallRepo.GetOnCallUser(ctx, teamID)
}

func (u *usecase) GetScheduleByID(ctx context.Context, id uuid.UUID) (*oncallModel.OnCallSchedule, error) {
	return u.oncallRepo.GetScheduleByID(ctx, id)
}

func (u *usecase) CreateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error {
	return u.oncallRepo.CreateSchedule(ctx, schedule)
}

func (u *usecase) UpdateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error {
	return u.oncallRepo.UpdateSchedule(ctx, schedule)
}

func (u *usecase) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	return u.oncallRepo.DeleteSchedule(ctx, id)
}

func (u *usecase) ListSchedules(ctx context.Context, teamID uuid.UUID) ([]oncallModel.OnCallSchedule, error) {
	return u.oncallRepo.ListSchedules(ctx, teamID)
}
