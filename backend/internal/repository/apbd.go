package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlBulanAPBD  = `SELECT coalesce(max(bulan_realisasi), 0) FROM apbd_akun WHERE tahun = $1`
	sqlHapusAPBD  = `DELETE FROM apbd_akun WHERE tahun = $1`
	sqlSimpanAPBD = `INSERT INTO apbd_akun (tahun, urutan, nama, anggaran, realisasi, bulan_realisasi, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
	sqlTahunAPBD = `SELECT CASE WHEN $1::int > 0 THEN $1::int ELSE coalesce(max(tahun), 0) END FROM apbd_akun`
	sqlBacaAPBD  = `SELECT nama, anggaran::float8, realisasi::float8, bulan_realisasi, provenance_id
FROM apbd_akun WHERE tahun = $1 ORDER BY urutan`
)

type APBDRepository struct {
	DB Basisdata
}

func (r APBDRepository) BulanRealisasiTersimpan(ctx context.Context, tahun int) (int, error) {
	var bulan int
	return bulan, r.DB.QueryRow(ctx, sqlBulanAPBD, tahun).Scan(&bulan)
}

func (r APBDRepository) SimpanAPBD(ctx context.Context, sumber usecase.APBDSumber) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, sumber.Provenance); err != nil {
			return err
		}
		apbd := sumber.APBD
		if _, err := tx.Exec(ctx, sqlHapusAPBD, apbd.Tahun); err != nil {
			return fmt.Errorf("hapus apbd %d: %w", apbd.Tahun, err)
		}
		for urutan, akun := range apbd.Akun {
			_, err := tx.Exec(ctx, sqlSimpanAPBD, apbd.Tahun, urutan, akun.Nama, akun.Anggaran, akun.Realisasi, apbd.BulanRealisasi, apbd.ProvenanceID)
			if err != nil {
				return fmt.Errorf("simpan akun apbd %s: %w", akun.Nama, err)
			}
		}
		return nil
	})
}

func (r APBDRepository) BacaAPBD(ctx context.Context, tahunDiminta int) (domain.APBDDaerah, error) {
	apbd := domain.APBDDaerah{}
	if err := r.DB.QueryRow(ctx, sqlTahunAPBD, tahunDiminta).Scan(&apbd.Tahun); err != nil {
		return apbd, fmt.Errorf("tahun apbd: %w", err)
	}
	akun, err := kumpulkan(ctx, r.DB, sqlBacaAPBD, func(b pgx.CollectableRow) (domain.AkunAPBD, error) {
		var a domain.AkunAPBD
		return a, b.Scan(&a.Nama, &a.Anggaran, &a.Realisasi, &apbd.BulanRealisasi, &apbd.ProvenanceID)
	}, apbd.Tahun)
	if err != nil {
		return apbd, fmt.Errorf("baca apbd %d: %w", apbd.Tahun, err)
	}
	if len(akun) == 0 {
		return apbd, domain.ErrAPBDTidakDitemukan
	}
	apbd.Akun = akun
	return apbd, nil
}
