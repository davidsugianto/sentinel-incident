# Development TODO List

## Status Legend

- [ ] Not Started
- [~] In Progress
- [x] Completed

---

## Core Features

### Incident Management

- [x] Implement full incident CRUD operations
  - [x] GET /v1/incidents - List incidents with pagination
  - [x] GET /v1/incidents/:id - Get incident by ID
  - [x] POST /v1/incidents - Create incident
  - [x] PUT /v1/incidents/:id - Update incident
  - [x] DELETE /v1/incidents/:id - Delete incident
- [x] Add incident status workflow (open → acknowledged → resolved)
- [x] Add severity levels support
- [ ] Implement incident search and filtering

### Alert Channels

- [x] Slack Integration
  - [x] Webhook configuration
  - [x] Message template support
  - [x] Rich message formatting (blocks, attachments)
- [x] Lark Integration
  - [x] Webhook configuration
  - [x] Message card templates
- [ ] Email Integration (future)
  - [ ] SMTP configuration
  - [ ] HTML template support

### On-Call Integration

- [x] PagerDuty Integration
  - [x] API client implementation
  - [x] Schedule sync
  - [x] Escalation policy mapping
- [ ] OpsGenie Integration (future)
- [ ] VictorOps Integration (future)

---

## Infrastructure

### Database

- [x] Add database migrations
  - [x] Create migrations table
  - [x] Incident table migration
  - [x] Alert table migration
  - [x] Team tables migration
- [x] Add connection pooling configuration
- [ ] Add database health check endpoint

### Authentication & Authorization

- [x] Implement JWT authentication middleware
- [ ] Add API key authentication
- [ ] Implement team-based authorization
- [ ] Add role-based access control (RBAC)

### Observability

- [ ] Structured logging implementation
- [ ] Metrics collection (Prometheus)
  - [ ] Request latency metrics
  - [ ] Error rate metrics
  - [ ] Alert delivery metrics
- [ ] Tracing support (OpenTelemetry)
- [ ] Health check endpoints
  - [ ] /health/live - Liveness probe
  - [ ] /health/ready - Readiness probe

---

## Quality Assurance

### Testing

- [x] Unit tests
  - [x] Usecase layer tests
  - [x] Repository layer tests
  - [x] Handler layer tests
- [ ] Integration tests
  - [ ] API endpoint tests
  - [ ] Database integration tests
- [x] Add test coverage reporting
- [x] Add mockery for interface mocking

### Code Quality

- [ ] Add golangci-lint configuration
- [ ] Add pre-commit hooks
- [ ] Add CI/CD pipeline
  - [ ] GitHub Actions workflow
  - [ ] Run tests on PR
  - [ ] Run linters on PR

---

## Documentation

- [x] API Documentation
  - [x] OpenAPI/Swagger specification
  - [x] Swagger UI endpoint
- [ ] Architecture decision records (ADRs)
- [ ] Runbook for operations
- [ ] Template documentation

---

## Deployment

- [ ] Docker optimization
  - [ ] Multi-stage build
  - [ ] Distroless image
- [ ] Kubernetes manifests
  - [ ] Deployment
  - [ ] Service
  - [ ] ConfigMap
  - [ ] Secret
  - [ ] HorizontalPodAutoscaler
- [ ] Helm chart
- [ ] Terraform modules

---

## Future Enhancements

- [ ] Webhook support for external integrations
- [ ] Incident templates
- [ ] Custom escalation policies
- [ ] SLA tracking and alerting
- [ ] Incident analytics dashboard
- [ ] Mobile push notifications
- [ ] SMS alerts (Twilio integration)

---

## Technical Debt

- [ ] Add proper error types and handling
- [ ] Implement graceful shutdown
- [ ] Add request validation
- [ ] Add rate limiting middleware
- [ ] Add request ID propagation

---

## Priority Order

### High Priority (MVP) - COMPLETED

- [x] Database migrations
- [x] Incident CRUD operations
- [x] Slack integration
- [x] Basic authentication

### Medium Priority - COMPLETED

1. [x] Lark integration
2. [x] PagerDuty integration
3. [x] Unit tests
4. [x] API documentation

### Low Priority

1. Email integration
2. Other on-call providers
3. Advanced analytics
4. Mobile push notifications
