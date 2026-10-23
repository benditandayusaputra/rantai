package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type KonfigurasiSekolahTerdekat struct {
	RadiusBawaanMeter   int `yaml:"radius_bawaan_m"`
	RadiusMaksimalMeter int `yaml:"radius_maksimal_m"`
}

type PembacaSekitarTitik interface {
	SekitarTitik(ctx context.Context, titik domain.Titik, radiusMeter int) (domain.SekitarTitik, error)
}

type CariSekolahTerdekat struct {
	Pembaca     PembacaSekitarTitik
	Konfigurasi KonfigurasiSekolahTerdekat
}

func (u CariSekolahTerdekat) Cari(ctx context.Context, titik domain.Titik, radiusMeter int) (domain.SekitarTitik, error) {
	if !isTitikValid(titik) {
		return domain.SekitarTitik{}, domain.InputTidakValid("koordinat tidak valid")
	}
	if radiusMeter == 0 {
		radiusMeter = u.Konfigurasi.RadiusBawaanMeter
	}
	if radiusMeter < 0 || radiusMeter > u.Konfigurasi.RadiusMaksimalMeter {
		return domain.SekitarTitik{}, domain.InputTidakValid(fmt.Sprintf("radius harus 1 sampai %d meter", u.Konfigurasi.RadiusMaksimalMeter))
	}
	hasil, err := u.Pembaca.SekitarTitik(ctx, titik, radiusMeter)
	if err != nil {
		return domain.SekitarTitik{}, fmt.Errorf("cari sekolah sekitar %v: %w", titik, err)
	}
	hasil.Titik, hasil.RadiusMeter = titik, radiusMeter
	return hasil, nil
}
