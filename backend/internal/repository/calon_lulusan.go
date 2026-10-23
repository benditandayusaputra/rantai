package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlSiswaSDPerKecamatan = `WITH terbaru AS (SELECT max(tahun_ajaran) AS tahun_ajaran FROM sekolah_rombel)
SELECT w.kode, coalesce(t.tahun_ajaran, ''),
    coalesce(sum(r.jumlah_siswa) FILTER (WHERE r.tingkat = $1), 0)::int,
    coalesce(sum(r.jumlah_siswa), 0)::int
FROM wilayah w
CROSS JOIN terbaru t
LEFT JOIN sekolah s ON s.kode_kecamatan = w.kode AND s.bentuk = 'SD'
LEFT JOIN sekolah_rombel r ON r.npsn = s.npsn AND r.tahun_ajaran = t.tahun_ajaran
WHERE w.tingkat = 'kecamatan'
GROUP BY w.kode, t.tahun_ajaran
ORDER BY w.kode`

	sqlSimpanCalonLulusan = `INSERT INTO calon_lulusan_kecamatan (kode_kecamatan, tahun_ajaran, jumlah, label, metode, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (kode_kecamatan, tahun_ajaran) DO UPDATE SET
    jumlah = EXCLUDED.jumlah, label = EXCLUDED.label, metode = EXCLUDED.metode,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (calon_lulusan_kecamatan.jumlah, calon_lulusan_kecamatan.label, calon_lulusan_kecamatan.metode)
    IS DISTINCT FROM (EXCLUDED.jumlah, EXCLUDED.label, EXCLUDED.metode)
RETURNING (xmax = 0) AS baru`
)

type CalonLulusanRepository struct {
	DB Basisdata
}

func (r CalonLulusanRepository) SiswaSDPerKecamatan(ctx context.Context) (string, []domain.SiswaSDKecamatan, error) {
	baris, err := r.DB.Query(ctx, sqlSiswaSDPerKecamatan, domain.TingkatSDTerakhir)
	if err != nil {
		return "", nil, fmt.Errorf("kueri siswa SD per kecamatan: %w", err)
	}
	tahunAjaran := ""
	daftar, err := pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.SiswaSDKecamatan, error) {
		var s domain.SiswaSDKecamatan
		err := b.Scan(&s.KodeKecamatan, &tahunAjaran, &s.SiswaKelas6, &s.TotalSiswa)
		return s, err
	})
	return tahunAjaran, daftar, err
}

func (r CalonLulusanRepository) SimpanCalonLulusan(ctx context.Context, data usecase.DataCalonLulusan) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi calon lulusan: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, p := range data.Provenance {
		if err := simpanProvenance(ctx, tx, p); err != nil {
			return usecase.HasilImpor{}, err
		}
	}
	hasil := usecase.HasilImpor{}
	for _, c := range data.Daftar {
		var baru bool
		err := tx.QueryRow(ctx, sqlSimpanCalonLulusan, c.KodeKecamatan, data.TahunAjaran, c.Jumlah, c.Label, c.Metode, c.ProvenanceID).Scan(&baru)
		_, ada, err := hasilSatuBaris(baru, err)
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan calon lulusan %s: %w", c.KodeKecamatan, err)
		}
		hasil.JumlahBaru += hitungJika(ada && baru)
		hasil.JumlahBerubah += hitungJika(ada && !baru)
	}
	return hasil, tx.Commit(ctx)
}
