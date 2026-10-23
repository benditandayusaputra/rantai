package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaTema interface {
	DaftarTemaAktif(ctx context.Context) ([]domain.TemaTersimpan, error)
	AmbilTemaAktif(ctx context.Context, kode string) (domain.TemaTersimpan, bool, error)
}

type BacaTema struct {
	Pembaca PembacaTema
}

func (u BacaTema) Daftar(ctx context.Context) ([]domain.TemaTersimpan, error) {
	daftar, err := u.Pembaca.DaftarTemaAktif(ctx)
	if err != nil {
		return nil, fmt.Errorf("daftar tema aktif: %w", err)
	}
	return daftar, nil
}

func (u BacaTema) Ambil(ctx context.Context, kode string) (domain.TemaTersimpan, error) {
	tema, ada, err := u.Pembaca.AmbilTemaAktif(ctx, kode)
	if err != nil {
		return domain.TemaTersimpan{}, fmt.Errorf("ambil tema %s: %w", kode, err)
	}
	if !ada {
		return domain.TemaTersimpan{}, domain.TidakDitemukan("tema tidak ditemukan")
	}
	return tema, nil
}
