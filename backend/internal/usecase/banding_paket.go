package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaPembanding interface {
	AmbilPaketDibandingkan(ctx context.Context, kode string) (domain.PaketDibandingkan, bool, error)
	NilaiPembanding(ctx context.Context, p domain.PaketDibandingkan) ([]domain.NilaiPembanding, error)
}

type BandingPaket struct {
	Pembaca       PembacaPembanding
	MinPembanding int
}

func (u BandingPaket) Bandingkan(ctx context.Context, kode string) (domain.PembandingPaket, error) {
	p, ada, err := u.Pembaca.AmbilPaketDibandingkan(ctx, kode)
	if err != nil {
		return domain.PembandingPaket{}, fmt.Errorf("ambil paket %s: %w", kode, err)
	}
	if !ada {
		return domain.PembandingPaket{}, domain.TidakDitemukan("paket tidak ditemukan")
	}
	pembanding := []domain.NilaiPembanding{}
	if p.Kategori != "" {
		if pembanding, err = u.Pembaca.NilaiPembanding(ctx, p); err != nil {
			return domain.PembandingPaket{}, err
		}
	}
	terbaru := p.ProvenanceID
	if len(pembanding) > 0 {
		terbaru = pembanding[0].ProvenanceID
	}
	return domain.BandingkanPaket(p, pembanding, terbaru, u.MinPembanding), nil
}
