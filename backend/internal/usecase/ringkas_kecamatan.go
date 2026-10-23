package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaIndikator interface {
	IndikatorTerbaru(ctx context.Context, kodeTema, kodeKecamatan string) (domain.IndikatorKecamatan, bool, error)
}

type PembacaAnggaranKecamatan interface {
	RingkasanAnggaran(ctx context.Context, kodeKecamatan string, tahun int) (domain.RingkasanAnggaran, bool, error)
}

type PembacaWakilKecamatan interface {
	WakilKecamatanKomisi(ctx context.Context, kodeKecamatan string, komisiPengawas []string) (domain.WakilKecamatan, error)
}

type RingkasKecamatan struct {
	Wilayah   PembacaWilayah
	Tema      PembacaDaftarTema
	Indikator PembacaIndikator
	Anggaran  PembacaAnggaranKecamatan
	Wakil     PembacaWakilKecamatan
	Sekarang  func() time.Time
}

func (u RingkasKecamatan) Ringkas(ctx context.Context, kode string) (domain.RingkasanKecamatan, error) {
	ringkasan, err := u.wilayah(ctx, kode)
	if err != nil {
		return domain.RingkasanKecamatan{}, err
	}
	daftarTema, err := u.Tema.Daftar(ctx)
	if err != nil {
		return domain.RingkasanKecamatan{}, fmt.Errorf("daftar tema aktif: %w", err)
	}
	kodeKecamatan := ringkasan.Kecamatan.Kode
	if ringkasan.Tema, err = u.temaKecamatan(ctx, daftarTema, kodeKecamatan); err != nil {
		return domain.RingkasanKecamatan{}, err
	}
	if ringkasan.Anggaran, err = u.anggaran(ctx, kodeKecamatan); err != nil {
		return domain.RingkasanKecamatan{}, err
	}
	ringkasan.Wakil, err = u.Wakil.WakilKecamatanKomisi(ctx, kodeKecamatan, domain.KomisiPengawasSemuaTema(daftarTema))
	return ringkasan, err
}

func (u RingkasKecamatan) wilayah(ctx context.Context, kode string) (domain.RingkasanKecamatan, error) {
	kodeKecamatan, err := kecamatanDariKode(kode)
	if err != nil {
		return domain.RingkasanKecamatan{}, err
	}
	ringkasan := domain.RingkasanKecamatan{}
	if kodeKecamatan != kode {
		kelurahan, err := CariWilayah{Pembaca: u.Wilayah}.Ambil(ctx, kode)
		if err != nil {
			return domain.RingkasanKecamatan{}, err
		}
		ringkasan.KelurahanDisorot = &kelurahan.RingkasanWilayah
	}
	ringkasan.Kecamatan, err = CariWilayah{Pembaca: u.Wilayah}.Ambil(ctx, kodeKecamatan)
	return ringkasan, err
}

func (u RingkasKecamatan) temaKecamatan(ctx context.Context, daftar []domain.TemaTersimpan, kodeKecamatan string) ([]domain.TemaKecamatan, error) {
	hasil := []domain.TemaKecamatan{}
	for _, t := range daftar {
		indikator, ada, err := u.Indikator.IndikatorTerbaru(ctx, t.Kode, kodeKecamatan)
		if err != nil {
			return nil, fmt.Errorf("indikator %s kecamatan %s: %w", t.Kode, kodeKecamatan, err)
		}
		temaKecamatan := domain.TemaKecamatan{Tema: t}
		if ada {
			temaKecamatan.Indikator = &indikator
		}
		hasil = append(hasil, temaKecamatan)
	}
	return hasil, nil
}

func (u RingkasKecamatan) anggaran(ctx context.Context, kodeKecamatan string) (*domain.RingkasanAnggaran, error) {
	anggaran, ada, err := u.Anggaran.RingkasanAnggaran(ctx, kodeKecamatan, u.Sekarang().Year())
	if err != nil {
		return nil, fmt.Errorf("ringkasan anggaran kecamatan %s: %w", kodeKecamatan, err)
	}
	if !ada {
		return nil, nil
	}
	return &anggaran, nil
}
