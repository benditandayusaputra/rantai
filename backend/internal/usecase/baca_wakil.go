package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaWakil interface {
	DapilKecamatan(ctx context.Context, kodeKecamatan string) (domain.RingkasanWilayah, domain.RingkasanDapil, bool, error)
	AmbilDapil(ctx context.Context, id int) (domain.RingkasanDapil, bool, error)
	IsiDapil(ctx context.Context, id int) ([]domain.RingkasanWilayah, []domain.Wakil, error)
	AmbilAnggota(ctx context.Context, id string) (domain.ProfilAnggota, bool, error)
}

type PencariTema interface {
	Ambil(ctx context.Context, kode string) (domain.TemaTersimpan, error)
}

type BacaWakil struct {
	Pembaca PembacaWakil
	Tema    PencariTema
}

func (u BacaWakil) WakilKecamatan(ctx context.Context, kode, kodeTema string) (domain.WakilKecamatan, error) {
	kodeKecamatan, err := kecamatanDariKode(kode)
	if err != nil {
		return domain.WakilKecamatan{}, err
	}
	komisiPengawas, err := u.komisiPengawas(ctx, kodeTema)
	if err != nil {
		return domain.WakilKecamatan{}, err
	}
	kecamatan, dapil, ada, err := u.Pembaca.DapilKecamatan(ctx, kodeKecamatan)
	if err != nil {
		return domain.WakilKecamatan{}, fmt.Errorf("ambil dapil kecamatan %s: %w", kodeKecamatan, err)
	}
	if !ada {
		return domain.WakilKecamatan{}, domain.TidakDitemukan("kecamatan tidak ditemukan atau belum punya data dapil")
	}
	detail, err := u.lengkapiDapil(ctx, dapil, komisiPengawas)
	if err != nil {
		return domain.WakilKecamatan{}, err
	}
	return domain.WakilKecamatan{Kecamatan: kecamatan, KomisiPengawas: komisiPengawas, DetailDapil: detail}, nil
}

func (u BacaWakil) Dapil(ctx context.Context, id int) (domain.DetailDapil, error) {
	dapil, ada, err := u.Pembaca.AmbilDapil(ctx, id)
	if err != nil {
		return domain.DetailDapil{}, fmt.Errorf("ambil dapil %d: %w", id, err)
	}
	if !ada {
		return domain.DetailDapil{}, domain.TidakDitemukan("dapil tidak ditemukan")
	}
	return u.lengkapiDapil(ctx, dapil, []string{})
}

func (u BacaWakil) Anggota(ctx context.Context, id string) (domain.ProfilAnggota, error) {
	profil, ada, err := u.Pembaca.AmbilAnggota(ctx, id)
	if err != nil {
		return domain.ProfilAnggota{}, fmt.Errorf("ambil anggota %s: %w", id, err)
	}
	if !ada {
		return domain.ProfilAnggota{}, domain.TidakDitemukan("anggota tidak ditemukan")
	}
	profil.Dapil, err = u.lengkapiDapil(ctx, profil.Dapil.RingkasanDapil, []string{})
	return profil, err
}

func (u BacaWakil) lengkapiDapil(ctx context.Context, dapil domain.RingkasanDapil, komisiPengawas []string) (domain.DetailDapil, error) {
	kecamatan, anggota, err := u.Pembaca.IsiDapil(ctx, dapil.ID)
	if err != nil {
		return domain.DetailDapil{}, fmt.Errorf("isi dapil %d: %w", dapil.ID, err)
	}
	return domain.DetailDapil{
		RingkasanDapil: dapil,
		Kecamatan:      kecamatan,
		Anggota:        domain.UrutkanWakil(anggota, komisiPengawas),
	}, nil
}

func (u BacaWakil) komisiPengawas(ctx context.Context, kodeTema string) ([]string, error) {
	if kodeTema == "" {
		return []string{}, nil
	}
	tema, err := u.Tema.Ambil(ctx, kodeTema)
	if err != nil {
		return nil, err
	}
	return domain.KomisiPengawas(tema.Tema), nil
}

func kecamatanDariKode(kode string) (string, error) {
	tingkat, err := domain.TingkatDariKode(kode)
	switch {
	case err != nil:
		return "", err
	case tingkat == domain.TingkatKelurahan:
		return domain.IndukKode(kode), nil
	case tingkat != domain.TingkatKecamatan:
		return "", domain.InputTidakValid("kode harus kode kecamatan atau kelurahan")
	default:
		return kode, nil
	}
}
