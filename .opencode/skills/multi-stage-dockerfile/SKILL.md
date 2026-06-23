---
name: multi-stage-dockerfile
description: 'Create optimized multi-stage Dockerfiles for any language or framework. Focus on smaller, more secure container images.'
---

# Multi-Stage Dockerfile

Create efficient multi-stage Dockerfiles that follow best practices for smaller, more secure container images.

## Structure
- Builder stage for compilation and dependency installation
- Runtime stage with only what's needed to run
- Copy only necessary artifacts from builder to runtime
- Use meaningful stage names with `AS` keyword

## Base Images
- Official, minimal base images with exact version tags
- Consider distroless for runtime stages
- Alpine-based for smaller footprint
- Ensure minimal runtime dependencies

## Layer Optimization
- Organize commands to maximize caching
- Frequent changes after infrequent changes
- Use `.dockerignore` to exclude unnecessary files
- Combine RUN commands with `&&`

## Security
- Use `USER` for non-root containers
- Remove build tools from final image
- Scan final image for vulnerabilities

## Muse-Specific
- Go API: `golang:1.22-alpine AS builder` → `alpine:3.19 AS runtime`
- Vue frontend: `node:20-alpine AS build` → `nginx:alpine AS serve`
- Python ML: `python:3.11-slim AS base` → `python:3.11-slim AS runtime`
