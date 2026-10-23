package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlNPSNTerdaftar = `SELECT npsn FROM sekolah`

	sqlSimpanDayaTampung = `INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (npsn, tahun_ajaran, jalur) DO UPDATE SET
    kursi = EXCLUDED.kursi, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE daya_tampung.kursi IS DISTINCT FROM EXCLUDED.kursi
RETURNING (xmax = 0) AS baru`
)

type DayaTampungRepository struct {
	DB Basisdata
}

func (r DayaTampungRepository) NPSNTerdaftar(ctx context.Context) (map[string]bool, error) {
	baris, err := r.DB.Query(ctx, sqlNPSNTerdaftar)
	if err != nil {
		return nil, fmt.Errorf("kueri NPSN sekolah: %w", err)
	}
	daftar, err := pgx.CollectRows(baris, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	terdaftar := make(map[string]bool, len(daftar))
	for _, npsn := range daftar {
		terdaftar[npsn] = true
	}
	return terdaftar, nil
}

func (r DayaTampungRepository) SimpanDayaTampung(ctx context.Context, semuaBerkas []usecase.BerkasDayaTampung) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi daya tampung: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	hasil := usecase.HasilImpor{}
	for _, berkas := range semuaBerkas {
		hasilBerkas, err := simpanBerkasDayaTampung(ctx, tx, berkas)
		if err != nil {
			return usecase.HasilImpor{}, err
		}
		hasil.JumlahBaru += hasilBerkas.JumlahBaru
		hasil.JumlahBerubah += hasilBerkas.JumlahBerubah
	}
	return hasil, tx.Commit(ctx)
}

func simpanBerkasDayaTampung(ctx context.Context, tx pgx.Tx, berkas usecase.BerkasDayaTampung) (usecase.HasilImpor, error) {
	if err := simpanProvenance(ctx, tx, berkas.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasil := usecase.HasilImpor{}
	for _, b := range berkas.Baris {
		var baru bool
		err := tx.QueryRow(ctx, sqlSimpanDayaTampung, b.NPSN, berkas.TahunAjaran, berkas.Jalur, b.Kursi, berkas.Provenance.ID).Scan(&baru)
		_, ada, err := hasilSatuBaris(baru, err)
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan daya tampung %s: %w", b.NPSN, err)
		}
		hasil.JumlahBaru += hitungJika(ada && baru)
		hasil.JumlahBerubah += hitungJika(ada && !baru)
	}
	return hasil, nil
}
