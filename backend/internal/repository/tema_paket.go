package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlPaketUntukTema = `SELECT p.kode, p.nama, i.nama
FROM paket p JOIN instansi i ON i.id = p.instansi_id
WHERE NOT p.is_simulation
ORDER BY p.kode`

	sqlSimpanTemaPaket = `UPDATE paket SET tema_kode = $2, diubah_pada = now()
WHERE kode = $1 AND tema_kode IS DISTINCT FROM $2`
)

func (r LokasiPaketRepository) DaftarPaketUntukTema(ctx context.Context) ([]domain.MasukanTema, error) {
	baris, err := r.DB.Query(ctx, sqlPaketUntukTema)
	if err != nil {
		return nil, fmt.Errorf("daftar paket untuk tema: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.MasukanTema, error) {
		var m domain.MasukanTema
		err := b.Scan(&m.KodePaket, &m.NamaPaket, &m.NamaInstansi)
		return m, err
	})
}

func (r LokasiPaketRepository) SimpanTemaPaket(ctx context.Context, temaPerPaket map[string][]string) (int, error) {
	berubah := 0
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		for kode, tema := range temaPerPaket {
			tag, err := tx.Exec(ctx, sqlSimpanTemaPaket, kode, tema)
			if err != nil {
				return fmt.Errorf("simpan tema paket %s: %w", kode, err)
			}
			berubah += int(tag.RowsAffected())
		}
		return nil
	})
	return berubah, err
}
