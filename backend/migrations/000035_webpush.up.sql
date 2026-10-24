ALTER TABLE kejadian ADD COLUMN didorong_pada timestamptz;

CREATE INDEX kejadian_belum_didorong_idx ON kejadian (dibuat_pada) WHERE didorong_pada IS NULL;
