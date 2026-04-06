package http

import (
	"github.com/davidsugianto/sentinel-incident/internal/handler/http/middleware"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authConfig *config.Auth
}

func NewAuthHandler(cfg *config.Auth) *AuthHandler {
	return &AuthHandler{authConfig: cfg}
}

type LoginRequest struct {
	UserID string `json:"user_id" binding:"required"`
	TeamID string `json:"team_id" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err)
		return
	}

	token, err := middleware.GenerateToken(h.authConfig, req.UserID, req.TeamID)
	if err != nil {
		response.Fail(c, 500, err)
		return
	}

	response.OK(c, gin.H{
		"token": token,
		"type":  "Bearer",
	})
}
