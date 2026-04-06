package main

import (
	"net/http"

	"github.com/davidsugianto/go-pkgs/grace"
	httpHandler "github.com/davidsugianto/sentinel-incident/internal/handler/http"
	"github.com/davidsugianto/sentinel-incident/internal/handler/http/middleware"
	"github.com/davidsugianto/sentinel-incident/internal/pkg/config"
	incident "github.com/davidsugianto/sentinel-incident/internal/usecase/incident"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
	handler    *httpHandler.Handler
	authHandler *httpHandler.AuthHandler
	config     *config.Config
}

type Dependencies struct {
	IncidentUseCase incident.Usecase
	Config          *config.Config
}

func New(deps Dependencies) *Server {
	return &Server{
		Server: &http.Server{},
		handler: httpHandler.New(httpHandler.Dependencies{
			IncidentUseCase: deps.IncidentUseCase,
		}),
		authHandler: httpHandler.NewAuthHandler(&deps.Config.Auth),
		config:      deps.Config,
	}
}

func (s *Server) v1Endpoint(r *gin.Engine) {
	g := r.Group("/v1")
	g.Use(gin.Recovery(), middleware.RequestID(), middleware.Logger())

	// health check (public)
	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// auth routes (public)
	auth := g.Group("/auth")
	auth.POST("/login", s.authHandler.Login)

	// incident routes (protected)
	incidents := g.Group("/incidents")
	incidents.Use(middleware.JWT(&s.config.Auth))
	incidents.GET("", s.handler.ListIncidents)
	incidents.GET("/:id", s.handler.GetIncident)
	incidents.POST("", s.handler.CreateIncident)
	incidents.PUT("/:id", s.handler.UpdateIncident)
	incidents.DELETE("/:id", s.handler.DeleteIncident)
}

func (s *Server) Run(port string) error {
	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     s.config.CORS.AllowedOrigins,
		AllowMethods:     s.config.CORS.AllowedMethods,
		AllowHeaders:     s.config.CORS.AllowedHeaders,
		AllowCredentials: s.config.CORS.AllowCredentials,
	}
	r.Use(cors.New(corsConfig))

	s.v1Endpoint(r)

	s.Addr = port
	s.Handler = r

	return grace.ServeHTTP(s.Addr, s.Handler)
}
