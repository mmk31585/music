APP_NAME=musicapp
DOCKER_COMPOSE=docker compose -f deployments/docker-compose.yml

run:
	@air -c .air.toml

run-worker:
	@air -c .air.worker.toml

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

build-worker:
	go build -o bin/worker ./cmd/worker

build-all: build build-worker

test:
	go test ./...

fmt:
	go fmt ./...

infra-up:
	$(DOCKER_COMPOSE) up -d

infra-down:
	$(DOCKER_COMPOSE) down

infra-restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up -d

infra-logs:
	$(DOCKER_COMPOSE) logs -f

migrate-up:
	./scripts/migrate.sh up

migrate-down:
	./scripts/migrate.sh down

migrate-status:
	./scripts/migrate.sh status

dev: infra-up
	go run ./cmd/api

dev-worker: infra-up
	go run ./cmd/worker

minio:
	@echo "MinIO API: http://localhost:9000"
	@echo "MinIO Console: http://localhost:9001"
	@echo "Username: minio"
	@echo "Password: minio123"
