"""
Celery application configuration for moja-ml-service.

Defines the Celery app instance with routing, serialization, and the
critical concurrency limit (see ``worker_concurrency`` doc below).

Queue layout::

    lyrics_queue       ←  all lyrics-generation tasks
    covers_queue       ←  all cover-optimization tasks
    embedding_queue    ←  all embedding-extraction tasks
    video_queue        ←  all video-processing tasks

Why separate queues for lyrics and embedding?
    Both openl3 (embedding) and faster-whisper (transcription) hold
    significant RAM when their models are loaded (~1-2 GB each).
    Running them as separate queues lets you choose between:
      a) Single combined worker (concurrency=1) listening to both
         queues — conservative default for a tight-RAM VPS.
      b) Two separate worker processes (each concurrency=1) on
         separate containers — if either queue backs up and enough
         RAM is available.
    Having explicit routing allows either model without code changes.
"""

from __future__ import annotations

from celery import Celery

from app.config import settings

celery_app = Celery(
    "moja_ml_service",
    broker=settings.celery_broker_url,
    backend=settings.celery_result_backend,
)

celery_app.conf.update(
    # ── Serialization ───────────────────────────────────────────────
    task_serializer="json",
    accept_content=["json"],
    result_serializer="json",
    timezone="UTC",
    enable_utc=True,
    # ── Concurrency ─────────────────────────────────────────────────
    # DO NOT increase worker_concurrency above 1 for this queue.
    #
    # faster-whisper "medium" model holds significant RAM during
    # inference (~2 GB working set). Running 2+ transcriptions in
    # parallel on a no-GPU VPS WILL cause OOM kills under realistic
    # load. If throughput becomes a bottleneck, scale by running
    # multiple WORKER PROCESSES (each still concurrency=1) behind
    # the same queue, not by raising this number. Confirm available
    # RAM before even considering that.  The same applies to openl3
    # which loads a TensorFlow model (~1 GB working set).
    worker_concurrency=1,
    worker_prefetch_multiplier=1,
    # ── Task acknowledgments ────────────────────────────────────────
    task_acks_late=True,
    task_reject_on_worker_lost=True,
    # ── Routing ─────────────────────────────────────────────────────
    # ── Queue routing ───────────────────────────────────────────────────
    #
    # Why a SEPARATE ``embedding_queue`` (not sharing ``lyrics_queue``)?
    #
    # Both embedding extraction (openl3 + tensorflow) and transcription
    # (faster-whisper "medium") have significant RAM footprints when their
    # respective models are loaded into memory (~1-2 GB working set each).
    # By routing them to separate queues you have two deployment choices:
    #
    #   A) SINGLE combined worker (concurrency=1) listening to ALL queues:
    #       celery worker --queues=lyrics_queue,embedding_queue,covers_queue
    #      Tasks run sequentially — safe on tight-RAM VPS, just slower at
    #      draining the backlog.  This is the default for development /
    #      low-traffic deployments (see docker-compose.yml).
    #
    #   B) Separate worker PROCESSES (each concurrency=1), each listening
    #      to one queue.  Needed if/when either queue backs up significantly
    #      under real load — each worker loads only its own model, so
    #      interleaving is not a concern.
    #
    # The routing MUST be explicit regardless of whether workers are
    # combined or split — implicit routing (publishing to the default
    # queue) would make it impossible to separate later without a
    # coordinated drain-and-redeploy.
    #
    # Note: covers_queue is intentionally in the same routing dict because
    # cover optimization (Pillow) is comparatively lightweight — it can
    # safely share a worker with either the lyrics or embedding queue
    # without memory pressure.
    task_routes={
        "app.workers.tasks.lyrics_tasks.*": {"queue": "lyrics_queue"},
        "app.workers.tasks.cover_tasks.*": {"queue": "covers_queue"},
        "app.workers.tasks.embedding_tasks.*": {"queue": "embedding_queue"},
        "app.workers.tasks.video_tasks.*": {"queue": "video_queue"},
    },
)

# ── Beat schedule (periodic tasks) ────────────────────────────────────
# Enables periodic re-optimization of covers and re-enrichment of metadata.
# Enable with: celery -A app.workers.celery_app beat
celery_app.conf.beat_schedule = {
    # Re-optimize covers that haven't been processed in 7 days
    "refresh-stale-covers-weekly": {
        "task": "app.workers.tasks.cover_tasks.refresh_stale_covers",
        "schedule": 604800,  # 7 days in seconds
        "kwargs": {"max_age_days": 7},
    },
}

# Import tasks so Celery auto-discovers them
import app.workers.tasks.cover_tasks  # noqa: E402, F401
import app.workers.tasks.embedding_tasks  # noqa: E402, F401
import app.workers.tasks.lyrics_tasks  # noqa: E402, F401
import app.workers.tasks.video_tasks  # noqa: E402, F401
