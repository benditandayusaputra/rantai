package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const sqlHakJawabTarget = `SELECT h.id, h.target_jenis, h.target_id, h.isi, coalesce(h.tautan_dokumen, ''), g.nama_tampilan,
    coalesce(i.nama, a.nama, ''), h.is_simulation, h.dibuat_pada
FROM hak_jawab h
JOIN pengguna g ON g.id = h.penjawab_id
LEFT JOIN instansi i ON i.id = h.instansi_id
LEFT JOIN anggota a ON a.id = h.anggota_id
WHERE h.target_jenis = $1 AND h.target_id = ANY ($2)
ORDER BY h.dibuat_pada, h.id`

func hakJawabUntuk(ctx context.Context, db Basisdata, jenis domain.JenisTargetHakJawab, id []string) (map[string][]domain.HakJawab, error) {
	hasil := map[string][]domain.HakJawab{}
	if len(id) == 0 {
		return hasil, nil
	}
	daftar, err := kumpulkan(ctx, db, sqlHakJawabTarget, pgx.RowToStructByPos[domain.HakJawab], jenis, id)
	if err != nil {
		return nil, fmt.Errorf("baca hak jawab %s: %w", jenis, err)
	}
	for _, h := range daftar {
		hasil[h.TargetID] = append(hasil[h.TargetID], h)
	}
	return hasil, nil
}

func hakJawabAtau(peta map[string][]domain.HakJawab, id string) []domain.HakJawab {
	if daftar, ada := peta[id]; ada {
		return daftar
	}
	return []domain.HakJawab{}
}
