package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobFotoWikimedia = "foto_wikimedia"

type PengambilFotoSekolah interface {
	Ambil(ctx context.Context) ([]domain.FotoSekolah, error)
}

type PenyimpanFotoCadangan interface {
	SimpanFotoCadangan(ctx context.Context, sumber domain.SumberFoto, daftar []domain.FotoSekolah) (HasilImpor, error)
}

type ImporFotoWikimedia struct {
	Pengambil PengambilFotoSekolah
	Penyimpan PenyimpanFotoCadangan
	Pencatat  PencatatIngest
}

func (u ImporFotoWikimedia) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobFotoWikimedia)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobFotoWikimedia, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporFotoWikimedia) impor(ctx context.Context) (HasilImpor, error) {
	foto, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca daftar foto wikimedia: %w", err)
	}
	return u.Penyimpan.SimpanFotoCadangan(ctx, domain.SumberFotoWikimedia, foto)
}
