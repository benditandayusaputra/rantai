CREATE TABLE dapil (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nomor integer NOT NULL,
    nama text NOT NULL,
    tahun_pemilu integer NOT NULL,
    jumlah_kursi integer NOT NULL CHECK (jumlah_kursi > 0),
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (tahun_pemilu, nomor)
);

CREATE TABLE dapil_wilayah (
    dapil_id integer NOT NULL REFERENCES dapil (id) ON DELETE CASCADE,
    kode_kecamatan text NOT NULL UNIQUE REFERENCES wilayah (kode),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dapil_id, kode_kecamatan)
);

CREATE TABLE anggota (
    id text PRIMARY KEY,
    nama text NOT NULL,
    nama_lain text[] NOT NULL DEFAULT '{}',
    partai text NOT NULL,
    fraksi text NOT NULL,
    dapil_id integer NOT NULL REFERENCES dapil (id),
    url_profil text NOT NULL,
    lhkpn_url text NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX anggota_dapil_id_idx ON anggota (dapil_id);
