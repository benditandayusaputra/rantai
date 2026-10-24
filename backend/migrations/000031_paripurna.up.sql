CREATE TYPE agenda_sikap_fraksi AS ENUM ('raperda', 'apbd');

CREATE TABLE sikap_fraksi (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    fraksi text NOT NULL,
    agenda agenda_sikap_fraksi NOT NULL,
    judul_agenda text NOT NULL,
    tanggal date NOT NULL,
    ringkasan text NOT NULL,
    url text NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (fraksi, url)
);

CREATE INDEX sikap_fraksi_fraksi_tanggal_idx ON sikap_fraksi (fraksi, tanggal DESC);

CREATE TABLE aspirasi_reses (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    masa_reses text NOT NULL,
    tanggal_paripurna date NOT NULL,
    bidang text NOT NULL CHECK (bidang <> '' AND bidang = lower(bidang)),
    uraian text NOT NULL,
    kode_wilayah text REFERENCES wilayah (kode),
    url text NOT NULL,
    provenance_id text NOT NULL REFERENCES provenance (id),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    UNIQUE (url, uraian)
);

CREATE INDEX aspirasi_reses_bidang_idx ON aspirasi_reses (bidang);
