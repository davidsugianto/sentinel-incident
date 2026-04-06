package http

import (
	"strconv"

	"github.com/davidsugianto/sentinel-incident/internal/model/incident"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/response"
	incidentUsecase "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateIncidentRequest struct {
	TeamID      string                 `json:"team_id" binding:"required"`
	Title       string                 `json:"title" binding:"required"`
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
	Severity    incident.Severity      `json:"severity"`
}

type UpdateIncidentRequest struct {
	Title       *string                 `json:"title,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Content     *map[string]interface{} `json:"content,omitempty"`
	Status      *incident.Status        `json:"status,omitempty"`
	Severity    *incident.Severity      `json:"severity,omitempty"`
}

func (h *Handler) CreateIncident(c *gin.Context) {
	var req CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	inc := &incident.Incident{
		TeamID:      req.TeamID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Severity:    req.Severity,
	}

	result, err := h.incidentUseCase.Create(c.Request.Context(), inc)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.Created(c, result)
}

func (h *Handler) GetIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	result, err := h.incidentUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, 404, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) ListIncidents(c *gin.Context) {
	var params incidentUsecase.ListParams

	params.TeamID = c.Query("team_id")
	params.Status = incident.Status(c.Query("status"))
	params.Severity = incident.Severity(c.Query("severity"))

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	if p, err := strconv.Atoi(page); err == nil {
		params.Page = p
	}
	if ps, err := strconv.Atoi(pageSize); err == nil {
		params.PageSize = ps
	}

	result, err := h.incidentUseCase.List(c.Request.Context(), params)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) UpdateIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	var req UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	result, err := h.incidentUseCase.Update(c.Request.Context(), id, &incidentUsecase.UpdateIncidentRequest{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Status:      req.Status,
		Severity:    req.Severity,
	})
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, result)
}

func (h *Handler) DeleteIncident(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.incidentUseCase.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, gin.H{"message": "incident deleted"})
}
