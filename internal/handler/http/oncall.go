package http

import (
	oncallModel "github.com/davidsugianto/sentinel-incident/internal/model/oncall"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	TeamID             string                 `json:"team_id" binding:"required"`
	Provider           oncallModel.Provider   `json:"provider" binding:"required"`
	ScheduleID         string                 `json:"schedule_id" binding:"required"`
	Config             map[string]interface{} `json:"config"`
	EscalationPolicyID string                 `json:"escalation_policy_id"`
}

type SyncSchedulesRequest struct {
	Provider oncallModel.Provider `json:"provider" binding:"required"`
}

func (h *Handler) GetOnCall(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	result, err := h.oncallUseCase.GetOnCallUser(c.Request.Context(), teamID)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	if result == nil {
		response.OK(c, gin.H{"message": "no on-call user found"})
		return
	}

	response.OK(c, result)
}

func (h *Handler) SyncSchedules(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	var req SyncSchedulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.oncallUseCase.SyncSchedules(c.Request.Context(), teamID, req.Provider); err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, gin.H{"message": "schedules synced successfully"})
}

func (h *Handler) ListSchedules(c *gin.Context) {
	teamIDStr := c.Param("team_id")
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	result, err := h.oncallUseCase.ListSchedules(c.Request.Context(), teamID)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) GetSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	result, err := h.oncallUseCase.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) CreateSchedule(c *gin.Context) {
	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	teamID, err := uuid.Parse(req.TeamID)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	// Add escalation policy ID to config if provided
	config := req.Config
	if config == nil {
		config = make(map[string]interface{})
	}
	if req.EscalationPolicyID != "" {
		config["escalation_policy_id"] = req.EscalationPolicyID
	}

	schedule := &oncallModel.OnCallSchedule{
		TeamID:     teamID,
		Provider:   req.Provider,
		ScheduleID: req.ScheduleID,
		Config:     config,
		IsActive:   true,
	}

	if err := h.oncallUseCase.CreateSchedule(c.Request.Context(), schedule); err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.Created(c, schedule)
}

func (h *Handler) UpdateSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	var req CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	schedule, err := h.oncallUseCase.GetScheduleByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err)
		return
	}

	if req.ScheduleID != "" {
		schedule.ScheduleID = req.ScheduleID
	}
	if req.Config != nil {
		schedule.Config = req.Config
	}
	if req.EscalationPolicyID != "" {
		if schedule.Config == nil {
			schedule.Config = make(map[string]interface{})
		}
		schedule.Config["escalation_policy_id"] = req.EscalationPolicyID
	}

	if err := h.oncallUseCase.UpdateSchedule(c.Request.Context(), schedule); err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, schedule)
}

func (h *Handler) DeleteSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.oncallUseCase.DeleteSchedule(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, gin.H{"message": "schedule deleted"})
}
