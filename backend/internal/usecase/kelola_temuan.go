package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type RepositoriTemuan interface {
	IsPaketAda(ctx context.Context, kode string) (bool, error)
	HitungTemuanSejak(ctx context.Context, penggunaID int64, sejak time.Time) (int, error)
	SimpanTemuan(ctx context.Context, masuk domain.TemuanMasuk) (int64, error)
	AmbilPengirimTemuan(ctx context.Context, id int64) (int64, bool, error)
	SimpanTanggapan(ctx context.Context, masuk domain.TanggapanMasuk) error
	DaftarTemuan(ctx context.Context, paketKode string, batas int) ([]domain.Temuan, error)
	SembunyikanTemuan(ctx context.Context, id int64) (bool, error)
	CatatModerasi(ctx context.Context, c domain.CatatanModerasi) error
}

type PermintaanModerasiTemuan struct {
	TemuanID    int64
	ModeratorID int64
	Alasan      string
}

type KelolaTemuan struct {
	Repo      RepositoriTemuan
	Terlarang domain.KosakataTerlarang
	Sekarang  func() time.Time
}

func (u KelolaTemuan) Daftar(ctx context.Context, paketKode string) ([]domain.Temuan, error) {
	daftar, err := u.Repo.DaftarTemuan(ctx, paketKode, domain.BatasTemuanPerPaket)
	if err != nil {
		return nil, fmt.Errorf("daftar temuan paket %s: %w", paketKode, err)
	}
	return daftar, nil
}

func (u KelolaTemuan) Kirim(ctx context.Context, masuk domain.TemuanMasuk) (int64, error) {
	masuk.Isi = strings.TrimSpace(masuk.Isi)
	masuk.TautanSumber = strings.TrimSpace(masuk.TautanSumber)
	if err := masuk.Validasi(u.Terlarang); err != nil {
		return 0, err
	}
	if err := u.pastikanBolehKirim(ctx, masuk); err != nil {
		return 0, err
	}
	id, err := u.Repo.SimpanTemuan(ctx, masuk)
	if err != nil {
		return 0, fmt.Errorf("simpan temuan paket %s: %w", masuk.PaketKode, err)
	}
	return id, nil
}

func (u KelolaTemuan) pastikanBolehKirim(ctx context.Context, masuk domain.TemuanMasuk) error {
	ada, err := u.Repo.IsPaketAda(ctx, masuk.PaketKode)
	if err != nil {
		return fmt.Errorf("cek paket %s: %w", masuk.PaketKode, err)
	}
	if !ada {
		return domain.TidakDitemukan("paket tidak ditemukan")
	}
	jumlah, err := u.Repo.HitungTemuanSejak(ctx, masuk.PenggunaID, u.Sekarang().Add(-domain.JendelaBatasTemuan))
	if err != nil {
		return fmt.Errorf("hitung temuan pengguna %d: %w", masuk.PenggunaID, err)
	}
	if jumlah >= domain.BatasTemuanPerHari {
		return domain.ErrBatasTemuanHarian
	}
	return nil
}

func (u KelolaTemuan) Tanggapi(ctx context.Context, masuk domain.TanggapanMasuk) error {
	masuk.TautanSumber = strings.TrimSpace(masuk.TautanSumber)
	if err := masuk.Validasi(); err != nil {
		return err
	}
	pengirim, ada, err := u.Repo.AmbilPengirimTemuan(ctx, masuk.TemuanID)
	if err != nil {
		return fmt.Errorf("ambil temuan %d: %w", masuk.TemuanID, err)
	}
	switch {
	case !ada:
		return domain.ErrTemuanTidakDitemukan
	case pengirim == masuk.PenggunaID:
		return domain.ErrTanggapiTemuanSendiri
	}
	if err := u.Repo.SimpanTanggapan(ctx, masuk); err != nil {
		return fmt.Errorf("simpan tanggapan temuan %d: %w", masuk.TemuanID, err)
	}
	return nil
}

func (u KelolaTemuan) Sembunyikan(ctx context.Context, p PermintaanModerasiTemuan) error {
	if strings.TrimSpace(p.Alasan) == "" {
		return domain.ErrAlasanModerasiKosong
	}
	ada, err := u.Repo.SembunyikanTemuan(ctx, p.TemuanID)
	if err != nil {
		return fmt.Errorf("sembunyikan temuan %d: %w", p.TemuanID, err)
	}
	if !ada {
		return domain.ErrTemuanTidakDitemukan
	}
	return u.Repo.CatatModerasi(ctx, domain.CatatanModerasi{
		ModeratorID: p.ModeratorID, TargetJenis: domain.TargetTemuan, TargetID: p.TemuanID,
		Aksi: domain.AksiSembunyikanTemuan, Alasan: p.Alasan,
	})
}
