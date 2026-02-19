package http

import (
	"github.com/davidsugianto/sentinel-incident/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateIncident(c *gin.Context) {
	data := map[string]string{"message": "CreateIncident"}
	response.OK(c, data)
}
