CREATE TYPE kondisi_laporan AS ENUM ('belum_dimulai', 'sedang_dikerjakan', 'terhenti', 'selesai_baik', 'selesai_rusak');

CREATE TYPE status_laporan AS ENUM ('DIKIRIM', 'MENUNGGU_SAKSI', 'TERVERIFIKASI', 'DITERUSKAN', 'DITANGGAPI', 'SELESAI', 'DITOLAK');

CREATE TABLE laporan (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id),
    paket_kode text REFERENCES paket (kode),
    kode_wilayah text NOT NULL REFERENCES wilayah (kode),
    lokasi_publik geometry(Point, 4326) NOT NULL,
    lokasi_asli geometry(Point, 4326) NOT NULL,
    akurasi_meter numeric(7, 1) NOT NULL CHECK (akurasi_meter >= 0),
    kondisi kondisi_laporan NOT NULL,
    catatan text NOT NULL DEFAULT '' CHECK (char_length(catatan) <= 500),
    status status_laporan NOT NULL DEFAULT 'DIKIRIM',
    is_simulation boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX laporan_pengguna_dibuat_idx ON laporan (pengguna_id, dibuat_pada);
CREATE INDEX laporan_paket_idx ON laporan (paket_kode);
CREATE INDEX laporan_lokasi_publik_idx ON laporan USING gist (lokasi_publik);

CREATE TABLE laporan_foto (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    laporan_id bigint NOT NULL REFERENCES laporan (id) ON DELETE CASCADE,
    urutan smallint NOT NULL,
    object_key text NOT NULL UNIQUE,
    thumbnail_key text NOT NULL UNIQUE,
    lebar integer NOT NULL,
    tinggi integer NOT NULL,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (laporan_id, urutan)
);

CREATE TABLE laporan_konfirmasi (
    laporan_id bigint NOT NULL REFERENCES laporan (id) ON DELETE CASCADE,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id),
    kondisi kondisi_laporan NOT NULL,
    lokasi_asli geometry(Point, 4326) NOT NULL,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (laporan_id, pengguna_id)
);

CREATE TABLE laporan_riwayat (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    laporan_id bigint NOT NULL REFERENCES laporan (id) ON DELETE CASCADE,
    status_lama status_laporan,
    status_baru status_laporan NOT NULL,
    oleh_pengguna_id bigint REFERENCES pengguna (id),
    alasan text NOT NULL DEFAULT '',
    dibuat_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX laporan_riwayat_laporan_idx ON laporan_riwayat (laporan_id, dibuat_pada);
