"""Abstract base class for audio storage resolution."""

from abc import ABC, abstractmethod


class AudioStorageClient(ABC):
    """Resolve an audio source identifier into a local file path.

    Two concrete implementations exist:

    * ``SharedVolumeStorageClient`` — the ``source`` is a relative path
      within a shared filesystem mount (shared_volume mode).
    * ``PresignedUrlStorageClient`` — the ``source`` is an HTTP(S) URL
      from which the file is downloaded (presigned_url mode).
    """

    @abstractmethod
    def fetch_to_local_path(self, source: str, dest_dir: str) -> str:
        """Download / resolve the audio source to a local file path.

        Args:
            source: Either a relative path (shared_volume) or a
                    presigned download URL (presigned_url).
            dest_dir: Directory where temp files may be written
                      (unused in shared_volume mode, required in
                      presigned_url mode).

        Returns:
            Absolute local path to the audio file.

        Raises:
            FileNotFoundError: The audio file does not exist or is
                unreachable.
            ValueError: Path traversal or other security violations.
            IOError: Download failures, timeouts, etc.
        """
        ...
