ALTER TABLE anggota ADD COLUMN foto_url text;

CREATE TABLE foto_sekolah (
    npsn text PRIMARY KEY REFERENCES sekolah (npsn) ON DELETE CASCADE,
    url text NOT NULL,
    sumber text NOT NULL CHECK (sumber IN ('sekolah_kita', 'wikimedia')),
    halaman text NOT NULL,
    kredit text NOT NULL DEFAULT '',
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);
