"""Tests for health check endpoints."""

import pytest
from httpx import AsyncClient


@pytest.mark.asyncio
async def test_liveness(async_client: AsyncClient) -> None:
    """GET /api/v1/health/live should return 200 with status ok."""
    response = await async_client.get("/api/v1/health/live")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}


@pytest.mark.asyncio
async def test_readiness_returns_ok(async_client: AsyncClient) -> None:
    """GET /api/v1/health/ready should return 200.

    In Phase 3 the FastAPI process doesn't load the Whisper model (that's
    the Celery worker's job), so readiness only checks DB + Redis.
    """
    response = await async_client.get("/api/v1/health/ready")
    assert response.status_code == 200
    data = response.json()
    assert "status" in data
    assert "checks" in data
    # Whisper should NOT be in the checks (it's a worker-only concern)
    assert "whisper" not in data["checks"]
