"""
API key authentication middleware for all ML service endpoints.

Reads the API key from env var ``ML_API_KEY``. All endpoints except health
checks require ``Authorization: Bearer <key>`` header.
"""

from __future__ import annotations

import os

from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer

_security = HTTPBearer(auto_error=False)


def _get_ml_api_key() -> str:
    return os.getenv("ML_API_KEY", "")


async def require_api_key(
    credentials: HTTPAuthorizationCredentials | None = Depends(_security),
) -> None:
    api_key = _get_ml_api_key()
    if not api_key:
        # No API key configured — allow all requests (dev mode)
        return
    if credentials is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Missing Authorization header",
            headers={"WWW-Authenticate": "Bearer"},
        )
    if credentials.credentials != api_key:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid API key",
            headers={"WWW-Authenticate": "Bearer"},
        )
