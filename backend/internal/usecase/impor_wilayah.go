package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobWilayah = "wilayah"

type DataSumberWilayah struct {
	Wilayah    []domain.Wilayah
	Provenance domain.Provenance
}

type DataSimpanWilayah struct {
	Provinsi   domain.Wilayah
	Wilayah    []domain.Wilayah
	Provenance domain.Provenance
}

type HasilImpor struct {
	JumlahBaru    int
	JumlahBerubah int
}

type PengambilWilayah interface {
	Ambil(ctx context.Context) (DataSumberWilayah, error)
}

type PenyimpanWilayah interface {
	SimpanWilayah(ctx context.Context, data DataSimpanWilayah) (HasilImpor, error)
}

type PencatatIngest interface {
	Mulai(ctx context.Context, job string) (int64, error)
	Selesai(ctx context.Context, id int64, hasil HasilImpor) error
	Gagal(ctx context.Context, id int64, sebab error) error
}

type IdentitasProvinsi struct {
	Kode string
	Nama string
}

type ImporWilayah struct {
	Pengambil   PengambilWilayah
	Penyimpan   PenyimpanWilayah
	Pencatat    PencatatIngest
	AturanAlias domain.AturanAlias
	Provinsi    IdentitasProvinsi
}

func (u ImporWilayah) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobWilayah)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobWilayah, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporWilayah) impor(ctx context.Context) (HasilImpor, error) {
	sumber, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil sumber wilayah: %w", err)
	}
	hasil, err := u.Penyimpan.SimpanWilayah(ctx, DataSimpanWilayah{
		Provinsi:   u.buatProvinsi(),
		Wilayah:    u.lengkapiAlias(sumber.Wilayah),
		Provenance: sumber.Provenance,
	})
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan wilayah: %w", err)
	}
	return hasil, nil
}

func (u ImporWilayah) buatProvinsi() domain.Wilayah {
	return domain.Wilayah{
		Kode:    u.Provinsi.Kode,
		Nama:    u.Provinsi.Nama,
		Tingkat: domain.TingkatProvinsi,
		Alias:   u.AturanAlias.BuatAlias(u.Provinsi.Nama, domain.TingkatProvinsi),
	}
}

func (u ImporWilayah) lengkapiAlias(daftar []domain.Wilayah) []domain.Wilayah {
	hasil := make([]domain.Wilayah, len(daftar))
	for i, w := range daftar {
		w.Alias = u.AturanAlias.BuatAlias(w.Nama, w.Tingkat)
		hasil[i] = w
	}
	return hasil
}
