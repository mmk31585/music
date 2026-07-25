# ── Moja ML Service — Dockerfile ───────────────────────────────────────
#
# Design decision: single-stage build.
#
# Reasoning:
#   - faster-whisper needs Cython-compiled wheels (ctranslate2, onnxruntime)
#     at runtime. A multi-stage build would copy those wheels + the venv
#     from builder → runtime, which adds complexity with uv's layout.
#     The image size difference (~300-400 MB) is not significant enough
#     to justify the extra maintenance surface for a single-service
#     deployment on a VPS.
#   - If image size ever becomes a concern, switch to multi-stage:
#     builder copies uv + deps, runtime copies /app/.venv from builder.

FROM python:3.12-slim

# ── System dependencies ───────────────────────────────────────────────
# faster-whisper / ctranslate2 need libgomp (OpenMP) at runtime.
# ffmpeg is needed for audio decoding.
RUN apt-get update && apt-get install -y --no-install-recommends \
    libgomp1 \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# ── Install uv (pinned version for reproducibility) ───────────────────
COPY --from=ghcr.io/astral-sh/uv:0.6.3 /uv /uvx /bin/

# ── Application code ──────────────────────────────────────────────────
WORKDIR /app

# Copy dependency manifests first for layer caching
COPY pyproject.toml uv.lock ./

# Install dependencies (no dev group in production)
RUN uv sync --frozen --no-dev

# Copy the rest of the application
COPY . .

# ── Non-root user for security ─────────────────────────────────────────
RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser && \
    mkdir -p /app/.cache && \
    chown -R appuser:appgroup /app
ENV HOME=/app UV_NO_SYNC=1
USER appuser

# ── Runtime ───────────────────────────────────────────────────────────
EXPOSE 8000
HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
  CMD python -c "import urllib.request; urllib.request.urlopen('http://localhost:8000/api/v1/health/live')" || exit 1

CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
