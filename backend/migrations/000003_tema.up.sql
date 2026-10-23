CREATE TABLE tema (
    kode text PRIMARY KEY,
    nama_id text NOT NULL,
    nama_en text NOT NULL,
    aktif boolean NOT NULL DEFAULT false,
    urutan integer NOT NULL,
    konfigurasi jsonb NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE indikator_kecamatan (
    tema_kode text NOT NULL REFERENCES tema (kode),
    kode_kecamatan text NOT NULL REFERENCES wilayah (kode),
    periode text NOT NULL,
    kebutuhan numeric NOT NULL,
    kapasitas numeric NOT NULL,
    kesenjangan numeric NOT NULL,
    rasio numeric,
    kelas text NOT NULL,
    label label_kepastian NOT NULL,
    rincian jsonb NOT NULL DEFAULT '{}',
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tema_kode, kode_kecamatan, periode)
);
