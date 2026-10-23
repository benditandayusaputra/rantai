CREATE TYPE jenis_instansi AS ENUM ('dinas', 'sudin', 'badan', 'kelurahan', 'lainnya');
CREATE TYPE jenis_sinyal AS ENUM ('S-01', 'S-02', 'S-03', 'S-04', 'S-05');

CREATE TABLE instansi (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama text NOT NULL UNIQUE,
    jenis jenis_instansi NOT NULL,
    kode_kota text REFERENCES wilayah (kode),
    induk_id bigint REFERENCES instansi (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE rup (
    kode text PRIMARY KEY,
    nama text NOT NULL,
    instansi_id bigint REFERENCES instansi (id),
    pagu numeric(18, 2),
    sumber_dana text,
    tahun_anggaran integer,
    metode text,
    jadwal_pelaksanaan_mulai date,
    jadwal_pelaksanaan_selesai date,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE penyedia (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nama text NOT NULL,
    nama_normal text NOT NULL UNIQUE,
    alias text[] NOT NULL DEFAULT '{}',
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE paket (
    kode text PRIMARY KEY,
    nama text NOT NULL,
    kode_rup text REFERENCES rup (kode),
    instansi_id bigint NOT NULL REFERENCES instansi (id),
    is_tender boolean NOT NULL,
    is_batal boolean NOT NULL DEFAULT false,
    jenis_pengadaan text NOT NULL,
    metode text NOT NULL,
    tahun_anggaran integer NOT NULL,
    pagu numeric(18, 2),
    hps numeric(18, 2),
    nilai_kontrak numeric(18, 2),
    jenis_kontrak text,
    tanggal_dibuat date,
    tanggal_kontrak date,
    tahap text NOT NULL,
    jumlah_peserta integer,
    penyedia_id bigint REFERENCES penyedia (id),
    jadwal_mulai date,
    jadwal_selesai date,
    tanggal_serah_terima date,
    kategori text,
    sektor text,
    tema_kode text[] NOT NULL DEFAULT '{}',
    lokasi_tingkat tingkat_wilayah,
    kode_wilayah text REFERENCES wilayah (kode),
    lokasi geometry(Point, 4326),
    rt integer,
    rw integer,
    is_simulation boolean NOT NULL DEFAULT false,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX paket_lokasi_idx ON paket USING gist (lokasi);
CREATE INDEX paket_kode_wilayah_idx ON paket (kode_wilayah);
CREATE INDEX paket_tahun_anggaran_idx ON paket (tahun_anggaran);
CREATE INDEX paket_tema_kode_idx ON paket USING gin (tema_kode);
CREATE INDEX paket_kategori_idx ON paket (kategori, jenis_pengadaan, tahun_anggaran);

CREATE TABLE paket_riwayat_tahap (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    paket_kode text NOT NULL REFERENCES paket (kode) ON DELETE CASCADE,
    tahap_lama text,
    tahap_baru text NOT NULL,
    dicatat_pada timestamptz NOT NULL,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX paket_riwayat_tahap_paket_idx ON paket_riwayat_tahap (paket_kode);

CREATE TABLE riwayat_penyedia (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    penyedia_id bigint NOT NULL REFERENCES penyedia (id),
    tahun integer NOT NULL,
    pembeli text NOT NULL,
    nama_paket text NOT NULL,
    nilai numeric(18, 2),
    sumber sumber_data NOT NULL,
    kode_sumber text NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (sumber, kode_sumber)
);

CREATE INDEX riwayat_penyedia_penyedia_idx ON riwayat_penyedia (penyedia_id, tahun);

CREATE TABLE sinyal (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    paket_kode text NOT NULL REFERENCES paket (kode) ON DELETE CASCADE,
    jenis jenis_sinyal NOT NULL,
    alasan text NOT NULL,
    pembanding jsonb NOT NULL DEFAULT '{}',
    aktif boolean NOT NULL DEFAULT true,
    dihitung_pada timestamptz NOT NULL,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (paket_kode, jenis)
);
