package usecase

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaAPBD interface {
	BacaAPBD(ctx context.Context, tahun int) (domain.APBDDaerah, error)
}

type BacaAPBD struct {
	Pembaca PembacaAPBD
	Aturan  domain.AturanAPBD
}

func (u BacaAPBD) RingkasBelanja(ctx context.Context, tahun int) (domain.RingkasanBelanjaAPBD, error) {
	apbd, err := u.Pembaca.BacaAPBD(ctx, tahun)
	if err != nil {
		return domain.RingkasanBelanjaAPBD{}, err
	}
	return apbd.RingkasBelanja(u.Aturan)
}
