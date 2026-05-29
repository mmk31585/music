APP_NAME=musicapp
DOCKER_COMPOSE=docker compose -f deployments/docker-compose.yml

run:
	@air

tidy:
	go mod tidy

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

fmt:
	go fmt ./...

infra-up:
	$(DOCKER_COMPOSE) up -d

infra-down:
	$(DOCKER_COMPOSE) down

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
