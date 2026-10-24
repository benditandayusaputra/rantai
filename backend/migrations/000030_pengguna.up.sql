CREATE TYPE peran_pengguna AS ENUM ('warga', 'anggota', 'instansi', 'moderator', 'admin');

CREATE TABLE pengguna (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text NOT NULL UNIQUE CHECK (email = lower(email)),
    nama_tampilan text NOT NULL DEFAULT '',
    peran peran_pengguna NOT NULL DEFAULT 'warga',
    instansi_id bigint REFERENCES instansi (id),
    anggota_id text REFERENCES anggota (id),
    is_simulation boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    CHECK ((peran = 'instansi') = (instansi_id IS NOT NULL)),
    CHECK (anggota_id IS NULL OR peran = 'anggota')
);

CREATE TABLE otp (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text NOT NULL,
    kode_hash bytea NOT NULL,
    kedaluwarsa timestamptz NOT NULL,
    percobaan integer NOT NULL DEFAULT 0,
    dipakai boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL
);

CREATE INDEX otp_email_dibuat_pada_idx ON otp (email, dibuat_pada DESC);

CREATE TABLE sesi_refresh (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id) ON DELETE CASCADE,
    token_hash bytea NOT NULL UNIQUE,
    kedaluwarsa timestamptz NOT NULL,
    dicabut boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX sesi_refresh_pengguna_id_idx ON sesi_refresh (pengguna_id);
