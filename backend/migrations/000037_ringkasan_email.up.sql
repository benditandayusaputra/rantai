CREATE TABLE ringkasan_email (
    pengguna_id bigint PRIMARY KEY REFERENCES pengguna (id) ON DELETE CASCADE,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    terakhir_dikirim timestamptz
);
