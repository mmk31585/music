"""Video processing domain — audio replacement, thumbnail extraction, metadata detection."""

from app.domain.video_processing.audio_replacer import (
    AudioReplacer,
    VideoProcessingError,
    VideoProcessingResult,
    get_audio_replacer,
)

__all__ = [
    "AudioReplacer",
    "VideoProcessingError",
    "VideoProcessingResult",
    "get_audio_replacer",
]
