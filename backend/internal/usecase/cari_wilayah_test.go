package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaTiruan struct {
	kueriDiterima string
	ada           bool
}

func (p *pembacaTiruan) CariWilayah(_ context.Context, kueri string, _ int) ([]domain.WilayahDitemukan, error) {
	p.kueriDiterima = kueri
	return []domain.WilayahDitemukan{}, nil
}

func (p *pembacaTiruan) KelurahanDiTitik(context.Context, domain.Titik) (domain.WilayahDitemukan, bool, error) {
	return domain.WilayahDitemukan{}, p.ada, nil
}

func (p *pembacaTiruan) AmbilWilayah(context.Context, string) (domain.DetailWilayah, bool, error) {
	return domain.DetailWilayah{}, p.ada, nil
}

func kategoriKesalahan(err error) domain.KategoriKesalahan {
	var kesalahan *domain.Kesalahan
	if errors.As(err, &kesalahan) {
		return kesalahan.Kategori
	}
	return 0
}

func TestCari_KueriBerantakan_DinormalkanSebelumDicari(t *testing.T) {
	pembaca := &pembacaTiruan{}

	if _, err := (CariWilayah{Pembaca: pembaca}).Cari(t.Context(), "  Kelapa   GADING "); err != nil {
		t.Fatal(err)
	}

	if pembaca.kueriDiterima != "kelapa gading" {
		t.Errorf("kueri = %q", pembaca.kueriDiterima)
	}
}

func TestCari_KueriTerlaluPendek_InputTidakValid(t *testing.T) {
	_, err := CariWilayah{Pembaca: &pembacaTiruan{}}.Cari(t.Context(), " ab ")

	if kategoriKesalahan(err) != domain.KategoriInputTidakValid {
		t.Errorf("err = %v", err)
	}
}

func TestDiTitik_Koordinat_KategoriKesalahanSesuai(t *testing.T) {
	kasus := []struct {
		nama     string
		titik    domain.Titik
		ada      bool
		kategori domain.KategoriKesalahan
	}{
		{"di dalam DKI", domain.Titik{Lat: -6.15, Lng: 106.9}, true, 0},
		{"di luar DKI", domain.Titik{Lat: -7.25, Lng: 112.75}, false, domain.KategoriTidakDitemukan},
		{"lintang mustahil", domain.Titik{Lat: -91, Lng: 106.9}, true, domain.KategoriInputTidakValid},
		{"bujur mustahil", domain.Titik{Lat: -6.15, Lng: 181}, true, domain.KategoriInputTidakValid},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			_, err := CariWilayah{Pembaca: &pembacaTiruan{ada: k.ada}}.DiTitik(t.Context(), k.titik)

			if kategoriKesalahan(err) != k.kategori {
				t.Errorf("err = %v, harap kategori %d", err, k.kategori)
			}
		})
	}
}

func TestAmbil_KodeDanKeberadaan_KategoriKesalahanSesuai(t *testing.T) {
	kasus := []struct {
		nama     string
		kode     string
		ada      bool
		kategori domain.KategoriKesalahan
	}{
		{"ada", "31.72.06", true, 0},
		{"tidak ada", "31.72.99", false, domain.KategoriTidakDitemukan},
		{"kode rusak", "31-72", true, domain.KategoriInputTidakValid},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			_, err := CariWilayah{Pembaca: &pembacaTiruan{ada: k.ada}}.Ambil(t.Context(), k.kode)

			if kategoriKesalahan(err) != k.kategori {
				t.Errorf("err = %v, harap kategori %d", err, k.kategori)
			}
		})
	}
}
