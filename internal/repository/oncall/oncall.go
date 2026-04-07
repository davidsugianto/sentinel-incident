package oncall

import (
	"context"
	"errors"

	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *repository) SyncSchedules(ctx context.Context, teamID uuid.UUID, provider oncallModel.Provider) error {
	if provider != oncallModel.ProviderPagerDuty {
		logger.Info(ctx, "Provider not supported for sync", map[string]interface{}{
			"provider": provider,
		})
		return nil
	}

	if r.pagerdutyClient == nil || !r.pagerdutyConfig.Enabled {
		logger.Info(ctx, "PagerDuty client not configured", nil)
		return nil
	}

	// Fetch schedules from PagerDuty
	schedules, err := r.pagerdutyClient.GetSchedules(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to fetch PagerDuty schedules", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Sync schedules to database
	for _, pdSchedule := range schedules {
		var existing oncallModel.OnCallSchedule
		err := r.db.Where("team_id = ? AND schedule_id = ?", teamID, pdSchedule.ID).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new schedule
			newSchedule := &oncallModel.OnCallSchedule{
				TeamID:     teamID,
				Provider:   provider,
				ScheduleID: pdSchedule.ID,
				Config: map[string]interface{}{
					"name": pdSchedule.Name,
				},
				IsActive: true,
			}
			if err := r.db.Create(newSchedule).Error; err != nil {
				logger.Error(ctx, "Failed to create schedule", map[string]interface{}{
					"schedule_id": pdSchedule.ID,
					"error":       err.Error(),
				})
				continue
			}
		} else if err == nil {
			// Update existing schedule
			existing.Config = map[string]interface{}{
				"name": pdSchedule.Name,
			}
			if err := r.db.Save(&existing).Error; err != nil {
				logger.Error(ctx, "Failed to update schedule", map[string]interface{}{
					"schedule_id": pdSchedule.ID,
					"error":       err.Error(),
				})
			}
		}
	}

	logger.Info(ctx, "Schedules synced successfully", map[string]interface{}{
		"team_id":   teamID,
		"provider":  provider,
		"schedules": len(schedules),
	})

	return nil
}

func (r *repository) GetOnCallUser(ctx context.Context, teamID uuid.UUID) (*OnCallUserInfo, error) {
	if r.pagerdutyClient == nil || !r.pagerdutyConfig.Enabled {
		logger.Info(ctx, "PagerDuty client not configured", nil)
		return nil, nil
	}

	// Get active schedules for team
	var schedules []oncallModel.OnCallSchedule
	if err := r.db.Where("team_id = ? AND is_active = ?", teamID, true).Find(&schedules).Error; err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, nil
	}

	// Get escalation policies from config
	var policyIDs []string
	for _, schedule := range schedules {
		if policyID, ok := schedule.Config["escalation_policy_id"].(string); ok && policyID != "" {
			policyIDs = append(policyIDs, policyID)
		}
	}

	if len(policyIDs) == 0 {
		return nil, nil
	}

	// Fetch on-call users from PagerDuty
	onCallUsers, err := r.pagerdutyClient.GetOnCallUsers(ctx, policyIDs)
	if err != nil {
		return nil, err
	}

	if len(onCallUsers) == 0 {
		return nil, nil
	}

	// Return the first on-call user
	user := onCallUsers[0]
	return &OnCallUserInfo{
		UserID:     user.User.ID,
		UserName:   user.User.Name,
		UserEmail:  user.User.Email,
		PolicyID:   user.EscalationPolicy.ID,
		PolicyName: user.EscalationPolicy.Name,
		ScheduleID: user.Schedule.ID,
	}, nil
}

func (r *repository) GetScheduleByID(ctx context.Context, id uuid.UUID) (*oncallModel.OnCallSchedule, error) {
	var schedule oncallModel.OnCallSchedule
	if err := r.db.First(&schedule, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

func (r *repository) CreateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error {
	return r.db.Create(schedule).Error
}

func (r *repository) UpdateSchedule(ctx context.Context, schedule *oncallModel.OnCallSchedule) error {
	return r.db.Save(schedule).Error
}

func (r *repository) DeleteSchedule(ctx context.Context, id uuid.UUID) error {
	result := r.db.Delete(&oncallModel.OnCallSchedule{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *repository) ListSchedules(ctx context.Context, teamID uuid.UUID) ([]oncallModel.OnCallSchedule, error) {
	var schedules []oncallModel.OnCallSchedule
	if err := r.db.Where("team_id = ?", teamID).Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}
