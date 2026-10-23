COMPOSE := docker compose
MIGRATE := $(COMPOSE) run --rm --build migrate

.PHONY: up down migrate migrate-down lint test

up:
	$(COMPOSE) up -d --wait db minio

down:
	$(COMPOSE) down

migrate:
	@ls backend/migrations/*.up.sql >/dev/null 2>&1 || { echo "belum ada migrasi"; exit 0; }; $(MIGRATE) up

migrate-down:
	$(MIGRATE) down -all

lint:
	cd backend && gofmt -l . | (! grep .) && go vet ./...

test:
	cd backend && go test ./...
