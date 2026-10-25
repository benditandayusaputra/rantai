ALTER TABLE sekolah DROP CONSTRAINT sekolah_bentuk_check;
ALTER TABLE sekolah ADD CONSTRAINT sekolah_bentuk_check CHECK (bentuk IN ('SD', 'MI', 'SMP', 'MTs', 'SMA', 'SMK', 'MA'));
ALTER TABLE sekolah DROP CONSTRAINT sekolah_jenjang_check;
ALTER TABLE sekolah ADD CONSTRAINT sekolah_jenjang_check CHECK (jenjang IN ('dasar', 'menengah_pertama', 'menengah_atas'));

ALTER TABLE sekolah_rombel DROP CONSTRAINT sekolah_rombel_tingkat_check;
ALTER TABLE sekolah_rombel ADD CONSTRAINT sekolah_rombel_tingkat_check CHECK (tingkat BETWEEN 1 AND 13);

ALTER TABLE calon_lulusan_kecamatan ADD COLUMN jenjang text NOT NULL DEFAULT 'dasar'
    CONSTRAINT calon_lulusan_kecamatan_jenjang_check CHECK (jenjang IN ('dasar', 'menengah_pertama'));
ALTER TABLE calon_lulusan_kecamatan ALTER COLUMN jenjang DROP DEFAULT;
ALTER TABLE calon_lulusan_kecamatan DROP CONSTRAINT calon_lulusan_kecamatan_pkey;
ALTER TABLE calon_lulusan_kecamatan ADD PRIMARY KEY (kode_kecamatan, tahun_ajaran, jenjang);
ALTER TABLE calon_lulusan_kecamatan DROP CONSTRAINT calon_lulusan_kecamatan_metode_check;
ALTER TABLE calon_lulusan_kecamatan ADD CONSTRAINT calon_lulusan_kecamatan_metode_check
    CHECK (metode IN ('per_sekolah', 'agregat_resmi', 'estimasi_total_sd', 'estimasi_total_smp'));
