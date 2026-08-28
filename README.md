# Kickbase - Football Team Management REST API ⚽

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![CI](https://github.com/widiskel/kickbase/actions/workflows/ci.yml/badge.svg)](https://github.com/widiskel/kickbase/actions/workflows/ci.yml)
[![Swagger Docs](https://img.shields.io/badge/Swagger-OpenAPI%202.0-85EA2D?logo=swagger)](http://localhost:8080/swagger/index.html)

A production-grade RESTful API backend for managing amateur football teams, players, match schedules, match results, and match reports for **Perusahaan XYZ**. Built with Go 1.22+, Gin, PostgreSQL 16, and GORM following **Strict Clean Architecture & SOLID Principles**.

---

## 🌟 Core Highlights & Engineering Rigor

- **Clean Architecture & SOLID**: Strict layer separation (`Handler` ➔ `Service Interface` ➔ `Repository Interface` ➔ `Database`). Dependency Inversion via interfaces with 100% isolated mock-based unit tests.
- **Go Style Standards**: Strict file limit `<500 LOC` per file, early returns, non-nil slice/map initialization, and clean DTO isolation.
- **Authentication & Security**: JWT Dual-Token lifecycle (1h Access Token + 7d Rotating Refresh Token with anti-replay revocation) and bcrypt password hashing.
- **Granular RBAC Domain Permissions**: 24 distinct permissions (`<domain>:<action>`: Create, Read, Update, Delete, Revert) across 6 domains (`teams`, `players`, `matches`, `results`, `reports`, `users`).
- **Pagination, Filtering & Sorting**: Dynamic query capabilities on all list endpoints (`page`, `limit`, `sort_by`, `order`, search, and domain filters).
- **Audit Trail & Revert**: Dedicated `_history` tables tracking JSONB changes and version increments with one-click snapshot rollback.
- **Full Observability**: Structured JSON logging (zerolog + Request ID tracing), Prometheus metrics (`/metrics`), and pre-provisioned Grafana dashboards.

---

## 💻 System Prerequisites by OS

| Operating System | Recommended Tools | Minimum Requirements |
|---|---|---|
| **Windows 10/11** | Docker Desktop / Podman Desktop (WSL2) | Go 1.22+, WSL2, Git |
| **macOS (Intel/Apple Silicon)** | Docker Desktop / Colima / Podman (`brew install podman podman-compose`) | Go 1.22+ (`brew install go`), Git |
| **Linux (Ubuntu/Debian/RHEL)** | Docker CE (`docker-compose-plugin`) / Podman (`podman-compose`) | Go 1.22+, PostgreSQL 16, Git |

---

## 🚀 Quick Start — Containerized Setup (Recommended)

Satu perintah untuk menjalankan seluruh 4 container sekaligus (**Kickbase API**, **PostgreSQL 16**, **Prometheus**, **Grafana**):

### 1. Menggunakan Docker Compose (Linux / macOS / Windows)

```bash
# Clone repository
git clone <repository-url>
cd kickbase

# Build and start all 4 services in background
docker compose up --build -d

# Stop all services (dan hapus volume data jika ingin reset bersih)
docker compose down -v
```

---

### 2. Menggunakan Podman Compose (Linux / macOS / Windows)

```bash
# Clone repository
git clone <repository-url>
cd kickbase

# Build and start all 4 services in background
podman compose up --build -d

# Khusus Pengguna Podman di Windows:
# Jika port container belum otomatis ter-route ke localhost Windows, jalankan helper ini:
powershell -ExecutionPolicy Bypass -File ./scripts/sync-portproxy.ps1

# Stop all services
podman compose down -v
```

---

### 3. Menggunakan Colima (Khusus macOS)

```bash
# Start colima engine
colima start --cpu 2 --memory 4

# Run compose
docker compose up --build -d
```

---

## 🛠️ Quick Start — Native Local Development (Non-Container)

Jika reviewer ingin menjalankan backend secara langsung di mesin lokal tanpa Docker/Podman:

```bash
# 1. Clone repository
git clone <repository-url>
cd kickbase

# 2. Setup Environment Variable
cp .env.example .env

# 3. Buat Database PostgreSQL lokal
# (Pastikan PostgreSQL 14+ aktif di port 5432)
psql -U postgres -c "CREATE DATABASE kickbase;"

# 4. (Opsional) Jalankan Migrasi DDL SQL manual
psql -U postgres -d kickbase -f migrations/000001_init_schema.up.sql

# 5. Install Dependencies & Generate Swagger Documentation
go mod tidy
swag init -g cmd/server/main.go -o docs

# 6. Jalankan Server API
# (Server otomatis melakukan auto-migration & auto-seeding data 44 pemain)
go run cmd/server/main.go
```

---

## 🌐 Service URLs & Access Summary

| Service | Local URL | Default Credentials | Description |
|---|---|---|---|
| **API Server** | [http://localhost:8080](http://localhost:8080) | - | Kickbase REST API Server |
| **Interactive Swagger UI** | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) | Bearer Token (See below) | OpenAPI 2.0 Interactive Testing |
| **OpenAPI Spec (JSON)** | [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json) | - | Raw OpenAPI 2.0 JSON Specification |
| **API Health Check** | [http://localhost:8080/api/health](http://localhost:8080/api/health) | - | Live DB & App Connectivity Health |
| **Prometheus Metrics** | [http://localhost:8080/metrics](http://localhost:8080/metrics) | - | Prometheus Raw Text Exposition Endpoint |
| **Prometheus Server UI** | [http://localhost:9090](http://localhost:9090) | - | PromQL Query Console & Metric Graphing |
| **Prometheus Pre-Loaded Graphs** | [http://localhost:9090/graph?g0.expr=kickbase%3Ahttp_requests%3Arate1m&g0.tab=0&g1.expr=kickbase%3Ahttp_errors%3Aratio_pct&g1.tab=0&g2.expr=kickbase%3Alatency%3Ap95_seconds&g2.tab=0&g3.expr=kickbase%3Amemory%3Aheap_mb&g3.tab=0](http://localhost:9090/graph?g0.expr=kickbase%3Ahttp_requests%3Arate1m&g0.tab=0&g1.expr=kickbase%3Ahttp_errors%3Aratio_pct&g1.tab=0&g2.expr=kickbase%3Alatency%3Ap95_seconds&g2.tab=0&g3.expr=kickbase%3Amemory%3Aheap_mb&g3.tab=0) | - | 4 Pre-loaded Charts (Rate, Errors, Latency, RAM) |
| **Prometheus Pre-Configured Rules** | [http://localhost:9090/rules](http://localhost:9090/rules) | - | 11 Live Evaluated Prometheus Recording Rules |
| **Prometheus Targets Status** | [http://localhost:9090/targets](http://localhost:9090/targets) | - | Live Target Scraping Health (`app:8080`) |
| **Grafana Dashboard** | [http://localhost:3000](http://localhost:3000) | `admin` / `admin` | Pre-configured Executive Dashboard |
| **PostgreSQL Database** | `localhost:5432` | `postgres` / `postgres` | Database `kickbase` (Auto-migrated & Seeded) |

---

## 🔑 Demo Credentials & RBAC Permission Matrix

Database seeder otomatis mengisi akun bawaan berikut saat pertama kali dijalankan:

| Role | Demo Username | Demo Password | Granted Permissions |
|---|---|---|---|
| **Admin** | `admin` | `password123` | **All 24 Permissions (`*`)**: Full Create, Read, Update, Delete, Revert on all domains |
| **Staff** | `staff` | `password123` | **Create, Read, Update**: Tidak memiliki permission `delete` & `revert` |
| **Viewer** | `viewer` | `password123` | **Read-Only (`*:read`)**: Hanya dapat membaca data publik |

### Cara Menggunakan Token di Swagger UI:
1. Buka [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).
2. Buka endpoint **`POST /api/auth/login`**, klik **Try it out**, dan masukkan payload:
   ```json
   {
     "username": "admin",
     "password": "password123"
   }
   ```
3. Salin nilai `access_token` dari respons JSON.
4. Klik tombol hijau **Authorize 🔓** di kanan atas halaman Swagger UI.
5. Masukkan format: `Bearer <access_token>` (atau tokennya saja) lalu klik **Authorize**.
6. Seluruh endpoint data mutasi (POST/PUT/DELETE/PATCH) sekarang siap diuji secara interaktif!

---

## 📡 REST API Specifications

### 1. Authentication (`/api/auth`)
- `POST /api/auth/register` — Registrasi user baru (Admin/Staff/Viewer).
- `POST /api/auth/login` — Login user (Mengembalikan Access Token, Refresh Token, dan 24 Permissions).
- `POST /api/auth/refresh` — Rotasi refresh token dan penerbitan token baru (Revocation protection).

### 2. Team Management (`/api/teams`)
- `GET /api/teams` — List tim dengan pagination, filter (`name`, `city`), dan sorting (`sort_by=founded_year&order=desc`). *(Public)*
- `GET /api/teams/:id` — Detail tim. *(Public)*
- `GET /api/teams/:id/history` — Riwayat perubahan data tim (Audit Trail). *(Public)*
- `POST /api/teams` — Tambah tim baru. *(Protected: `teams:create`)*
- `PUT /api/teams/:id` — Update tim dengan optimistic locking (`version`). *(Protected: `teams:update`)*
- `DELETE /api/teams/:id` — Soft delete tim (Dilarang jika tim masih memiliki pemain). *(Protected: `teams:delete`)*
- `POST /api/teams/:id/revert` — Rollback ke versi sebelumnya. *(Protected: `teams:revert`)*

### 3. Player Management (`/api/players`)
- `GET /api/players` — List pemain dengan filter (`team_id`, `position`, `name`) & sorting (`sort_by=jersey_number&order=asc`). *(Public)*
- `GET /api/players/:id` — Detail pemain. *(Public)*
- `GET /api/players/:id/history` — Audit trail riwayat pemain. *(Public)*
- `POST /api/players` — Tambah pemain baru. *(Protected: `players:create`)*
  - **Validasi eFootball 15 Posisi**: `CF, SS, LWF, RWF, AMF, CMF, DMF, LMF, RMF, CB, LB, RB, GK`.
  - **22 Playstyle Opsional**: `Goal Poacher, Prolific Winger, Box-to-Box, Build Up`, dll.
  - **Unique Jersey Number**: Nomor punggung 1-99 unik per tim.
  - **Rentang Fisik**: Tinggi 150-220 cm, Berat 40-150 kg.
- `PUT /api/players/:id` — Update pemain. *(Protected: `players:update`)*
- `DELETE /api/players/:id` — Soft delete pemain (Dilarang jika pemain memiliki riwayat gol). *(Protected: `players:delete`)*
- `POST /api/players/:id/revert` — Rollback ke versi sebelumnya. *(Protected: `players:revert`)*

### 4. Match Scheduling (`/api/matches`)
- `GET /api/matches` — List pertandingan dengan filter status (`scheduled, completed, cancelled, deferred`), tim, dan rentang tanggal. *(Public)*
- `GET /api/matches/:id` — Detail jadwal pertandingan. *(Public)*
- `GET /api/matches/:id/history` & `POST /api/matches/:id/revert`. *(Public/Protected)*
- `POST /api/matches` — Jadwalkan pertandingan (`home_team_id != away_team_id`). *(Protected: `matches:create`)*
- `PATCH /api/matches/:id/status` — Update status siklus pertandingan. *(Protected: `matches:update`)*

### 5. Match Results & Goals (`/api/results`)
- `POST /api/results` — Input hasil skor dan pencetak gol. Status pertandingan otomatis berubah menjadi `completed`. *(Protected: `results:create`)*
- `GET /api/results/:matchId` — Detail hasil pertandingan dan daftar gol. *(Public)*

### 6. Match Reports (`/api/reports`)
- `GET /api/reports/matches` — List seluruh ringkasan laporan pertandingan. *(Public)*
- `GET /api/reports/matches/:id` — Laporan pertandingan publik:
  - Skor akhir & Status pemenang (`Tim Home Menang`, `Tim Away Menang`, `Draw`).
  - Daftar **Top Scorer** pertandingan beserta jumlah gol yang dicetak.
  - Total **Akumulasi Kemenangan** tim Home dan tim Away sepanjang masa.

---

## 📊 Observability, Logging & Tracing

Sistem Kickbase mengimplementasikan standar observabilitas enterprise tiga pilar (Logs, Metrics, Traces):

### 1. Structured JSON Logging & Request Tracing
- **Library**: `github.com/rs/zerolog`
- **Request Tracing**: Setiap request yang masuk secara otomatis diberi UUID unik melalui middleware dan diteruskan di header `X-Request-ID`.
- **Level Differentiation**: `INFO` untuk 2xx/3xx, `WARN` untuk 4xx (client errors), dan `ERROR` untuk 5xx (server faults).
- **Log Payload Format**:
  ```json
  {
    "level": "info",
    "request_id": "e8c5a652-4ca8-4609-87e9-9b2eb77c5140",
    "method": "POST",
    "path": "/api/auth/login",
    "status": 200,
    "latency": 49.16,
    "client_ip": "192.0.2.1",
    "time": "2026-08-28T10:32:53+07:00",
    "message": "request processed"
  }
  ```

### 2. Prometheus Time-Series Metrics (`/metrics`)
- **Library**: `github.com/prometheus/client_golang`
- **Application & HTTP Metrics**:
  - `http_requests_total{method, path, status}` (Counter) — Menghitung total incoming request per endpoint dan status code.
  - `http_request_duration_seconds{method, path}` (Histogram) — Mengukur distribusi latensi request (p50, p90, p99).
- **Live Dynamic Business Gauges** (Terkoneksi real-time ke Database):
  - `teams_total` — Jumlah tim aktif saat ini.
  - `players_total` — Jumlah pemain aktif di seluruh klub (44 pemain).
  - `matches_total` — Total seluruh jadwal pertandingan.
  - `matches_completed_total` — Total pertandingan yang telah selesai dengan laporan hasil.

### 3. Grafana Pre-Configured Dashboards
- **URL**: [http://localhost:3000](http://localhost:3000) (Login default: `admin` / `admin`, atau Anonymous Admin enabled).
- **Dashboard**: `Kickbase API - Executive Production Dashboard` (Pre-provisioned otomatis):
  - 🟢 **Row 1**: Application Health (`HEALTHY (UP)`), Service Uptime, RAM/Heap Usage (MiB), Active Goroutines, CPU Core Utilization Rate.
  - 🚨 **Row 2**: Error Tracking Breakdown (Total 4xx, Total 5xx, Overall Failure Rate %, Top Failing Endpoints Trace Table).
  - ⚽ **Row 3**: Football Domain Metrics (Teams, Players, Matches, Completed).
  - 🚀 **Row 4**: HTTP Traffic & Latency percentiles.
  - 📊 **Row 5**: Runtime Resource Trends (Heap allocations over time, Goroutines growth).

---

## 🏗️ Architecture & Code Organization

```
cmd/
└── server/
    └── main.go                # Application bootstrap & Swagger annotations (<65 lines)

internal/
├── config/                    # Environment variable loader
├── database/                  # PostgreSQL connection, auto-migration & 44-player seeder
├── domain/                    # Pure domain models (User, Team, Player, Match, Permission, etc.)
├── interfaces/                # Service & Repository interfaces (Dependency Inversion)
├── repository/                # GORM repository implementations with dynamic queries
├── service/                   # Core business logic, validation, & audit trail
├── handler/                   # HTTP handlers per domain (<250 lines/file) + DTO definitions
├── middleware/                # CORS, Security Headers, JWT Auth, RBAC Guard, Prometheus Metrics
└── router/                    # Route wiring & dependency injection

migrations/                    # Explicit SQL DDL Migration and Seed Scripts
├── 000001_init_schema.up.sql
├── 000001_init_schema.down.sql
└── 000002_seed_initial_data.up.sql

deploy/
├── prometheus.yml             # Prometheus scrape configuration
└── grafana/                   # Datasources & pre-provisioned dashboards

test/
├── mocks/                     # Repository mock structs for isolated testing
├── unit/service/              # 100% Mock-based unit test suites (42 test cases)
└── integration/               # Full database integration tests (6 test suites)
```

---

## 🧪 Automated Testing

### Menjalankan Unit & Integration Tests:

```bash
# Jalankan seluruh test suite
go test ./... -v

# Jalankan unit test service terisolasi (Mock-based, ~80% coverage)
go test ./test/unit/service/... -v

# Jalankan integration test (End-to-End dengan PostgreSQL)
go test ./test/integration/... -v
```

---

## 📄 License

Project ini dilisensikan di bawah [MIT License](LICENSE).
