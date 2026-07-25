APP_NAME=musicapp
DOCKER_COMPOSE=docker compose --env-file .env -f deployments/docker-compose.yml

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
	go test ./... -race -count=1

test-coverage:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic -count=1

cover-html: test-coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: file://$(PWD)/coverage.html"

cover-ci:
	go test ./... -race -coverprofile=coverage.out -covermode=atomic -count=1
	go tool cover -func=coverage.out | tail -1 | awk '{print $$NF}' | sed 's/%//' | xargs -I{} bash -c 'if [ "{}" \< "30.0" ]; then echo "FAIL: Coverage {}% < 30.0% threshold"; exit 1; else echo "PASS: Coverage {}% >= 30.0%"; fi'

lint-go:
	golangci-lint run ./... --timeout 5m

lint-frontend:
	cd frontend && npm run lint

fmt:
	go fmt ./...

infra-up:
	$(DOCKER_COMPOSE) up -d postgres redis minio

infra-down:
	$(DOCKER_COMPOSE) down

infra-restart:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up -d postgres redis

infra-logs:
	$(DOCKER_COMPOSE) logs -f

docker-up:
	$(DOCKER_COMPOSE) up -d --build

docker-down:
	$(DOCKER_COMPOSE) down

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-rebuild:
	$(DOCKER_COMPOSE) build --no-cache api
	$(DOCKER_COMPOSE) build --no-cache frontend

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
