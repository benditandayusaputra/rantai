DELETE FROM calon_lulusan_kecamatan WHERE jenjang <> 'dasar' OR metode = 'estimasi_total_smp';
ALTER TABLE calon_lulusan_kecamatan DROP CONSTRAINT calon_lulusan_kecamatan_metode_check;
ALTER TABLE calon_lulusan_kecamatan ADD CONSTRAINT calon_lulusan_kecamatan_metode_check
    CHECK (metode IN ('per_sekolah', 'agregat_resmi', 'estimasi_total_sd'));
ALTER TABLE calon_lulusan_kecamatan DROP CONSTRAINT calon_lulusan_kecamatan_pkey;
ALTER TABLE calon_lulusan_kecamatan ADD PRIMARY KEY (kode_kecamatan, tahun_ajaran);
ALTER TABLE calon_lulusan_kecamatan DROP COLUMN jenjang;

DELETE FROM sekolah_rombel WHERE tingkat > 9;
ALTER TABLE sekolah_rombel DROP CONSTRAINT sekolah_rombel_tingkat_check;
ALTER TABLE sekolah_rombel ADD CONSTRAINT sekolah_rombel_tingkat_check CHECK (tingkat BETWEEN 1 AND 9);

DELETE FROM daya_tampung WHERE npsn IN (SELECT npsn FROM sekolah WHERE jenjang = 'menengah_atas');
DELETE FROM sekolah_rombel WHERE npsn IN (SELECT npsn FROM sekolah WHERE jenjang = 'menengah_atas');
DELETE FROM sekolah WHERE jenjang = 'menengah_atas';
ALTER TABLE sekolah DROP CONSTRAINT sekolah_jenjang_check;
ALTER TABLE sekolah ADD CONSTRAINT sekolah_jenjang_check CHECK (jenjang IN ('dasar', 'menengah_pertama'));
ALTER TABLE sekolah DROP CONSTRAINT sekolah_bentuk_check;
ALTER TABLE sekolah ADD CONSTRAINT sekolah_bentuk_check CHECK (bentuk IN ('SD', 'MI', 'SMP', 'MTs'));
