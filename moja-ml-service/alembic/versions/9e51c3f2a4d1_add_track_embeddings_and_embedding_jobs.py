"""add track_embeddings and embedding_jobs tables with pgvector

Revision ID: 9e51c3f2a4d1
Revises: 77a349a83b71
Create Date: 2026-06-21 14:00:00.000000
"""
from collections.abc import Sequence

import sqlalchemy as sa
from pgvector.sqlalchemy import Vector

from alembic import op

revision: str = "9e51c3f2a4d1"
down_revision: str | Sequence[str] | None = "77a349a83b71"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    """Upgrade schema."""
    # Enable pgvector extension (must be before any table using Vector type).
    # This is idempotent — safe to run if already enabled.
    op.execute("CREATE EXTENSION IF NOT EXISTS vector")

    # ── track_embeddings table ──────────────────────────────────────────
    # Stores the 512-dim OpenL3 embedding per track plus metadata for
    # similarity boosting/filtering.  The HNSW index below enables fast
    # approximate nearest-neighbor search using cosine distance.
    op.create_table(
        "track_embeddings",
        sa.Column("track_id", sa.String(), nullable=False),
        sa.Column("embedding", Vector(512), nullable=False),
        sa.Column("embedding_model", sa.String(), nullable=False),
        sa.Column("tempo_bpm", sa.Float(), nullable=True),
        sa.Column("genre_ids", sa.ARRAY(sa.String()), nullable=False, server_default="{}"),
        sa.Column("release_year", sa.Integer(), nullable=True),
        sa.Column("extracted_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.PrimaryKeyConstraint("track_id"),
    )

    # HNSW index for approximate nearest-neighbor search using cosine
    # distance.  Cosine is the standard similarity metric for audio
    # embeddings from models like OpenL3 — direction (spectral profile)
    # matters more than magnitude (loudness), which makes cosine the
    # natural choice over Euclidean (L2) for this embedding family.
    # HNSW is preferred over IVFFlat for this scale (thousands to low
    # millions of tracks): it offers better recall at the same query
    # speed and doesn't need a training step.
    op.execute(
        "CREATE INDEX track_embeddings_hnsw_idx ON track_embeddings "
        "USING hnsw (embedding vector_cosine_ops)"
    )

    # ── embedding_jobs table ────────────────────────────────────────────
    op.create_table(
        "embedding_jobs",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("track_id", sa.String(), nullable=False),
        sa.Column("idempotency_key", sa.String(), nullable=False),
        sa.Column("status", sa.String(), nullable=False, server_default="queued"),
        sa.Column("audio_file_path", sa.String(), nullable=True),
        sa.Column("audio_download_url", sa.String(), nullable=True),
        sa.Column("genre_ids", sa.ARRAY(sa.String()), nullable=False, server_default="{}"),
        sa.Column("release_year", sa.Integer(), nullable=True),
        sa.Column("error_message", sa.String(), nullable=True),
        sa.Column("retry_count", sa.Integer(), nullable=False, server_default="0"),
        sa.Column("created_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(
        op.f("ix_embedding_jobs_idempotency_key"),
        "embedding_jobs",
        ["idempotency_key"],
        unique=True,
    )
    op.create_index(
        op.f("ix_embedding_jobs_track_id"),
        "embedding_jobs",
        ["track_id"],
        unique=False,
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_index(op.f("ix_embedding_jobs_track_id"), table_name="embedding_jobs")
    op.drop_index(op.f("ix_embedding_jobs_idempotency_key"), table_name="embedding_jobs")
    op.drop_table("embedding_jobs")
    op.execute("DROP INDEX IF EXISTS track_embeddings_hnsw_idx")
    op.drop_table("track_embeddings")
    op.execute("DROP EXTENSION IF EXISTS vector")
