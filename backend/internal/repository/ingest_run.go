package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlMulaiIngest = `INSERT INTO ingest_run (sumber, mulai, status) VALUES ($1, now(), 'berjalan') RETURNING id`

	sqlSelesaiIngest = `UPDATE ingest_run
SET selesai = now(), status = 'berhasil', jumlah_baru = $2, jumlah_berubah = $3, diubah_pada = now()
WHERE id = $1`

	sqlGagalIngest = `UPDATE ingest_run
SET selesai = now(), status = 'gagal', pesan_error = $2, diubah_pada = now()
WHERE id = $1`

	sqlPengambilanTerakhir = `SELECT DISTINCT ON (a.sumber) a.sumber, a.status::text, a.mulai,
    (SELECT max(b.selesai) FROM ingest_run b WHERE b.sumber = a.sumber AND b.status = 'berhasil')
FROM ingest_run a
ORDER BY a.sumber, a.mulai DESC`
)

type IngestRunRepository struct {
	DB Basisdata
}

func (r IngestRunRepository) Mulai(ctx context.Context, job string) (int64, error) {
	var id int64
	if err := r.DB.QueryRow(ctx, sqlMulaiIngest, job).Scan(&id); err != nil {
		return 0, fmt.Errorf("catat mulai ingest %s: %w", job, err)
	}
	return id, nil
}

func (r IngestRunRepository) Selesai(ctx context.Context, id int64, hasil usecase.HasilImpor) error {
	if _, err := r.DB.Exec(ctx, sqlSelesaiIngest, id, hasil.JumlahBaru, hasil.JumlahBerubah); err != nil {
		return fmt.Errorf("catat selesai ingest %d: %w", id, err)
	}
	return nil
}

func (r IngestRunRepository) Gagal(ctx context.Context, id int64, sebab error) error {
	if _, err := r.DB.Exec(ctx, sqlGagalIngest, id, sebab.Error()); err != nil {
		return fmt.Errorf("catat gagal ingest %d: %w", id, err)
	}
	return nil
}

func (r IngestRunRepository) DaftarPengambilanTerakhir(ctx context.Context) ([]domain.PengambilanSumber, error) {
	baris, err := r.DB.Query(ctx, sqlPengambilanTerakhir)
	if err != nil {
		return nil, fmt.Errorf("daftar pengambilan terakhir: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.PengambilanSumber, error) {
		var p domain.PengambilanSumber
		err := b.Scan(&p.Job, &p.StatusTerakhir, &p.DijalankanPada, &p.BerhasilPada)
		return p, err
	})
}
