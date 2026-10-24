package usecase

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PenyimpanPeran interface {
	AturPeran(ctx context.Context, email string, peran domain.Peran) error
}

type BeriPeran struct {
	Penyimpan PenyimpanPeran
}

func (u BeriPeran) Jalankan(ctx context.Context, email string, peran domain.Peran) error {
	bersih, err := domain.NormalisasiEmail(email)
	if err != nil {
		return err
	}
	if !peran.IsDibuatLewatPerintah() {
		return domain.ErrPeranBukanUntukPerintah
	}
	return u.Penyimpan.AturPeran(ctx, bersih, peran)
}
