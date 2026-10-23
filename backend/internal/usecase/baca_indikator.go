package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaTampilanPendidikan interface {
	FiturPetaTema(ctx context.Context, temaKode string) ([]domain.FiturPetaIndikator, error)
	DetailPendidikanKecamatan(ctx context.Context, kodeKecamatan string) (domain.DetailPendidikanKecamatan, bool, error)
	SekolahMenengahPertama(ctx context.Context, kodeKecamatan string) ([]domain.SekolahDenganKursi, error)
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

func (u BacaIndikator) DetailPendidikan(ctx context.Context, kodeKecamatan string) (domain.DetailPendidikanKecamatan, error) {
	if err := validasiKodeKecamatan(kodeKecamatan); err != nil {
		return domain.DetailPendidikanKecamatan{}, err
	}
	detail, ada, err := u.Pembaca.DetailPendidikanKecamatan(ctx, kodeKecamatan)
	if err != nil {
		return domain.DetailPendidikanKecamatan{}, fmt.Errorf("ambil detail pendidikan %s: %w", kodeKecamatan, err)
	}
	if !ada {
		return domain.DetailPendidikanKecamatan{}, domain.TidakDitemukan("data pendidikan kecamatan tidak ditemukan")
	}
	return detail, nil
}

func (u BacaIndikator) SekolahKecamatan(ctx context.Context, kodeKecamatan string) ([]domain.SekolahDenganKursi, error) {
	if err := validasiKodeKecamatan(kodeKecamatan); err != nil {
		return nil, err
	}
	daftar, err := u.Pembaca.SekolahMenengahPertama(ctx, kodeKecamatan)
	if err != nil {
		return nil, fmt.Errorf("ambil sekolah kecamatan %s: %w", kodeKecamatan, err)
	}
	return daftar, nil
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
