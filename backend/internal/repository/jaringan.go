package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const sqlRelasiPaket = `SELECT p.kode, p.nama, coalesce(p.nilai_kontrak, p.hps, p.pagu, 0)::float8, p.provenance_id,
    i.id, i.nama, p.penyedia_id, pe.nama
FROM paket p
JOIN instansi i ON i.id = p.instansi_id
JOIN penyedia pe ON pe.id = p.penyedia_id
WHERE NOT p.is_simulation AND NOT p.is_batal
    AND ($1::text = '' OR p.kode_wilayah = $1::text OR p.kode_wilayah LIKE $1::text || '.%')
    AND ($2::bigint = 0 OR p.penyedia_id = $2::bigint)
    AND ($3::int = 0 OR p.tahun_anggaran = $3::int)
ORDER BY 3 DESC, p.kode
LIMIT $4`

type JaringanRepository struct {
	DB Basisdata
}

func (r JaringanRepository) BacaRelasiPaket(ctx context.Context, f domain.FilterJaringan, batas int) ([]domain.RelasiPaket, error) {
	return kumpulkan(ctx, r.DB, sqlRelasiPaket, func(b pgx.CollectableRow) (domain.RelasiPaket, error) {
		var rp domain.RelasiPaket
		return rp, b.Scan(&rp.PaketKode, &rp.PaketNama, &rp.Nilai, &rp.ProvenanceID,
			&rp.InstansiID, &rp.InstansiNama, &rp.PenyediaID, &rp.PenyediaNama)
	}, f.KodeWilayah, f.PenyediaID, f.Tahun, batas)
}
