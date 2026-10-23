package usecase

import (
	"context"
	"fmt"
	"time"

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

type PenyimpanTemaPaket interface {
	DaftarPaketUntukTema(ctx context.Context) ([]domain.MasukanTema, error)
	SimpanTemaPaket(ctx context.Context, temaPerPaket map[string][]string) (int, error)
}

type TurunanPaket struct {
	Gazetir        PembacaGazetir
	Lokasi         PenyimpanLokasiPaket
	Tema           PembacaTema
	TemaPaket      PenyimpanTemaPaket
	KategoriSinyal PenyimpanKategoriSinyal
	Pencatat       PencatatIngest
	AturanLokasi   domain.AturanLokasi
	Anggaran       AturanTurunanAnggaran
	Sekarang       func() time.Time
}

func (u TurunanPaket) Jalankan(ctx context.Context) error {
	return jalankanTercatat(ctx, u.Pencatat, JobTurunan, u.hitung)
}

func (u TurunanPaket) hitung(ctx context.Context) (HasilImpor, error) {
	total := 0
	for _, langkah := range []func(context.Context) (int, error){u.tentukanLokasi, u.tandaiTema, u.golongkanKategori, u.hitungSinyal} {
		berubah, err := langkah(ctx)
		if err != nil {
			return HasilImpor{}, err
		}
		total += berubah
	}
	return HasilImpor{JumlahBerubah: total}, nil
}

func (u TurunanPaket) tandaiTema(ctx context.Context) (int, error) {
	daftarTema, err := u.Tema.DaftarTemaAktif(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca tema aktif: %w", err)
	}
	penanda, err := domain.PenandaTemaBaru(temaDari(daftarTema))
	if err != nil {
		return 0, err
	}
	masukan, err := u.TemaPaket.DaftarPaketUntukTema(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca paket untuk tema: %w", err)
	}
	temaPerPaket := make(map[string][]string, len(masukan))
	for _, m := range masukan {
		temaPerPaket[m.KodePaket] = penanda.Tandai(m)
	}
	return u.TemaPaket.SimpanTemaPaket(ctx, temaPerPaket)
}

func temaDari(daftar []domain.TemaTersimpan) []domain.Tema {
	hasil := make([]domain.Tema, len(daftar))
	for i, t := range daftar {
		hasil[i] = t.Tema
	}
	return hasil
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
