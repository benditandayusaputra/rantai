package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlKandidatWilayahSekolah = `WITH titik AS (SELECT ST_SetSRID(ST_MakePoint($1::float8, $2::float8), 4326) AS g)
SELECT
    coalesce((SELECT w.kode FROM wilayah w, titik WHERE w.tingkat = 'kecamatan' AND ST_Covers(w.geom, titik.g) ORDER BY w.kode LIMIT 1), ''),
    coalesce((SELECT w.kode FROM wilayah w, titik WHERE w.tingkat = 'kelurahan' AND ST_Covers(w.geom, titik.g) ORDER BY w.kode LIMIT 1), ''),
    coalesce((SELECT min(kode) FROM wilayah
        WHERE tingkat = 'kecamatan' AND regexp_replace(lower(nama), '[^a-z]', '', 'g') = $3 HAVING count(*) = 1), '')`

	sqlSimpanSekolah = `INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, kode_kelurahan, lokasi, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, NULLIF($8, ''),
    CASE WHEN $9::float8 IS NULL THEN NULL ELSE ST_SetSRID(ST_MakePoint($9::float8, $10::float8), 4326) END, $11)
ON CONFLICT (npsn) DO UPDATE SET
    nama = EXCLUDED.nama, bentuk = EXCLUDED.bentuk, jenjang = EXCLUDED.jenjang, status = EXCLUDED.status,
    is_madrasah = EXCLUDED.is_madrasah, kode_kecamatan = EXCLUDED.kode_kecamatan, kode_kelurahan = EXCLUDED.kode_kelurahan,
    lokasi = EXCLUDED.lokasi, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (sekolah.nama, sekolah.bentuk, sekolah.status, sekolah.kode_kecamatan, sekolah.kode_kelurahan)
        IS DISTINCT FROM (EXCLUDED.nama, EXCLUDED.bentuk, EXCLUDED.status, EXCLUDED.kode_kecamatan, EXCLUDED.kode_kelurahan)
    OR sekolah.lokasi IS DISTINCT FROM EXCLUDED.lokasi
RETURNING (xmax = 0) AS baru`

	sqlSimpanRombel = `INSERT INTO sekolah_rombel (npsn, tahun_ajaran, tingkat, jumlah_rombel, jumlah_siswa, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (npsn, tahun_ajaran, tingkat) DO UPDATE SET
    jumlah_rombel = EXCLUDED.jumlah_rombel, jumlah_siswa = EXCLUDED.jumlah_siswa,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (sekolah_rombel.jumlah_rombel, sekolah_rombel.jumlah_siswa) IS DISTINCT FROM (EXCLUDED.jumlah_rombel, EXCLUDED.jumlah_siswa)`

	sqlSimpanFotoResmi = `INSERT INTO foto_sekolah (npsn, url, sumber, halaman, kredit)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (npsn) DO UPDATE SET
    url = EXCLUDED.url, sumber = EXCLUDED.sumber, halaman = EXCLUDED.halaman, kredit = EXCLUDED.kredit, diubah_pada = now()
WHERE (foto_sekolah.url, foto_sekolah.sumber, foto_sekolah.halaman, foto_sekolah.kredit)
    IS DISTINCT FROM (EXCLUDED.url, EXCLUDED.sumber, EXCLUDED.halaman, EXCLUDED.kredit)`

	sqlHapusFotoResmi = `DELETE FROM foto_sekolah WHERE npsn = $1 AND sumber = $2`

	sqlKecamatanTanpaSekolah = `SELECT w.kode || ' ' || w.nama FROM wilayah w
WHERE w.tingkat = 'kecamatan' AND NOT EXISTS (SELECT 1 FROM sekolah s WHERE s.kode_kecamatan = w.kode)
ORDER BY w.kode`
)

type SekolahRepository struct {
	DB Basisdata
}

func (r SekolahRepository) KandidatWilayah(ctx context.Context, lokasi *domain.Titik, namaKecamatanNormal string) (domain.KandidatWilayah, error) {
	lng, lat := koordinatOpsional(lokasi)
	var k domain.KandidatWilayah
	err := r.DB.QueryRow(ctx, sqlKandidatWilayahSekolah, lng, lat, namaKecamatanNormal).
		Scan(&k.KecamatanDariTitik, &k.KelurahanDariTitik, &k.KecamatanDariNama)
	if err != nil {
		return domain.KandidatWilayah{}, fmt.Errorf("cari kandidat wilayah sekolah: %w", err)
	}
	return k, nil
}

func (r SekolahRepository) SimpanSekolah(ctx context.Context, data usecase.DataSimpanSekolah) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi sekolah: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasil := usecase.HasilImpor{}
	for _, s := range data.Sekolah {
		baru, berubah, err := simpanSatuSekolah(ctx, tx, s, data.Provenance.ID)
		if err != nil {
			return usecase.HasilImpor{}, err
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, tx.Commit(ctx)
}

func simpanSatuSekolah(ctx context.Context, tx pgx.Tx, s usecase.SekolahTerpetakan, provenanceID string) (baru, berubah bool, err error) {
	lng, lat := koordinatOpsional(s.Lokasi)
	err = tx.QueryRow(ctx, sqlSimpanSekolah, s.NPSN, s.Nama, s.Bentuk, s.Bentuk.Jenjang(), s.Status, s.Bentuk.IsMadrasah(),
		s.KodeKecamatan, s.KodeKelurahan, lng, lat, provenanceID).Scan(&baru)
	_, ada, err := hasilSatuBaris(baru, err)
	if err != nil {
		return false, false, fmt.Errorf("simpan sekolah %s: %w", s.NPSN, err)
	}
	if err := simpanFotoResmi(ctx, tx, s.NPSN, s.Foto); err != nil {
		return false, false, err
	}
	for _, r := range s.Rombel {
		if _, err := tx.Exec(ctx, sqlSimpanRombel, s.NPSN, s.TahunAjaran, r.Tingkat, r.JumlahRombel, r.JumlahSiswa, provenanceID); err != nil {
			return false, false, fmt.Errorf("simpan rombel %s tingkat %d: %w", s.NPSN, r.Tingkat, err)
		}
	}
	return ada && baru, ada && !baru, nil
}

func simpanFotoResmi(ctx context.Context, tx pgx.Tx, npsn string, foto *domain.Foto) error {
	var err error
	if foto == nil {
		_, err = tx.Exec(ctx, sqlHapusFotoResmi, npsn, domain.SumberFotoSekolahKita)
	} else {
		_, err = tx.Exec(ctx, sqlSimpanFotoResmi, npsn, foto.URL, foto.Sumber, foto.Halaman, foto.Kredit)
	}
	if err != nil {
		return fmt.Errorf("simpan foto sekolah %s: %w", npsn, err)
	}
	return nil
}

func (r SekolahRepository) KecamatanTanpaSekolah(ctx context.Context) ([]string, error) {
	baris, err := r.DB.Query(ctx, sqlKecamatanTanpaSekolah)
	if err != nil {
		return nil, fmt.Errorf("kueri kecamatan tanpa sekolah: %w", err)
	}
	return pgx.CollectRows(baris, pgx.RowTo[string])
}

func koordinatOpsional(lokasi *domain.Titik) (lng, lat *float64) {
	if lokasi == nil {
		return nil, nil
	}
	return &lokasi.Lng, &lokasi.Lat
}
