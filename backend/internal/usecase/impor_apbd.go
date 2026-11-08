package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobAPBD = "apbd_djpk"

type APBDSumber struct {
	APBD       domain.APBDDaerah
	Provenance domain.Provenance
}

type PengambilAPBD interface {
	AmbilAPBD(ctx context.Context, tahun, bulanTerakhir int) (APBDSumber, error)
}

type PenyimpanAPBD interface {
	BulanRealisasiTersimpan(ctx context.Context, tahun int) (int, error)
	SimpanAPBD(ctx context.Context, sumber APBDSumber) error
}

type ImporAPBD struct {
	Pengambil PengambilAPBD
	Penyimpan PenyimpanAPBD
	Pencatat  PencatatIngest
	Tahun     []int
	Sekarang  func() time.Time
}

func (u ImporAPBD) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobAPBD, u.impor)
}

func (u ImporAPBD) impor(ctx context.Context) (HasilImpor, error) {
	hasil := HasilImpor{}
	for _, tahun := range u.Tahun {
		berubah, err := u.imporTahun(ctx, tahun)
		if err != nil {
			return hasil, fmt.Errorf("impor apbd %d: %w", tahun, err)
		}
		if berubah {
			hasil.JumlahBerubah++
		}
	}
	return hasil, nil
}

func (u ImporAPBD) imporTahun(ctx context.Context, tahun int) (bool, error) {
	bulanTerakhir := bulanTerakhirAPBD(tahun, u.Sekarang())
	tersimpan, err := u.Penyimpan.BulanRealisasiTersimpan(ctx, tahun)
	if err != nil || bulanTerakhir == 0 || tersimpan == domain.BulanPenuhAPBD {
		return false, err
	}
	sumber, err := u.Pengambil.AmbilAPBD(ctx, tahun, bulanTerakhir)
	if err != nil {
		return false, err
	}
	return true, u.Penyimpan.SimpanAPBD(ctx, sumber)
}

func bulanTerakhirAPBD(tahun int, sekarang time.Time) int {
	switch {
	case tahun < sekarang.Year():
		return domain.BulanPenuhAPBD
	case tahun == sekarang.Year():
		return int(sekarang.Month())
	default:
		return 0
	}
}
