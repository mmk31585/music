"""Audio storage resolution clients (shared_volume / presigned_url)."""

from app.infra.storage.base import AudioStorageClient
from app.infra.storage.factory import get_storage_client

__all__ = ["AudioStorageClient", "get_storage_client"]
