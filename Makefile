-include .env

DB_LOKAL := postgres://rantai:rantai@localhost:5433/rantai?sslmode=disable
DB_UJI := postgres://rantai:rantai@localhost:5433/rantai_uji?sslmode=disable
DATABASE_URL ?= $(DB_LOKAL)
TEST_DATABASE_URL ?= $(DB_UJI)
MINIO_ROOT_USER ?= rantai
MINIO_ROOT_PASSWORD ?= rantai-rahasia
export DATABASE_URL DATABASE_URL_MIGRASI TEST_DATABASE_URL MINIO_ROOT_USER MINIO_ROOT_PASSWORD

COMPOSE := docker compose
MIGRATE := $(COMPOSE) run --rm --build migrate
KLIEN_PG := $(COMPOSE) run --rm -T klien-pg sh -c

.PHONY: up down migrate siapkan-uji e2e migrate-down lint test impor-wilayah impor-tema impor-kpu-dprd impor-lpse impor-sirup impor-opentender turunan laporan-geocode snapshot-paket muat-snapshot-paket impor-umpan-dprd impor-umpan-media ringkas-umpan impor-sekolah impor-daya-tampung hitung-calon-lulusan hitung-indikator impor-paripurna kunci-vapid ringkasan-mingguan snapshot-demo seed-demo tanya-demo

up:
	$(COMPOSE) up -d --wait db minio

down:
	$(COMPOSE) down

migrate:
	@ls backend/migrations/*.up.sql >/dev/null 2>&1 || { echo "belum ada migrasi"; exit 0; }; $(MIGRATE) up

siapkan-uji:
	$(COMPOSE) exec -T db sh -c 'createdb -U rantai rantai_uji 2>/dev/null || true'
	$(MAKE) migrate DATABASE_URL_MIGRASI='postgres://rantai:rantai@db:5432/rantai_uji?sslmode=disable'

migrate-down:
	$(MIGRATE) down -all

lint:
	cd backend && golangci-lint run ./... && go run ./tools/cekkomentar
	cd frontend && npm run lint

test:
	cd backend && go test ./...
	cd frontend && npm test

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
	$(KLIEN_PG) 'pg_dump "$$DATABASE_URL_MIGRASI" --data-only --inserts --on-conflict-do-nothing $(TABEL_PAKET)' | gzip > $(SNAPSHOT_PAKET)

muat-snapshot-paket:
	gunzip -c $(SNAPSHOT_PAKET) | $(KLIEN_PG) 'psql "$$DATABASE_URL_MIGRASI" -q -v ON_ERROR_STOP=1'

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

hitung-calon-lulusan:
	cd backend && go run ./cmd/worker calon_lulusan

hitung-indikator:
	cd backend && go run ./cmd/worker indikator

impor-paripurna:
	cd backend && go run ./cmd/worker paripurna

kunci-vapid:
	cd backend && go run ./tools/kuncivapid

ringkasan-mingguan:
	cd backend && go run ./cmd/worker ringkasan_mingguan

SNAPSHOT_DEMO := data/snapshot/rantai.sql.gz
TABEL_AKUN_DAN_WARGA := -T pengguna -T otp -T sesi_refresh -T langganan -T notifikasi -T ringkasan_email -T moderasi_log \
	-T laporan -T laporan_foto -T laporan_konfirmasi -T laporan_riwayat -T hak_jawab

snapshot-demo:
	$(KLIEN_PG) 'pg_dump "$$DATABASE_URL_MIGRASI" --data-only --inserts --on-conflict-do-nothing -n public -T schema_migrations -T spatial_ref_sys $(TABEL_AKUN_DAN_WARGA)' | gzip > $(SNAPSHOT_DEMO)

tanya-demo:
	cd backend && go run ./tools/rekamtanya

seed-demo:
	@test -f $(SNAPSHOT_DEMO) || { echo "$(SNAPSHOT_DEMO) belum ada; buat dengan make snapshot-demo di mesin yang datanya lengkap"; exit 1; }
	{ echo 'SET session_replication_role = replica;'; gunzip -c $(SNAPSHOT_DEMO) | sed '/^SET transaction_timeout/d'; } | \
		$(KLIEN_PG) 'psql "$$DATABASE_URL_MIGRASI" -q -v ON_ERROR_STOP=1'
	cd backend && go run ./cmd/worker skenario_demo

e2e:
	cd frontend && npm run e2e
