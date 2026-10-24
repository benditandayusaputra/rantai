package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type RepositoriPantau interface {
	IsWilayahAda(ctx context.Context, kode string) (bool, error)
	SimpanLangganan(ctx context.Context, penggunaID int64, kodeKecamatan string) (domain.Langganan, error)
	DaftarLangganan(ctx context.Context, penggunaID int64) ([]domain.Langganan, error)
	HapusLangganan(ctx context.Context, penggunaID, id int64) (bool, error)
	DaftarNotifikasi(ctx context.Context, penggunaID int64, batas int) (domain.HalamanNotifikasi, error)
	TandaiDibaca(ctx context.Context, penggunaID, id int64) (bool, error)
}

type PembuatKejadian interface {
	BuatKejadian(ctx context.Context, sejak time.Time) (int, error)
}

type PendorongKejadian interface {
	Dorong(ctx context.Context) (int, error)
}

type Pantau struct {
	Repo RepositoriPantau
}

func (u Pantau) Langgani(ctx context.Context, penggunaID int64, kode string) (domain.Langganan, error) {
	kodeKecamatan, err := kecamatanDariKode(kode)
	if err != nil {
		return domain.Langganan{}, err
	}
	ada, err := u.Repo.IsWilayahAda(ctx, kodeKecamatan)
	if err != nil {
		return domain.Langganan{}, fmt.Errorf("periksa wilayah %s: %w", kodeKecamatan, err)
	}
	if !ada {
		return domain.Langganan{}, domain.TidakDitemukan("kecamatan tidak ditemukan")
	}
	return u.Repo.SimpanLangganan(ctx, penggunaID, kodeKecamatan)
}

func (u Pantau) Daftar(ctx context.Context, penggunaID int64) ([]domain.Langganan, error) {
	return u.Repo.DaftarLangganan(ctx, penggunaID)
}

func (u Pantau) Berhenti(ctx context.Context, penggunaID, id int64) error {
	ada, err := u.Repo.HapusLangganan(ctx, penggunaID, id)
	if err != nil {
		return fmt.Errorf("hapus langganan %d: %w", id, err)
	}
	if !ada {
		return domain.ErrLanggananTidakDitemukan
	}
	return nil
}

func (u Pantau) Notifikasi(ctx context.Context, penggunaID int64) (domain.HalamanNotifikasi, error) {
	return u.Repo.DaftarNotifikasi(ctx, penggunaID, domain.NotifikasiPerHalaman)
}

func (u Pantau) TandaiDibaca(ctx context.Context, penggunaID, id int64) error {
	ada, err := u.Repo.TandaiDibaca(ctx, penggunaID, id)
	if err != nil {
		return fmt.Errorf("tandai notifikasi %d: %w", id, err)
	}
	if !ada {
		return domain.TidakDitemukan("notifikasi tidak ditemukan")
	}
	return nil
}
