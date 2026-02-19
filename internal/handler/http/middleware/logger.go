package middleware

import (
	"time"

	"github.com/davidsugianto/sentinel-incident/internal/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info(c.Request.Context(), "incoming request", map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   status,
			"latency":  latency.String(),
			"clientIP": c.ClientIP(),
		})
	}
}
