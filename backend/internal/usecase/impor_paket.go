package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	JobLPSEDaftar = "lpse_daftar"
	JobLPSEDetail = "lpse_detail"
)

type DaftarPaketSumber struct {
	Paket      []domain.PaketDaftar
	Provenance domain.Provenance
}

type DetailPaketSumber struct {
	Detail     domain.DetailPaket
	Provenance domain.Provenance
}

type DataSimpanDetailPaket struct {
	Detail     domain.DetailPaket
	Instansi   domain.Instansi
	Provenance domain.Provenance
}

type PengambilDaftarPaket interface {
	AmbilDaftarPaket(ctx context.Context, tahun int) (DaftarPaketSumber, error)
}

type PenyimpanDaftarPaket interface {
	SimpanDaftarPaket(ctx context.Context, data DaftarPaketSumber) (HasilImpor, error)
}

type PengambilDetailPaket interface {
	AmbilDetailPaket(ctx context.Context, kodeSPSE string) (DetailPaketSumber, error)
}

type PenyimpanDetailPaket interface {
	DaftarPaketPerluDetail(ctx context.Context, jenisPengadaan []string) ([]string, error)
	SimpanDetailPaket(ctx context.Context, data DataSimpanDetailPaket) error
}

type ImporDaftarPaket struct {
	Pengambil PengambilDaftarPaket
	Penyimpan PenyimpanDaftarPaket
	Pencatat  PencatatIngest
	Tahun     []int
}

type ImporDetailPaket struct {
	Pengambil      PengambilDetailPaket
	Penyimpan      PenyimpanDetailPaket
	Pencatat       PencatatIngest
	AturanInstansi domain.AturanInstansi
	JenisPengadaan []string
}

func (u ImporDaftarPaket) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobLPSEDaftar, u.impor)
}

func (u ImporDaftarPaket) impor(ctx context.Context) (HasilImpor, error) {
	total := HasilImpor{}
	for _, tahun := range u.Tahun {
		sumber, err := u.Pengambil.AmbilDaftarPaket(ctx, tahun)
		if err != nil {
			return total, fmt.Errorf("ambil daftar paket %d: %w", tahun, err)
		}
		hasil, err := u.Penyimpan.SimpanDaftarPaket(ctx, sumber)
		if err != nil {
			return total, fmt.Errorf("simpan daftar paket %d: %w", tahun, err)
		}
		total = total.tambah(hasil)
	}
	return total, nil
}

func (u ImporDetailPaket) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobLPSEDetail, u.impor)
}

func (u ImporDetailPaket) impor(ctx context.Context) (HasilImpor, error) {
	daftarKodeSPSE, err := u.Penyimpan.DaftarPaketPerluDetail(ctx, u.JenisPengadaan)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("daftar paket perlu detail: %w", err)
	}
	hasil := HasilImpor{}
	for _, kodeSPSE := range daftarKodeSPSE {
		if err := u.imporSatu(ctx, kodeSPSE); err != nil {
			return hasil, err
		}
		hasil.JumlahBerubah++
	}
	return hasil, nil
}

func (u ImporDetailPaket) imporSatu(ctx context.Context, kodeSPSE string) error {
	sumber, err := u.Pengambil.AmbilDetailPaket(ctx, kodeSPSE)
	if err != nil {
		return fmt.Errorf("ambil detail paket non-tender %s: %w", kodeSPSE, err)
	}
	err = u.Penyimpan.SimpanDetailPaket(ctx, DataSimpanDetailPaket{
		Detail:     sumber.Detail,
		Instansi:   u.AturanInstansi.Klasifikasi(sumber.Detail.SatuanKerja),
		Provenance: sumber.Provenance,
	})
	if err != nil {
		return fmt.Errorf("simpan detail paket non-tender %s: %w", kodeSPSE, err)
	}
	return nil
}

func (h HasilImpor) tambah(lain HasilImpor) HasilImpor {
	return HasilImpor{JumlahBaru: h.JumlahBaru + lain.JumlahBaru, JumlahBerubah: h.JumlahBerubah + lain.JumlahBerubah}
}

func jalankanTercatat(ctx context.Context, pencatat PencatatIngest, job string, impor func(context.Context) (HasilImpor, error)) error {
	id, err := pencatat.Mulai(ctx, job)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", job, err)
	}
	hasil, err := impor(ctx)
	if err != nil {
		return errors.Join(err, pencatat.Gagal(ctx, id, err))
	}
	return pencatat.Selesai(ctx, id, hasil)
}
