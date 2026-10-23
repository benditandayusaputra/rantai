package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobOpentender = "opentender"

type RiwayatSumber struct {
	Riwayat    []domain.RiwayatPenyedia
	Provenance domain.Provenance
}

type PengambilRiwayat interface {
	AmbilRiwayat(ctx context.Context, tahun int) (RiwayatSumber, error)
}

type PenyimpanRiwayat interface {
	SimpanRiwayat(ctx context.Context, sumber RiwayatSumber, penyedia map[string]domain.Penyedia) (HasilImpor, error)
}

type ImporRiwayatPenyedia struct {
	Pengambil PengambilRiwayat
	Penyimpan PenyimpanRiwayat
	Pencatat  PencatatIngest
	Tahun     []int
}

func (u ImporRiwayatPenyedia) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobOpentender, u.impor)
}

func (u ImporRiwayatPenyedia) impor(ctx context.Context) (HasilImpor, error) {
	total := HasilImpor{}
	for _, tahun := range u.Tahun {
		sumber, err := u.Pengambil.AmbilRiwayat(ctx, tahun)
		if err != nil {
			return total, fmt.Errorf("ambil riwayat opentender %d: %w", tahun, err)
		}
		hasil, err := u.Penyimpan.SimpanRiwayat(ctx, sumber, penyediaRiwayat(sumber.Riwayat))
		if err != nil {
			return total, fmt.Errorf("simpan riwayat opentender %d: %w", tahun, err)
		}
		total = total.tambah(hasil)
	}
	return total, nil
}

func penyediaRiwayat(riwayat []domain.RiwayatPenyedia) map[string]domain.Penyedia {
	penyedia := map[string]domain.Penyedia{}
	for _, r := range riwayat {
		penyedia[r.NamaPenyedia] = domain.Penyedia{Nama: r.NamaPenyedia, NamaNormal: domain.NormalisasiNamaPenyedia(r.NamaPenyedia)}
	}
	return penyedia
}
