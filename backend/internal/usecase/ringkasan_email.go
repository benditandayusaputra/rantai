package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobRingkasanMingguan = "ringkasan_mingguan"

type RepositoriRingkasanEmail interface {
	IsRingkasanEmailAktif(ctx context.Context, penggunaID int64) (bool, error)
	AktifkanRingkasanEmail(ctx context.Context, penggunaID int64) error
	MatikanRingkasanEmail(ctx context.Context, penggunaID int64) error
	PenerimaRingkasan(ctx context.Context, terkirimSebelum time.Time) ([]domain.PenerimaRingkasan, error)
	KabarRingkasan(ctx context.Context, penggunaID int64, sejak time.Time, batas int) ([]domain.Notifikasi, error)
	TandaiRingkasanTerkirim(ctx context.Context, penggunaID int64, pada time.Time) error
}

type PengirimEmail interface {
	Kirim(ctx context.Context, email, subjek, isi string) error
}

type RingkasanEmail struct {
	Repo RepositoriRingkasanEmail
}

func (u RingkasanEmail) IsAktif(ctx context.Context, penggunaID int64) (bool, error) {
	return u.Repo.IsRingkasanEmailAktif(ctx, penggunaID)
}

func (u RingkasanEmail) Aktifkan(ctx context.Context, penggunaID int64) error {
	return u.Repo.AktifkanRingkasanEmail(ctx, penggunaID)
}

func (u RingkasanEmail) Matikan(ctx context.Context, penggunaID int64) error {
	return u.Repo.MatikanRingkasanEmail(ctx, penggunaID)
}

type KirimRingkasanMingguan struct {
	Repo        RepositoriRingkasanEmail
	Pengirim    PengirimEmail
	AlamatSitus string
	Sekarang    func() time.Time
}

func (u KirimRingkasanMingguan) Jalankan(ctx context.Context) (int, error) {
	sekarang := u.Sekarang()
	sejak := sekarang.Add(-domain.JendelaKejadian)
	penerima, err := u.Repo.PenerimaRingkasan(ctx, sejak)
	if err != nil {
		return 0, fmt.Errorf("baca penerima ringkasan: %w", err)
	}
	terkirim := 0
	for _, p := range penerima {
		isTerkirim, err := u.kirimSatu(ctx, p, sejak)
		if err != nil {
			return terkirim, err
		}
		if isTerkirim {
			terkirim++
		}
		if err := u.Repo.TandaiRingkasanTerkirim(ctx, p.PenggunaID, sekarang); err != nil {
			return terkirim, fmt.Errorf("tandai ringkasan pengguna %d: %w", p.PenggunaID, err)
		}
	}
	return terkirim, nil
}

func (u KirimRingkasanMingguan) kirimSatu(ctx context.Context, p domain.PenerimaRingkasan, sejak time.Time) (bool, error) {
	kabar, err := u.Repo.KabarRingkasan(ctx, p.PenggunaID, sejak, domain.BatasKabarRingkasan)
	if err != nil {
		return false, fmt.Errorf("baca kabar ringkasan pengguna %d: %w", p.PenggunaID, err)
	}
	if len(kabar) == 0 {
		return false, nil
	}
	email := domain.SusunEmailRingkasan(kabar, u.AlamatSitus)
	if err := u.Pengirim.Kirim(ctx, p.Email, email.Subjek, email.Isi); err != nil {
		return false, fmt.Errorf("kirim ringkasan pengguna %d: %w", p.PenggunaID, err)
	}
	return true, nil
}
