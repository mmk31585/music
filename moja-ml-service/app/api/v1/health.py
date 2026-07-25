"""
Health check endpoints.

Liveness  → process is alive (used by orchestrator restart policies).
Readiness → service can actually do work.

Architecture note (Phase 3)
---------------------------
The Whisper model is loaded ONLY in the Celery worker process (via a
``worker_process_init`` signal in ``app.workers.tasks.lyrics_tasks``).
The FastAPI web process never loads the model — it only enqueues jobs
and queries their status. Therefore ``/health/ready`` in the FastAPI
process checks **database + Redis connectivity** only, not Whisper.

If you need to verify the Whisper model is loaded, check the Celery
worker's startup logs for the ``warm_up_whisper_model`` signal output.
"""

from typing import Any

import redis.asyncio as aioredis
from fastapi import APIRouter
from sqlalchemy import text

from app.config import settings
from app.infra.db.session import async_session_maker

router = APIRouter()


@router.get("/health/live")
async def liveness() -> dict[str, str]:
    """Process is up and responding. Used by orchestrator restart policies."""
    return {"status": "ok"}


@router.get("/health/ready")
async def readiness() -> dict[str, Any]:
    """
    Service readiness probe for the FastAPI web process.

    Checks:
      1. Database connectivity (SELECT 1).
      2. Redis connectivity (PING).

    NOTE: Whisper model readiness is NOT checked here — the model lives
    in the Celery worker process, not the FastAPI process. See module
    docstring for details.
    """
    checks: dict[str, Any] = {
        "status": "ok",
        "checks": {
            "database": {"status": "unknown"},
            "redis": {"status": "unknown"},
        },
    }

    # ── Database check ──────────────────────────────────────────────
    try:
        async with async_session_maker() as session:
            result = await session.execute(text("SELECT 1"))
            row = result.scalar_one()
            if row == 1:
                checks["checks"]["database"]["status"] = "ok"
            else:
                checks["checks"]["database"]["status"] = "error"
                checks["checks"]["database"]["reason"] = "unexpected_select_1_result"
    except Exception as exc:
        checks["checks"]["database"]["status"] = "error"
        checks["checks"]["database"]["reason"] = str(exc)

    # ── Redis check ─────────────────────────────────────────────────
    try:
        r = aioredis.from_url(settings.redis_url)
        pong = await r.ping()
        await r.aclose()
        if pong is True:
            checks["checks"]["redis"]["status"] = "ok"
        else:
            checks["checks"]["redis"]["status"] = "error"
            checks["checks"]["redis"]["reason"] = "ping_returned_non_true"
    except Exception as exc:
        checks["checks"]["redis"]["status"] = "error"
        checks["checks"]["redis"]["reason"] = str(exc)

    # ── Overall status ──────────────────────────────────────────────
    all_ok = all(
        check.get("status") == "ok"
        for check in checks["checks"].values()
    )
    if not all_ok:
        checks["status"] = "degraded"

    return checks
