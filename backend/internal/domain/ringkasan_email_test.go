package domain

import (
	"strings"
	"testing"
)

func TestSusunEmailRingkasan_DuaKecamatan_DikelompokkanDenganTautanDanLabelSimulasi(t *testing.T) {
	kabar := []Notifikasi{
		{Jenis: KejadianPaketBaru, KodeKecamatan: "31.72.06", NamaKecamatan: "Kelapa Gading", Judul: "Rehab Ruang Kelas", Tautan: "/paket/a"},
		{Jenis: KejadianLaporanTerverifikasi, KodeKecamatan: "31.72.06", NamaKecamatan: "Kelapa Gading", Judul: "Laporan atap", Tautan: "/laporan/1", IsSimulation: true},
		{Jenis: KejadianHakJawab, KodeKecamatan: "31.72.01", NamaKecamatan: "Penjaringan", Judul: "Tanggapan Sudin", Tautan: "/paket/b"},
	}

	email := SusunEmailRingkasan(kabar, "https://rantai.uji")

	if !strings.Contains(email.Subjek, "3 kabar") {
		t.Errorf("subjek = %q", email.Subjek)
	}
	harap := []string{
		"\nKelapa Gading\n- [Paket baru] Rehab Ruang Kelas\n  https://rantai.uji/id/paket/a\n",
		"- [Laporan terverifikasi] Laporan atap (Data simulasi)\n",
		"\nPenjaringan\n- [Hak jawab baru] Tanggapan Sudin\n",
		"https://rantai.uji/id/pantau",
	}
	for _, potongan := range harap {
		if !strings.Contains(email.Isi, potongan) {
			t.Errorf("isi tidak memuat %q:\n%s", potongan, email.Isi)
		}
	}
	if strings.Count(email.Isi, "Kelapa Gading\n") != 1 {
		t.Errorf("kecamatan tertulis lebih dari sekali:\n%s", email.Isi)
	}
}
