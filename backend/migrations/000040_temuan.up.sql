CREATE TYPE sikap_temuan AS ENUM ('mendukung', 'membantah');

CREATE TABLE temuan (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id),
    paket_kode text NOT NULL REFERENCES paket (kode),
    isi text NOT NULL CHECK (char_length(isi) BETWEEN 20 AND 1000),
    tautan_sumber text NOT NULL CHECK (char_length(tautan_sumber) <= 500),
    disembunyikan boolean NOT NULL DEFAULT false,
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX temuan_paket_idx ON temuan (paket_kode, dibuat_pada);
CREATE INDEX temuan_pengguna_dibuat_idx ON temuan (pengguna_id, dibuat_pada);

CREATE TABLE temuan_tanggapan (
    temuan_id bigint NOT NULL REFERENCES temuan (id) ON DELETE CASCADE,
    pengguna_id bigint NOT NULL REFERENCES pengguna (id),
    sikap sikap_temuan NOT NULL,
    tautan_sumber text NOT NULL DEFAULT '' CHECK (char_length(tautan_sumber) <= 500),
    dibuat_pada timestamptz NOT NULL DEFAULT now(),
    diubah_pada timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (temuan_id, pengguna_id)
);

ALTER TABLE moderasi_log DROP CONSTRAINT moderasi_log_target_jenis_check;
ALTER TABLE moderasi_log ADD CONSTRAINT moderasi_log_target_jenis_check
    CHECK (target_jenis IN ('laporan', 'laporan_foto', 'umpan', 'temuan'));
