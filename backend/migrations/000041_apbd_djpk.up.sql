ALTER TYPE sumber_data ADD VALUE IF NOT EXISTS 'djpk';

CREATE TABLE apbd_akun (
    tahun int NOT NULL,
    urutan int NOT NULL,
    nama text NOT NULL,
    anggaran numeric(20, 2) NOT NULL,
    realisasi numeric(20, 2) NOT NULL,
    bulan_realisasi int NOT NULL CHECK (bulan_realisasi BETWEEN 1 AND 12),
    provenance_id text NOT NULL REFERENCES provenance (id),
    PRIMARY KEY (tahun, urutan)
);
