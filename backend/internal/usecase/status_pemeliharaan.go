package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type JadwalPaketTersimpan struct {
	JenisPengadaan string
	Jadwal         domain.JadwalPaket
}

type PembacaJadwalPaket interface {
	AmbilJadwalPaket(ctx context.Context, kode string) (JadwalPaketTersimpan, bool, error)
}

type KonfigurasiPemeliharaan struct {
	domain.AturanPemeliharaan `yaml:",inline"`
	JenisPengadaan            []string `yaml:"jenis_pengadaan"`
}

type StatusPemeliharaan struct {
	Pembaca     PembacaJadwalPaket
	Konfigurasi KonfigurasiPemeliharaan
	Sekarang    func() time.Time
}

type HasilPemeliharaan struct {
	IsBerlaku bool
	Estimasi  domain.EstimasiPemeliharaan
}

func (u StatusPemeliharaan) Hitung(ctx context.Context, kode string) (HasilPemeliharaan, error) {
	tersimpan, ada, err := u.Pembaca.AmbilJadwalPaket(ctx, kode)
	if err != nil {
		return HasilPemeliharaan{}, fmt.Errorf("ambil jadwal paket %s: %w", kode, err)
	}
	if !ada {
		return HasilPemeliharaan{}, domain.TidakDitemukan("paket tidak ditemukan")
	}
	if !slices.Contains(u.Konfigurasi.JenisPengadaan, tersimpan.JenisPengadaan) {
		return HasilPemeliharaan{}, nil
	}
	return HasilPemeliharaan{IsBerlaku: true, Estimasi: u.Konfigurasi.Estimasi(tersimpan.Jadwal, u.Sekarang())}, nil
}
