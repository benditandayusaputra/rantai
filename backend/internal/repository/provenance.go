package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const sqlAmbilProvenance = `SELECT id, sumber, coalesce(url, ''), diambil_pada, label, coalesce(hash_konten, ''), coalesce(catatan, '')
FROM provenance
WHERE id = $1`

type ProvenanceRepository struct {
	DB Basisdata
}

func (r ProvenanceRepository) AmbilProvenance(ctx context.Context, id string) (domain.Provenance, bool, error) {
	var p domain.Provenance
	err := r.DB.QueryRow(ctx, sqlAmbilProvenance, id).Scan(
		&p.ID, &p.Sumber, &p.URL, &p.DiambilPada, &p.Label, &p.HashKonten, &p.Catatan,
	)
	return hasilSatuBaris(p, err)
}
