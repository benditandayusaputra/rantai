-include .env

DATABASE_URL ?= postgres://rantai:rantai@localhost:5433/rantai?sslmode=disable
TEST_DATABASE_URL ?= $(DATABASE_URL)
export DATABASE_URL TEST_DATABASE_URL

COMPOSE := docker compose
MIGRATE := $(COMPOSE) run --rm --build migrate

.PHONY: up down migrate migrate-down lint test impor-wilayah impor-tema impor-kpu-dprd impor-lpse impor-sirup impor-opentender turunan laporan-geocode snapshot-paket muat-snapshot-paket impor-umpan-dprd impor-umpan-media ringkas-umpan impor-sekolah impor-daya-tampung

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

impor-lpse:
	cd backend && go run ./cmd/worker lpse_daftar && go run ./cmd/worker lpse_detail

impor-sirup:
	cd backend && go run ./cmd/worker sirup

impor-opentender:
	cd backend && go run ./cmd/worker opentender

turunan:
	cd backend && go run ./cmd/worker turunan

laporan-geocode:
	cd backend && go run ./cmd/worker laporan_geocode

SNAPSHOT_PAKET := data/snapshot/paket.sql.gz
TABEL_PAKET := -t provenance -t instansi -t rup -t penyedia -t paket -t paket_riwayat_tahap

snapshot-paket:
	$(COMPOSE) exec -T db pg_dump -U rantai -d rantai --data-only --inserts --on-conflict-do-nothing $(TABEL_PAKET) | gzip > $(SNAPSHOT_PAKET)

muat-snapshot-paket:
	gunzip -c $(SNAPSHOT_PAKET) | $(COMPOSE) exec -T db psql -q -v ON_ERROR_STOP=1 -U rantai -d rantai

impor-umpan-dprd:
	cd backend && go run ./cmd/worker umpan_dprd

impor-umpan-media:
	cd backend && go run ./cmd/worker umpan_media

ringkas-umpan:
	cd backend && go run ./cmd/worker ringkas_umpan

impor-sekolah:
	cd backend && go run ./cmd/worker dapodik

impor-daya-tampung:
	cd backend && go run ./cmd/worker daya_tampung
