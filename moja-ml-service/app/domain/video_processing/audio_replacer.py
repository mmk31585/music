"""
Audio replacement for user-uploaded video edits.

Replaces a video file's audio track with a specified audio file from Moja's
master track, starting at a configurable offset within the audio. Uses FFmpeg
for the mux operation -- the video stream is copied (no re-encode), only the
audio stream is re-encoded into AAC.

Domain isolation
----------------
``AudioReplacer`` knows nothing about HTTP, Celery, S3, or shared volumes.
It only knows "here is a video file on disk, here is an audio file on disk".
The outer layer (worker / storage) is responsible for getting the files to
local paths before calling this service.
"""

from __future__ import annotations

import logging
import os
import shutil
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path

import ffmpeg

logger = logging.getLogger(__name__)


# ── Custom exceptions ──────────────────────────────────────────────────


class VideoProcessingError(Exception):
    """Raised when the audio-replacement / thumbnail-extraction pipeline fails.

    Wraps any internal FFmpeg crash, I/O error, etc. so callers see a single
    exception type that maps to a permanent job failure (not retried).
    """


# ── Result dataclass ───────────────────────────────────────────────────


@dataclass
class VideoProcessingResult:
    """The final output of the video processing pipeline."""

    output_video_path: str = ""
    thumbnail_path: str = ""
    duration_ms: int = 0
    aspect_ratio: str = ""  # e.g. "9:16", "16:9", "1:1"
    file_size_bytes: int = 0


# ── Service ───────────────────────────────────────────────────────────


class AudioReplacer:
    """Replaces a video file's audio track with a specified audio file.

    Processing steps:
    1. Probe the input video to determine duration and aspect ratio.
    2. Trim the input audio starting at ``track_start_ms``, for the
       duration of the video (video dictates length, audio follows).
    3. Mux: video stream from input_video + trimmed audio stream.
       Video is **copied** (``-c:v copy``), audio is re-encoded as AAC.
    4. Extract a thumbnail frame at 10% of the video duration.
    5. All work happens in a temp subdirectory; results are atomically
       moved to *output_dir* only on success.
    6. Raises ``VideoProcessingError`` on any FFmpeg non-zero exit.
    """

    def process(
        self,
        video_path: str,
        audio_path: str,
        track_start_ms: int,
        output_dir: str,
    ) -> VideoProcessingResult:
        """Execute the full audio-replacement pipeline.

        Args:
            video_path: Absolute path to the uploaded video file on disk.
            audio_path: Absolute path to the master track audio file on disk.
            track_start_ms: Offset into the track's audio to start from (ms).
            output_dir: Directory where the final output files will be placed.

        Returns:
            A ``VideoProcessingResult`` with paths to the output video and
            thumbnail, plus duration and aspect ratio metadata.

        Raises:
            VideoProcessingError: If any FFmpeg step fails.
        """
        # Validate inputs exist
        for label, path in [("video", video_path), ("audio", audio_path)]:
            if not os.path.isfile(path):
                raise VideoProcessingError(
                    f"Input {label} file not found: {path}"
                )

        os.makedirs(output_dir, exist_ok=True)

        # Work in a temp directory to avoid leaving partial files
        with tempfile.TemporaryDirectory(
            prefix="video_proc_", dir=output_dir
        ) as tmp_dir:
            try:
                # ── Step 1: Probe input video ─────────────────────────
                logger.info("Probing video: %s", video_path)
                probe = ffmpeg.probe(video_path)
                video_stream = self._find_video_stream(probe)
                if video_stream is None:
                    raise VideoProcessingError(
                        "No video stream found in input file"
                    )

                duration_s = float(probe.get("format", {}).get("duration", 0))
                duration_ms = int(duration_s * 1000)
                aspect_ratio = self._detect_aspect_ratio(video_stream)

                logger.info(
                    "Video: duration=%.2fs, aspect=%s, codec=%s",
                    duration_s,
                    aspect_ratio,
                    video_stream.get("codec_name", "unknown"),
                )

                # ── Step 2: Mux video + trimmed audio ────────────────
                # Temporary paths in the working directory
                temp_output = os.path.join(tmp_dir, "output.mp4")
                start_s = track_start_ms / 1000.0

                logger.info(
                    "Muxing: video=%s, audio=%s, audio_offset=%.2fs",
                    video_path,
                    audio_path,
                    start_s,
                )

                input_video = ffmpeg.input(video_path)
                input_audio = ffmpeg.input(audio_path, ss=start_s)

                stream = ffmpeg.output(
                    input_video["v"],
                    input_audio["a"],
                    temp_output,
                    # Video: copy (no re-encode)
                    vcodec="copy",
                    # Audio: re-encode as AAC 192k
                    acodec="aac",
                    audio_bitrate="192k",
                    # Shortest: stop when the shortest input ends (the video)
                    shortest=None,
                )

                # Suppress ffmpeg-python's default verbose output --
                # we handle errors via the return code.
                stream = ffmpeg.overwrite_output(stream)

                # ffmpeg-python builds a command but we run it via
                # subprocess for explicit error handling.
                cmd = stream.compile()
                logger.debug("FFmpeg mux command: %s", " ".join(cmd))

                self._run_ffmpeg(cmd, "Muxing audio onto video")

                # Verify output was created
                if not os.path.isfile(temp_output):
                    raise VideoProcessingError(
                        "FFmpeg mux completed but output file not found"
                    )

                file_size = os.path.getsize(temp_output)
                logger.info(
                    "Mux complete: %s (%d bytes)",
                    temp_output,
                    file_size,
                )

                # ── Step 3: Extract thumbnail ─────────────────────────
                # Pick a frame at 10% of the video duration to avoid
                # black frames at 0% (fade-in) or end credits.
                thumb_offset_s = duration_s * 0.10
                temp_thumb = os.path.join(tmp_dir, "thumbnail.jpg")

                logger.info(
                    "Extracting thumbnail at %.2fs (10%% of duration)",
                    thumb_offset_s,
                )

                thumb_cmd = [
                    "ffmpeg",
                    "-y",
                    "-ss", str(thumb_offset_s),
                    "-i", video_path,
                    "-vframes", "1",
                    "-q:v", "2",  # high-quality JPEG
                    temp_thumb,
                ]
                self._run_ffmpeg(thumb_cmd, "Thumbnail extraction")

                if not os.path.isfile(temp_thumb):
                    raise VideoProcessingError(
                        "Thumbnail extraction completed but file not found"
                    )

                logger.info("Thumbnail extracted: %s", temp_thumb)

                # ── Step 4: Move results to final output dir ──────────
                output_video_name = (
                    f"video_{Path(video_path).stem}_"
                    f"{Path(audio_path).stem}_processed.mp4"
                )
                final_video_path = os.path.join(
                    output_dir, output_video_name
                )
                shutil.move(temp_output, final_video_path)

                output_thumb_name = (
                    f"thumb_{Path(video_path).stem}_"
                    f"{Path(audio_path).stem}.jpg"
                )
                final_thumb_path = os.path.join(
                    output_dir, output_thumb_name
                )
                shutil.move(temp_thumb, final_thumb_path)

                logger.info(
                    "Results moved to: video=%s, thumb=%s",
                    final_video_path,
                    final_thumb_path,
                )

                return VideoProcessingResult(
                    output_video_path=final_video_path,
                    thumbnail_path=final_thumb_path,
                    duration_ms=duration_ms,
                    aspect_ratio=aspect_ratio,
                    file_size_bytes=file_size,
                )

            except VideoProcessingError:
                raise
            except ffmpeg.Error as exc:
                stderr = (
                    exc.stderr.decode("utf-8", errors="replace")[:500]
                    if exc.stderr
                    else ""
                )
                raise VideoProcessingError(
                    f"FFmpeg failed: {stderr}"
                ) from exc
            except Exception as exc:
                logger.exception("Unexpected error during video processing")
                raise VideoProcessingError(
                    f"Video processing failed: {exc}"
                ) from exc

    # ── Internal helpers ───────────────────────────────────────────────

    @staticmethod
    def _find_video_stream(probe: dict) -> dict | None:
        """Return the first video stream from the FFprobe output."""
        streams = probe.get("streams", [])
        for s in streams:
            if s.get("codec_type") == "video":
                return s
        return None

    @staticmethod
    def _detect_aspect_ratio(video_stream: dict) -> str:
        """Detect display aspect ratio from the video stream metadata.

        Returns a string like ``"16:9"``, ``"9:16"``, ``"4:3"``, ``"1:1"``.
        Falls back to deriving from width/height if ``display_aspect_ratio``
        is not present.
        """
        dar = video_stream.get("display_aspect_ratio")
        if dar and ":" in dar:
            return dar

        width = video_stream.get("width", 1)
        height = video_stream.get("height", 1)

        # Reduce to lowest terms using GCD
        def _gcd(a: int, b: int) -> int:
            while b:
                a, b = b, a % b
            return a

        g = _gcd(width, height)
        return f"{width // g}:{height // g}"

    @staticmethod
    def _run_ffmpeg(cmd: list[str], step_label: str) -> None:
        """Execute an FFmpeg command and raise on failure."""
        logger.info("Running FFmpeg step [%s]...", step_label)
        try:
            subprocess.run(
                cmd,
                capture_output=True,
                text=False,
                check=True,
            )
        except subprocess.CalledProcessError as exc:
            stderr = (
                exc.stderr.decode("utf-8", errors="replace")[:500]
                if exc.stderr
                else ""
            )
            raise VideoProcessingError(
                f"FFmpeg [{step_label}] failed (exit {exc.returncode}): {stderr}"
            ) from exc

        logger.debug("FFmpeg [%s] completed successfully", step_label)


# ── Singleton factory (established convention) ─────────────────────────


_instance: AudioReplacer | None = None


def get_audio_replacer() -> AudioReplacer:
    """Return the singleton ``AudioReplacer`` instance.

    ``AudioReplacer`` holds no state, so this is trivial, but the singleton
    factory pattern is the established convention in this codebase (matching
    ``get_transcriber()``, ``get_embedding_extractor()``, etc.).
    """
    global _instance  # noqa: PLW0603
    if _instance is None:
        _instance = AudioReplacer()
    return _instance
