from typing import Literal

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    # ── App ──────────────────────────────────────────────────────────────
    app_name: str = "moja-ml-service"
    environment: str = "development"
    debug: bool = False

    # ── Database (this service's OWN database — separate from Go's catalog DB) ──
    database_url: str = "postgresql+asyncpg://moja_ml:changeme_in_production@localhost:5433/moja_ml"

    # ── Redis / Celery ───────────────────────────────────────────────────
    redis_url: str = "redis://localhost:6379/0"
    celery_broker_url: str = "redis://localhost:6379/1"
    celery_result_backend: str = "redis://localhost:6379/2"

    # ── Storage mode ─────────────────────────────────────────────────────
    # "shared_volume"  — Go's uploads dir mounted at audio_shared_path (local/VPS)
    # "presigned_url"  — Go passes a presigned S3 download URL in the job payload
    storage_mode: Literal["shared_volume", "presigned_url"] = "shared_volume"
    audio_shared_path: str = "/mnt/audio-uploads"

    # ── Whisper ──────────────────────────────────────────────────────────
    whisper_model_size: str = "medium"
    whisper_compute_type: str = "int8"
    whisper_device: str = "cpu"
    whisper_language_hint: str = "fa"

    # ── Audio Embedding (openl3) ─────────────────────────────────────────
    audio_embedding_model_variant: str = "mel256"
    audio_embedding_size: int = 512
    audio_embedding_device: str = "cpu"

    # ── Go backend callback ──────────────────────────────────────────────
    go_backend_callback_url: str = "http://localhost:8080/api/v1"
    webhook_hmac_secret: str = "dev-secret-change-in-production"

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        env_prefix="",
    )


settings = Settings()
