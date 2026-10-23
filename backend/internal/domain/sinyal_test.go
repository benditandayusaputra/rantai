package domain

import (
	"fmt"
	"testing"
	"time"
)

var hariIniSinyal = time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC)

func aturanSinyalUji() AturanSinyal {
	a := AturanSinyal{MinPembanding: 8}
	a.S01.Metode = []string{"Tender", "Seleksi"}
	a.S02.PengaliIQR = 1.5
	a.S03.Metode = "Pengadaan Langsung"
	a.S03.JenisPengadaan = "Pekerjaan Konstruksi"
	a.S03.JumlahMin = 3
	a.S03.JendelaHari = 30
	a.S03.PorsiAmbang = 0.85
	a.S03.Ambang = []AmbangPengadaanLangsung{
		{BerlakuMulai: time.Date(2021, time.February, 2, 0, 0, 0, 0, time.UTC), Nilai: 200e6, Dasar: "Perpres 12/2021"},
		{BerlakuMulai: time.Date(2025, time.April, 30, 0, 0, 0, 0, time.UTC), Nilai: 400e6, Dasar: "Perpres 46/2025"},
	}
	a.S04.PorsiNilaiMin = 0.4
	a.S04.JumlahPaketMin = 5
	a.S05.TenggangHari = 30
	return a
}

func angka[T any](n T) *T { return &n }

func jenisPerPaket(daftar []Sinyal, jenis JenisSinyal) map[string]bool {
	hasil := map[string]bool{}
	for _, s := range daftar {
		if s.Jenis == jenis {
			hasil[s.KodePaket] = true
		}
	}
	return hasil
}

func TestS01_TenderSatuPeserta_PositifDanPengadaanLangsungNegatif(t *testing.T) {
	daftar := []PaketSinyal{
		{Kode: "tender", Metode: "Tender", JumlahPeserta: angka(1)},
		{Kode: "pl", Metode: "Pengadaan Langsung", JumlahPeserta: angka(1)},
		{Kode: "tender-ramai", Metode: "Tender", JumlahPeserta: angka(4)},
	}

	hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalSatuPeserta)

	if !hasil["tender"] || hasil["pl"] || hasil["tender-ramai"] {
		t.Errorf("S-01 = %v", hasil)
	}
}

func paketSejenis(n int, hps float64) []PaketSinyal {
	daftar := []PaketSinyal{}
	for i := range n {
		daftar = append(daftar, PaketSinyal{Kode: fmt.Sprint(i), Kategori: "saluran", JenisPengadaan: "Pekerjaan Konstruksi", TahunAnggaran: 2026, HPS: angka(hps + float64(i)*1e6)})
	}
	return daftar
}

func TestS02_HPSJauhDiAtasSejenis_PositifDanDiBawahMinPembandingNegatif(t *testing.T) {
	cukup := append(paketSejenis(8, 100e6), PaketSinyal{Kode: "tinggi", Kategori: "saluran", JenisPengadaan: "Pekerjaan Konstruksi", TahunAnggaran: 2026, HPS: angka(900e6)})
	kurang := append(paketSejenis(6, 100e6), PaketSinyal{Kode: "tinggi", Kategori: "saluran", JenisPengadaan: "Pekerjaan Konstruksi", TahunAnggaran: 2026, HPS: angka(900e6)})

	positif := jenisPerPaket(aturanSinyalUji().Hitung(cukup, hariIniSinyal), SinyalNilaiDiAtasSejenis)
	negatif := jenisPerPaket(aturanSinyalUji().Hitung(kurang, hariIniSinyal), SinyalNilaiDiAtasSejenis)

	if !positif["tinggi"] || len(positif) != 1 || len(negatif) != 0 {
		t.Errorf("positif %v, negatif %v", positif, negatif)
	}
}

func paketPL(kode string, hari int, hps float64) PaketSinyal {
	dibuat := time.Date(2026, time.March, hari, 0, 0, 0, 0, time.UTC)
	return PaketSinyal{Kode: kode, Metode: "Pengadaan Langsung", JenisPengadaan: "Pekerjaan Konstruksi", Kategori: "saluran", InstansiID: 1, HPS: angka(hps), TanggalDibuat: &dibuat}
}

func TestS03_TigaPaketMepetBatasDalam30Hari_Positif(t *testing.T) {
	daftar := []PaketSinyal{paketPL("a", 1, 390e6), paketPL("b", 10, 395e6), paketPL("c", 25, 399e6)}

	hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalMepetBatasPL)

	if len(hasil) != 3 {
		t.Errorf("S-03 = %v", hasil)
	}
}

func TestS03_JaraknyaLebihDari30HariAtauNilaiRendah_Negatif(t *testing.T) {
	berjauhan := []PaketSinyal{paketPL("a", 1, 390e6), paketPL("b", 10, 395e6), paketPL("c", 1, 399e6)}
	berjauhan[2].TanggalDibuat = angka(time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC))
	rendah := []PaketSinyal{paketPL("a", 1, 390e6), paketPL("b", 10, 395e6), paketPL("c", 25, 200e6)}

	for nama, daftar := range map[string][]PaketSinyal{"berjauhan": berjauhan, "rendah": rendah} {
		if hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalMepetBatasPL); len(hasil) != 0 {
			t.Errorf("%s: S-03 = %v", nama, hasil)
		}
	}
}

func TestS03_SebelumPerpres46_AmbangLama200Juta(t *testing.T) {
	daftar := []PaketSinyal{paketPL("a", 1, 190e6), paketPL("b", 2, 195e6), paketPL("c", 3, 199e6)}
	for i := range daftar {
		daftar[i].TanggalDibuat = angka(daftar[i].TanggalDibuat.AddDate(-1, 0, 0))
	}

	if hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalMepetBatasPL); len(hasil) != 3 {
		t.Errorf("S-03 ambang lama = %v", hasil)
	}
}

func paketMenang(kode string, penyedia int64, nilai float64) PaketSinyal {
	return PaketSinyal{Kode: kode, Kategori: "jalan", InstansiID: 7, TahunAnggaran: 2026, PenyediaID: &penyedia, NilaiKontrak: angka(nilai)}
}

func TestS04_PenyediaMenangLebihDari40Persen_Positif(t *testing.T) {
	daftar := []PaketSinyal{paketMenang("a", 1, 50), paketMenang("b", 1, 50), paketMenang("c", 2, 40), paketMenang("d", 3, 30), paketMenang("e", 4, 30)}

	hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalPenyediaDominan)

	if !hasil["a"] || !hasil["b"] || len(hasil) != 2 {
		t.Errorf("S-04 = %v", hasil)
	}
}

func TestS04_KurangDariLimaPaketAtauTersebar_Negatif(t *testing.T) {
	sedikit := []PaketSinyal{paketMenang("a", 1, 50), paketMenang("b", 1, 50), paketMenang("c", 2, 10), paketMenang("d", 3, 10)}
	tersebar := []PaketSinyal{paketMenang("a", 1, 30), paketMenang("b", 2, 30), paketMenang("c", 3, 30), paketMenang("d", 4, 30), paketMenang("e", 5, 30)}

	for nama, daftar := range map[string][]PaketSinyal{"sedikit": sedikit, "tersebar": tersebar} {
		if hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalPenyediaDominan); len(hasil) != 0 {
			t.Errorf("%s: S-04 = %v", nama, hasil)
		}
	}
}

func TestS05_LewatJadwalDenganLaporanTerverifikasi_PositifTanpaLaporanNegatif(t *testing.T) {
	selesai := time.Date(2026, time.August, 1, 0, 0, 0, 0, time.UTC)
	daftar := []PaketSinyal{
		{Kode: "dilaporkan", JadwalSelesai: &selesai, LaporanTerverifikasiTertahan: 1},
		{Kode: "tanpa-laporan", JadwalSelesai: &selesai},
		{Kode: "masih-tenggang", JadwalSelesai: angka(time.Date(2026, time.October, 1, 0, 0, 0, 0, time.UTC)), LaporanTerverifikasiTertahan: 2},
	}

	hasil := jenisPerPaket(aturanSinyalUji().Hitung(daftar, hariIniSinyal), SinyalMelewatiJadwal)

	if !hasil["dilaporkan"] || len(hasil) != 1 {
		t.Errorf("S-05 = %v", hasil)
	}
}

func TestHitungStatistik_KuartilInterpolasiLinear(t *testing.T) {
	s := HitungStatistik([]float64{4, 1, 3, 2, 5})

	if s.Q1 != 2 || s.Median != 3 || s.Q3 != 4 || s.Jumlah != 5 || s.BatasAtas(1.5) != 7 {
		t.Errorf("statistik = %+v", s)
	}
}

func TestFormatRupiah(t *testing.T) {
	if hasil := FormatRupiah(185101703.6); hasil != "Rp185.101.704" {
		t.Errorf("format = %s", hasil)
	}
}

func TestBandingkanPaket_KurangDariMinPembanding_TidakTersedia(t *testing.T) {
	p := PaketDibandingkan{Kode: "x", Kategori: "saluran", HPS: angka(100.0)}
	tujuh := make([]NilaiPembanding, 7)
	delapan := make([]NilaiPembanding, 8)

	if BandingkanPaket(p, tujuh, "prv", 8).IsTersedia || !BandingkanPaket(p, delapan, "prv", 8).IsTersedia {
		t.Error("batas min_pembanding tidak dipatuhi")
	}
	if BandingkanPaket(PaketDibandingkan{Kode: "y", HPS: angka(1.0)}, delapan, "prv", 8).IsTersedia {
		t.Error("paket tanpa kategori tetap dibandingkan")
	}
}
