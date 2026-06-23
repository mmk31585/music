"""add video_processing_jobs table

Revision ID: f0644ae1851b
Revises: 9e51c3f2a4d1
Create Date: 2026-06-22 00:00:00.000000
"""
from collections.abc import Sequence

import sqlalchemy as sa

from alembic import op

revision: str = "f0644ae1851b"
down_revision: str | Sequence[str] | None = "9e51c3f2a4d1"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    """Upgrade schema."""
    op.create_table(
        "video_processing_jobs",
        sa.Column("id", sa.Uuid(), nullable=False),
        sa.Column("video_id", sa.String(), nullable=False),
        sa.Column("idempotency_key", sa.String(), nullable=False),
        sa.Column("status", sa.String(), nullable=False, server_default="queued"),
        sa.Column("raw_video_path", sa.String(), nullable=True),
        sa.Column("raw_video_download_url", sa.String(), nullable=True),
        sa.Column("audio_file_path", sa.String(), nullable=True),
        sa.Column("audio_download_url", sa.String(), nullable=True),
        sa.Column("track_start_ms", sa.Integer(), nullable=False, server_default="0"),
        sa.Column("result_video_path", sa.String(), nullable=True),
        sa.Column("result_thumbnail_path", sa.String(), nullable=True),
        sa.Column("result_duration_ms", sa.Integer(), nullable=True),
        sa.Column("result_aspect_ratio", sa.String(), nullable=True),
        sa.Column("error_message", sa.String(), nullable=True),
        sa.Column("retry_count", sa.Integer(), nullable=False, server_default="0"),
        sa.Column("callback_delivered", sa.Boolean(), nullable=False, server_default="false"),
        sa.Column("created_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.Column("updated_at", sa.DateTime(), nullable=False, server_default=sa.func.now()),
        sa.PrimaryKeyConstraint("id"),
    )
    op.create_index(
        op.f("ix_video_processing_jobs_video_id"),
        "video_processing_jobs",
        ["video_id"],
        unique=False,
    )
    op.create_index(
        op.f("ix_video_processing_jobs_idempotency_key"),
        "video_processing_jobs",
        ["idempotency_key"],
        unique=True,
    )


def downgrade() -> None:
    """Downgrade schema."""
    op.drop_index(
        op.f("ix_video_processing_jobs_idempotency_key"),
        table_name="video_processing_jobs",
    )
    op.drop_index(
        op.f("ix_video_processing_jobs_video_id"),
        table_name="video_processing_jobs",
    )
    op.drop_table("video_processing_jobs")
