package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrLanggananPushKedaluwarsa = errors.New("langganan push sudah tidak berlaku")

type RepositoriPush interface {
	SimpanLanggananPush(ctx context.Context, penggunaID int64, l domain.LanggananPush) (int, error)
	HapusLanggananPush(ctx context.Context, penggunaID int64, endpoint string) error
	KejadianBelumDidorong(ctx context.Context) ([]domain.KejadianUntukPush, error)
	TujuanPush(ctx context.Context, kodeKecamatan string) ([]domain.TujuanPush, error)
	HapusLanggananPushByID(ctx context.Context, id int64) error
	TandaiDidorong(ctx context.Context, kejadianID int64) error
}

type PengirimPush interface {
	Kirim(ctx context.Context, l domain.LanggananPush, pesan domain.PesanPush) error
}

type LayananPush struct {
	Repo        RepositoriPush
	KunciPublik string
}

func (u LayananPush) Aktifkan(ctx context.Context, penggunaID int64, l domain.LanggananPush) (int, error) {
	if err := l.Validasi(); err != nil {
		return 0, err
	}
	jumlah, err := u.Repo.SimpanLanggananPush(ctx, penggunaID, l)
	if err != nil {
		return 0, fmt.Errorf("simpan langganan push: %w", err)
	}
	return jumlah, nil
}

func (u LayananPush) KunciPublikPush() string {
	return u.KunciPublik
}

func (u LayananPush) Matikan(ctx context.Context, penggunaID int64, endpoint string) error {
	return u.Repo.HapusLanggananPush(ctx, penggunaID, endpoint)
}

type DorongPush struct {
	Repo     RepositoriPush
	Pengirim PengirimPush
}

func (u DorongPush) Dorong(ctx context.Context) (int, error) {
	kejadian, err := u.Repo.KejadianBelumDidorong(ctx)
	if err != nil {
		return 0, fmt.Errorf("baca kejadian belum didorong: %w", err)
	}
	terkirim := 0
	for _, k := range kejadian {
		jumlah, err := u.dorongSatu(ctx, k)
		if err != nil {
			return terkirim, err
		}
		terkirim += jumlah
	}
	return terkirim, nil
}

func (u DorongPush) dorongSatu(ctx context.Context, k domain.KejadianUntukPush) (int, error) {
	tujuan, err := u.Repo.TujuanPush(ctx, k.KodeKecamatan)
	if err != nil {
		return 0, fmt.Errorf("baca tujuan push kejadian %d: %w", k.ID, err)
	}
	terkirim := 0
	pesan := domain.SusunPesanPush(k)
	for _, t := range tujuan {
		err := u.Pengirim.Kirim(ctx, t.Langganan, pesan)
		switch {
		case errors.Is(err, ErrLanggananPushKedaluwarsa):
			err = u.Repo.HapusLanggananPushByID(ctx, t.LanggananID)
		case err == nil:
			terkirim++
		}
		if err != nil {
			return terkirim, fmt.Errorf("dorong kejadian %d: %w", k.ID, err)
		}
	}
	return terkirim, u.Repo.TandaiDidorong(ctx, k.ID)
}
