package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaTampilanPendidikan interface {
	FiturPetaTema(ctx context.Context, temaKode string) ([]domain.FiturPetaIndikator, error)
	DetailPendidikanKecamatan(ctx context.Context, peralihan domain.Peralihan, kodeKecamatan string) (domain.DetailPendidikanKecamatan, bool, error)
	TrenRasio(ctx context.Context, temaKode, kodeKecamatan string) ([]domain.TitikTren, error)
	SekolahJenjang(ctx context.Context, jenjang domain.Jenjang, kodeKecamatan string) ([]domain.SekolahDenganKursi, error)
}

type BacaIndikator struct {
	Pembaca PembacaTampilanPendidikan
	Tema    PembacaTema
}

func (u BacaIndikator) PetaTema(ctx context.Context, temaKode string) (domain.PetaTema, error) {
	tema, err := BacaTema{Pembaca: u.Tema}.Ambil(ctx, temaKode)
	if err != nil {
		return domain.PetaTema{}, err
	}
	fitur, err := u.Pembaca.FiturPetaTema(ctx, temaKode)
	if err != nil {
		return domain.PetaTema{}, fmt.Errorf("ambil peta tema %s: %w", temaKode, err)
	}
	return domain.PetaTema{Tema: tema, Fitur: fitur}, nil
}

func (u BacaIndikator) DetailPendidikan(ctx context.Context, temaKode, kodeKecamatan string) (domain.DetailPendidikanKecamatan, error) {
	peralihan, err := peralihanKecamatan(temaKode, kodeKecamatan)
	if err != nil {
		return domain.DetailPendidikanKecamatan{}, err
	}
	tema, err := BacaTema{Pembaca: u.Tema}.Ambil(ctx, temaKode)
	if err != nil {
		return domain.DetailPendidikanKecamatan{}, err
	}
	detail, ada, err := u.Pembaca.DetailPendidikanKecamatan(ctx, peralihan, kodeKecamatan)
	if err != nil {
		return domain.DetailPendidikanKecamatan{}, fmt.Errorf("ambil detail %s %s: %w", temaKode, kodeKecamatan, err)
	}
	if !ada {
		return domain.DetailPendidikanKecamatan{}, domain.TidakDitemukan("data pendidikan kecamatan tidak ditemukan")
	}
	tren, err := u.Pembaca.TrenRasio(ctx, temaKode, kodeKecamatan)
	if err != nil {
		return domain.DetailPendidikanKecamatan{}, fmt.Errorf("ambil tren %s %s: %w", temaKode, kodeKecamatan, err)
	}
	detail.Peralihan, detail.Ambang, detail.Tren = peralihan, tema.Indikator.AmbangRasio, domain.TrenLayakTampil(tren)
	return detail, nil
}

func (u BacaIndikator) SekolahKecamatan(ctx context.Context, temaKode, kodeKecamatan string) ([]domain.SekolahDenganKursi, error) {
	peralihan, err := peralihanKecamatan(temaKode, kodeKecamatan)
	if err != nil {
		return nil, err
	}
	daftar, err := u.Pembaca.SekolahJenjang(ctx, peralihan.JenjangTujuan, kodeKecamatan)
	if err != nil {
		return nil, fmt.Errorf("ambil sekolah %s kecamatan %s: %w", peralihan.JenjangTujuan, kodeKecamatan, err)
	}
	return daftar, nil
}

func peralihanKecamatan(temaKode, kodeKecamatan string) (domain.Peralihan, error) {
	if err := validasiKodeKecamatan(kodeKecamatan); err != nil {
		return domain.Peralihan{}, err
	}
	return domain.PeralihanTema(temaKode)
}

func validasiKodeKecamatan(kode string) error {
	tingkat, err := domain.TingkatDariKode(kode)
	if err != nil {
		return err
	}
	if tingkat != domain.TingkatKecamatan {
		return domain.InputTidakValid("kode harus kode kecamatan: " + kode)
	}
	return nil
}
