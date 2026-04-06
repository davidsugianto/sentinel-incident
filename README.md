# Sentinel Incident

Incident management tool that supports alerting across multiple channels with easy custom messaging and on-call integrations.

## Features

- **Multi-channel Alerts**: Send incident notifications to Slack, Lark (more platforms coming)
- **Custom Templates**: Define your own alert messages using Go templates
- **Easy Configuration**: YAML-based configuration with environment variables support
- **On-call Integrations**: PagerDuty integration for on-call rotations

## Architecture

![Sentinel Incident Architecture](./docs/images/sentinel-incident.svg)
https://excalidraw.com/#json=4DCgEjDo_3VNjWlv2PxbS,_DraZtMnOr-C8GnPH6akOw

```
┌─────────────────────────────────────────────────────────────────┐
│                        cmd/http                                  │
│                    (Entry Point & Server)                        │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                    internal/handler/http                         │
│                    (HTTP Handlers / Controllers)                 │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                    internal/usecase                              │
│                       (Business Logic)                           │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                   internal/repository                            │
│                     (Data Access Layer)                          │
├─────────────────────┬─────────────────────┬─────────────────────┤
│        alert        │      incident       │        oncall        │
│   (Alert Channels)  │    (Incident DB)    │   (PagerDuty API)   │
└─────────────────────┴─────────────────────┴─────────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.25.7+
- PostgreSQL 14+
- Docker (optional)

### Local Development

```bash
# Clone the repository
git clone https://github.com/davidsugianto/sentinel-incident.git
cd sentinel-incident

# Install dependencies
go mod download

# Copy and configure settings
cp configs/config.yaml configs/config.local.yaml
# Edit config.yaml with your database credentials

# Run the application
go run ./cmd/http
```

### Docker Development

```bash
# Build and run with hot-reload
docker build -t sentinel-incident .
docker run -p 8080:8080 -v $(pwd):/app sentinel-incident
```

### With Air (Hot Reload)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot-reload
air -c .air.toml
```

## API Endpoints

| Method | Endpoint          | Description        |
|--------|-------------------|--------------------|
| GET    | `/v1/ping`        | Health check       |
| POST   | `/v1/incidents`   | Create incident    |

## Configuration

Configuration is loaded from `configs/config.yaml`. The application supports:

- **Hot-reload**: Configuration changes are applied automatically
- **Environment variables**: Override any config via env vars

### Configuration Structure

```yaml
server:
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: sentinel_incident
  sslmode: disable

auth:
  jwt_secret: "your-secret-key"

cors:
  allowed_origins:
    - "http://localhost:3000"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  allow_credentials: true
```

## Project Structure

```
sentinel-incident/
├── cmd/
│   └── http/              # Application entry point
├── internal/
│   ├── handler/http/      # HTTP handlers and routes
│   │   └── middleware/    # HTTP middlewares
│   ├── model/             # Domain models
│   ├── pkg/               # Shared utilities
│   │   ├── config/        # Configuration loader
│   │   ├── db/            # Database connection
│   │   ├── logger/        # Logging utilities
│   │   └── response/      # HTTP response helpers
│   ├── repository/        # Data access layer
│   │   ├── alert/         # Alert channel integrations
│   │   ├── incident/      # Incident persistence
│   │   └── oncall/        # On-call integrations
│   └── usecase/           # Business logic
│       └── incident/      # Incident usecases
├── configs/               # Configuration files
├── docs/                  # Documentation
├── Dockerfile
├── go.mod
└── README.md
```

## Tech Stack

- **Language**: Go 1.25.7
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Configuration**: Viper
- **Logging**: Logrus

## License

MIT License
