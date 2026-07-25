# Muse — Premium Persian Music Ecosystem

> A full-stack, immersive music streaming platform — blending a Go/Gin API backend with a Vue 3 SPA frontend, ML-powered features, and a rich social/creator economy.

---

## Table of Contents

- [Stack](#stack)
- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Manual Setup (step by step)](#manual-setup-step-by-step)
  - [1. Infrastructure](#1-infrastructure)
  - [2. Backend](#2-backend)
  - [3. Frontend](#3-frontend)
  - [4. Database Migrations](#4-database-migrations)
- [Unified Dev Script](#unified-dev-script)
- [Docker Deployment](#docker-deployment)
- [Project Structure](#project-structure)
- [Available Commands](#available-commands)
- [Architecture](#architecture)
- [Documentation](#documentation)

---

## Stack

| Layer | Technology |
|-------|-----------|
| **Frontend** | Vue 3 (Composition API, `<script setup>`), TypeScript (strict), Vite 7, Pinia 3, PrimeVue 4, Tailwind CSS v4, Vue Router 4, Axios, Zod 4 |
| **Backend** | Go 1.25, Gin 1.12, Goose (migrations), Zap (logging) |
| **Infrastructure** | Postgres 17, Redis 7.4, OpenSearch 3.2, MinIO (S3-compatible storage) |
| **ML Service** | FastAPI (Python), Whisper (lyrics transcription), Celery |
| **Auth** | JWT (access + refresh tokens) |
| **Storage** | Local filesystem or S3/MinIO |

---

## Prerequisites

| Tool | Purpose |
|------|---------|
| **Go** 1.25+ | Backend API |
| **Node.js** ^20.19 or >=22.12 | Frontend dev server |
| **Docker** & **Docker Compose** | Database, Redis, search, storage |
| **Air** | Go hot-reload (optional but recommended) |
| **Goose** | Database migrations |
| **Python** 3.10+ | ML microservice (Whisper) |
| **Make** | Build automation |

---

## Installation — Linux

### 1. System packages

```bash
sudo apt update
sudo apt install -y build-essential git curl wget unzip make
```

### 2. Go

```bash
# Download latest stable
curl -OL https://go.dev/dl/go1.25.2.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.25.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### 3. Node.js (via nvm)

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.1/install.sh | bash
source ~/.bashrc
nvm install 22
node -v && npm -v
```

### 4. Docker & Docker Compose

```bash
# Docker
curl -fsSL https://get.docker.com | sudo sh
sudo usermod -aG docker $USER
# Log out and back in for group to take effect

# Verify
docker --version
docker compose version
```

### 5. Air (Go hot-reload)

```bash
go install github.com/air-verse/air@latest
air -v
```

### 6. Goose (migrations)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -version
```

### 7. Python & ML dependencies

```bash
sudo apt install -y python3 python3-pip python3-venv
python3 -m venv ~/muse-ml-venv
source ~/muse-ml-venv/bin/activate
cd moja-ml-service
pip install -r requirements.txt
```

---

## Installation — Windows

### 1. Go

Download the installer from <https://go.dev/dl/> (`.msi` for Windows amd64). Run it. Verify:

```powershell
# Open a new terminal after install
go version
```

### 2. Node.js (via nvm-windows)

Download and install **nvm-windows** from <https://github.com/coreybutler/nvm-windows/releases>. Then:

```powershell
nvm install 22
nvm use 22
node -v && npm -v
```

Alternatively, download the LTS installer directly from <https://nodejs.org/>.

### 3. Docker Desktop

Download **Docker Desktop for Windows** from <https://www.docker.com/products/docker-desktop/>. Install and enable **WSL 2 backend** (Settings > General > Use WSL 2 based engine).

```powershell
docker --version
docker compose version
```

### 4. Air (Go hot-reload)

```powershell
go install github.com/air-verse/air@latest
air -v
```

### 5. Goose (migrations)

```powershell
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -version
```

### 6. Python & ML dependencies

Download Python 3.10+ from <https://www.python.org/downloads/> (check "Add to PATH" during install).

```powershell
python -m venv %USERPROFILE%\muse-ml-venv
%USERPROFILE%\muse-ml-venv\Scripts\activate
cd moja-ml-service
pip install -r requirements.txt
```

### 7. Make (optional)

Windows does not ship `make`. Install via [Chocolatey](https://chocolatey.org/) or [Scoop](https://scoop.sh/):

```powershell
# Chocolatey
choco install make

# Scoop
scoop install make
```

If you skip `make`, you can run the underlying commands directly (e.g. `go run ./cmd/api`, `docker compose up -d`).

---

## Quick Start

```bash
# 1. Clone & enter
git clone <repo-url> && cd music

# 2. Copy environment files
cp .env.example .env
cp frontend/.env.example frontend/.env

# 3. Start infrastructure (Docker)
make infra-up

# 4. Run database migrations
./scripts/migrate.sh up

# 5. Start the Go API (in one terminal)
make run              # with Air (hot reload)
# or
go run ./cmd/api      # without Air

# 6. Start the frontend (in another terminal)
cd frontend && npm install && npm run dev

# 7. Open in browser
# Frontend → http://localhost:3000
# API      → http://localhost:8080/api/v1/health
```

> **Or use the unified dev script** — `bash scripts/dev.sh start` (starts everything including infra, API, frontend, and ML service).

---

## Manual Setup (step by step)

### 1. Infrastructure

Start Postgres, Redis, OpenSearch, and MinIO via Docker:

```bash
make infra-up
```

Verify they are healthy:

```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
```

Expected: `musicapp_postgres`, `musicapp_redis`, `musicapp_opensearch`, `musicapp_minio` all in `healthy` or `Up` state.

To stop infrastructure:

```bash
make infra-down
```

### 2. Backend

```bash
# Copy environment (edit if needed — defaults work for local dev)
cp .env.example .env

# Run with Air (hot reload on file changes)
make run

# Or run directly
go run ./cmd/api
```

The API starts on **`http://localhost:8080`**. Health check: `http://localhost:8080/api/v1/health`.

Environment variables are documented in [`.env.example`](.env.example):

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_PORT` | `8080` | API server port |
| `POSTGRES_URL` | `postgres://postgres:postgres@localhost:5432/musicapp?sslmode=disable` | Postgres connection |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `JWT_ACCESS_SECRET` | — | JWT signing key (change in production!) |
| `STORAGE_DRIVER` | `local` | `local` or `s3` |

### 3. Frontend

```bash
cd frontend

# Copy environment
cp .env.example .env

# Install dependencies
npm install

# Start dev server (with HMR)
npm run dev            # → http://localhost:3000

# Or expose to network
npm run dev:expose     # → http://0.0.0.0:3000
```

The Vite dev server proxies `/api/v1` requests to `localhost:8080` (configured in [`vite.config.ts`](frontend/vite.config.ts)).

**Available npm scripts:**

| Command | Description |
|---------|-------------|
| `npm run dev` | Vite dev server with HMR (`:3000`) |
| `npm run dev:expose` | Dev server exposed on network |
| `npm run build` | Type-check + production build |
| `npm run build-only` | Production build only |
| `npm run type-check` | vue-tsc type checking |
| `npm run preview` | Preview production build |
| `npm run lint` | ESLint with auto-fix |
| `npm run format` | Prettier formatting |

### 4. Database Migrations

Migrations use [Goose](https://github.com/pressly/goose). SQL files are in [`migrations/`](migrations/).

```bash
# Run all pending migrations
./scripts/migrate.sh up

# Rollback last migration
./scripts/migrate.sh down

# Check migration status
./scripts/migrate.sh status
```

The `POSTGRES_URL` environment variable must be set (defaults to `postgres://postgres:postgres@localhost:5432/musicapp?sslmode=disable`).

---

## Unified Dev Script

[`scripts/dev.sh`](scripts/dev.sh) manages all services in parallel with a single command:

```bash
bash scripts/dev.sh start      # Start everything
bash scripts/dev.sh stop       # Stop everything gracefully
bash scripts/dev.sh restart    # Restart everything
bash scripts/dev.sh status     # Show health of all services
bash scripts/dev.sh logs       # Tail all logs (follow)
bash scripts/dev.sh logs api   # Tail specific service log
```

Services managed:

| Service | Technology | Port | Description |
|---------|-----------|------|-------------|
| `infra` | Docker Compose | — | Postgres, Redis, MinIO, OpenSearch |
| `api` | Go/Gin | `:8080` | REST API |
| `web` | Vite | `:3000` | Frontend dev server |
| `ml` | FastAPI | `:8000` | ML service (Whisper lyrics) |
| `worker` | Celery | — | Async ML inference |

Logs are written to `logs/{api,web,ml-api,worker}.log`.

---

## Docker Deployment

Full production stack via Docker Compose:

```bash
# Build and start all services
make docker-up

# View logs
make docker-logs

# Rebuild specific services
make docker-rebuild

# Stop everything
make docker-down
```

The production compose file is at [`deployments/docker-compose.yml`](deployments/docker-compose.yml). It builds:

- **`musicapp_api`** — Go API from [`Dockerfile`](Dockerfile) (multi-stage, ~15 MB)
- **`musicapp_frontend`** — Nginx-served SPA from [`frontend/Dockerfile`](frontend/Dockerfile)
- **Infra containers** — Postgres 17, Redis 7, OpenSearch 3, MinIO

**Important:** Set secure `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` in production (passed via environment or `.env`).

---

## Project Structure

```
music/
├── cmd/
│   ├── api/              # Go API entrypoint
│   └── worker/           # Background worker entrypoint
├── internal/
│   └── app/              # Go application core (handlers, services, models)
├── frontend/
│   ├── src/
│   │   ├── assets/       # CSS, fonts, images
│   │   ├── components/   # Vue components (admin/, auth/, common/, music/, social/)
│   │   ├── composables/  # Shared logic (useAuth, usePlayer, useCatalogSearch, etc.)
│   │   ├── layouts/      # Layout shells (LayoutMusicApp, LayoutAuth, etc.)
│   │   ├── pages/        # Route pages (app/, admin/, auth/, errors/)
│   │   ├── plugins/      # Axios client, request factory
│   │   ├── router/       # Vue Router configuration
│   │   ├── services/     # API modules, audio engine, WebSocket, storage
│   │   ├── stores/       # Pinia stores (player, user-auth, feature-flags)
│   │   ├── types/        # TypeScript definitions
│   │   └── utils/        # PrimeVue preset, utility functions
│   ├── vite.config.ts    # Vite configuration
│   └── package.json
├── migrations/           # SQL migrations (Goose)
├── scripts/
│   ├── dev.sh            # Unified dev server manager
│   ├── migrate.sh        # Database migration helper
│   └── watchdog.sh       # Server health monitor
├── deployments/          # Docker Compose + systemd units
│   └── docker-compose.yml
├── moja-ml-service/      # ML service (FastAPI + Whisper)
├── cmd/                  # Go entrypoints
├── internal/             # Go internal packages
├── Makefile             # Top-level commands
├── Dockerfile           # API Docker image
└── Dockerfile.worker    # Worker Docker image
```

---

## Available Commands

### Makefile (`make <target>`)

| Target | Description |
|--------|-------------|
| `run` | Start Go API with Air (hot reload) |
| `run-worker` | Start worker with Air (hot reload) |
| `build` | Build Go API binary → `bin/api` |
| `build-worker` | Build worker binary → `bin/worker` |
| `build-all` | Build both API + worker |
| `test` | Run Go tests |
| `fmt` | Format Go code |
| `infra-up` | Start Docker infra (Postgres, Redis, OpenSearch, MinIO) |
| `infra-down` | Stop Docker infra |
| `infra-restart` | Restart Docker infra |
| `infra-logs` | Follow Docker logs |
| `docker-up` | Full production stack via Docker Compose |
| `docker-down` | Stop production stack |
| `docker-rebuild` | Rebuild API + frontend images |
| `migrate-up` | Run DB migrations |
| `migrate-down` | Rollback last migration |
| `migrate-status` | Check migration status |
| `dev` | Start infra + Go API |
| `dev-worker` | Start infra + Go worker |

---

## Architecture

### Frontend Agent Architecture

The frontend follows a parallel agent architecture (see [`frontend/AGENTS.md`](frontend/AGENTS.md)):

| Agent | Domain | Key Files |
|-------|--------|-----------|
| `@auth` | Login, JWT, guards | `stores/user-auth.ts`, `composables/auth/*` |
| `@player` | Audio engine, queue, radio | `stores/player.ts`, `composables/player/*` |
| `@catalog` | Search, tracks, albums, artists | `pages/app/Page*.vue`, `services/api/catalog/*` |
| `@social` | Rooms, clubs, notifications | `components/social/*`, `services/api/social/*` |
| `@creator` | Creator dash, badges, XP | `components/creator/*`, `services/api/creator/*` |
| `@admin` | Admin panel, moderation | `pages/admin/*`, `components/admin/*` |
| `@ui` | Layout, common components, styling | `layouts/*`, `components/common/*` |
| `@infra` | API barrel, router, stores | `services/api/index.ts`, `router/index.ts` |

### Backend

Go/Gin REST API serving under `/api/v1`, with JWT authentication, Postgres persistence, Redis caching, and OpenSearch for full-text search. File storage supports local filesystem and S3/MinIO.

### ML Service

FastAPI-based service in `moja-ml-service/` providing:
- Whisper-based Persian lyrics transcription
- Audio feature extraction
- Celery-backed async processing

---

## Documentation

| Document | Description |
|----------|-------------|
| [`frontend/AGENTS.md`](frontend/AGENTS.md) | Agent team architecture & file ownership |
| [`frontend/PERFORMANCE.md`](frontend/PERFORMANCE.md) | Performance-sensitive code guidelines |
| [`FRONTEND_ARCHITECTURE.md`](FRONTEND_ARCHITECTURE.md) | Product plan & frontend architecture |
| [`TODO.md`](TODO.md) | Upcoming features & known issues |
| [docs/](docs/) | Additional project documentation |
| [`.env.example`](.env.example) | Environment variable reference |
| [`frontend/.env.example`](frontend/.env.example) | Frontend environment variable reference |
| [`migrations/`](migrations/) | Database migration SQL files |
