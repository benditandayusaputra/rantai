package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlPaketDibandingkan = `SELECT kode, coalesce(kategori, ''), jenis_pengadaan, tahun_anggaran, hps, provenance_id FROM paket WHERE kode = $1`

	sqlNilaiPembanding = `SELECT p.kode, p.hps, p.provenance_id
FROM paket p JOIN provenance pr ON pr.id = p.provenance_id
WHERE p.kategori = $1 AND p.jenis_pengadaan = $2 AND p.tahun_anggaran = $3 AND p.kode <> $4
    AND p.hps IS NOT NULL AND NOT p.is_batal AND NOT p.is_simulation
ORDER BY pr.diambil_pada DESC, p.kode`
)

func (r PaketRepository) AmbilPaketDibandingkan(ctx context.Context, kode string) (domain.PaketDibandingkan, bool, error) {
	var p domain.PaketDibandingkan
	err := r.DB.QueryRow(ctx, sqlPaketDibandingkan, kode).Scan(&p.Kode, &p.Kategori, &p.JenisPengadaan, &p.TahunAnggaran, &p.HPS, &p.ProvenanceID)
	return hasilSatuBaris(p, err)
}

func (r PaketRepository) NilaiPembanding(ctx context.Context, p domain.PaketDibandingkan) ([]domain.NilaiPembanding, error) {
	baris, err := r.DB.Query(ctx, sqlNilaiPembanding, p.Kategori, p.JenisPengadaan, p.TahunAnggaran, p.Kode)
	if err != nil {
		return nil, fmt.Errorf("nilai pembanding paket %s: %w", p.Kode, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.NilaiPembanding, error) {
		var n domain.NilaiPembanding
		err := b.Scan(&n.Kode, &n.HPS, &n.ProvenanceID)
		return n, err
	})
}
