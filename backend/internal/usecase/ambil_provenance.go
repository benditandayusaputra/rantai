package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaProvenance interface {
	AmbilProvenance(ctx context.Context, id string) (domain.Provenance, bool, error)
}

type AmbilProvenance struct {
	Pembaca PembacaProvenance
}

func (u AmbilProvenance) Ambil(ctx context.Context, id string) (domain.Provenance, error) {
	provenance, ada, err := u.Pembaca.AmbilProvenance(ctx, id)
	if err != nil {
		return domain.Provenance{}, fmt.Errorf("ambil provenance %s: %w", id, err)
	}
	if !ada {
		return domain.Provenance{}, domain.TidakDitemukan("sumber data tidak ditemukan")
	}
	return provenance, nil
}
