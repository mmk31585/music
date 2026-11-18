# Root image builds the backend API (see docker-compose for full stack)
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /out/api ./api
EXPOSE 8080
CMD ["./api"]
EOF
echo done