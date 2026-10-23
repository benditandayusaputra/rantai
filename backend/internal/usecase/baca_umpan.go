package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const UmpanPerHalaman = 20

type PembacaUmpan interface {
	IsAnggotaAda(ctx context.Context, id string) (bool, error)
	UmpanAnggota(ctx context.Context, filter domain.FilterUmpan) (domain.HalamanUmpan, error)
}

type BacaUmpan struct {
	Pembaca PembacaUmpan
}

func (u BacaUmpan) UmpanAnggota(ctx context.Context, filter domain.FilterUmpan) (domain.HalamanUmpan, error) {
	if err := validasiFilterUmpan(filter); err != nil {
		return domain.HalamanUmpan{}, err
	}
	ada, err := u.Pembaca.IsAnggotaAda(ctx, filter.AnggotaID)
	if err != nil {
		return domain.HalamanUmpan{}, fmt.Errorf("periksa anggota %s: %w", filter.AnggotaID, err)
	}
	if !ada {
		return domain.HalamanUmpan{}, domain.TidakDitemukan("anggota tidak ditemukan")
	}
	filter.PerHalaman = UmpanPerHalaman
	filter.Tag = strings.ToLower(strings.TrimSpace(filter.Tag))
	halaman, err := u.Pembaca.UmpanAnggota(ctx, filter)
	if err != nil {
		return domain.HalamanUmpan{}, fmt.Errorf("baca umpan anggota %s: %w", filter.AnggotaID, err)
	}
	halaman.PerHalaman = filter.PerHalaman
	return halaman, nil
}

func validasiFilterUmpan(filter domain.FilterUmpan) error {
	switch {
	case filter.Halaman < 1:
		return domain.InputTidakValid("halaman harus angka 1 atau lebih")
	case filter.Jenis != "" && filter.Jenis != domain.SumberUmpanResmiDPRD && filter.Jenis != domain.SumberUmpanMedia:
		return domain.InputTidakValid("jenis harus resmi_dprd atau media")
	default:
		return nil
	}
}
