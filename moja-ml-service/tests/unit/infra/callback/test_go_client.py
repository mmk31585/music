"""Tests for ``GoCallbackClient``.

Uses ``respx`` to mock HTTP responses.
"""

from __future__ import annotations

import hashlib
import hmac
import json

import pytest
import respx

from app.infra.callback.go_client import (
    CallbackRejectedError,
    CallbackTransientError,
    GoCallbackClient,
)

# Shared test data
TEST_CALLBACK_URL = "http://go-backend:8080/api/v1"
TEST_HMAC_SECRET = "test-secret-key"
TRACK_ID = "track-abc-123"
LRC_CONTENT = "[00:01.00]Hello world\n[00:04.50]Test line"
PLAIN_TEXT = "Hello world\nTest line"
CONFIDENCE = 0.95
DETECTED_LANG = "fa"
WHISPER_VERSION = "medium-int8"

# The expected callback path (appended by the client)
EXPECTED_PATH = "/internal/v1/lyrics/callback"


# ── Test helpers ────────────────────────────────────────────────────────


def _expected_endpoint() -> str:
    return f"{TEST_CALLBACK_URL}{EXPECTED_PATH}"


def _build_client(
    callback_url: str = TEST_CALLBACK_URL,
    hmac_secret: str = TEST_HMAC_SECRET,
) -> GoCallbackClient:
    return GoCallbackClient(callback_url, hmac_secret)


# ── HMAC signing tests ─────────────────────────────────────────────────


class TestSigning:
    """Tests for the HMAC signing method."""

    def test_signature_deterministic(self) -> None:
        """Same payload + secret + timestamp always produces the same sig."""
        client = _build_client()
        sig1 = client._sign_payload("1712345678", b'{"key":"value"}')
        sig2 = client._sign_payload("1712345678", b'{"key":"value"}')
        assert sig1 == sig2

    def test_signature_changes_with_different_secret(self) -> None:
        """Different secret produces a different signature."""
        client_a = _build_client(hmac_secret="secret_a")
        client_b = _build_client(hmac_secret="secret_b")
        sig_a = client_a._sign_payload("1712345678", b"test")
        sig_b = client_b._sign_payload("1712345678", b"test")
        assert sig_a != sig_b

    def test_signature_changes_with_different_timestamp(self) -> None:
        """Different timestamp produces a different signature (replay protection)."""
        client = _build_client()
        sig1 = client._sign_payload("1712345678", b"test")
        sig2 = client._sign_payload("1712345679", b"test")
        assert sig1 != sig2

    def test_signature_changes_with_different_body(self) -> None:
        """Different body produces a different signature."""
        client = _build_client()
        sig1 = client._sign_payload("1712345678", b'{"a":1}')
        sig2 = client._sign_payload("1712345678", b'{"a":2}')
        assert sig1 != sig2

    def test_signature_matches_external_computation(self) -> None:
        """The signing algorithm matches a reference HMAC computation."""
        client = _build_client()
        timestamp = "1712345678"
        body = b'{"msg":"test"}'
        sig = client._sign_payload(timestamp, body)

        # Recompute via pure stdlib — should match.
        message = f"{timestamp}.{body.decode()}".encode()
        expected = hmac.new(
            TEST_HMAC_SECRET.encode(),
            message,
            hashlib.sha256,
        ).hexdigest()

        assert sig == expected


# ── Delivery tests (with respx) ─────────────────────────────────────────


class TestDeliverLyricsResult:
    """Tests for the full callback delivery using mocked HTTP."""

    def test_returns_true_on_200(self) -> None:
        """A 200 response returns True."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            route = respx.post(endpoint).respond(200)

            result = client.deliver_lyrics_result(
                track_id=TRACK_ID,
                lrc_content=LRC_CONTENT,
                plain_text=PLAIN_TEXT,
                confidence=CONFIDENCE,
                detected_language=DETECTED_LANG,
                whisper_model_version=WHISPER_VERSION,
            )

        assert result is True
        assert route.called

    def test_returns_true_on_204(self) -> None:
        """A 204 response also returns True."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            route = respx.post(endpoint).respond(204)

            result = client.deliver_lyrics_result(
                track_id=TRACK_ID,
                lrc_content=LRC_CONTENT,
                plain_text=PLAIN_TEXT,
                confidence=CONFIDENCE,
                detected_language=DETECTED_LANG,
                whisper_model_version=WHISPER_VERSION,
            )

        assert result is True
        assert route.called

    def test_sends_hmac_and_timestamp_headers(self) -> None:
        """Request includes X-Moja-Signature and X-Moja-Timestamp headers."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            route = respx.post(endpoint).respond(200)
            client.deliver_lyrics_result(
                track_id=TRACK_ID,
                lrc_content=LRC_CONTENT,
                plain_text=PLAIN_TEXT,
                confidence=CONFIDENCE,
                detected_language=DETECTED_LANG,
                whisper_model_version=WHISPER_VERSION,
            )

        assert route.called
        request = route.calls[0].request
        assert "X-Moja-Signature" in request.headers
        assert "X-Moja-Timestamp" in request.headers
        # Timestamp should be a numeric string (unix epoch seconds).
        assert request.headers["X-Moja-Timestamp"].isdigit()

    def test_sends_expected_body(self) -> None:
        """Request body contains the expected JSON fields."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            route = respx.post(endpoint).respond(200)
            client.deliver_lyrics_result(
                track_id=TRACK_ID,
                lrc_content=LRC_CONTENT,
                plain_text=PLAIN_TEXT,
                confidence=CONFIDENCE,
                detected_language=DETECTED_LANG,
                whisper_model_version=WHISPER_VERSION,
            )

        request = route.calls[0].request
        body = json.loads(request.content)
        assert body["track_id"] == TRACK_ID
        assert body["lrc_content"] == LRC_CONTENT
        assert body["confidence"] == CONFIDENCE
        assert body["detected_language"] == DETECTED_LANG
        assert body["whisper_model_version"] == WHISPER_VERSION

    def test_raises_rejected_on_400(self) -> None:
        """A 400 response raises CallbackRejectedError."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            respx.post(endpoint).respond(400, text="Bad request")

            with pytest.raises(CallbackRejectedError, match="400"):
                client.deliver_lyrics_result(
                    track_id=TRACK_ID,
                    lrc_content=LRC_CONTENT,
                    plain_text=PLAIN_TEXT,
                    confidence=CONFIDENCE,
                    detected_language=DETECTED_LANG,
                    whisper_model_version=WHISPER_VERSION,
                )

    def test_raises_rejected_on_401(self) -> None:
        """A 401 (bad HMAC) raises CallbackRejectedError."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            respx.post(endpoint).respond(401, text="Invalid signature")

            with pytest.raises(CallbackRejectedError, match="401"):
                client.deliver_lyrics_result(
                    track_id=TRACK_ID,
                    lrc_content=LRC_CONTENT,
                    plain_text=PLAIN_TEXT,
                    confidence=CONFIDENCE,
                    detected_language=DETECTED_LANG,
                    whisper_model_version=WHISPER_VERSION,
                )

    def test_raises_transient_on_500(self) -> None:
        """A 500 response raises CallbackTransientError."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            respx.post(endpoint).respond(500)

            with pytest.raises(CallbackTransientError, match="500"):
                client.deliver_lyrics_result(
                    track_id=TRACK_ID,
                    lrc_content=LRC_CONTENT,
                    plain_text=PLAIN_TEXT,
                    confidence=CONFIDENCE,
                    detected_language=DETECTED_LANG,
                    whisper_model_version=WHISPER_VERSION,
                )

    def test_raises_transient_on_503(self) -> None:
        """A 503 response raises CallbackTransientError."""
        client = _build_client()
        endpoint = _expected_endpoint()

        with respx.mock:
            respx.post(endpoint).respond(503)

            with pytest.raises(CallbackTransientError, match="503"):
                client.deliver_lyrics_result(
                    track_id=TRACK_ID,
                    lrc_content=LRC_CONTENT,
                    plain_text=PLAIN_TEXT,
                    confidence=CONFIDENCE,
                    detected_language=DETECTED_LANG,
                    whisper_model_version=WHISPER_VERSION,
                )
