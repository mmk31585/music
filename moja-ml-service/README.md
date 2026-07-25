# Moja ML Service

**Audio/ML processing microservice for the Moja (موجا) Persian music platform.**

This service handles CPU-bound ML workloads that the Go backend offloads via
Celery tasks: speech-to-text transcription (faster-whisper), audio feature
extraction, recommendation inference, and lyrics processing. It runs as an
internal microservice alongside the main Go API — never exposed to the public
internet.

---

## Storage Integration Decision

> **Problem:** How does this Python service access audio files uploaded by
> users through the Go backend?

The Go backend stores uploaded audio files. This ML service needs to read them
for processing (transcription, feature extraction, etc.). We cannot assume the
same storage topology in every deployment — a solo developer on a single VPS
and a team running on Kubernetes with S3 have very different needs.

Two modes are supported, selected via the `STORAGE_MODE` environment variable:

### Mode A: `shared_volume` (default, recommended for local dev)

```
┌──────────────────┐     mount (ro)     ┌──────────────────┐
│  Go backend      │ ──────────────────▶│  ML service      │
│  uploads/        │   -v ./uploads:     │  /mnt/audio-     │
│  track-audio/    │     /mnt/audio-     │  uploads/        │
│  *.mp3           │     uploads:ro      │  track-audio/    │
└──────────────────┘                    └──────────────────┘
                                           audio_file_path:
                                           "track-audio/abc.mp3"
```

**How it works:**
- The Go backend's `uploads/` directory is bind-mounted (read-only) into the
  ML service container at `AUDIO_SHARED_PATH`.
- Celery job payloads contain `audio_file_path` — a relative path from the
  uploads root.
- The ML service constructs the full path as
  `{AUDIO_SHARED_PATH}/{audio_file_path}`.

**When to use:**
- Single-VPS deployments where everything runs on one machine.
- Local development with Docker Compose.
- No object storage available.

### Mode B: `presigned_url` (recommended for production with S3/MinIO)

```
┌──────────────────┐  job payload        ┌──────────────────┐
│  Go backend      │ ──────────────────▶│  ML service      │
│  (has S3 creds)  │  audio_download_url │  (no creds)     │
│                  │  "https://s3...."   │                  │
│  ──► generates   │                     │  ──► HTTP GET    │
│      presigned   │                     │      download    │
│      download URL│                     │      → temp file │
└──────────────────┘                     └──────────────────┘
                                           ➜ process ➜ delete
```

**How it works:**
- Go generates a presigned/signed download URL (S3 or MinIO) and includes it
  as `audio_download_url` in the Celery job payload.
- The ML service downloads the file via HTTP(S) to a temp directory, processes
  it, and deletes the temp file.
- The ML service **never needs AWS credentials or S3 access keys**.

**When to use:**
- Multi-server deployments.
- S3/MinIO object storage is already in use.
- You want strict credential isolation (Go has S3 access, ML service has none).

### Job payload schema

Both modes are supported by a single job payload schema. Exactly one of the
two audio source fields will be present:

```python
class TranscriptionJobPayload(BaseModel):
    """Payload for a transcription (or other audio-processing) task.

    Exactly one of audio_file_path or audio_download_url must be provided,
    depending on STORAGE_MODE. The service checks which field is set and
    resolves the audio source accordingly.
    """

    # Shared volume mode
    audio_file_path: str | None = None
    # e.g. "track-audio/abc123.mp3"
    # Full path: /mnt/audio-uploads/track-audio/abc123.mp3

    # Presigned URL mode
    audio_download_url: str | None = None
    # e.g. "https://minio.example.com/moja-uploads/track-audio/abc123.mp3?X-Amz-Signature=..."

    # Common fields
    track_id: str
    job_id: str
    webhook_callback_url: str
    # Where to POST results. Overrides the default go_backend_callback_url
    # when the Go backend wants a specific endpoint per task type.
    language: str | None = None  # Override whisper_language_hint
```

This schema is defined here for documentation. The actual implementation
(Pydantic model) will be created in Phase 3 when the Celery worker is built.

---

## Local Development

### Prerequisites

- Python 3.11+
- [uv](https://docs.astral.sh/uv/) (package manager)
- Docker & Docker Compose (for Redis + Postgres)

### Setup

```bash
# 1. Clone and enter the project
cd moja-ml-service

# 2. Create virtual environment and install dependencies
uv sync

# 3. Copy environment file and edit as needed
cp .env.example .env

# 4. Start infrastructure (Redis + Postgres)
docker compose up -d

# 5. Run database migrations
uv run alembic upgrade head

# 6. Start the FastAPI dev server
uv run uvicorn app.main:app --reload
```

The service will be available at http://localhost:8000.

### Health checks

```bash
# Liveness — is the process up?
curl http://localhost:8000/api/v1/health/live

# Readiness — can the service actually do work?
curl http://localhost:8000/api/v1/health/ready
```

### Running tests

```bash
uv run pytest
```

With coverage:

```bash
uv run pytest --cov=app
```

### Linting and type checking

```bash
uv run ruff check .
uv run mypy app/
```

---

## Project Structure

```
moja-ml-service/
├── app/
│   ├── main.py                 # FastAPI application entry point
│   ├── config.py               # pydantic-settings configuration
│   ├── api/
│   │   └── v1/
│   │       ├── router.py       # API router aggregation
│   │       └── health.py       # Liveness + readiness probes
│   ├── domain/
│   │   ├── lyrics/             # Lyrics processing (Phase 2+)
│   │   ├── recommendation/     # Recommendation engine (Phase 4+)
│   │   └── audio_features/     # Feature extraction (Phase 3+)
│   ├── workers/
│   │   └── __init__.py         # Celery workers (Phase 3+)
│   ├── infra/
│   │   ├── storage/            # Audio file resolution (shared_volume / presigned_url)
│   │   ├── db/                 # SQLAlchemy setup + session dependency
│   │   └── callback/           # Webhook callback to Go backend
│   └── schemas/
│       └── __init__.py         # Pydantic schemas (Phase 3+)
├── alembic/                    # Database migrations
├── tests/
│   ├── unit/
│   └── integration/
├── Dockerfile
├── docker-compose.yml
├── pyproject.toml
└── README.md
```

---

## Infrastructure Constraints

| Constraint          | Value                                       |
|---------------------|---------------------------------------------|
| CPU                 | x86_64, no GPU                              |
| RAM                 | 4-8 GB (assumed; flag any requirement >8GB) |
| Python              | 3.11+                                       |
| Package manager     | uv                                          |
| Message broker      | Redis (via Celery)                          |
| Database            | PostgreSQL 16 (async via asyncpg)           |
| ASGI server         | Uvicorn                                     |
| ML inference        | faster-whisper (CPU int8 quantized)         |

---

## Deployment Notes

- **Whisper model:** On first run, faster-whisper downloads the model from
  Hugging Face Hub (~1.5 GB for `medium`). Ensure the VPS has internet access
  for the initial download, or pre-download and cache the model.
- **RAM:** The `medium` Whisper model with `int8` quantization uses ~2 GB RAM
  during transcription. A 4 GB VPS is sufficient for one concurrent job.
  For higher throughput, add more Celery workers or increase VPS memory.
- **Concurrency:** The FastAPI app handles HTTP requests. Heavy ML work runs
  in Celery workers (Phase 3), keeping the API responsive for health checks
  and light queries.
