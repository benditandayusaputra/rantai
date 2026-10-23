ALTER TYPE sumber_data ADD VALUE IF NOT EXISTS 'disdik';

CREATE TABLE daya_tampung (
    npsn text NOT NULL REFERENCES sekolah (npsn),
    tahun_ajaran text NOT NULL,
    jalur text NOT NULL CHECK (jalur IN ('negeri', 'spmb_bersama')),
    kursi integer NOT NULL CHECK (kursi >= 0),
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (npsn, tahun_ajaran, jalur)
);
