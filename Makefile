-include .env

DATABASE_URL ?= postgres://rantai:rantai@localhost:5433/rantai?sslmode=disable
TEST_DATABASE_URL ?= $(DATABASE_URL)
export DATABASE_URL TEST_DATABASE_URL

COMPOSE := docker compose
MIGRATE := $(COMPOSE) run --rm --build migrate

.PHONY: up down migrate migrate-down lint test impor-wilayah impor-tema impor-kpu-dprd

up:
	$(COMPOSE) up -d --wait db minio

down:
	$(COMPOSE) down

migrate:
	@ls backend/migrations/*.up.sql >/dev/null 2>&1 || { echo "belum ada migrasi"; exit 0; }; $(MIGRATE) up

migrate-down:
	$(MIGRATE) down -all

lint:
	cd backend && golangci-lint run ./... && go run ./tools/cekkomentar

test:
	cd backend && go test ./...

impor-wilayah:
	cd backend && go run ./cmd/worker wilayah

impor-tema:
	cd backend && go run ./cmd/worker tema

impor-kpu-dprd:
	cd backend && go run ./cmd/worker kpu_dprd
