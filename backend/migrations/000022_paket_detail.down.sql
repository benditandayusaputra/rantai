ALTER TABLE paket DROP CONSTRAINT IF EXISTS paket_jenis_kode_spse_unik;
ALTER TABLE paket DROP COLUMN IF EXISTS detail_diambil_pada;
ALTER TABLE paket DROP COLUMN IF EXISTS kode_spse;
