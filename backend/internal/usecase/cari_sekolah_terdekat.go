package usecase

import (
	"cmp"
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type KonfigurasiSekolahTerdekat struct {
	RadiusBawaanMeter   int `yaml:"radius_bawaan_m"`
	RadiusMaksimalMeter int `yaml:"radius_maksimal_m"`
}

type PembacaSekitarTitik interface {
	SekitarTitik(ctx context.Context, kueri domain.KueriSekitarTitik) (domain.SekitarTitik, error)
}

type PermintaanSekolahTerdekat struct {
	Titik       domain.Titik
	RadiusMeter int
	TemaKode    string
}

type CariSekolahTerdekat struct {
	Pembaca     PembacaSekitarTitik
	Konfigurasi KonfigurasiSekolahTerdekat
}

func (u CariSekolahTerdekat) Cari(ctx context.Context, permintaan PermintaanSekolahTerdekat) (domain.SekitarTitik, error) {
	kueri, err := u.susunKueri(permintaan)
	if err != nil {
		return domain.SekitarTitik{}, err
	}
	hasil, err := u.Pembaca.SekitarTitik(ctx, kueri)
	if err != nil {
		return domain.SekitarTitik{}, fmt.Errorf("cari sekolah sekitar %v: %w", kueri.Titik, err)
	}
	hasil.KueriSekitarTitik = kueri
	return hasil, nil
}

func (u CariSekolahTerdekat) susunKueri(p PermintaanSekolahTerdekat) (domain.KueriSekitarTitik, error) {
	if !isTitikValid(p.Titik) {
		return domain.KueriSekitarTitik{}, domain.InputTidakValid("koordinat tidak valid")
	}
	radius := cmp.Or(p.RadiusMeter, u.Konfigurasi.RadiusBawaanMeter)
	if radius < 0 || radius > u.Konfigurasi.RadiusMaksimalMeter {
		return domain.KueriSekitarTitik{}, domain.InputTidakValid(fmt.Sprintf("radius harus 1 sampai %d meter", u.Konfigurasi.RadiusMaksimalMeter))
	}
	peralihan, err := domain.PeralihanTema(cmp.Or(p.TemaKode, domain.TemaPendidikan))
	if err != nil {
		return domain.KueriSekitarTitik{}, domain.InputTidakValid("tema sekolah terdekat tidak dikenal: " + p.TemaKode)
	}
	return domain.KueriSekitarTitik{Titik: p.Titik, RadiusMeter: radius, Peralihan: peralihan}, nil
}
