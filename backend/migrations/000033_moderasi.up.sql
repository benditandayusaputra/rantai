ALTER TABLE laporan_foto ADD COLUMN disembunyikan boolean NOT NULL DEFAULT false;

CREATE TABLE moderasi_log (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    moderator_id bigint NOT NULL REFERENCES pengguna (id),
    target_jenis text NOT NULL CHECK (target_jenis IN ('laporan', 'laporan_foto', 'umpan')),
    target_id bigint NOT NULL,
    aksi text NOT NULL,
    alasan text NOT NULL DEFAULT '',
    dibuat_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX moderasi_log_target_idx ON moderasi_log (target_jenis, target_id, dibuat_pada);
