package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	sqlHapusFotoCadanganLama = `DELETE FROM foto_sekolah WHERE sumber = $1 AND NOT (npsn = ANY ($2))`

	sqlSimpanFotoCadangan = `INSERT INTO foto_sekolah (npsn, url, sumber, halaman, kredit)
SELECT s.npsn, $2, $3, $4, $5 FROM sekolah s WHERE s.npsn = $1
ON CONFLICT (npsn) DO UPDATE SET
    url = EXCLUDED.url, halaman = EXCLUDED.halaman, kredit = EXCLUDED.kredit, diubah_pada = now()
WHERE foto_sekolah.sumber = EXCLUDED.sumber
    AND (foto_sekolah.url, foto_sekolah.halaman, foto_sekolah.kredit) IS DISTINCT FROM (EXCLUDED.url, EXCLUDED.halaman, EXCLUDED.kredit)
RETURNING (xmax = 0) AS baru`
)

type FotoSekolahRepository struct {
	DB Basisdata
}

func (r FotoSekolahRepository) SimpanFotoCadangan(ctx context.Context, sumber domain.SumberFoto, daftar []domain.FotoSekolah) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi foto sekolah: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	npsn := make([]string, len(daftar))
	for i, f := range daftar {
		npsn[i] = f.NPSN
	}
	if _, err := tx.Exec(ctx, sqlHapusFotoCadanganLama, sumber, npsn); err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("hapus foto %s lama: %w", sumber, err)
	}
	hasil := usecase.HasilImpor{}
	for _, f := range daftar {
		baru, berubah, err := tulisBaris(tx.QueryRow(ctx, sqlSimpanFotoCadangan, f.NPSN, f.URL, sumber, f.Halaman, f.Kredit))
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan foto sekolah %s: %w", f.NPSN, err)
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, tx.Commit(ctx)
}
