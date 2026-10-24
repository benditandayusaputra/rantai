package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type RepositoriModerasi interface {
	AmbilUntukVerifikasi(ctx context.Context, id int64) (domain.LaporanUntukVerifikasi, bool, error)
	UbahStatus(ctx context.Context, perubahan domain.PerubahanStatus) error
	AntreanModerasi(ctx context.Context) (domain.AntreanModerasi, error)
	SembunyikanFoto(ctx context.Context, laporanID, fotoID int64) (bool, error)
	SembunyikanUmpan(ctx context.Context, id int64, alasan string) (bool, error)
	CatatModerasi(ctx context.Context, c domain.CatatanModerasi) error
}

type FotoModerasiBertanda struct {
	domain.FotoModerasi
	URL          string
	URLThumbnail string
}

type LaporanAntrean struct {
	domain.LaporanModerasi
	FotoBertanda []FotoModerasiBertanda
}

type AntreanModerasi struct {
	Laporan []LaporanAntrean
	Umpan   []domain.UmpanModerasi
}

type PermintaanModerasiUmpan struct {
	UmpanID     int64
	ModeratorID int64
	Alasan      string
}

type Moderasi struct {
	Repo    RepositoriModerasi
	URLFoto PembuatURLFoto
}

func (u Moderasi) Antrean(ctx context.Context) (AntreanModerasi, error) {
	antrean, err := u.Repo.AntreanModerasi(ctx)
	if err != nil {
		return AntreanModerasi{}, fmt.Errorf("baca antrean moderasi: %w", err)
	}
	laporan := make([]LaporanAntrean, len(antrean.Laporan))
	for i, l := range antrean.Laporan {
		foto, err := u.tandatangani(ctx, l.Foto)
		if err != nil {
			return AntreanModerasi{}, err
		}
		laporan[i] = LaporanAntrean{LaporanModerasi: l, FotoBertanda: foto}
	}
	return AntreanModerasi{Laporan: laporan, Umpan: antrean.Umpan}, nil
}

func (u Moderasi) tandatangani(ctx context.Context, foto []domain.FotoModerasi) ([]FotoModerasiBertanda, error) {
	hasil := make([]FotoModerasiBertanda, len(foto))
	for i, f := range foto {
		utama, err := u.URLFoto.URLBertanda(ctx, f.ObjectKey)
		if err != nil {
			return nil, err
		}
		thumbnail, err := u.URLFoto.URLBertanda(ctx, f.ThumbnailKey)
		if err != nil {
			return nil, err
		}
		hasil[i] = FotoModerasiBertanda{FotoModerasi: f, URL: utama, URLThumbnail: thumbnail}
	}
	return hasil, nil
}

func (u Moderasi) ModerasiLaporan(ctx context.Context, p domain.PermintaanModerasiLaporan) error {
	if err := p.Validasi(); err != nil {
		return err
	}
	if p.Aksi == domain.AksiSembunyikanFoto {
		return u.sembunyikanFoto(ctx, p)
	}
	laporan, err := ambilUntukVerifikasi(ctx, u.Repo, p.LaporanID)
	if err != nil {
		return err
	}
	if err := domain.Transisi(laporan.Status, p.StatusTujuan()); err != nil {
		return err
	}
	moderator := p.ModeratorID
	err = u.Repo.UbahStatus(ctx, domain.PerubahanStatus{
		LaporanID: laporan.ID, Dari: laporan.Status, Ke: p.StatusTujuan(), OlehPenggunaID: &moderator, Alasan: p.AlasanRiwayat(),
	})
	if err != nil {
		return fmt.Errorf("moderasi laporan %d: %w", p.LaporanID, err)
	}
	return u.catat(ctx, domain.CatatanModerasi{
		ModeratorID: p.ModeratorID, TargetJenis: domain.TargetLaporan, TargetID: p.LaporanID, Aksi: p.Aksi, Alasan: p.AlasanRiwayat(),
	})
}

func (u Moderasi) sembunyikanFoto(ctx context.Context, p domain.PermintaanModerasiLaporan) error {
	ada, err := u.Repo.SembunyikanFoto(ctx, p.LaporanID, p.FotoID)
	if err != nil {
		return fmt.Errorf("sembunyikan foto %d: %w", p.FotoID, err)
	}
	if !ada {
		return domain.TidakDitemukan("foto laporan tidak ditemukan")
	}
	return u.catat(ctx, domain.CatatanModerasi{
		ModeratorID: p.ModeratorID, TargetJenis: domain.TargetLaporanFoto, TargetID: p.FotoID, Aksi: p.Aksi, Alasan: p.Alasan,
	})
}

func (u Moderasi) SembunyikanUmpan(ctx context.Context, p PermintaanModerasiUmpan) error {
	if p.Alasan == "" {
		return domain.ErrAlasanModerasiKosong
	}
	ada, err := u.Repo.SembunyikanUmpan(ctx, p.UmpanID, p.Alasan)
	if err != nil {
		return fmt.Errorf("sembunyikan umpan %d: %w", p.UmpanID, err)
	}
	if !ada {
		return domain.TidakDitemukan("item umpan tidak ditemukan")
	}
	return u.catat(ctx, domain.CatatanModerasi{
		ModeratorID: p.ModeratorID, TargetJenis: domain.TargetUmpan, TargetID: p.UmpanID, Aksi: domain.AksiSembunyikanUmpan, Alasan: p.Alasan,
	})
}

func (u Moderasi) catat(ctx context.Context, c domain.CatatanModerasi) error {
	if err := u.Repo.CatatModerasi(ctx, c); err != nil {
		return fmt.Errorf("catat moderasi %s %d: %w", c.TargetJenis, c.TargetID, err)
	}
	return nil
}
