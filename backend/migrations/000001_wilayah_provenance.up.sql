CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TYPE tingkat_wilayah AS ENUM ('provinsi', 'kota', 'kecamatan', 'kelurahan');
CREATE TYPE label_kepastian AS ENUM ('resmi', 'proksi', 'estimasi', 'simulasi');
CREATE TYPE sumber_data AS ENUM ('kemendagri', 'lpse', 'sirup', 'dapodik', 'kpu', 'dprd', 'opentender', 'juknis_spmb', 'rantai', 'simulasi');
CREATE TYPE status_ingest AS ENUM ('berjalan', 'berhasil', 'gagal');

CREATE TABLE provenance (
    id text PRIMARY KEY,
    sumber sumber_data NOT NULL,
    url text,
    diambil_pada timestamptz NOT NULL,
    label label_kepastian NOT NULL,
    hash_konten text,
    catatan text,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE wilayah (
    kode text PRIMARY KEY,
    tingkat tingkat_wilayah NOT NULL,
    nama text NOT NULL,
    induk_kode text REFERENCES wilayah (kode),
    alias text[] NOT NULL DEFAULT '{}',
    geom geometry(MultiPolygon, 4326) NOT NULL,
    titik_tengah geometry(Point, 4326) NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX wilayah_geom_idx ON wilayah USING gist (geom);
CREATE INDEX wilayah_induk_kode_idx ON wilayah (induk_kode);

CREATE TABLE ingest_run (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    sumber text NOT NULL,
    mulai timestamptz NOT NULL,
    selesai timestamptz,
    status status_ingest NOT NULL,
    jumlah_baru integer NOT NULL DEFAULT 0,
    jumlah_berubah integer NOT NULL DEFAULT 0,
    pesan_error text,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);
