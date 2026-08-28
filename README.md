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
- **Full Observability**: Structured JSON logging (zerolog + Request ID tracking), Prometheus metrics (`/metrics`), and pre-provisioned Grafana dashboards.

---

## 🚀 Quick Start — One-Command Setup

### Option A: Docker / Podman Compose (Recommended)

Menjalankan seluruh 4 container sekaligus (**Kickbase API**, **PostgreSQL 16**, **Prometheus**, **Grafana**):

```bash
# 1. Clone repository
git clone <repository-url>
cd kickbase

# 2. Start full stack
docker compose up --build -d
# Atau jika menggunakan Podman:
podman compose up --build -d
```

> 💡 **Khusus Pengguna Podman di Windows**: Jika mengakses dari browser Windows, jalankan helper berikut untuk mensinkronkan port container ke `localhost`:
> ```powershell
> powershell -ExecutionPolicy Bypass -File ./scripts/sync-portproxy.ps1
> ```

---

### Option B: Local Setup (Native Go + PostgreSQL)

```bash
# 1. Clone repository
git clone <repository-url>
cd kickbase

# 2. Setup Environment
cp .env.example .env

# 3. Create PostgreSQL Database
psql -U postgres -c "CREATE DATABASE kickbase;"

# 4. Install Dependencies & Generate Swagger Docs
go mod tidy
swag init -g cmd/server/main.go -o docs

# 5. Run Server
go run cmd/server/main.go
```

---

## 🌐 Service URLs & Access Summary

| Service | Local URL | Default Credentials | Description |
|---|---|---|---|
| **API Server** | [http://localhost:8080](http://localhost:8080) | - | Kickbase REST API |
| **Interactive Swagger UI** | [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) | Bearer Token (See below) | OpenAPI 2.0 Web Interface |
| **API Health Check** | [http://localhost:8080/api/health](http://localhost:8080/api/health) | - | Live DB & Service health |
| **Prometheus Live Metrics** | [http://localhost:8080/metrics](http://localhost:8080/metrics) | - | Prometheus exposition endpoint |
| **Prometheus Server** | [http://localhost:9090](http://localhost:9090) | - | Target scraping & metrics query |
| **Grafana Dashboard** | [http://localhost:3000](http://localhost:3000) | `admin` / `admin` | Pre-configured metrics visualizer |

---

## 🔑 Demo Credentials & RBAC Permission Matrix

Kamu dapat langsung mencoba kredensial demo berikut atau mendaftarkan akun baru melalui endpoint `POST /api/auth/register`:

| Role | Demo Username | Demo Password | Granted Permissions |
|---|---|---|---|
| **Admin** | `admin_demo` | `password123` | **All 24 Permissions (`*`)**: Full Create, Read, Update, Delete, Revert on all domains |
| **Staff** | `staff_demo` | `password123` | **Create, Read, Update**: Tidak memiliki permission `delete` & `revert` |
| **Viewer** | `viewer_demo` | `password123` | **Read-Only (`*:read`)**: Hanya dapat membaca data publik |

### Cara Menggunakan Token di Swagger UI:
1. Buka [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).
2. Buka endpoint **`POST /api/auth/login`**, klik **Try it out**, dan masukkan payload:
   ```json
   {
     "username": "admin_demo",
     "password": "password123"
   }
   ```
   *(Jika user belum terdaftar, panggil `POST /api/auth/register` terlebih dahulu dengan role `admin`)*.
3. Salin nilai `access_token` dari respons JSON.
4. Klik tombol hijau **Authorize 🔓** di kanan atas halaman Swagger UI.
5. Masukkan format: `Bearer <access_token>` (atau cukup tokennya saja) lalu klik **Authorize**.
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

## 🏗️ Architecture & Code Organization

```
cmd/
└── server/
    └── main.go                # Application bootstrap & Swagger annotations (<65 lines)

internal/
├── config/                    # Environment variable loader
├── database/                  # PostgreSQL connection & auto-migration
├── domain/                    # Pure domain models (User, Team, Player, Match, Permission, etc.)
├── interfaces/                # Service & Repository interfaces (Dependency Inversion)
├── repository/                # GORM repository implementations with dynamic queries
├── service/                   # Core business logic, validation, & audit trail
├── handler/                   # HTTP handlers per domain (<250 lines/file) + DTO definitions
├── middleware/                # CORS, Security Headers, JWT Auth, RBAC Guard, Prometheus Metrics
└── router/                    # Route wiring & dependency injection

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

## 📊 Observability & Metrics

- **Structured Logging**: Seluruh HTTP request di-log dalam format JSON terstruktur menggunakan `zerolog` menyertakan `request_id`, `latency_ms`, `status_code`, dan `client_ip`.
- **Prometheus Metrics**:
  - Counter: `http_requests_total{method, path, status}`
  - Histogram: `http_request_duration_seconds{method, path}`
  - Dynamic Gauges: `kickbase_teams_total`, `kickbase_players_total`, `kickbase_matches_total`, `kickbase_matches_completed_total` (terkoneksi real-time ke DB).

---

## 📄 License

Project ini dilisensikan di bawah [MIT License](LICENSE).
