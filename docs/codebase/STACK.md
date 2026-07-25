# Muse — Technology Stack

## Frontend

| Layer | Technology | Version | Notes |
|-------|-----------|---------|-------|
| **Framework** | Vue 3 (Composition API) | ^3.5.28 | `<script setup lang="ts">` everywhere |
| **Language** | TypeScript | ~5.9.0 | Strict mode |
| **Build** | Vite | ^7.3.1 | Dev on :3000, proxy `/api/v1` → :8080 |
| **State** | Pinia | ^3.0.4 | Setup function pattern, 5 stores |
| **Router** | Vue Router | ^4.6.3 | `createWebHistory`, lazy routes |
| **UI Components** | PrimeVue | ^4.5.4 | Auto-imported, custom dark preset |
| **CSS** | Tailwind CSS | ^4.2.0 | Via `@tailwindcss/vite` plugin |
| **HTTP** | Axios | ^1.13.5 | Custom wrapper with interceptor chain |
| **Icons** | Lucide Vue Next | ^0.552.0 | |
| **Validation** | Zod | ^4.3.6 | Runtime API response validation |
| **Testing** | Vitest | ^4.1.9 | + Vue Test Utils |
| **Lint** | ESLint | ^10.0 | + Prettier |
| **Drag** | VueDraggable | ^4.1.0 | Sortable lists |

## Backend

| Layer | Technology | Version | Notes |
|-------|-----------|---------|-------|
| **Language** | Go | 1.25.9 | |
| **Framework** | Gin | v1.12.0 | HTTP routing + middleware |
| **DB Driver** | pgx + sqlx | v5.10.0 / v1.4.0 | Connection pool + query helpers |
| **Migrations** | Goose | v3.27.1 | SQL files in `migrations/` |
| **Auth** | golang-jwt | v5.3.1 | Access + refresh tokens |
| **Validation** | go-playground/validator | v10.30.2 | |
| **Logging** | zap | v1.28.0 | Structured, global singleton |
| **Cache** | go-redis | v9.19.0 | |
| **Search** | OpenSearch Go client | v1.1.0 | |
| **Events** | In-process bus | Custom | + transactional outbox (pg) |
| **Metrics** | Prometheus | v1.23.2 | `/metrics` endpoint |
| **Docs** | Swaggo (Swagger) | v1.16.6 | `/swagger/*` |
| **Storage** | Local / S3 (MinIO) | Custom | Switchable via config |
| **WebSocket** | gorilla/websocket | v1.5.3 | Real-time features |
| **Testing** | testify | v1.11.1 | + go-sqlmock |
| **Media** | dhowden/tag | v0.0.0... | Audio metadata parsing |

## Infrastructure

| Service | Image/Version | Purpose |
|---------|--------------|---------|
| **PostgreSQL** | postgres:17 | Primary database |
| **Redis** | redis:8.8.0 | Cache + session + pub/sub |
| **OpenSearch** | opensearch:3.2.0 | Full-text search |
| **MinIO** | minio/minio | S3-compatible object storage |
| **Docker** | Compose | Local dev orchestration |

## ML/AI Service (Python)

| Layer | Technology | Notes |
|-------|-----------|-------|
| **Framework** | FastAPI | REST API for ML model inference |
| **Task Queue** | Celery + Redis | Async batch processing |
| **Speech-to-Text** | OpenAI Whisper | Lyrics transcription |
| **LLM Client** | OpenRouter | AI chat integration |
| **Embeddings** | OpenAI API | Track/lyric embeddings |
| **ORM** | SQLAlchemy + Alembic | DB migrations |
| **Mood Analysis** | Custom | OpenAI-based mood extraction |

## Dev Tooling

| Tool | Purpose |
|------|---------|
| **Air** | Live-reload for Go (`.air.toml`) |
| **golangci-lint** | Go linting |
| **ESLint + Prettier** | Frontend linting/formatting |
| **Commitlint** | Conventional commit enforcement |
| **dev.sh** | Unified dev server manager |
| **watchdog.sh** | Automon service health monitor |
