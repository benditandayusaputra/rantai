package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const alasanTerverifikasiSaksi = "dikonfirmasi %d akun dalam 200 m dan 14 hari"

type RepositoriVerifikasi interface {
	AmbilUntukVerifikasi(ctx context.Context, id int64) (domain.LaporanUntukVerifikasi, bool, error)
	SimpanKonfirmasi(ctx context.Context, laporanID int64, k domain.Konfirmasi) error
	UbahStatus(ctx context.Context, perubahan domain.PerubahanStatus) error
}

type HasilKonfirmasi struct {
	Status      domain.StatusLaporan
	JumlahSaksi int
}

type KonfirmasiLaporan struct {
	Repo     RepositoriVerifikasi
	Sekarang func() time.Time
}

func (u KonfirmasiLaporan) Konfirmasi(ctx context.Context, masuk domain.KonfirmasiMasuk) (HasilKonfirmasi, error) {
	if err := masuk.Validasi(); err != nil {
		return HasilKonfirmasi{}, err
	}
	laporan, err := ambilUntukVerifikasi(ctx, u.Repo, masuk.LaporanID)
	if err != nil {
		return HasilKonfirmasi{}, err
	}
	if err := laporan.PeriksaBisaDikonfirmasi(masuk.PenggunaID); err != nil {
		return HasilKonfirmasi{}, err
	}
	k := domain.Konfirmasi{PenggunaID: masuk.PenggunaID, Kondisi: masuk.Kondisi, Titik: masuk.Titik, DibuatPada: u.Sekarang()}
	if err := u.Repo.SimpanKonfirmasi(ctx, laporan.ID, k); err != nil {
		return HasilKonfirmasi{}, err
	}
	laporan.Konfirmasi = append(laporan.Konfirmasi, k)
	return u.verifikasiBilaCukup(ctx, laporan)
}

func (u KonfirmasiLaporan) verifikasiBilaCukup(ctx context.Context, laporan domain.LaporanUntukVerifikasi) (HasilKonfirmasi, error) {
	hasil := HasilKonfirmasi{Status: laporan.Status, JumlahSaksi: laporan.JumlahSaksiSah()}
	if !laporan.IsCukupSaksi() {
		return hasil, nil
	}
	if err := domain.Transisi(laporan.Status, domain.StatusTerverifikasi); err != nil {
		return hasil, err
	}
	err := u.Repo.UbahStatus(ctx, domain.PerubahanStatus{
		LaporanID: laporan.ID, Dari: laporan.Status, Ke: domain.StatusTerverifikasi,
		Alasan: fmt.Sprintf(alasanTerverifikasiSaksi, hasil.JumlahSaksi),
	})
	if err != nil {
		return hasil, fmt.Errorf("verifikasi laporan %d: %w", laporan.ID, err)
	}
	hasil.Status = domain.StatusTerverifikasi
	return hasil, nil
}

type pengambilVerifikasi interface {
	AmbilUntukVerifikasi(ctx context.Context, id int64) (domain.LaporanUntukVerifikasi, bool, error)
}

func ambilUntukVerifikasi(ctx context.Context, repo pengambilVerifikasi, id int64) (domain.LaporanUntukVerifikasi, error) {
	laporan, ada, err := repo.AmbilUntukVerifikasi(ctx, id)
	if err != nil {
		return laporan, fmt.Errorf("ambil laporan %d: %w", id, err)
	}
	if !ada {
		return laporan, domain.TidakDitemukan("laporan tidak ditemukan")
	}
	return laporan, nil
}
