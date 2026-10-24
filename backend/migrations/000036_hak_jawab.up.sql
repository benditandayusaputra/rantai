ALTER TYPE sumber_data ADD VALUE IF NOT EXISTS 'instansi';

CREATE TYPE target_hak_jawab AS ENUM ('sinyal', 'laporan', 'paket', 'umpan');

CREATE TABLE hak_jawab (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    target_jenis target_hak_jawab NOT NULL,
    target_id text NOT NULL,
    instansi_id bigint REFERENCES instansi (id),
    anggota_id text REFERENCES anggota (id),
    isi text NOT NULL CHECK (char_length(isi) BETWEEN 1 AND 2000),
    tautan_dokumen text,
    penjawab_id bigint NOT NULL REFERENCES pengguna (id),
    kode_kecamatan text REFERENCES wilayah (kode),
    is_simulation boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    CHECK ((instansi_id IS NULL) <> (anggota_id IS NULL))
);

CREATE INDEX hak_jawab_target_idx ON hak_jawab (target_jenis, target_id, dibuat_pada);
CREATE INDEX hak_jawab_dibuat_pada_idx ON hak_jawab (dibuat_pada);

CREATE FUNCTION tolak_ubah_hak_jawab() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'hak jawab tidak bisa diubah atau dihapus';
END;
$$;

CREATE TRIGGER hak_jawab_tetap BEFORE UPDATE OR DELETE ON hak_jawab
    FOR EACH ROW EXECUTE FUNCTION tolak_ubah_hak_jawab();

ALTER TABLE paket
    ADD COLUMN nomor_berita_acara text CHECK (char_length(nomor_berita_acara) <= 100),
    ADD COLUMN serah_terima_provenance_id text REFERENCES provenance (id);
