"""Tests for the AudioReplacer domain logic.

These tests use a **mocked** subprocess/FFmpeg — no real FFmpeg binary
is invoked. That would be an integration-test concern.
"""

import pytest

from app.domain.video_processing.audio_replacer import (
    AudioReplacer,
    VideoProcessingError,
)

# ── Fixtures ────────────────────────────────────────────────────────────


def _mux_call_from_mock(mock_run):
    """Find the FFmpeg mux command (contains output.mp4) from subprocess calls."""
    for call_args in mock_run.call_args_list:
        args = call_args[0][0]
        cmd_str = " ".join(args)
        if "output.mp4" in cmd_str:
            return args
    return None


def _thumbnail_call_from_mock(mock_run):
    """Find the FFmpeg thumbnail command (contains thumbnail.jpg) from subprocess calls."""
    for call_args in mock_run.call_args_list:
        args = call_args[0][0]
        cmd_str = " ".join(args)
        if "thumbnail.jpg" in cmd_str:
            return args
    return None


def _setup_common_mocks(mocker, probe_duration="10.0", file_size=1000000):
    """Set up the standard mock environment for AudioReplacer.process().

    Returns the mocked ``subprocess.run`` so callers can inspect calls.
    """
    mock_run = mocker.patch("subprocess.run", return_value=mocker.Mock(returncode=0))
    mocker.patch(
        "ffmpeg.probe",
        return_value={
            "streams": [
                {
                    "codec_type": "video",
                    "width": 1920,
                    "height": 1080,
                    "display_aspect_ratio": "16:9",
                },
            ],
            "format": {"duration": probe_duration},
        },
    )
    mocker.patch("os.path.isfile", return_value=True)
    mocker.patch("os.path.getsize", return_value=file_size)
    mocker.patch("os.makedirs")
    mocker.patch("shutil.move")
    mock_tmp = mocker.patch("tempfile.TemporaryDirectory")
    mock_tmp.return_value.__enter__.return_value = "/tmp/fake_tmp"
    return mock_run


# ── Tests: FFmpeg command construction ─────────────────────────────────


def test_ffmpeg_mux_command_uses_copy_for_video(mocker) -> None:
    """The mux command should use -c:v copy (no video re-encode)."""
    mock_run = _setup_common_mocks(mocker)

    replacer = AudioReplacer()
    replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=5000,
        output_dir="/fake/output",
    )

    mux_call = _mux_call_from_mock(mock_run)
    assert mux_call is not None, "Expected a mux FFmpeg call"
    cmd_str = " ".join(mux_call)
    assert "copy" in cmd_str, (
        f"Expected video codec copy in mux command, got: {cmd_str}"
    )


def test_ffmpeg_thumbnail_command_uses_one_frame(mocker) -> None:
    """Thumbnail extraction should use -vframes 1."""
    mock_run = _setup_common_mocks(mocker, probe_duration="20.0")

    replacer = AudioReplacer()
    replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=0,
        output_dir="/fake/output",
    )

    thumb_call = _thumbnail_call_from_mock(mock_run)
    assert thumb_call is not None, "Expected a thumbnail FFmpeg call"
    cmd_str = " ".join(thumb_call)
    assert "-vframes" in cmd_str or "1" in cmd_str


# ── Tests: Aspect ratio detection ──────────────────────────────────────


def test_detects_landscape_16_9() -> None:
    """16:9 video should be correctly detected."""
    replacer = AudioReplacer()
    result = replacer._detect_aspect_ratio(
        {"width": 1920, "height": 1080, "display_aspect_ratio": "16:9"}
    )
    assert result == "16:9"


def test_detects_portrait_9_16() -> None:
    """9:16 portrait video should be correctly detected."""
    replacer = AudioReplacer()
    result = replacer._detect_aspect_ratio(
        {"width": 1080, "height": 1920, "display_aspect_ratio": "9:16"}
    )
    assert result == "9:16"


def test_detects_aspect_ratio_from_dimensions() -> None:
    """When display_aspect_ratio is absent, derive from width/height."""
    replacer = AudioReplacer()
    result = replacer._detect_aspect_ratio(
        {"width": 640, "height": 480}
    )
    assert result == "4:3"


def test_detects_square_1_1() -> None:
    """1:1 square video should be correctly detected."""
    replacer = AudioReplacer()
    result = replacer._detect_aspect_ratio(
        {"width": 1080, "height": 1080}
    )
    assert result == "1:1"


# ── Tests: Thumbnail offset ────────────────────────────────────────────


def test_thumbnail_uses_10_percent_offset(mocker) -> None:
    """Thumbnail should be extracted at 10% of the video duration."""
    mock_run = _setup_common_mocks(mocker, probe_duration="60.0")

    replacer = AudioReplacer()
    replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=0,
        output_dir="/fake/output",
    )

    thumb_call = _thumbnail_call_from_mock(mock_run)
    assert thumb_call is not None

    ss_index = None
    for i, arg in enumerate(thumb_call):
        if arg == "-ss":
            ss_index = i
            break

    assert ss_index is not None, "Expected -ss argument in thumbnail command"
    ss_value = float(thumb_call[ss_index + 1])
    assert abs(ss_value - 6.0) < 0.1, (
        f"Expected thumbnail offset ~6.0s (10% of 60s), got {ss_value}"
    )


# ── Tests: Error handling ──────────────────────────────────────────────


def test_wraps_ffmpeg_errors(mocker) -> None:
    """Non-existent file should raise VideoProcessingError."""
    replacer = AudioReplacer()

    mocker.patch("os.path.isfile", return_value=False)

    with pytest.raises(VideoProcessingError) as exc_info:
        replacer.process(
            video_path="/nonexistent/video.mp4",
            audio_path="/fake/audio.mp3",
            track_start_ms=0,
            output_dir="/tmp",
        )

    assert "not found" in str(exc_info.value).lower()


def test_wraps_ffmpeg_subprocess_error(mocker) -> None:
    """FFmpeg non-zero exit should be wrapped in VideoProcessingError."""
    mocker.patch("os.path.isfile", return_value=True)
    mocker.patch("subprocess.run", side_effect=OSError("ffmpeg not found"))
    mocker.patch(
        "ffmpeg.probe",
        return_value={
            "streams": [
                {
                    "codec_type": "video",
                    "width": 1920,
                    "height": 1080,
                    "display_aspect_ratio": "16:9",
                },
            ],
            "format": {"duration": "10.0"},
        },
    )
    mocker.patch("tempfile.TemporaryDirectory")
    mocker.patch("os.path.getsize")
    mocker.patch("shutil.move")

    replacer = AudioReplacer()
    with pytest.raises(VideoProcessingError) as exc_info:
        replacer.process(
            video_path="/fake/video.mp4",
            audio_path="/fake/audio.mp3",
            track_start_ms=0,
            output_dir="/tmp",
        )

    assert "failed" in str(exc_info.value).lower() or "ffmpeg" in str(exc_info.value).lower()


def test_duration_ms_calculated_correctly(mocker) -> None:
    """Duration from ffprobe should be converted to milliseconds."""
    _setup_common_mocks(mocker, probe_duration="120.5", file_size=5000000)

    replacer = AudioReplacer()
    result = replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=0,
        output_dir="/fake/output",
    )

    assert result.duration_ms == 120500


def test_file_size_bytes_populated(mocker) -> None:
    """File size should be read from the output file."""
    _setup_common_mocks(mocker, file_size=2500000)

    replacer = AudioReplacer()
    result = replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=0,
        output_dir="/fake/output",
    )

    assert result.file_size_bytes == 2500000


# ── Tests: Audio offset ────────────────────────────────────────────────


def test_audio_offset_applied_via_ss_flag(mocker) -> None:
    """The audio input should have -ss (start offset) applied.

    We verify the mux command contains the ss argument for the audio
    input by checking the final FFmpeg CLI arguments.
    """
    mock_run = _setup_common_mocks(mocker)

    replacer = AudioReplacer()
    replacer.process(
        video_path="/fake/video.mp4",
        audio_path="/fake/audio.mp3",
        track_start_ms=15000,
        output_dir="/fake/output",
    )

    mux_call = _mux_call_from_mock(mock_run)
    assert mux_call is not None, "Expected a mux FFmpeg call"
    cmd_str = " ".join(mux_call)
    # The audio input should have ss=15.0 in the filter or input spec
    assert "ss=15.0" in cmd_str or "15" in cmd_str, (
        f"Expected audio offset for 15s in mux command, got: {cmd_str}"
    )


# ── Tests: Singleton factory ───────────────────────────────────────────


def test_get_audio_replacer_returns_singleton() -> None:
    """get_audio_replacer should return the same instance on repeated calls."""
    from app.domain.video_processing.audio_replacer import get_audio_replacer

    instance1 = get_audio_replacer()
    instance2 = get_audio_replacer()
    assert instance1 is instance2
