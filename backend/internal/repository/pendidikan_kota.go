package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const sqlKecamatanKota = `SELECT kode, nama, tingkat FROM wilayah
WHERE induk_kode = $1 AND tingkat = 'kecamatan' ORDER BY nama, kode`

const sqlIndikatorKota = `WITH terbaru AS (
    SELECT max(periode) AS periode FROM indikator_kecamatan WHERE tema_kode = $2
)
SELECT ` + sqlKolomIndikator + `
FROM wilayah w JOIN indikator_kecamatan i ON i.kode_kecamatan = w.kode CROSS JOIN terbaru t
WHERE w.induk_kode = $1 AND i.tema_kode = $2 AND i.periode = t.periode ORDER BY w.kode`

const sqlSekolahKota = `SELECT s.kode_kecamatan, s.bentuk,
    count(*) FILTER (WHERE s.status = 'negeri')::int,
    count(*) FILTER (WHERE s.status = 'swasta')::int,
    p.label::text, p.id, p.diambil_pada
FROM sekolah s JOIN wilayah w ON w.kode = s.kode_kecamatan
JOIN provenance p ON p.id = s.provenance_id
WHERE w.induk_kode = $1 AND s.bentuk IN ('SD', 'SMP', 'SMA', 'SMK')
GROUP BY s.kode_kecamatan, s.bentuk, p.label, p.id, p.diambil_pada
ORDER BY s.kode_kecamatan, s.bentuk, p.id`

type PendidikanKotaRepository struct {
	DB Basisdata
}

func (r PendidikanKotaRepository) KecamatanKota(ctx context.Context, kode string) ([]domain.RingkasanWilayah, error) {
	baris, err := r.DB.Query(ctx, sqlKecamatanKota, kode)
	if err != nil {
		return nil, fmt.Errorf("kueri kecamatan kabupaten/kota: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.RingkasanWilayah, error) {
		var w domain.RingkasanWilayah
		err := b.Scan(&w.Kode, &w.Nama, &w.Tingkat)
		return w, err
	})
}

func (r PendidikanKotaRepository) IndikatorKota(ctx context.Context, kode, tema string) ([]domain.IndikatorTersimpan, error) {
	baris, err := r.DB.Query(ctx, sqlIndikatorKota, kode, tema)
	if err != nil {
		return nil, fmt.Errorf("kueri indikator kabupaten/kota: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.IndikatorTersimpan, error) {
		var i domain.IndikatorTersimpan
		err := b.Scan(tujuanIndikator(&i)...)
		return i, err
	})
}

func (r PendidikanKotaRepository) SekolahKota(ctx context.Context, kode string) ([]domain.JumlahSekolahKecamatan, error) {
	baris, err := r.DB.Query(ctx, sqlSekolahKota, kode)
	if err != nil {
		return nil, fmt.Errorf("kueri sebaran sekolah kabupaten/kota: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.JumlahSekolahKecamatan, error) {
		var s domain.JumlahSekolahKecamatan
		err := b.Scan(&s.KodeKecamatan, &s.Bentuk, &s.Negeri, &s.Swasta, &s.Label, &s.ProvenanceID, &s.DiambilPada)
		return s, err
	})
}

func (r PendidikanKotaRepository) DaftarKota(ctx context.Context) ([]domain.RingkasanWilayah, error) {
	baris, err := r.DB.Query(ctx, `SELECT kode, nama, tingkat FROM wilayah WHERE tingkat = 'kota' AND kode LIKE '31.%' ORDER BY kode`)
	if err != nil {
		return nil, fmt.Errorf("kueri daftar kabupaten/kota: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.RingkasanWilayah, error) {
		var w domain.RingkasanWilayah
		err := b.Scan(&w.Kode, &w.Nama, &w.Tingkat)
		return w, err
	})
}
