package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	prefiksKunciFoto      = "laporan/"
	akhiranKunciThumbnail = "-480.jpg"
	akhiranKunciUtama     = ".jpg"
)

type PemrosesFoto interface {
	Proses(isi []byte) (domain.FotoDiproses, error)
}

type PenyimpanFoto interface {
	SimpanJPEG(ctx context.Context, kunci string, isi []byte) error
}

type RepositoriLaporan interface {
	JumlahLaporanSejak(ctx context.Context, penggunaID int64, sejak time.Time) (int, error)
	JarakPaketMeter(ctx context.Context, kode string, titik domain.Titik) (*float64, bool, error)
	WilayahDiTitik(ctx context.Context, titik domain.Titik) (string, bool, error)
	SimpanLaporan(ctx context.Context, laporan domain.LaporanBaru) (int64, error)
}

type LaporanDiterima struct {
	ID     int64
	Status domain.StatusLaporan
}

type TerimaLaporan struct {
	Pemroses  PemrosesFoto
	Penyimpan PenyimpanFoto
	Repo      RepositoriLaporan
	Sekarang  func() time.Time
	BuatKunci func() string
}

func (u TerimaLaporan) Terima(ctx context.Context, masuk domain.LaporanMasuk) (LaporanDiterima, error) {
	if err := masuk.Validasi(); err != nil {
		return LaporanDiterima{}, err
	}
	if err := u.periksaBatasHarian(ctx, masuk.PenggunaID); err != nil {
		return LaporanDiterima{}, err
	}
	baru, err := u.susunLaporan(ctx, masuk)
	if err != nil {
		return LaporanDiterima{}, err
	}
	if baru.Foto, err = u.simpanFoto(ctx, masuk.Foto); err != nil {
		return LaporanDiterima{}, err
	}
	id, err := u.Repo.SimpanLaporan(ctx, baru)
	if err != nil {
		return LaporanDiterima{}, fmt.Errorf("simpan laporan: %w", err)
	}
	return LaporanDiterima{ID: id, Status: domain.StatusMenungguSaksi}, nil
}

func (u TerimaLaporan) periksaBatasHarian(ctx context.Context, penggunaID int64) error {
	jumlah, err := u.Repo.JumlahLaporanSejak(ctx, penggunaID, u.Sekarang().Add(-domain.JendelaBatasLaporan))
	if err != nil {
		return fmt.Errorf("hitung laporan pengguna %d: %w", penggunaID, err)
	}
	if jumlah >= domain.BatasLaporanPerHari {
		return domain.ErrBatasLaporanHarian
	}
	return nil
}

func (u TerimaLaporan) susunLaporan(ctx context.Context, masuk domain.LaporanMasuk) (domain.LaporanBaru, error) {
	if err := u.periksaPaket(ctx, masuk); err != nil {
		return domain.LaporanBaru{}, err
	}
	kodeWilayah, ada, err := u.Repo.WilayahDiTitik(ctx, masuk.Titik)
	if err != nil {
		return domain.LaporanBaru{}, fmt.Errorf("cari wilayah titik laporan: %w", err)
	}
	if !ada {
		return domain.LaporanBaru{}, domain.ErrLokasiDiLuarDKI
	}
	return domain.LaporanBaru{LaporanMasuk: masuk, KodeWilayah: kodeWilayah, TitikPublik: domain.BulatkanLokasiPublik(masuk.Titik)}, nil
}

func (u TerimaLaporan) periksaPaket(ctx context.Context, masuk domain.LaporanMasuk) error {
	if masuk.PaketKode == "" {
		return nil
	}
	jarak, ada, err := u.Repo.JarakPaketMeter(ctx, masuk.PaketKode, masuk.Titik)
	switch {
	case err != nil:
		return fmt.Errorf("periksa paket %s: %w", masuk.PaketKode, err)
	case !ada:
		return domain.TidakDitemukan("paket tidak ditemukan")
	case jarak != nil && *jarak > domain.RadiusProyekLaporanMeter:
		return domain.ErrPaketLaporanJauh
	default:
		return nil
	}
}

func (u TerimaLaporan) simpanFoto(ctx context.Context, semua [][]byte) ([]domain.FotoTersimpan, error) {
	diproses := make([]domain.FotoDiproses, len(semua))
	for i, isi := range semua {
		hasil, err := u.Pemroses.Proses(isi)
		if err != nil {
			return nil, fmt.Errorf("foto %d: %w", i+1, err)
		}
		diproses[i] = hasil
	}
	tersimpan := make([]domain.FotoTersimpan, len(diproses))
	for i, foto := range diproses {
		tersimpan[i] = kunciFoto(prefiksKunciFoto+u.BuatKunci(), foto)
		if err := unggahFoto(ctx, u.Penyimpan, tersimpan[i], foto); err != nil {
			return nil, err
		}
	}
	return tersimpan, nil
}

func kunciFoto(dasar string, foto domain.FotoDiproses) domain.FotoTersimpan {
	return domain.FotoTersimpan{
		ObjectKey: dasar + akhiranKunciUtama, ThumbnailKey: dasar + akhiranKunciThumbnail, Lebar: foto.Lebar, Tinggi: foto.Tinggi,
	}
}

func unggahFoto(ctx context.Context, penyimpan PenyimpanFoto, kunci domain.FotoTersimpan, foto domain.FotoDiproses) error {
	if err := penyimpan.SimpanJPEG(ctx, kunci.ObjectKey, foto.Utama); err != nil {
		return err
	}
	return penyimpan.SimpanJPEG(ctx, kunci.ThumbnailKey, foto.Thumbnail)
}
