---
applyTo: '**/Dockerfile,**/Dockerfile.*,**/*.dockerfile,**/docker-compose*.yml,**/compose*.yml'
description: 'Best practices for creating optimized, secure Docker images for the Muse Go/Vue stack'
---

# Containerization & Docker Best Practices for Muse

## Core Principles

### 1. Multi-Stage Builds
Use multiple `FROM` instructions to separate build-time dependencies from runtime dependencies.

```dockerfile
# Stage 1: Build Go binary
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# Stage 2: Runtime
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
```

### 2. Choose Minimal Base Images
- Use `alpine` variants for Go services
- Use `node:20-alpine` for frontend builds
- Avoid `latest` tag; pin specific versions
- For frontend production: use `nginx:alpine` to serve static files

### 3. Optimize Image Layers
- Order Dockerfile instructions from least to most frequently changing
- Combine `RUN` commands with `&&` to minimize layers
- Clean up package manager caches in the same `RUN` command

```dockerfile
# GOOD
RUN apt-get update && \
    apt-get install -y --no-install-recommends some-package && \
    rm -rf /var/lib/apt/lists/*
```

### 4. Use .dockerignore
```dockerignore
.git*
node_modules
.env
*.log
dist/
tmp/
.git/
.idea/
.vscode/
```

## Security

### Non-Root User
```dockerfile
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
```

### Health Checks
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
```

### No Secrets in Layers
- Never `COPY` .env files into images
- Use Docker secrets or environment variables at runtime
- Build args for non-sensitive config only

## Muse-Specific Patterns

### Go API Service
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /api ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]
```

### Vue Frontend Service
```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### ML Worker Service
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["celery", "-A", "tasks", "worker", "--loglevel=info"]
```
