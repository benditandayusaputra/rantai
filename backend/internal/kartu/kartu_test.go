package kartu

import (
	"bytes"
	"image/png"
	"os"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func konfigurasiUji(t *testing.T) Konfigurasi {
	t.Helper()
	var k Konfigurasi
	if err := config.BacaYAML("../../config/kartu.yaml", &k); err != nil {
		t.Fatal(err)
	}
	return k
}

func ringkasanUji() domain.RingkasanKecamatan {
	rasio := 1.834
	tema := domain.TemaTersimpan{Tema: domain.Tema{Kode: domain.TemaPendidikan}}
	tema.Indikator.Judul = domain.Teks{ID: "Kursi SMP negeri dibanding lulusan SD", EN: "Public junior high seats vs primary graduates"}
	var r domain.RingkasanKecamatan
	r.Kecamatan.Kode = "31.72.06"
	r.Kecamatan.Nama = "Kelapa Gading"
	r.Kecamatan.Hierarki = []domain.RingkasanWilayah{{Kode: "31.72", Nama: "Kota Administrasi Jakarta Utara", Tingkat: domain.TingkatKota}}
	r.Tema = []domain.TemaKecamatan{{Tema: tema, Indikator: &domain.IndikatorKecamatan{Periode: "2026/2027", Rasio: &rasio, Label: domain.LabelResmi}}}
	r.Anggaran = &domain.RingkasanAnggaran{Tahun: 2026, JumlahPaket: 12, TotalHPS: 4_250_000_000}
	r.Wakil.Nomor = 2
	r.Wakil.JumlahKursi = 9
	r.Wakil.KomisiPengawas = []string{"E"}
	return r
}

func TestSusunIsi_BahasaIndonesia_AngkaRingkasDanTautanPendek(t *testing.T) {
	isi := susunIsi(ringkasanUji(), konfigurasiUji(t), "id")

	if isi.Judul != "Kecamatan Kelapa Gading" || isi.Kota != "Kota Administrasi Jakarta Utara" {
		t.Errorf("judul %q, kota %q", isi.Judul, isi.Kota)
	}
	if isi.Blok[0].Nilai != "1,8" || isi.Blok[0].Label != domain.LabelResmi || isi.Blok[1].Nilai != "Rp4,3 miliar" {
		t.Errorf("blok = %+v", isi.Blok)
	}
	if isi.Blok[1].Label != domain.LabelProksi || isi.Blok[2].Nilai != "9 kursi DPRD" {
		t.Errorf("blok = %+v", isi.Blok)
	}
	if isi.Tautan != "rantai-civic.vercel.app/id/kecamatan/31.72.06" {
		t.Errorf("tautan = %q", isi.Tautan)
	}
}

func TestSusunIsi_BahasaInggris_PemisahTitikDanJudulInggris(t *testing.T) {
	isi := susunIsi(ringkasanUji(), konfigurasiUji(t), "en")

	if isi.Judul != "Kelapa Gading District" || isi.Blok[0].Nilai != "1.8" || isi.Blok[1].Nilai != "Rp4.3 billion" {
		t.Errorf("isi = %+v", isi)
	}
	if isi.Blok[0].Judul != "Public junior high seats vs primary graduates" {
		t.Errorf("judul tema = %q", isi.Blok[0].Judul)
	}
}

func TestSusunIsi_TanpaIndikatorDanPaket_KeteranganKosongTanpaLabel(t *testing.T) {
	r := ringkasanUji()
	r.Tema[0].Indikator = nil
	r.Anggaran = nil

	isi := susunIsi(r, konfigurasiUji(t), "xx")

	if isi.Blok[0].Nilai != "Belum dihitung" || isi.Blok[0].Label != "" || isi.Blok[1].Nilai != "Belum ada paket berlokasi" {
		t.Errorf("blok = %+v", isi.Blok)
	}
}

func TestGambar_RingkasanLengkap_PNGUkuranKartu(t *testing.T) {
	pembuat, err := Baru(konfigurasiUji(t))
	if err != nil {
		t.Fatal(err)
	}
	isi, err := pembuat.Gambar(ringkasanUji(), "id")
	if err != nil {
		t.Fatal(err)
	}
	gambar, err := png.Decode(bytes.NewReader(isi))
	if err != nil {
		t.Fatal(err)
	}
	if ukuran := gambar.Bounds().Size(); ukuran.X != lebarKartu || ukuran.Y != tinggiKartu {
		t.Errorf("ukuran = %v", ukuran)
	}
	if tujuan := os.Getenv("KARTU_CONTOH"); tujuan != "" {
		_ = os.WriteFile(tujuan, isi, 0o600)
	}
}
