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
	UmpanBertagAnggota(ctx context.Context, anggotaID string) ([]domain.UmpanTampil, error)
}

type BacaUmpan struct {
	Pembaca PembacaUmpan
}

func (u BacaUmpan) IsuAnggota(ctx context.Context, anggotaID string) ([]domain.IsuDisuarakan, error) {
	if err := u.pastikanAnggotaAda(ctx, anggotaID); err != nil {
		return nil, err
	}
	umpan, err := u.Pembaca.UmpanBertagAnggota(ctx, anggotaID)
	if err != nil {
		return nil, fmt.Errorf("baca umpan bertag anggota %s: %w", anggotaID, err)
	}
	return domain.KelompokkanIsu(umpan), nil
}

func (u BacaUmpan) pastikanAnggotaAda(ctx context.Context, anggotaID string) error {
	ada, err := u.Pembaca.IsAnggotaAda(ctx, anggotaID)
	if err != nil {
		return fmt.Errorf("periksa anggota %s: %w", anggotaID, err)
	}
	if !ada {
		return domain.TidakDitemukan("anggota tidak ditemukan")
	}
	return nil
}

func (u BacaUmpan) UmpanAnggota(ctx context.Context, filter domain.FilterUmpan) (domain.HalamanUmpan, error) {
	if err := validasiFilterUmpan(filter); err != nil {
		return domain.HalamanUmpan{}, err
	}
	if err := u.pastikanAnggotaAda(ctx, filter.AnggotaID); err != nil {
		return domain.HalamanUmpan{}, err
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
