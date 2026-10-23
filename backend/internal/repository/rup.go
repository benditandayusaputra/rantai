package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlRUPBelumLengkap = `SELECT r.kode FROM rup r JOIN provenance p ON p.id = r.provenance_id
WHERE p.sumber <> 'sirup'
ORDER BY r.tahun_anggaran DESC NULLS LAST, r.kode`

	sqlSimpanRUP = `UPDATE rup SET nama = $2, instansi_id = $3, pagu = $4, sumber_dana = $5, tahun_anggaran = $6, metode = $7,
    jadwal_pelaksanaan_mulai = $8, jadwal_pelaksanaan_selesai = $9, provenance_id = $10, diubah_pada = now()
WHERE kode = $1`
)

type RUPRepository struct {
	DB Basisdata
}

func (r RUPRepository) DaftarRUPBelumLengkap(ctx context.Context) ([]string, error) {
	baris, err := r.DB.Query(ctx, sqlRUPBelumLengkap)
	if err != nil {
		return nil, fmt.Errorf("daftar rup belum lengkap: %w", err)
	}
	return pgx.CollectRows(baris, pgx.RowTo[string])
}

func (r RUPRepository) SimpanRUP(ctx context.Context, sumber usecase.RUPSumber, instansi domain.Instansi) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, sumber.Provenance); err != nil {
			return err
		}
		idInstansi, err := simpanInstansi(ctx, tx, instansi)
		if err != nil {
			return err
		}
		u := sumber.RUP
		_, err = tx.Exec(ctx, sqlSimpanRUP, u.Kode, u.Nama, idInstansi, u.Pagu, u.SumberDana, u.TahunAnggaran, u.Metode,
			u.PelaksanaanMulai, u.PelaksanaanSelesai, sumber.Provenance.ID)
		if err != nil {
			return fmt.Errorf("simpan rup %s: %w", u.Kode, err)
		}
		return nil
	})
}
