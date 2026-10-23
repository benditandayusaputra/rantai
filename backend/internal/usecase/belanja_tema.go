package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaBelanja interface {
	BelanjaKecamatan(ctx context.Context, kodeTema string, tahun int) ([]domain.BelanjaKecamatan, error)
	BelanjaLuarKecamatan(ctx context.Context, kodeTema string, tahun int) ([]domain.BelanjaLuarKecamatan, error)
}

type BandingBelanja struct {
	Tema     BacaTema
	Pembaca  PembacaBelanja
	Sekarang func() time.Time
}

func (u BandingBelanja) Hitung(ctx context.Context, kodeTema string, tahun int) (domain.BelanjaTema, error) {
	tema, err := u.Tema.Ambil(ctx, kodeTema)
	if err != nil {
		return domain.BelanjaTema{}, err
	}
	if tahun == 0 {
		tahun = u.Sekarang().Year()
	}
	kecamatan, err := u.Pembaca.BelanjaKecamatan(ctx, kodeTema, tahun)
	if err != nil {
		return domain.BelanjaTema{}, fmt.Errorf("belanja kecamatan: %w", err)
	}
	luar, err := u.Pembaca.BelanjaLuarKecamatan(ctx, kodeTema, tahun)
	if err != nil {
		return domain.BelanjaTema{}, fmt.Errorf("belanja luar kecamatan: %w", err)
	}
	return domain.BelanjaTema{Tema: tema, Tahun: tahun, Kecamatan: kecamatan, LuarKecamatan: luar}, nil
}
