# syntax=docker/dockerfile:1.7
FROM golang:1.25.9-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/api ./cmd/api

FROM alpine:3.23
RUN apk add --no-cache ca-certificates ffmpeg tzdata
COPY --from=builder /bin/api /app/api
COPY --from=builder /app/migrations /app/migrations
WORKDIR /app
EXPOSE 8080
ENTRYPOINT ["./api"]
