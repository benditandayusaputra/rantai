package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaPendidikanKota interface {
	DaftarKota(ctx context.Context) ([]domain.RingkasanWilayah, error)
	KecamatanKota(ctx context.Context, kode string) ([]domain.RingkasanWilayah, error)
	IndikatorKota(ctx context.Context, kode, tema string) ([]domain.IndikatorTersimpan, error)
	SekolahKota(ctx context.Context, kode string) ([]domain.JumlahSekolahKecamatan, error)
}

type BacaPendidikanKota struct {
	Pembaca PembacaPendidikanKota
	Wilayah PembacaWilayah
}

func (u BacaPendidikanKota) Baca(ctx context.Context, tema, kode string) (domain.PendidikanKota, error) {
	if err := validasiKodeKota(kode); err != nil {
		return domain.PendidikanKota{}, err
	}
	if _, err := domain.PeralihanTema(tema); err != nil {
		return domain.PendidikanKota{}, err
	}
	wilayah, ada, err := u.Wilayah.AmbilWilayah(ctx, kode)
	if err != nil {
		return domain.PendidikanKota{}, err
	}
	if !ada {
		return domain.PendidikanKota{}, domain.TidakDitemukan("kabupaten/kota tidak ditemukan")
	}
	d := domain.PendidikanKota{Wilayah: wilayah.RingkasanWilayah, Tema: tema}
	if err := u.isi(ctx, &d); err != nil {
		return domain.PendidikanKota{}, fmt.Errorf("baca pendidikan kabupaten/kota: %w", err)
	}
	d.Ringkasan = domain.RingkasKapasitasKota(d.Kecamatan, d.Indikator)
	return d, nil
}

func (u BacaPendidikanKota) isi(ctx context.Context, d *domain.PendidikanKota) error {
	var err error
	if d.KabupatenKota, err = u.Pembaca.DaftarKota(ctx); err != nil {
		return err
	}
	if d.Kecamatan, err = u.Pembaca.KecamatanKota(ctx, d.Wilayah.Kode); err != nil {
		return err
	}
	if d.Indikator, err = u.Pembaca.IndikatorKota(ctx, d.Wilayah.Kode, d.Tema); err != nil {
		return err
	}
	if len(d.Indikator) > 0 {
		d.Periode = d.Indikator[0].Periode
	}
	d.Sekolah, err = u.Pembaca.SekolahKota(ctx, d.Wilayah.Kode)
	return err
}

func validasiKodeKota(kode string) error {
	tingkat, err := domain.TingkatDariKode(kode)
	if err != nil {
		return err
	}
	if tingkat != domain.TingkatKota {
		return domain.InputTidakValid("kode harus kode kabupaten/kota: " + kode)
	}
	return nil
}
