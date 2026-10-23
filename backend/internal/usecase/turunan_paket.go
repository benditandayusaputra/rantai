package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobTurunan = "turunan"

type LokasiTerhitung struct {
	KodePaket string
	Lokasi    domain.LokasiPaket
}

type PembacaGazetir interface {
	BacaGazetir(ctx context.Context) ([]domain.EntriGazetir, error)
}

type PenyimpanLokasiPaket interface {
	DaftarPaketUntukLokasi(ctx context.Context) ([]domain.MasukanLokasi, error)
	SimpanLokasiPaket(ctx context.Context, daftar []LokasiTerhitung) (int, error)
}

type TurunanPaket struct {
	Gazetir      PembacaGazetir
	Lokasi       PenyimpanLokasiPaket
	Pencatat     PencatatIngest
	AturanLokasi domain.AturanLokasi
}

func (u TurunanPaket) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobTurunan, u.hitung)
}

func (u TurunanPaket) hitung(ctx context.Context) (HasilImpor, error) {
	berubah, err := u.tentukanLokasi(ctx)
	if err != nil {
		return HasilImpor{}, err
	}
	return HasilImpor{JumlahBerubah: berubah}, nil
}

func (u TurunanPaket) tentukanLokasi(ctx context.Context) (int, error) {
	gazetir, err := u.Gazetir.BacaGazetir(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca gazetir: %w", err)
	}
	masukan, err := u.Lokasi.DaftarPaketUntukLokasi(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca paket untuk lokasi: %w", err)
	}
	penentu := domain.PenentuLokasiBaru(u.AturanLokasi, gazetir)
	hasil := make([]LokasiTerhitung, len(masukan))
	for i, m := range masukan {
		hasil[i] = LokasiTerhitung{KodePaket: m.KodePaket, Lokasi: penentu.Tentukan(m)}
	}
	return u.Lokasi.SimpanLokasiPaket(ctx, hasil)
}
