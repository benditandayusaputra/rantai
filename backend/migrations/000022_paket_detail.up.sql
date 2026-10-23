ALTER TABLE paket ADD COLUMN kode_spse text NOT NULL;
ALTER TABLE paket ADD COLUMN detail_diambil_pada timestamptz;
ALTER TABLE paket ADD CONSTRAINT paket_jenis_kode_spse_unik UNIQUE (is_tender, kode_spse);
