package domain

import (
	"errors"
	"testing"
	"time"
)

func TestTransisi_SemuaPerpindahanSah_Diterima(t *testing.T) {
	sah := [][2]StatusLaporan{
		{StatusDikirim, StatusMenungguSaksi}, {StatusDikirim, StatusDitolak},
		{StatusMenungguSaksi, StatusTerverifikasi}, {StatusMenungguSaksi, StatusDitolak},
		{StatusTerverifikasi, StatusDiteruskan}, {StatusDiteruskan, StatusDitanggapi}, {StatusDitanggapi, StatusSelesai},
	}
	for _, s := range sah {
		if err := Transisi(s[0], s[1]); err != nil {
			t.Errorf("%s ke %s: %v", s[0], s[1], err)
		}
	}
}

func TestTransisi_PerpindahanTidakSah_ErrTransisiTidakSah(t *testing.T) {
	tidakSah := [][2]StatusLaporan{
		{StatusDikirim, StatusTerverifikasi}, {StatusMenungguSaksi, StatusSelesai}, {StatusTerverifikasi, StatusDitolak},
		{StatusDitolak, StatusMenungguSaksi}, {StatusSelesai, StatusDitanggapi}, {StatusTerverifikasi, StatusMenungguSaksi},
		{StatusDiteruskan, StatusSelesai}, {StatusMenungguSaksi, StatusMenungguSaksi},
	}
	for _, s := range tidakSah {
		if err := Transisi(s[0], s[1]); !errors.Is(err, ErrTransisiTidakSah) {
			t.Errorf("%s ke %s: err = %v", s[0], s[1], err)
		}
	}
}

var waktuLaporanUji = time.Date(2026, time.October, 1, 9, 0, 0, 0, time.UTC)

func laporanVerifikasiUji(konfirmasi ...Konfirmasi) LaporanUntukVerifikasi {
	return LaporanUntukVerifikasi{
		PelaporID: 1, Status: StatusMenungguSaksi, Kondisi: KondisiTerhenti,
		Titik: Titik{Lat: -6.12, Lng: 106.90}, DibuatPada: waktuLaporanUji, Konfirmasi: konfirmasi,
	}
}

func konfirmasiUji(pengguna int64, geserUtaraMeter float64, setelah time.Duration) Konfirmasi {
	return Konfirmasi{
		PenggunaID: pengguna, Kondisi: KondisiTerhenti,
		Titik: Titik{Lat: -6.12 + geserUtaraMeter/meterPerDerajatLintang, Lng: 106.90}, DibuatPada: waktuLaporanUji.Add(setelah),
	}
}

func TestJumlahSaksiSah(t *testing.T) {
	sehari := 24 * time.Hour
	kasus := []struct {
		nama       string
		konfirmasi []Konfirmasi
		harap      int
	}{
		{nama: "tanpa konfirmasi, pelapor saja", harap: 1},
		{nama: "dua saksi dekat dan baru", konfirmasi: []Konfirmasi{konfirmasiUji(2, 50, sehari), konfirmasiUji(3, 190, 13*sehari)}, harap: 3},
		{nama: "saksi di luar 200 m tidak dihitung", konfirmasi: []Konfirmasi{konfirmasiUji(2, 50, sehari), konfirmasiUji(3, 210, sehari)}, harap: 2},
		{nama: "saksi lewat 14 hari tidak dihitung", konfirmasi: []Konfirmasi{konfirmasiUji(2, 50, sehari), konfirmasiUji(3, 10, 14*sehari+time.Minute)}, harap: 2},
		{nama: "kondisi berbeda tidak dihitung", konfirmasi: []Konfirmasi{{PenggunaID: 2, Kondisi: KondisiSelesaiBaik, Titik: Titik{Lat: -6.12, Lng: 106.90}, DibuatPada: waktuLaporanUji}}, harap: 1},
		{nama: "konfirmasi pelapor tidak menambah hitungan", konfirmasi: []Konfirmasi{konfirmasiUji(1, 0, 0), konfirmasiUji(2, 0, 0)}, harap: 2},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if jumlah := laporanVerifikasiUji(k.konfirmasi...).JumlahSaksiSah(); jumlah != k.harap {
				t.Errorf("jumlah = %d, harap %d", jumlah, k.harap)
			}
		})
	}
}

func TestPeriksaBisaDikonfirmasi(t *testing.T) {
	kasus := []struct {
		nama     string
		laporan  LaporanUntukVerifikasi
		pengguna int64
		harap    error
	}{
		{nama: "saksi baru", laporan: laporanVerifikasiUji(), pengguna: 2},
		{nama: "pelapor sendiri", laporan: laporanVerifikasiUji(), pengguna: 1, harap: ErrKonfirmasiOlehPelapor},
		{nama: "sudah mengonfirmasi", laporan: laporanVerifikasiUji(konfirmasiUji(2, 0, 0)), pengguna: 2, harap: ErrKonfirmasiSudahTercatat},
		{nama: "laporan sudah terverifikasi", laporan: LaporanUntukVerifikasi{Status: StatusTerverifikasi}, pengguna: 2, harap: ErrLaporanTidakMenunggu},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.laporan.PeriksaBisaDikonfirmasi(k.pengguna); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
