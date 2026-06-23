"""
Synchronous SQLAlchemy session factory for Celery worker processes.

Celery tasks are synchronous, so they cannot use the async session from
``app.infra.db.session``. This module provides an equivalent sync factory
that reads the same ``DATABASE_URL`` from ``app.config.settings``.
"""

from __future__ import annotations

from collections.abc import Generator
from contextlib import contextmanager

from sqlalchemy import create_engine
from sqlalchemy.orm import Session, sessionmaker

from app.config import settings

# Convert the async connection string to a sync one for Celery workers.
#   postgresql+asyncpg://user:pass@host/db  →  postgresql://user:pass@host/db
_sync_url = settings.database_url.replace("+asyncpg", "")

_sync_engine = create_engine(
    _sync_url,
    pool_size=2,
    max_overflow=4,
    echo=settings.debug,
)

_sync_session_factory = sessionmaker(
    bind=_sync_engine,
    class_=Session,
    expire_on_commit=False,
)


@contextmanager
def get_sync_db() -> Generator[Session, None, None]:
    """Yield a sync DB session, closing it when the context exits.

    Usage in Celery tasks::

        with get_sync_db() as session:
            job = session.get(LyricsJob, job_id)
            ...
    """
    session = _sync_session_factory()
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()
