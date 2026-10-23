package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobRingkasUmpan = "ringkas_umpan"

const RingkasanDitolak = ""

type PenyimpanRingkasan interface {
	UmpanBelumDiringkas(ctx context.Context, batas int) ([]domain.UmpanUntukRingkas, error)
	SimpanRingkasan(ctx context.Context, id int64, ringkasan string) error
}

type PenulisRingkasan interface {
	Ringkas(ctx context.Context, judul, cuplikan string) (string, error)
}

type RingkasUmpan struct {
	Umpan    PenyimpanRingkasan
	Penulis  PenulisRingkasan
	Pencatat PencatatIngest
	Aturan   domain.AturanRingkasan
	Batas    int
	Laporkan func(pesan string)
}

func (u RingkasUmpan) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobRingkasUmpan)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobRingkasUmpan, err)
	}
	hasil, err := u.ringkasSemua(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u RingkasUmpan) ringkasSemua(ctx context.Context) (HasilImpor, error) {
	daftar, err := u.Umpan.UmpanBelumDiringkas(ctx, u.Batas)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca umpan belum diringkas: %w", err)
	}
	hasil := HasilImpor{}
	for _, umpan := range daftar {
		ringkasan, err := u.ringkasSatu(ctx, umpan)
		if err != nil {
			return hasil, err
		}
		if err := u.Umpan.SimpanRingkasan(ctx, umpan.ID, ringkasan); err != nil {
			return hasil, fmt.Errorf("simpan ringkasan umpan %d: %w", umpan.ID, err)
		}
		hasil.JumlahBerubah++
	}
	return hasil, nil
}

func (u RingkasUmpan) ringkasSatu(ctx context.Context, umpan domain.UmpanUntukRingkas) (string, error) {
	ringkasan, err := u.Penulis.Ringkas(ctx, umpan.Judul, umpan.Cuplikan)
	if err != nil {
		return "", fmt.Errorf("ringkas umpan %d: %w", umpan.ID, err)
	}
	if err := u.Aturan.Validasi(ringkasan, umpan); err != nil {
		u.laporkan(fmt.Sprintf("%s umpan %d: ringkasan ditolak, hanya judul yang tampil: %v", JobRingkasUmpan, umpan.ID, err))
		return RingkasanDitolak, nil
	}
	return ringkasan, nil
}

func (u RingkasUmpan) laporkan(pesan string) {
	if u.Laporkan != nil {
		u.Laporkan(pesan)
	}
}
