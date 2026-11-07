package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaRelasiPaket interface {
	BacaRelasiPaket(ctx context.Context, f domain.FilterJaringan, batas int) ([]domain.RelasiPaket, error)
}

type BacaJaringan struct {
	Pembaca PembacaRelasiPaket
}

func (u BacaJaringan) Susun(ctx context.Context, f domain.FilterJaringan) (domain.JaringanRelasi, error) {
	if err := f.Validasi(); err != nil {
		return domain.JaringanRelasi{}, err
	}
	relasi, err := u.Pembaca.BacaRelasiPaket(ctx, f, domain.BatasPaketJaringan)
	if err != nil {
		return domain.JaringanRelasi{}, fmt.Errorf("baca relasi paket: %w", err)
	}
	jaringan := domain.SusunJaringan(relasi)
	jaringan.IsTerpotong = len(relasi) >= domain.BatasPaketJaringan
	return jaringan, nil
}
