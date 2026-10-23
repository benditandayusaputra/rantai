CREATE TABLE calon_lulusan_kecamatan (
    kode_kecamatan text NOT NULL REFERENCES wilayah (kode),
    tahun_ajaran text NOT NULL,
    jumlah integer NOT NULL CHECK (jumlah >= 0),
    label label_kepastian NOT NULL,
    metode text NOT NULL CHECK (metode IN ('per_sekolah', 'agregat_resmi', 'estimasi_total_sd')),
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (kode_kecamatan, tahun_ajaran)
);
