package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PenyimpanKategoriSinyal interface {
	SimpanKategoriPaket(ctx context.Context, kategoriPerPaket map[string]domain.Kategori) (int, error)
	DaftarPaketSinyal(ctx context.Context) ([]domain.PaketSinyal, error)
	SimpanSinyal(ctx context.Context, daftar []domain.Sinyal, dihitungPada time.Time) (int, error)
}

type AturanTurunanAnggaran struct {
	Kategori []domain.AturanKategori
	Sinyal   domain.AturanSinyal
}

func (u TurunanPaket) golongkanKategori(ctx context.Context) (int, error) {
	penggolong, err := domain.PenggolongKategoriBaru(u.Anggaran.Kategori)
	if err != nil {
		return 0, err
	}
	masukan, err := u.TemaPaket.DaftarPaketUntukTema(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca paket untuk kategori: %w", err)
	}
	kategori := make(map[string]domain.Kategori, len(masukan))
	for _, m := range masukan {
		kategori[m.KodePaket], _ = penggolong.Golongkan(m.NamaPaket)
	}
	return u.KategoriSinyal.SimpanKategoriPaket(ctx, kategori)
}

func (u TurunanPaket) hitungSinyal(ctx context.Context) (int, error) {
	paket, err := u.KategoriSinyal.DaftarPaketSinyal(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca paket untuk sinyal: %w", err)
	}
	sekarang := u.Sekarang()
	return u.KategoriSinyal.SimpanSinyal(ctx, u.Anggaran.Sinyal.Hitung(paket, sekarang), sekarang)
}
