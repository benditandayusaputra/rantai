package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaPaketKosong struct {
	radius float64
}

func (p *pembacaPaketKosong) AmbilPaket(context.Context, string) (domain.PaketLengkap, bool, error) {
	return domain.PaketLengkap{}, false, nil
}

func (p *pembacaPaketKosong) DaftarPaketKecamatan(context.Context, FilterPaketKecamatan) ([]domain.PaketRingkas, int, error) {
	return []domain.PaketRingkas{}, 0, nil
}

func (p *pembacaPaketKosong) PetaPaket(context.Context, FilterPetaPaket) ([]domain.TitikPaket, []domain.AreaPaket, error) {
	return []domain.TitikPaket{}, []domain.AreaPaket{}, nil
}

func (p *pembacaPaketKosong) PaketTerdekat(_ context.Context, _ domain.Titik, radius float64, _ int) ([]domain.TitikPaket, error) {
	p.radius = radius
	return []domain.TitikPaket{}, nil
}

func bacaPaketUji(pembaca *pembacaPaketKosong) BacaPaket {
	return BacaPaket{Pembaca: pembaca, Aturan: AturanBacaPaket{RadiusBawaanM: 500, RadiusMaksM: 3000, LuasKotakMaks: 1}}
}

func TestTerdekat_RadiusKosongDanBerlebih(t *testing.T) {
	pembaca := &pembacaPaketKosong{}
	u := bacaPaketUji(pembaca)

	if _, err := u.Terdekat(t.Context(), domain.Titik{}, 0); err != nil || pembaca.radius != 500 {
		t.Errorf("radius bawaan = %v, err = %v", pembaca.radius, err)
	}
	var kesalahan *domain.Kesalahan
	if _, err := u.Terdekat(t.Context(), domain.Titik{}, 5000); !errors.As(err, &kesalahan) {
		t.Errorf("radius berlebih diterima: %v", err)
	}
}

func TestPeta_KotakTerbalikAtauTerlaluLuas_Ditolak(t *testing.T) {
	u := bacaPaketUji(&pembacaPaketKosong{})
	for _, kotak := range []domain.Kotak{{Barat: 107, Selatan: -6, Timur: 106, Utara: -5}, {Barat: 100, Selatan: -10, Timur: 110, Utara: 0}} {
		if _, err := u.Peta(t.Context(), FilterPetaPaket{Kotak: kotak}); err == nil {
			t.Errorf("kotak %+v diterima", kotak)
		}
	}
}

func TestAmbil_PaketTidakAda_TidakDitemukan(t *testing.T) {
	_, err := bacaPaketUji(&pembacaPaketKosong{}).Ambil(t.Context(), "x")

	var kesalahan *domain.Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriTidakDitemukan {
		t.Errorf("err = %v", err)
	}
}
