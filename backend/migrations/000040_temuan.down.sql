DELETE FROM moderasi_log WHERE target_jenis = 'temuan';
ALTER TABLE moderasi_log DROP CONSTRAINT moderasi_log_target_jenis_check;
ALTER TABLE moderasi_log ADD CONSTRAINT moderasi_log_target_jenis_check
    CHECK (target_jenis IN ('laporan', 'laporan_foto', 'umpan'));
DROP TABLE temuan_tanggapan;
DROP TABLE temuan;
DROP TYPE sikap_temuan;
