CREATE TABLE komisi (
    kode text PRIMARY KEY CHECK (kode IN ('A', 'B', 'C', 'D', 'E')),
    nama text NOT NULL,
    bidang text[] NOT NULL CHECK (cardinality(bidang) > 0),
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE anggota ADD COLUMN komisi_kode text REFERENCES komisi (kode);

CREATE TABLE anggota_akd (
    anggota_id text NOT NULL REFERENCES anggota (id) ON DELETE CASCADE,
    alat_kelengkapan text NOT NULL,
    jabatan text NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (anggota_id, alat_kelengkapan)
);
