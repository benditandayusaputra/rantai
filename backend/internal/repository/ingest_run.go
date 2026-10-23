package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	sqlMulaiIngest = `INSERT INTO ingest_run (sumber, mulai, status) VALUES ($1, now(), 'berjalan') RETURNING id`

	sqlSelesaiIngest = `UPDATE ingest_run
SET selesai = now(), status = 'berhasil', jumlah_baru = $2, jumlah_berubah = $3, diubah_pada = now()
WHERE id = $1`

	sqlGagalIngest = `UPDATE ingest_run
SET selesai = now(), status = 'gagal', pesan_error = $2, diubah_pada = now()
WHERE id = $1`
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
