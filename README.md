# Kickbase - Football Team Management API

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/widiskel/kickbase/actions/workflows/ci.yml/badge.svg)](https://github.com/widiskel/kickbase/actions/workflows/ci.yml)

A REST API backend for managing amateur football teams, players, match schedules, match results, and match reports for Perusahaan XYZ.

## Features

- **Team Management**: CRUD operations for football teams with soft delete
- **Player Management**: Manage players per team with jersey number uniqueness
- **Match Scheduling**: Schedule matches between teams with status lifecycle
- **Match Results**: Report match results with goal scorers and times
- **Match Reports**: View match reports with cumulative wins and top scorers
- **Audit Trail**: Track all entity changes with history tables and revert capability
- **Observability**: Structured logging, Prometheus metrics, Grafana dashboards

## Prerequisites

- **Go 1.22+** installed ([download](https://go.dev/dl/))
- **PostgreSQL 14+** running locally ([download](https://www.postgresql.org/download/))
- **swag CLI** installed (for Swagger docs):
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```

## Quick Start — Docker (Recommended)

```bash
git clone <repository-url>
cd kickbase
docker compose up
```

### Access

| Service | URL | Credentials |
|---------|-----|-------------|
| API | http://localhost:8080 | - |
| Swagger UI | http://localhost:8080/swagger/index.html | - |
| Prometheus | http://localhost:9090 | - |
| Grafana | http://localhost:3000 | admin / admin |
| Health Check | http://localhost:8080/api/health | - |

## Quick Start — Manual Setup

### 1. Clone and configure

```bash
git clone <repository-url>
cd kickbase
cp .env.example .env
```

### 2. Create database

```bash
psql -U postgres -c "CREATE DATABASE kickbase;"
```

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Generate Swagger docs

```bash
swag init -g cmd/server/main.go
```

### 5. Run the server

```bash
go run cmd/server/main.go
```

### Access

| Service | URL |
|---------|-----|
| API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/swagger/index.html |
| Health Check | http://localhost:8080/api/health |

## API Endpoints

### Teams

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/teams | Create a new team |
| GET | /api/teams | List all teams (paginated) |
| GET | /api/teams/:id | Get team by ID |
| PUT | /api/teams/:id | Update team |
| DELETE | /api/teams/:id | Soft delete team |
| GET | /api/teams/:id/history | Get team history |
| POST | /api/teams/:id/revert | Revert team to version |

### Players

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/teams/:teamId/players | Add player to team |
| GET | /api/teams/:teamId/players | List players by team |
| GET | /api/players/:id | Get player by ID |
| PUT | /api/players/:id | Update player |
| DELETE | /api/players/:id | Soft delete player |
| GET | /api/players/:id/history | Get player history |
| POST | /api/players/:id/revert | Revert player to version |

### Matches

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/matches | Schedule a match |
| GET | /api/matches | List all matches |
| GET | /api/matches/:id | Get match by ID |
| PATCH | /api/matches/:id/status | Update match status |
| GET | /api/matches/:id/history | Get match history |
| POST | /api/matches/:id/revert | Revert match to version |
| POST | /api/matches/:matchId/result | Report match result |
| GET | /api/matches/:matchId/result | Get match result |

### Reports

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/reports/matches | List all match reports |
| GET | /api/reports/matches/:id | Get single match report |

### Health

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/health | Health check |

## Architecture

```
cmd/server/main.go          # Application entry point

internal/
├── config/                  # Environment configuration
├── domain/                  # Domain entities (Team, Player, Match, etc.)
├── handler/                 # HTTP handlers with Swagger annotations
├── service/                 # Business logic
├── repository/              # Database operations
├── middleware/               # CORS, security, logging, metrics
├── router/                  # Route registration
└── database/                # DB connection and migration

deploy/                      # Deployment configs (Prometheus, Grafana)
docs/                        # Generated Swagger documentation
```

### Clean Architecture

```
Handler → Service → Repository → Database
```

- **Handler**: Receives HTTP requests, validates input, calls service
- **Service**: Contains business rules and orchestration
- **Repository**: Handles database operations
- **Domain**: Pure structs with no framework dependency

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DB_HOST | Yes | localhost | PostgreSQL host |
| DB_PORT | Yes | 5432 | PostgreSQL port |
| DB_USER | Yes | postgres | PostgreSQL user |
| DB_PASSWORD | Yes | postgres | PostgreSQL password |
| DB_NAME | Yes | kickbase | Database name |
| DB_SSLMODE | Yes | disable | SSL mode |
| SERVER_PORT | Yes | 8080 | Server port |
| LOG_LEVEL | No | info | Log level (debug, info, warn, error) |

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific layer tests
go test ./internal/service/...
go test ./internal/repository/...
go test ./internal/handler/...
```

### Coverage Requirements

- Service layer: ≥80%
- Repository layer: ≥70%
- Handler layer: ≥60%
- Overall: ≥70%

## Observability

### Structured Logging

All requests are logged with structured JSON including:
- Request ID (X-Request-ID header)
- Method, path, status code
- Latency, client IP

### Prometheus Metrics

Access at http://localhost:9090

Metrics exposed:
- `http_requests_total` - Request counts
- `http_request_duration_seconds` - Latency histograms
- `teams_total`, `players_total`, `matches_total` - Business metrics

### Grafana Dashboards

Access at http://localhost:3000 (admin/admin)

Pre-configured datasource: Prometheus

## License

MIT
