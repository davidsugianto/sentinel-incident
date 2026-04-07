# Development Guidelines

This document provides guidelines and best practices for developing Sentinel Incident.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Git Workflow](#git-workflow)
- [Testing Guidelines](#testing-guidelines)
- [API Design Guidelines](#api-design-guidelines)
- [Database Guidelines](#database-guidelines)
- [Error Handling](#error-handling)
- [Logging Guidelines](#logging-guidelines)
- [Security Guidelines](#security-guidelines)

---

## Getting Started

### Prerequisites

- **Go**: 1.24+ (required for `crypto/sha3` in standard library)
- **Docker**: Latest version
- **Docker Compose**: v2+
- **Make**: Build automation

### Required Tools

```bash
# Install Go tools
go install github.com/cosmtrek/air@latest           # Hot-reload
go install github.com/swaggo/swag/cmd/swag@latest   # Swagger docs
go install github.com/vektra/mockery/v2@latest      # Mock generation
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linting
go install golang.org/x/tools/cmd/goimports@latest  # Import organization
```

### Quick Start

```bash
# Clone and start development environment
git clone <repository-url>
cd sentinel-incident

# Copy environment file
cp .env.example .env

# Start all services with hot-reload
make docker-dev
```

---

## Development Environment

### Local Development (without Docker)

```bash
# Run PostgreSQL locally (or use Docker)
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=pass \
  -e POSTGRES_DB=sentinel_incident \
  postgres:15

# Run migrations
make db-migrate

# Start with hot-reload
make dev
```

### Docker Development (Recommended)

```bash
# Start all services
make docker-dev

# Or start in background
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Logging level | `debug` |
| `GIN_MODE` | Gin framework mode | `debug` |

See `.env.example` for complete list.

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
│   ├── usecase/           # Business logic layer
│   ├── repository/        # Data access layer
│   │   ├── alert/         # Alert repositories (Slack, Lark)
│   │   └── incident/      # Incident repositories
│   ├── model/             # Domain models and entities
│   └── pkg/               # Internal utilities
│       ├── config/        # Configuration management
│       ├── db/            # Database connection
│       ├── logger/        # Logging utilities
│       └── response/      # HTTP response helpers
├── migrations/            # Database migrations
├── configs/               # Configuration files
│   └── config.yaml        # Default configuration
├── docs/                  # Documentation
│   ├── swagger.json       # OpenAPI specification
│   └── swagger.yaml       # OpenAPI specification
├── test/                  # Integration tests
├── Dockerfile             # Container definition
├── docker-compose.yml     # Service orchestration
├── Makefile               # Build automation
└── go.mod                 # Go module definition
```

### Layer Responsibilities

| Layer | Responsibility | Dependencies |
|-------|---------------|--------------|
| **Handler** | HTTP request/response, validation | Usecase |
| **Usecase** | Business logic, orchestration | Repository |
| **Repository** | Data access, external integrations | Database, APIs |
| **Model** | Domain entities, DTOs | None |
| **Pkg** | Shared utilities | None |

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
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting (no code change)
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance

**Examples:**
```
feat(alert): add Slack webhook integration
fix(incident): resolve pagination offset issue
docs(readme): update installation instructions
test(usecase): add incident creation tests
```

### Pull Request Process

1. Create feature branch from `master`
2. Make changes with clear commits
3. Run tests: `make test`
4. Run linter: `make lint`
5. Update documentation if needed
6. Create PR with description of changes
7. Request review
8. Address review comments
9. Squash and merge

---

## Testing Guidelines

### Test Structure

```
internal/
├── usecase/
│   └── incident/
│       ├── usecase.go
│       └── usecase_test.go
└── repository/
    └── incident/
        ├── repository.go
        └── repository_test.go
```

### Test Naming

```go
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {}

// Examples
func TestCreateIncident_ValidInput_ReturnsIncident(t *testing.T) {}
func TestCreateIncident_EmptyContent_ReturnsError(t *testing.T) {}
func TestGetIncident_NotFound_ReturnsError(t *testing.T) {}
```

### Test Pattern

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
    assert.True(t, mockRepo.CreateCalled)
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

# Verbose
go test -v ./...
```

### Mock Generation

```bash
# Generate mocks
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

---

## Database Guidelines

### Migrations

Create migrations in `migrations/` directory:

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

-- migrations/20240101000000_create_incidents_table.down.sql
DROP TABLE IF EXISTS incidents;
```

### Migration Commands

```bash
# Run migrations
make db-migrate

# Rollback last migration
make db-rollback

# Reset database
make db-reset
```

### GORM Best Practices

```go
// Use struct for model
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

| Level | Usage |
|-------|-------|
| `DEBUG` | Development details, temporary debugging |
| `INFO` | Normal operations, business events |
| `WARN` | Unexpected but handled situations |
| `ERROR` | Errors that affect operations |
| `FATAL` | Unrecoverable errors, app termination |

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

### Sensitive Data

- Never log sensitive data (passwords, tokens, PII)
- Use environment variables for secrets
- Encrypt sensitive data at rest
- Use HTTPS in production

---

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
lsof -i :8080
kill -9 <PID>
```

**Database connection failed:**
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U user -d sentinel_incident
```

**Hot-reload not working:**
```bash
# Check Air is installed
which air

# Check .air.toml configuration
cat .air.toml
```

**Module dependency issues:**
```bash
go mod tidy
go mod download
```

---

## Useful Commands

```bash
# Development
make dev              # Run with hot-reload
make docker-dev       # Run with Docker + hot-reload

# Build
make build            # Build binary
make docker-build     # Build Docker image

# Database
make db-migrate       # Run migrations
make db-reset         # Reset database

# Testing
make test             # Run tests
make test-coverage    # Run with coverage

# Code Quality
make lint             # Run linter
make fmt              # Format code
make tidy             # Tidy modules

# Utilities
make swagger          # Generate Swagger docs
make mocks            # Generate mocks
make clean            # Clean build artifacts
```
