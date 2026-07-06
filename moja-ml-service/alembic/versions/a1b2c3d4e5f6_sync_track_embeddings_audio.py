"""sync track_embeddings → track_embeddings_audio

The Go migration 000045 renamed track_embeddings → track_embeddings_audio
to separate audio embeddings (OpenL3) from text embeddings (OpenAI).
The Python code already uses track_embeddings_audio everywhere, but the
Alembic migration 9e51c3f2a4d1 still creates the old table name. This
migration:

1. Drops the stale track_embeddings table (if it exists — created by
   the old Alembic migration when it ran after Go 000045)
2. Ensures the HNSW index exists on track_embeddings_audio for fast
   approximate nearest-neighbor search (cosine distance)

Revision ID: a1b2c3d4e5f6
Revises: f0644ae1851b
Create Date: 2026-07-06 00:00:00.000000
"""
from collections.abc import Sequence

import sqlalchemy as sa
from pgvector.sqlalchemy import Vector

from alembic import op

revision: str = "a1b2c3d4e5f6"
down_revision: str | Sequence[str] | None = "f0644ae1851b"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


INDEX_NAME = "track_embeddings_audio_hnsw_idx"
TABLE_NAME = "track_embeddings_audio"
STALE_TABLE = "track_embeddings"


def upgrade() -> None:
    """Upgrade schema."""

    # ── Drop stale track_embeddings table ────────────────────────────────
    # Go migration 000045 renamed this to track_embeddings_audio, so if
    # the old table still exists it means Alembic created it after the Go
    # migration ran (dual-schema-ownership hazard).
    op.execute(f"DROP TABLE IF EXISTS {STALE_TABLE} CASCADE")

    # ── Create HNSW index for fast ANN search using cosine distance ──────
    # pgvector's HNSW index is the recommended index type for our scale
    # (thousands to low millions of tracks): better recall than IVFFlat at
    # the same query speed, no training step needed.
    #
    # We use a DO block to check existence because PostgreSQL does not
    # support CREATE INDEX IF NOT EXISTS with the USING hnsw clause before
    # pgvector 0.7+.
    op.execute(
        f"""
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1
                FROM pg_class c
                JOIN pg_namespace n ON n.oid = c.relnamespace
                WHERE c.relname = '{INDEX_NAME}'
                  AND n.nspname = current_schema
            ) THEN
                CREATE INDEX {INDEX_NAME}
                    ON {TABLE_NAME}
                    USING hnsw (embedding vector_cosine_ops);
            END IF;
        END
        $$;
        """
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.execute(f"DROP INDEX IF EXISTS {INDEX_NAME}")
    # No need to recreate the stale table on downgrade — it was a no-op
    # if the Go migration had already run.
