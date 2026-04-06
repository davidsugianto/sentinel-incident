# Entity Relationship Diagram Documentation

## Overview

This document describes the database entities and their relationships for Sentinel Incident.

## ERD Diagram

```
┌─────────────────────────────────────┐
│              incidents               │
├─────────────────────────────────────┤
│ id          │ UUID     │ PK         │
│ team_id     │ VARCHAR  │ NOT NULL   │
│ title       │ VARCHAR  │            │
│ description │ TEXT     │            │
│ content     │ JSONB    │            │
│ status      │ VARCHAR  │            │
│ severity    │ VARCHAR  │            │
│ resolved    │ BOOLEAN  │ DEFAULT F  │
│ created_at  │ TIMESTAMP│            │
│ updated_at  │ TIMESTAMP│            │
│ resolved_at │ TIMESTAMP│            │
└─────────────────────────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────────────────────────┐
│              alerts                  │
├─────────────────────────────────────┤
│ id           │ UUID    │ PK         │
│ incident_id  │ UUID    │ FK         │
│ channel      │ VARCHAR │ NOT NULL   │
│ status       │ VARCHAR │            │
│ sent_at      │ TIMESTAMP│           │
│ error_message│ TEXT    │            │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│              teams                   │
├─────────────────────────────────────┤
│ id          │ UUID     │ PK         │
│ name        │ VARCHAR  │ NOT NULL   │
│ slug        │ VARCHAR  │ UNIQUE     │
│ created_at  │ TIMESTAMP│            │
│ updated_at  │ TIMESTAMP│            │
└─────────────────────────────────────┘
         │
         │ 1:N
         ▼
┌─────────────────────────────────────┐
│           team_channels              │
├─────────────────────────────────────┤
│ id          │ UUID     │ PK         │
│ team_id     │ UUID     │ FK         │
│ channel_type│ VARCHAR  │ NOT NULL   │
│ config      │ JSONB    │            │
│ is_active   │ BOOLEAN  │ DEFAULT T  │
└─────────────────────────────────────┘

┌─────────────────────────────────────┐
│           oncall_schedules           │
├─────────────────────────────────────┤
│ id          │ UUID     │ PK         │
│ team_id     │ UUID     │ FK         │
│ provider    │ VARCHAR  │ NOT NULL   │
│ schedule_id │ VARCHAR  │ NOT NULL   │
│ config      │ JSONB    │            │
│ is_active   │ BOOLEAN  │ DEFAULT T  │
└─────────────────────────────────────┘
```

## Entity Descriptions

### incidents

Core entity representing an incident event.

| Column      | Type      | Description                                    |
|-------------|-----------|------------------------------------------------|
| id          | UUID      | Unique identifier                              |
| team_id     | VARCHAR   | Reference to the team that owns this incident  |
| title       | VARCHAR   | Incident title                                 |
| description | TEXT      | Detailed description                           |
| content     | JSONB     | Flexible content for custom incident data      |
| status      | VARCHAR   | Current status (open, acknowledged, resolved)  |
| severity    | VARCHAR   | Severity level (critical, high, medium, low)   |
| resolved    | BOOLEAN   | Whether the incident has been resolved         |
| created_at  | TIMESTAMP | Record creation timestamp                      |
| updated_at  | TIMESTAMP | Last update timestamp                          |
| resolved_at | TIMESTAMP | When the incident was resolved                 |

### alerts

Tracks alert delivery attempts and status.

| Column       | Type      | Description                                    |
|--------------|-----------|------------------------------------------------|
| id           | UUID      | Unique identifier                              |
| incident_id  | UUID      | Reference to the incident                      |
| channel      | VARCHAR   | Alert channel (slack, lark, email, etc.)       |
| status       | VARCHAR   | Delivery status (pending, sent, failed)        |
| sent_at      | TIMESTAMP | When the alert was sent                        |
| error_message| TEXT      | Error message if delivery failed               |

### teams

Team management for multi-tenant support.

| Column     | Type      | Description                                    |
|------------|-----------|------------------------------------------------|
| id         | UUID      | Unique identifier                              |
| name       | VARCHAR   | Team name                                      |
| slug       | VARCHAR   | URL-friendly identifier                        |
| created_at | TIMESTAMP | Record creation timestamp                      |
| updated_at | TIMESTAMP | Last update timestamp                          |

### team_channels

Configures alert channels per team.

| Column       | Type      | Description                                    |
|--------------|-----------|------------------------------------------------|
| id           | UUID      | Unique identifier                              |
| team_id      | UUID      | Reference to the team                          |
| channel_type | VARCHAR   | Channel type (slack, lark, webhook)            |
| config       | JSONB     | Channel-specific configuration (webhooks, etc.)|
| is_active    | BOOLEAN   | Whether this channel is active                 |

### oncall_schedules

On-call schedule integrations.

| Column      | Type      | Description                                    |
|-------------|-----------|------------------------------------------------|
| id          | UUID      | Unique identifier                              |
| team_id     | UUID      | Reference to the team                          |
| provider    | VARCHAR   | Provider name (pagerduty, opsgenie, etc.)      |
| schedule_id | VARCHAR   | Provider's schedule identifier                 |
| config      | JSONB     | Provider-specific configuration                |
| is_active   | BOOLEAN   | Whether this schedule is active                |

## Relationships

| Relationship              | Type | Description                                    |
|---------------------------|------|------------------------------------------------|
| incidents → alerts        | 1:N  | An incident can have multiple alerts           |
| teams → incidents         | 1:N  | A team can have multiple incidents             |
| teams → team_channels     | 1:N  | A team can have multiple alert channels        |
| teams → oncall_schedules  | 1:N  | A team can have multiple on-call schedules     |

## Indexes

Recommended indexes for optimal query performance:

```sql
-- incidents
CREATE INDEX idx_incidents_team_id ON incidents(team_id);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_created_at ON incidents(created_at);

-- alerts
CREATE INDEX idx_alerts_incident_id ON alerts(incident_id);
CREATE INDEX idx_alerts_status ON alerts(status);

-- team_channels
CREATE INDEX idx_team_channels_team_id ON team_channels(team_id);

-- oncall_schedules
CREATE INDEX idx_oncall_schedules_team_id ON oncall_schedules(team_id);
```

## Migration Notes

When implementing migrations, consider:

1. Use GORM's AutoMigrate for initial schema creation
2. Add custom migrations for indexes
3. JSONB columns allow flexible schema evolution without migrations
4. Consider partitioning for high-volume tables (incidents, alerts)
