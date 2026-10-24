package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var waktuLaporanVerifikasi = time.Date(2026, time.October, 1, 9, 0, 0, 0, time.UTC)

type repoVerifikasiTiruan struct {
	laporan   domain.LaporanUntukVerifikasi
	perubahan []domain.PerubahanStatus
}

func (r *repoVerifikasiTiruan) AmbilUntukVerifikasi(_ context.Context, id int64) (domain.LaporanUntukVerifikasi, bool, error) {
	return r.laporan, id == r.laporan.ID, nil
}

func (r *repoVerifikasiTiruan) SimpanKonfirmasi(_ context.Context, _ int64, k domain.Konfirmasi) error {
	r.laporan.Konfirmasi = append(r.laporan.Konfirmasi, k)
	return nil
}

func (r *repoVerifikasiTiruan) UbahStatus(_ context.Context, p domain.PerubahanStatus) error {
	r.perubahan = append(r.perubahan, p)
	return nil
}

func konfirmasiMasukUji(pengguna int64) domain.KonfirmasiMasuk {
	return domain.KonfirmasiMasuk{LaporanID: 9, PenggunaID: pengguna, Titik: domain.Titik{Lat: -6.12, Lng: 106.9}, AkurasiMeter: 10, Kondisi: domain.KondisiTerhenti}
}

func konfirmasiLaporanUji() (KonfirmasiLaporan, *repoVerifikasiTiruan) {
	repo := &repoVerifikasiTiruan{laporan: domain.LaporanUntukVerifikasi{
		ID: 9, PelaporID: 1, Status: domain.StatusMenungguSaksi, Kondisi: domain.KondisiTerhenti,
		Titik: domain.Titik{Lat: -6.12, Lng: 106.9}, DibuatPada: waktuLaporanVerifikasi,
	}}
	return KonfirmasiLaporan{Repo: repo, Sekarang: func() time.Time { return waktuLaporanVerifikasi.Add(time.Hour) }}, repo
}

func TestKonfirmasiLaporan_SaksiKetiga_Terverifikasi(t *testing.T) {
	u, repo := konfirmasiLaporanUji()

	kedua, errKedua := u.Konfirmasi(context.Background(), konfirmasiMasukUji(2))
	ketiga, errKetiga := u.Konfirmasi(context.Background(), konfirmasiMasukUji(3))

	if errKedua != nil || kedua.Status != domain.StatusMenungguSaksi || kedua.JumlahSaksi != 2 {
		t.Errorf("kedua %+v, err %v", kedua, errKedua)
	}
	if errKetiga != nil || ketiga.Status != domain.StatusTerverifikasi || len(repo.perubahan) != 1 || repo.perubahan[0].Ke != domain.StatusTerverifikasi {
		t.Errorf("ketiga %+v, err %v, perubahan %+v", ketiga, errKetiga, repo.perubahan)
	}
}

func TestKonfirmasiLaporan_LaporanBerpaket_DiteruskanKeInstansi(t *testing.T) {
	u, repo := konfirmasiLaporanUji()
	repo.laporan.PaketKode = "simulasi-ruang-kelas-1"

	_, _ = u.Konfirmasi(context.Background(), konfirmasiMasukUji(2))
	ketiga, err := u.Konfirmasi(context.Background(), konfirmasiMasukUji(3))

	if err != nil || ketiga.Status != domain.StatusDiteruskan || len(repo.perubahan) != 2 ||
		repo.perubahan[1].Dari != domain.StatusTerverifikasi || repo.perubahan[1].Ke != domain.StatusDiteruskan {
		t.Errorf("ketiga %+v, err %v, perubahan %+v", ketiga, err, repo.perubahan)
	}
}

func TestKonfirmasiLaporan_Penolakan(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*domain.KonfirmasiMasuk)
		harap error
	}{
		{nama: "akurasi buruk", ubah: func(k *domain.KonfirmasiMasuk) { k.AkurasiMeter = 150 }, harap: domain.ErrAkurasiLokasiBuruk},
		{nama: "pelapor sendiri", ubah: func(k *domain.KonfirmasiMasuk) { k.PenggunaID = 1 }, harap: domain.ErrKonfirmasiOlehPelapor},
		{nama: "laporan tidak ada", ubah: func(k *domain.KonfirmasiMasuk) { k.LaporanID = 404 }},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			u, repo := konfirmasiLaporanUji()
			masuk := konfirmasiMasukUji(2)
			k.ubah(&masuk)
			_, err := u.Konfirmasi(context.Background(), masuk)
			var kesalahan *domain.Kesalahan
			if err == nil || (k.harap != nil && !errors.Is(err, k.harap)) || !errors.As(err, &kesalahan) || len(repo.laporan.Konfirmasi) != 0 {
				t.Errorf("err %v, konfirmasi %v", err, repo.laporan.Konfirmasi)
			}
		})
	}
}
