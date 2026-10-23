package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobSIRUP = "sirup"

type RUPSumber struct {
	RUP        domain.RUP
	Provenance domain.Provenance
}

type PengambilRUP interface {
	AmbilRUP(ctx context.Context, kode string) (RUPSumber, error)
}

type PenyimpanRUP interface {
	DaftarRUPBelumLengkap(ctx context.Context) ([]string, error)
	SimpanRUP(ctx context.Context, sumber RUPSumber, instansi domain.Instansi) error
}

type ImporRUP struct {
	Pengambil      PengambilRUP
	Penyimpan      PenyimpanRUP
	Pencatat       PencatatIngest
	AturanInstansi domain.AturanInstansi
}

func (u ImporRUP) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobSIRUP, u.impor)
}

func (u ImporRUP) impor(ctx context.Context) (HasilImpor, error) {
	daftarKode, err := u.Penyimpan.DaftarRUPBelumLengkap(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("daftar rup belum lengkap: %w", err)
	}
	hasil := HasilImpor{}
	for _, kode := range daftarKode {
		sumber, err := u.Pengambil.AmbilRUP(ctx, kode)
		if errors.Is(err, domain.ErrRUPTidakTersedia) {
			continue
		}
		if err != nil {
			return hasil, fmt.Errorf("ambil rup %s: %w", kode, err)
		}
		if err := u.Penyimpan.SimpanRUP(ctx, sumber, u.AturanInstansi.Klasifikasi(sumber.RUP.SatuanKerja)); err != nil {
			return hasil, fmt.Errorf("simpan rup %s: %w", kode, err)
		}
		hasil.JumlahBerubah++
	}
	return hasil, nil
}
