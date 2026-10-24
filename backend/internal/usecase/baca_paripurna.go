package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaParipurna interface {
	IsFraksiAda(ctx context.Context, fraksi string) (bool, error)
	SikapFraksi(ctx context.Context, fraksi string) ([]domain.SikapFraksi, error)
	IsDapilAda(ctx context.Context, id int) (bool, error)
	AspirasiReses(ctx context.Context, filter domain.FilterAspirasiReses) ([]domain.AspirasiReses, error)
}

type BacaParipurna struct {
	Pembaca PembacaParipurna
}

func (u BacaParipurna) SikapFraksi(ctx context.Context, fraksi string) ([]domain.SikapFraksi, error) {
	ada, err := u.Pembaca.IsFraksiAda(ctx, fraksi)
	if err != nil {
		return nil, fmt.Errorf("periksa fraksi %s: %w", fraksi, err)
	}
	if !ada {
		return nil, domain.TidakDitemukan("fraksi tidak ditemukan")
	}
	sikap, err := u.Pembaca.SikapFraksi(ctx, fraksi)
	if err != nil {
		return nil, fmt.Errorf("baca sikap fraksi %s: %w", fraksi, err)
	}
	return sikap, nil
}

func (u BacaParipurna) AspirasiReses(ctx context.Context, filter domain.FilterAspirasiReses) ([]domain.AspirasiReses, error) {
	if err := u.pastikanDapilAda(ctx, filter.DapilID); err != nil {
		return nil, err
	}
	filter.Bidang = strings.ToLower(strings.TrimSpace(filter.Bidang))
	aspirasi, err := u.Pembaca.AspirasiReses(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("baca aspirasi reses: %w", err)
	}
	return aspirasi, nil
}

func (u BacaParipurna) pastikanDapilAda(ctx context.Context, id int) error {
	if id == 0 {
		return nil
	}
	ada, err := u.Pembaca.IsDapilAda(ctx, id)
	if err != nil {
		return fmt.Errorf("periksa dapil %d: %w", id, err)
	}
	if !ada {
		return domain.TidakDitemukan("dapil tidak ditemukan")
	}
	return nil
}
