package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const sqlSimpanRiwayat = `INSERT INTO riwayat_penyedia (penyedia_id, tahun, pembeli, nama_paket, nilai, sumber, kode_sumber, provenance_id)
VALUES ($1, $2, $3, $4, $5, 'opentender', $6, $7)
ON CONFLICT (sumber, kode_sumber) DO UPDATE SET
    penyedia_id = EXCLUDED.penyedia_id, tahun = EXCLUDED.tahun, pembeli = EXCLUDED.pembeli, nama_paket = EXCLUDED.nama_paket,
    nilai = EXCLUDED.nilai, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (riwayat_penyedia.penyedia_id, riwayat_penyedia.tahun, riwayat_penyedia.pembeli, riwayat_penyedia.nama_paket, riwayat_penyedia.nilai)
    IS DISTINCT FROM (EXCLUDED.penyedia_id, EXCLUDED.tahun, EXCLUDED.pembeli, EXCLUDED.nama_paket, EXCLUDED.nilai)
RETURNING (xmax = 0)`

type RiwayatPenyediaRepository struct {
	DB Basisdata
}

func (r RiwayatPenyediaRepository) SimpanRiwayat(ctx context.Context, sumber usecase.RiwayatSumber, penyedia map[string]domain.Penyedia) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, sumber.Provenance); err != nil {
			return err
		}
		idPenyedia := map[string]int64{}
		for nama, p := range penyedia {
			id, err := simpanPenyedia(ctx, tx, &p)
			if err != nil {
				return err
			}
			idPenyedia[nama] = *id
		}
		for _, riwayat := range sumber.Riwayat {
			if err := simpanSatuRiwayat(ctx, tx, riwayat, idPenyedia[riwayat.NamaPenyedia], sumber.Provenance.ID, &hasil); err != nil {
				return err
			}
		}
		return nil
	})
	return hasil, err
}

func simpanSatuRiwayat(ctx context.Context, tx pgx.Tx, r domain.RiwayatPenyedia, idPenyedia int64, provenanceID string, hasil *usecase.HasilImpor) error {
	var baru bool
	err := tx.QueryRow(ctx, sqlSimpanRiwayat, idPenyedia, r.Tahun, r.Pembeli, r.NamaPaket, r.Nilai, r.KodeSumber, provenanceID).Scan(&baru)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	case err != nil:
		return fmt.Errorf("simpan riwayat %s: %w", r.KodeSumber, err)
	case baru:
		hasil.JumlahBaru++
	default:
		hasil.JumlahBerubah++
	}
	return nil
}
