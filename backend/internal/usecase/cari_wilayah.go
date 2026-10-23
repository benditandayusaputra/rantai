package usecase

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	PanjangMinimalKueri  = 3
	PanjangMaksimalKueri = 200
	BatasHasilCari       = 20
	batasLintang         = 90
	batasBujur           = 180
)

type PembacaWilayah interface {
	CariWilayah(ctx context.Context, kueri string, batas int) ([]domain.WilayahDitemukan, error)
	KelurahanDiTitik(ctx context.Context, titik domain.Titik) (domain.WilayahDitemukan, bool, error)
	AmbilWilayah(ctx context.Context, kode string) (domain.DetailWilayah, bool, error)
}

type CariWilayah struct {
	Pembaca PembacaWilayah
}

func (u CariWilayah) Cari(ctx context.Context, kueri string) ([]domain.WilayahDitemukan, error) {
	kueriNormal := normalkanKueri(kueri)
	panjang := utf8.RuneCountInString(kueriNormal)
	if panjang < PanjangMinimalKueri || panjang > PanjangMaksimalKueri {
		return nil, domain.InputTidakValid(fmt.Sprintf("kata pencarian harus %d sampai %d huruf", PanjangMinimalKueri, PanjangMaksimalKueri))
	}
	hasil, err := u.Pembaca.CariWilayah(ctx, kueriNormal, BatasHasilCari)
	if err != nil {
		return nil, fmt.Errorf("cari wilayah %q: %w", kueriNormal, err)
	}
	return hasil, nil
}

func (u CariWilayah) DiTitik(ctx context.Context, titik domain.Titik) (domain.WilayahDitemukan, error) {
	if !isTitikValid(titik) {
		return domain.WilayahDitemukan{}, domain.InputTidakValid("koordinat tidak valid")
	}
	wilayah, ada, err := u.Pembaca.KelurahanDiTitik(ctx, titik)
	if err != nil {
		return domain.WilayahDitemukan{}, fmt.Errorf("cari kelurahan di titik %v: %w", titik, err)
	}
	if !ada {
		return domain.WilayahDitemukan{}, domain.TidakDitemukan("lokasi berada di luar wilayah DKI Jakarta")
	}
	return wilayah, nil
}

func (u CariWilayah) Ambil(ctx context.Context, kode string) (domain.DetailWilayah, error) {
	if _, err := domain.TingkatDariKode(kode); err != nil {
		return domain.DetailWilayah{}, err
	}
	wilayah, ada, err := u.Pembaca.AmbilWilayah(ctx, kode)
	if err != nil {
		return domain.DetailWilayah{}, fmt.Errorf("ambil wilayah %s: %w", kode, err)
	}
	if !ada {
		return domain.DetailWilayah{}, domain.TidakDitemukan("wilayah tidak ditemukan")
	}
	return wilayah, nil
}

func normalkanKueri(kueri string) string {
	return strings.Join(strings.Fields(strings.ToLower(kueri)), " ")
}

func isTitikValid(titik domain.Titik) bool {
	return titik.Lat >= -batasLintang && titik.Lat <= batasLintang &&
		titik.Lng >= -batasBujur && titik.Lng <= batasBujur
}
