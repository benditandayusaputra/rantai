package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaLaporan interface {
	RepositoriVerifikasi
	AmbilLaporanPublik(ctx context.Context, id int64) (domain.LaporanPublik, bool, error)
	PetaLaporan(ctx context.Context, kotak domain.Kotak) ([]domain.TitikLaporan, error)
}

type PembuatURLFoto interface {
	URLBertanda(ctx context.Context, kunci string) (string, error)
}

type FotoBertanda struct {
	domain.FotoTersimpan
	URL          string
	URLThumbnail string
}

type DetailLaporan struct {
	domain.LaporanPublik
	FotoBertanda []FotoBertanda
}

type BacaLaporan struct {
	Pembaca PembacaLaporan
	URLFoto PembuatURLFoto
}

func (u BacaLaporan) Ambil(ctx context.Context, id int64) (DetailLaporan, error) {
	laporan, ada, err := u.Pembaca.AmbilLaporanPublik(ctx, id)
	if err != nil {
		return DetailLaporan{}, fmt.Errorf("ambil laporan %d: %w", id, err)
	}
	if !ada {
		return DetailLaporan{}, domain.TidakDitemukan("laporan tidak ditemukan")
	}
	verifikasi, err := ambilUntukVerifikasi(ctx, u.Pembaca, id)
	if err != nil {
		return DetailLaporan{}, err
	}
	laporan.JumlahSaksi = verifikasi.JumlahSaksiSah()
	foto, err := u.tandatanganiFoto(ctx, laporan.Foto)
	return DetailLaporan{LaporanPublik: laporan, FotoBertanda: foto}, err
}

func (u BacaLaporan) tandatanganiFoto(ctx context.Context, foto []domain.FotoTersimpan) ([]FotoBertanda, error) {
	hasil := make([]FotoBertanda, len(foto))
	for i, f := range foto {
		utama, err := u.URLFoto.URLBertanda(ctx, f.ObjectKey)
		if err != nil {
			return nil, err
		}
		thumbnail, err := u.URLFoto.URLBertanda(ctx, f.ThumbnailKey)
		if err != nil {
			return nil, err
		}
		hasil[i] = FotoBertanda{FotoTersimpan: f, URL: utama, URLThumbnail: thumbnail}
	}
	return hasil, nil
}

func (u BacaLaporan) Peta(ctx context.Context, kotak domain.Kotak) ([]domain.TitikLaporan, error) {
	if err := kotak.Periksa(domain.LuasKotakPetaLaporanMaks); err != nil {
		return nil, err
	}
	titik, err := u.Pembaca.PetaLaporan(ctx, kotak)
	if err != nil {
		return nil, fmt.Errorf("baca peta laporan: %w", err)
	}
	return titik, nil
}
