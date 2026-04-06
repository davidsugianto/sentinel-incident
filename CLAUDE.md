# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sentinel Incident is an incident management tool that supports multi-channel alerting (Slack, Lark) with custom Go templates and on-call integrations (PagerDuty).

## Build and Run Commands

```bash
# Build the application
go build -o ./bin/sentinel-incident ./cmd/http

# Run directly
go run ./cmd/http

# Run with hot-reload (requires air: go install github.com/cosmtrek/air@latest)
air -c .air.toml

# Run via Docker (development with hot-reload)
docker build -t sentinel-incident . && docker run -p 8080:8080 sentinel-incident
```

## Testing

No tests exist yet. When adding tests, use Go's standard testing pattern:

```bash
go test ./...
go test -v ./internal/usecase/incident/...
```

## Architecture

This project follows clean architecture with dependency injection:

```
cmd/http/           → Entry point, server setup, route registration
internal/
├── handler/http/   → HTTP handlers (controllers) - receives requests, calls usecases
├── usecase/        → Business logic - orchestrates repositories
├── repository/     → Data access layer - database and external integrations
├── model/          → Domain models
└── pkg/            → Shared utilities (config, db, logger, response)
```

### Dependency Injection Pattern

Each layer uses a `Dependencies` struct for constructor injection:

```go
// Repository
type Dependencies struct {
    Database *gorm.DB
}

// Usecase
type Dependencies struct {
    AlertRepo alertRepo.Repository
}

// Handler
type Dependencies struct {
    IncidentUseCase incident.Usecase
}
```

### Key Interfaces

- `alert.Repository` - SendAlert(ctx, *Incident) - handles alert delivery to channels
- `incident.Usecase` - CreateIncident(ctx, teamID, content, params) - business logic

## Configuration

YAML-based config at `configs/config.yaml`. Uses Viper with:
- Hot-reload on file changes
- Environment variable support via `AutomaticEnv()`

Config is loaded once as a singleton and accessed via `config.GetConfig()`.

## API Endpoints

All routes are versioned under `/v1`:

- `GET /v1/ping` - Health check
- `POST /v1/incidents` - Create incident

## Database

PostgreSQL with GORM. Connection configured via `configs/config.yaml` under `database` key.
