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
COPY --from=builder /bin/server /app/server
COPY --from=builder /app/migrations /app/migrations
WORKDIR /app
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
ENTRYPOINT ["./server"]
