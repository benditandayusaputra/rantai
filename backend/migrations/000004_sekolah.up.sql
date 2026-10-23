CREATE TABLE sekolah (
    npsn text PRIMARY KEY,
    nama text NOT NULL,
    bentuk text NOT NULL CHECK (bentuk IN ('SD', 'MI', 'SMP', 'MTs')),
    jenjang text NOT NULL CHECK (jenjang IN ('dasar', 'menengah_pertama')),
    status text NOT NULL CHECK (status IN ('negeri', 'swasta')),
    is_madrasah boolean NOT NULL,
    kode_kecamatan text NOT NULL REFERENCES wilayah (kode),
    kode_kelurahan text REFERENCES wilayah (kode),
    lokasi geometry(Point, 4326),
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX sekolah_lokasi_idx ON sekolah USING gist (lokasi);
CREATE INDEX sekolah_kode_kecamatan_idx ON sekolah (kode_kecamatan);

CREATE TABLE sekolah_rombel (
    npsn text NOT NULL REFERENCES sekolah (npsn),
    tahun_ajaran text NOT NULL,
    tingkat smallint NOT NULL CHECK (tingkat BETWEEN 1 AND 9),
    jumlah_rombel integer NOT NULL,
    jumlah_siswa integer NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (npsn, tahun_ajaran, tingkat)
);
