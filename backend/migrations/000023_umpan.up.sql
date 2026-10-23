CREATE TYPE jenis_sumber_umpan AS ENUM ('resmi_dprd', 'media');
CREATE TYPE dasar_cocok_umpan AS ENUM ('nama_dan_dprd', 'nama_dan_fraksi', 'situs_resmi');

CREATE TABLE umpan_kegiatan (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    jenis_sumber jenis_sumber_umpan NOT NULL,
    media_nama text NOT NULL,
    media_domain text NOT NULL,
    judul text NOT NULL,
    cuplikan text NOT NULL DEFAULT '' CHECK (char_length(cuplikan) <= 300),
    ringkasan text CHECK (char_length(ringkasan) <= 160),
    url text NOT NULL UNIQUE,
    terbit_pada timestamptz NOT NULL,
    tag_isu text[] NOT NULL DEFAULT '{}',
    disembunyikan boolean NOT NULL DEFAULT false,
    alasan_disembunyikan text,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX umpan_kegiatan_jenis_terbit_idx ON umpan_kegiatan (jenis_sumber, terbit_pada DESC);

CREATE TABLE umpan_anggota (
    umpan_id bigint NOT NULL REFERENCES umpan_kegiatan (id) ON DELETE CASCADE,
    anggota_id text NOT NULL REFERENCES anggota (id) ON DELETE CASCADE,
    dasar_cocok dasar_cocok_umpan NOT NULL,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (umpan_id, anggota_id)
);

CREATE INDEX umpan_anggota_anggota_id_idx ON umpan_anggota (anggota_id);
