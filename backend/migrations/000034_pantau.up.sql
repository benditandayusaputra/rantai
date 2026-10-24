CREATE TYPE kanal_langganan AS ENUM ('feed', 'webpush', 'email');

CREATE TYPE jenis_kejadian AS ENUM (
    'paket_baru', 'melewati_jadwal', 'pemeliharaan_berakhir', 'laporan_terverifikasi', 'hak_jawab', 'indikator_diperbarui'
);

CREATE TABLE langganan (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id) ON DELETE CASCADE,
    kode_wilayah text NOT NULL REFERENCES wilayah (kode),
    kanal kanal_langganan NOT NULL DEFAULT 'feed',
    webpush_endpoint jsonb,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (pengguna_id, kode_wilayah, kanal)
);

CREATE TABLE kejadian (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    jenis jenis_kejadian NOT NULL,
    kode_kecamatan text NOT NULL REFERENCES wilayah (kode),
    judul text NOT NULL,
    tautan text NOT NULL,
    kunci text NOT NULL UNIQUE,
    terjadi_pada timestamptz NOT NULL,
    is_simulation boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX kejadian_kecamatan_idx ON kejadian (kode_kecamatan, terjadi_pada DESC);

CREATE TABLE notifikasi (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id) ON DELETE CASCADE,
    kejadian_id bigint NOT NULL REFERENCES kejadian (id) ON DELETE CASCADE,
    dibaca_pada timestamptz,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (pengguna_id, kejadian_id)
);

CREATE INDEX notifikasi_pengguna_idx ON notifikasi (pengguna_id, dibuat_pada DESC);
