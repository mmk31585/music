# syntax=docker/dockerfile:1.7
FROM golang:1.25.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Build arg selects which binary to compile: "api" (default) or "worker".
ARG CMD_DIR=api
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/server ./cmd/${CMD_DIR}

FROM alpine:3.23
RUN apk add --no-cache ca-certificates ffmpeg tzdata curl \
    && adduser -D -u 1001 appuser
# Internet imports shell out to yt-dlp. Use the self-contained musllinux build
# (bundles Python) so the runtime does not need python3/apk at build time.
RUN wget -q -O /usr/local/bin/yt-dlp \
        "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_musllinux" \
    && chmod +x /usr/local/bin/yt-dlp
# Node is yt-dlp's JS runtime for YouTube. Copy it from the official image
# (plus the two C++ libs it links against) to avoid depending on apk.
COPY --from=node:22-alpine /usr/local/bin/node /usr/local/bin/node
COPY --from=node:22-alpine /usr/lib/libstdc++.so.6* /usr/lib/
COPY --from=node:22-alpine /usr/lib/libgcc_s.so.1 /usr/lib/
COPY --from=builder /bin/server /app/server
COPY --from=builder /app/migrations /app/migrations
WORKDIR /app
# /app/uploads exists in the image so that when a named volume is mounted there,
# Docker initializes it with appuser ownership (the app runs as appuser).
RUN mkdir -p /app/uploads && chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -sf http://localhost:8080/api/v1/health || exit 1
ENTRYPOINT ["./server"]
