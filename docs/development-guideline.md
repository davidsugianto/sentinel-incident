# Development Guidelines

This document provides guidelines and best practices for developing Sentinel Incident.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
  - [Local Development](#local-development-without-docker)
  - [Docker Development](#docker-development-recommended)
  - [Docker Services](#docker-services)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Git Workflow](#git-workflow)
- [Testing Guidelines](#testing-guidelines)
- [API Design Guidelines](#api-design-guidelines)
- [Database Guidelines](#database-guidelines)
- [Error Handling](#error-handling)
- [Logging Guidelines](#logging-guidelines)
- [Security Guidelines](#security-guidelines)
- [Makefile Reference](#makefile-reference)
- [Troubleshooting](#troubleshooting)

---

## Getting Started

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.24+ | Runtime (required for `crypto/sha3`) |
| **Docker** | Latest | Containerization |
| **Docker Compose** | v2+ | Service orchestration |
| **Make** | Any | Build automation |

### Required Tools

```bash
# Install Go tools
go install github.com/cosmtrek/air@latest           # Hot-reload
go install github.com/swaggo/swag/cmd/swag@latest   # Swagger docs
go install github.com/vektra/mockery/v2@latest      # Mock generation
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
go install golang.org/x/tools/cmd/goimports@latest  # Import organization

# Optional: Pre-commit framework
pip install pre-commit  # or: brew install pre-commit
```

### Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd sentinel-incident

# Copy environment file
cp .env.example .env

# Install git hooks
make hooks-install

# Start all services with hot-reload
make docker-dev
```

### Pre-commit Hooks

Pre-commit hooks ensure code quality before each commit.

**Install hooks:**
```bash
# Option 1: Using Makefile (recommended)
make hooks-install

# Option 2: Using pre-commit framework
pre-commit install

# Option 3: Manual installation
cp .githooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# Option 4: Configure git hooks path
git config core.hooksPath .githooks
```

**Run checks manually:**
```bash
make pre-commit
```

**Hooks run automatically on commit:**
| Check | Description |
|-------|-------------|
| `go fmt` | Format Go code |
| `go mod tidy` | Tidy Go modules |
| `go test ./...` | Run all unit tests |
| `go build` | Verify build passes |

**Uninstall hooks:**
```bash
make hooks-uninstall
```

---

## Development Environment

### Local Development (without Docker)

```bash
# Run PostgreSQL locally
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=pass \
  -e POSTGRES_DB=sentinel_incident \
  postgres:15

# Run Redis (optional)
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Run migrations
make db-migrate

# Start with hot-reload (requires air installed)
make dev
```

### Docker Development (Recommended)

```bash
# Start all services (foreground with logs)
make docker-dev

# Start in background
make docker-up

# View logs
make docker-logs

# View specific service logs
docker compose logs -f api
docker compose logs -f postgres

# Stop services
make docker-down

# Clean all (containers, volumes, images)
make docker-clean
```

### Docker Services

| Service | Port | Description | UI |
|---------|------|-------------|-----|
| **api** | 8080 | Application server | http://localhost:8080 |
| **postgres** | 5432 | PostgreSQL database | - |
| **redis** | 6379 | Redis cache | - |
| **adminer** | 8081 | Database UI | http://localhost:8081 |
| **redis-commander** | 8082 | Redis UI | http://localhost:8082 |

> Note: Adminer and Redis Commander are defined in `docker-compose.override.yml` and loaded automatically.

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `GIN_MODE` | Gin framework mode | `debug` | No |
| `DATABASE_HOST` | Database host | `localhost` | Yes |
| `DATABASE_PORT` | Database port | `5432` | Yes |
| `DATABASE_USER` | Database user | `user` | Yes |
| `DATABASE_PASSWORD` | Database password | `pass` | Yes |
| `DATABASE_NAME` | Database name | `sentinel_incident` | Yes |
| `DATABASE_URL` | Full connection URL (overrides individual vars) | - | No |
| `REDIS_ADDR` | Redis address | `localhost:6379` | No |
| `JWT_SECRET` | JWT signing secret | - | **Yes** |
| `LOG_LEVEL` | Logging level | `debug` | No |
| `LOG_FORMAT` | Log format (json/text) | `json` | No |
| `ENV` | Environment name | `development` | No |
| `SLACK_WEBHOOK_URL` | Slack webhook URL | - | No |
| `LARK_WEBHOOK_URL` | Lark webhook URL | - | No |
| `PAGERDUTY_API_TOKEN` | PagerDuty API token | - | No |

### Configuration Files

| File | Purpose |
|------|---------|
| `configs/config.yaml` | Default application configuration |
| `.env` | Environment variables (git-ignored) |
| `.env.example` | Environment variables template |
| `docker-compose.yml` | Docker services definition |
| `docker-compose.override.yml` | Local Docker overrides (auto-loaded) |
| `.air.toml` | Air hot-reload configuration |

---

## Project Structure

```
sentinel-incident/
├── cmd/
│   └── http/              # Application entry point
│       └── main.go        # Main entry file
├── internal/
│   ├── handler/           # HTTP handlers (controllers)
│   │   └── http/          # REST API handlers
│   │       └── middleware/ # HTTP middleware
│   ├── usecase/           # Business logic layer
│   ├── repository/        # Data access layer
│   │   ├── alert/         # Alert repositories (Slack, Lark)
│   │   └── incident/      # Incident repositories
│   ├── model/             # Domain models and entities
│   │   ├── alert/
│   │   ├── incident/
│   │   ├── oncall/
│   │   └── team/
│   └── pkg/               # Internal utilities
│       ├── config/        # Configuration management
│       ├── db/            # Database connection
│       ├── logger/        # Logging utilities
│       ├── pagerduty/     # PagerDuty client
│       └── response/      # HTTP response helpers
├── migrations/            # Database migrations
├── configs/               # Configuration files
├── docs/                  # Documentation
│   ├── swagger.json       # OpenAPI specification
│   ├── swagger.yaml       # OpenAPI specification
│   ├── development-guideline.md
│   └── development-todolist.md
├── .githooks/             # Git hooks
│   └── pre-commit         # Pre-commit hook script
├── Dockerfile             # Container definition
├── docker-compose.yml     # Service orchestration
├── docker-compose.override.yml  # Local overrides
├── .air.toml              # Air configuration
├── Makefile               # Build automation
└── go.mod                 # Go module definition
```

### Layer Responsibilities

| Layer | Responsibility | Dependencies | Example |
|-------|---------------|--------------|---------|
| **Handler** | HTTP request/response, validation | Usecase | `internal/handler/http/` |
| **Usecase** | Business logic, orchestration | Repository | `internal/usecase/` |
| **Repository** | Data access, external integrations | Database, APIs | `internal/repository/` |
| **Model** | Domain entities, DTOs | None | `internal/model/` |
| **Pkg** | Shared utilities | None | `internal/pkg/` |

### Dependency Flow

```
Handler → Usecase → Repository → Database/API
   ↓         ↓          ↓
 Model    Model      Model
```

---

## Coding Standards

### Go Style Guide

Follow [Effective Go](https://golang.org/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Naming Conventions

```go
// Packages: lowercase, single word preferred
package incident

// Interfaces: verb or noun ending with -er
type Repository interface {}
type AlertSender interface {}

// Structs: PascalCase
type IncidentService struct {}

// Functions/Methods: PascalCase (exported) or camelCase (unexported)
func CreateIncident() {}
func (s *Service) createAlert() {}

// Constants: PascalCase or UPPER_SNAKE_CASE
const MaxRetryAttempts = 3
const DEFAULT_TIMEOUT = 30 * time.Second

// Errors: descriptive, start with context
var ErrIncidentNotFound = errors.New("incident not found")
```

### Dependency Injection

Use constructor injection with Dependencies struct:

```go
// Repository
type Dependencies struct {
    Database *gorm.DB
}

func NewRepository(deps Dependencies) Repository {
    return &repository{db: deps.Database}
}

// Usecase
type Dependencies struct {
    IncidentRepo incident.Repository
    AlertRepo    alert.Repository
}

func NewUsecase(deps Dependencies) Usecase {
    return &usecase{
        incidentRepo: deps.IncidentRepo,
        alertRepo:    deps.AlertRepo,
    }
}

// Handler
type Dependencies struct {
    IncidentUsecase incident.Usecase
}

func NewHandler(deps Dependencies) Handler {
    return &handler{incidentUsecase: deps.IncidentUsecase}
}
```

### Error Handling

```go
// Return errors with context
if err != nil {
    return fmt.Errorf("failed to create incident: %w", err)
}

// Use custom error types for domain errors
var ErrIncidentNotFound = errors.New("incident not found")

func (r *repository) GetByID(id string) (*model.Incident, error) {
    var incident model.Incident
    if err := r.db.First(&incident, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrIncidentNotFound
        }
        return nil, fmt.Errorf("failed to get incident: %w", err)
    }
    return &incident, nil
}
```

---

## Git Workflow

### Branch Naming

```
f/feature-name      # Feature branch
b/bugfix-name       # Bug fix branch
h/hotfix-name       # Hot fix branch
r/release-x.x.x     # Release branch
```

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
| Type | Description |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation |
| `style` | Formatting (no code change) |
| `refactor` | Code refactoring |
| `test` | Adding tests |
| `chore` | Maintenance |

**Examples:**
```
feat(alert): add Slack webhook integration
fix(incident): resolve pagination offset issue
docs(readme): update installation instructions
test(usecase): add incident creation tests
chore(deps): update dependencies
```

### Pull Request Process

1. Create feature branch from `master`
2. Make changes with clear commits
3. **Pre-commit hooks run automatically**
4. Run tests manually: `make test`
5. Run linter: `make lint`
6. Update documentation if needed
7. Create PR with description of changes
8. Request review
9. Address review comments
10. Squash and merge

---

## Testing Guidelines

### Test Structure

```
internal/
├── usecase/
│   └── incident/
│       ├── usecase.go
│       └── usecase_test.go
├── repository/
│   └── incident/
│       ├── repository.go
│       └── repository_test.go
└── handler/
    └── http/
        ├── incident.go
        └── incident_test.go
```

### Test Naming

```go
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {}

// Examples
func TestCreateIncident_ValidInput_ReturnsIncident(t *testing.T) {}
func TestCreateIncident_EmptyContent_ReturnsError(t *testing.T) {}
func TestGetIncident_NotFound_ReturnsError(t *testing.T) {}
```

### Test Pattern (Arrange-Act-Assert)

```go
func TestCreateIncident_ValidInput_ReturnsIncident(t *testing.T) {
    // Arrange
    mockRepo := &MockRepository{
        CreateFunc: func(ctx context.Context, incident *model.Incident) error {
            incident.ID = "test-id"
            return nil
        },
    }
    deps := usecase.Dependencies{IncidentRepo: mockRepo}
    uc := usecase.NewUsecase(deps)

    // Act
    result, err := uc.CreateIncident(context.Background(), "team-1", "content", nil)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test-id", result.ID)
}
```

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/usecase/incident/...

# With coverage
make test-coverage

# Verbose output
go test -v ./...

# Run specific test
go test -run TestCreateIncident ./...
```

### Mock Generation

```bash
# Generate all mocks
make mocks

# Or with mockery directly
mockery --name=Repository --dir=./internal/repository/incident --output=./internal/mocks
```

---

## API Design Guidelines

### RESTful Conventions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/incidents` | List incidents |
| `GET` | `/v1/incidents/:id` | Get incident |
| `POST` | `/v1/incidents` | Create incident |
| `PUT` | `/v1/incidents/:id` | Update incident |
| `DELETE` | `/v1/incidents/:id` | Delete incident |
| `GET` | `/v1/ping` | Health check |
| `GET` | `/swagger/index.html` | API documentation |

### Request/Response Format

```json
// Success Response
{
  "code": 200,
  "message": "Success",
  "data": { ... }
}

// Error Response
{
  "code": 400,
  "message": "Validation error",
  "errors": [
    {"field": "content", "message": "required"}
  ]
}

// Paginated Response
{
  "code": 200,
  "message": "Success",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

### HTTP Status Codes

| Code | Usage |
|------|-------|
| `200` | Success |
| `201` | Created |
| `400` | Bad request / Validation error |
| `401` | Unauthorized |
| `404` | Not found |
| `500` | Internal server error |

### Swagger Documentation

```go
// CreateIncident godoc
// @Summary Create a new incident
// @Description Create a new incident with the provided details
// @Tags incidents
// @Accept json
// @Produce json
// @Param request body CreateIncidentRequest true "Incident data"
// @Success 201 {object} response.Response{data=model.Incident}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/incidents [post]
// @Security BearerAuth
func (h *handler) Create(c *gin.Context) { ... }
```

**Generate Swagger docs:**
```bash
make swagger
```

---

## Database Guidelines

### Migrations

Create migrations in `migrations/` directory with naming pattern:
```
{timestamp}_{description}.up.sql
{timestamp}_{description}.down.sql
```

Example:
```sql
-- migrations/20240101000000_create_incidents_table.up.sql
CREATE TABLE incidents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL,
    content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    severity VARCHAR(20) NOT NULL DEFAULT 'medium',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_incidents_team_id ON incidents(team_id);
CREATE INDEX idx_incidents_status ON incidents(status);

-- migrations/20240101000000_create_incidents_table.down.sql
DROP TABLE IF EXISTS incidents;
```

### Migration Commands

```bash
make db-migrate    # Run all pending migrations
make db-rollback   # Rollback last migration
make db-reset      # Drop all and re-run migrations
make db-drop       # Drop all migrations
```

### GORM Best Practices

```go
// Model definition with proper tags
type Incident struct {
    ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
    TeamID    string    `gorm:"type:uuid;not null;index"`
    Content   string    `gorm:"type:text;not null"`
    Status    string    `gorm:"type:varchar(20);not null;default:'open'"`
    Severity  string    `gorm:"type:varchar(20);not null;default:'medium'"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Use transactions for multiple operations
err := r.db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&incident).Error; err != nil {
        return err
    }
    if err := tx.Create(&alert).Error; err != nil {
        return err
    }
    return nil
})

// Use context for timeout/cancellation
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := r.db.WithContext(ctx).First(&incident, "id = ?", id).Error
```

---

## Error Handling

### Custom Error Types

```go
// Define domain errors
var (
    ErrIncidentNotFound = errors.New("incident not found")
    ErrInvalidStatus    = errors.New("invalid incident status")
    ErrUnauthorized     = errors.New("unauthorized access")
)

// Use errors.Is for comparison
if errors.Is(err, ErrIncidentNotFound) {
    return response.NotFound(c, "Incident not found")
}
```

### HTTP Error Responses

```go
// Handler error handling
func (h *handler) GetByID(c *gin.Context) {
    incident, err := h.uc.GetByID(c.Request.Context(), id)
    if err != nil {
        switch {
        case errors.Is(err, incident.ErrNotFound):
            response.NotFound(c, "Incident not found")
        case errors.Is(err, ErrUnauthorized):
            response.Unauthorized(c, "Unauthorized")
        default:
            response.InternalError(c, "Failed to get incident")
        }
        return
    }
    response.Success(c, incident)
}
```

---

## Logging Guidelines

### Log Levels

| Level | Usage | Example |
|-------|-------|---------|
| `DEBUG` | Development details, temporary debugging | Variable values, flow tracing |
| `INFO` | Normal operations, business events | Incident created, alert sent |
| `WARN` | Unexpected but handled situations | Retry attempt, deprecated API use |
| `ERROR` | Errors that affect operations | Database error, API failure |
| `FATAL` | Unrecoverable errors | Config load failure, startup error |

### Structured Logging

```go
import log "github.com/sirupsen/logrus"

// With fields
log.WithFields(log.Fields{
    "incident_id": id,
    "team_id":     teamID,
}).Info("Incident created successfully")

// Error logging
log.WithError(err).WithField("incident_id", id).Error("Failed to send alert")

// Debug (development only)
log.Debugf("Processing incident: %+v", incident)

// Warning
log.WithFields(log.Fields{
    "retry_count": count,
    "max_retries": maxRetries,
}).Warn("Retrying API call")
```

---

## Security Guidelines

### Authentication

```go
// JWT middleware
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            response.Unauthorized(c, "Missing authorization token")
            c.Abort()
            return
        }
        // Validate token...
        c.Next()
    }
}
```

### Input Validation

```go
type CreateIncidentRequest struct {
    TeamID   string `json:"team_id" binding:"required,uuid"`
    Content  string `json:"content" binding:"required,min=10,max=5000"`
    Severity string `json:"severity" binding:"omitempty,oneof=low medium high critical"`
}

func (h *handler) Create(c *gin.Context) {
    var req CreateIncidentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request", err.Error())
        return
    }
    // Process request...
}
```

### Security Checklist

- [ ] Never log sensitive data (passwords, tokens, PII)
- [ ] Use environment variables for secrets
- [ ] Validate all user input
- [ ] Use parameterized queries (GORM handles this)
- [ ] Enable CORS with specific origins only
- [ ] Use HTTPS in production
- [ ] Rotate JWT secrets periodically
- [ ] Implement rate limiting for production

---

## Makefile Reference

### Development Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run with Air hot-reload (local) |
| `make docker-dev` | Run with Docker Compose + Air |
| `make build` | Build the binary |
| `make run` | Build and run the binary |

### Docker Commands

| Command | Description |
|---------|-------------|
| `make docker-build` | Build Docker image |
| `make docker-up` | Start all services in background |
| `make docker-down` | Stop all services |
| `make docker-logs` | View container logs (follow) |
| `make docker-clean` | Remove containers, volumes, images |

### Database Commands

| Command | Description |
|---------|-------------|
| `make db-migrate` | Run database migrations |
| `make db-rollback` | Rollback last migration |
| `make db-reset` | Reset database (drop and recreate) |
| `make db-drop` | Drop all migrations |

### Testing & Quality Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make lint` | Run golangci-lint |
| `make fmt` | Format code with go fmt |
| `make imports` | Organize imports with goimports |
| `make tidy` | Run go mod tidy |

### Utility Commands

| Command | Description |
|---------|-------------|
| `make clean` | Clean build artifacts |
| `make deps` | Download dependencies |
| `make swagger` | Generate Swagger documentation |
| `make mocks` | Generate mocks with mockery |

### Git Hooks Commands

| Command | Description |
|---------|-------------|
| `make hooks-install` | Install pre-commit hooks |
| `make hooks-uninstall` | Remove pre-commit hooks |
| `make pre-commit` | Run pre-commit checks manually |

---

## Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port
lsof -i :8080
lsof -i :5432

# Kill process
kill -9 <PID>
```

#### Database Connection Failed

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U user -d sentinel_incident

# Restart database
docker compose restart postgres
```

#### Hot-reload Not Working

```bash
# Check Air is installed
which air

# Check .air.toml configuration
cat .air.toml

# Reinstall Air
go install github.com/cosmtrek/air@latest
```

#### Module Dependency Issues

```bash
# Clean and re-download
go clean -modcache
go mod download
go mod tidy
```

#### Docker Issues

```bash
# View logs
docker compose logs -f

# Restart all services
docker compose restart

# Clean rebuild
make docker-clean
make docker-dev
```

#### Pre-commit Hook Fails

```bash
# Run manually to see errors
make pre-commit

# Format code
make fmt

# Fix imports
make imports

# Run tests
make test
```

### Debug Mode

```bash
# Enable debug logging
export LOG_LEVEL=debug
export GIN_MODE=debug

# Run with verbose output
go run ./cmd/http
```

### Health Check

```bash
# Check if API is running
curl http://localhost:8080/v1/ping

# Check Swagger docs
open http://localhost:8080/swagger/index.html
```

---

## Quick Reference

```bash
# Start development
make docker-dev

# Run tests
make test

# Format and lint
make fmt && make lint

# Database operations
make db-migrate

# Generate docs
make swagger

# Install git hooks
make hooks-install
```
