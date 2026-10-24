package domain

import (
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

var polaDigit = regexp.MustCompile(`\d`)

func bacaYAMLUji[T any](t *testing.T, path string) T {
	t.Helper()
	isi, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var hasil T
	if err := yaml.Unmarshal(isi, &hasil); err != nil {
		t.Fatal(err)
	}
	return hasil
}

func temaPendidikanUji(t *testing.T) TemaTersimpan {
	return TemaTersimpan{Tema: bacaYAMLUji[Tema](t, "../../config/tema/pendidikan.yaml")}
}

func bahanAspirasiLengkap(t *testing.T) BahanAspirasi {
	tema := temaPendidikanUji(t)
	rasio := 0.75
	diambil := time.Date(2026, time.October, 23, 9, 0, 0, 0, time.UTC)
	return BahanAspirasi{
		Tema: tema, Masalah: tema.TanggungJawab[0], Kecamatan: "Kelapa Gading", Kota: "Kota Administrasi Jakarta Utara",
		TahunAnggaran: 2026,
		Indikator: &DataBersumber[IndikatorKecamatan]{
			Nilai:      IndikatorKecamatan{Periode: "2026/2027", Kebutuhan: 1240, Kapasitas: 930, Kesenjangan: 310, Rasio: &rasio, Label: LabelResmi},
			Provenance: Provenance{Sumber: SumberRantai, URL: "https://contoh.test/rumus", DiambilPada: diambil},
		},
		Anggaran: &DataBersumber[RingkasanAnggaran]{
			Nilai:      RingkasanAnggaran{Tahun: 2026, JumlahPaket: 3, TotalHPS: 4500000000},
			Provenance: Provenance{Sumber: "lpse", URL: "https://contoh.test/lpse", DiambilPada: diambil},
		},
		NamaDapil:     "Dapil 2",
		AnggotaKomisi: []Wakil{{Nama: "Anggota Uji", Fraksi: "Fraksi Uji"}},
		Jalur:         bacaYAMLUji[KonfigurasiJalurAspirasi](t, "../../config/jalur_aspirasi.yaml").Untuk(map[string]string{"kecamatan": "Kelapa Gading", "kota": "Kota Administrasi Jakarta Utara"}),
	}
}

func susunAspirasiUji(t *testing.T, bahan BahanAspirasi, bahasa string) Dokumen {
	return SusunDokumenAspirasi(bahan, bacaYAMLUji[KosakataAspirasi](t, "../../config/aspirasi.yaml"), bahasa)
}

func TestSusunDokumenAspirasi_SetiapAngkaDataBersumberDanBertanggal(t *testing.T) {
	for _, bahasa := range []string{BahasaIndonesia, BahasaInggris} {
		dokumen := susunAspirasiUji(t, bahanAspirasiLengkap(t), bahasa)
		tanggal := FormatTanggal(time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC), bahasa)

		for _, baris := range dokumen.Bagian[0].Baris {
			if polaDigit.MatchString(baris) && !strings.Contains(baris, tanggal) {
				t.Errorf("%s: angka tanpa sumber bertanggal: %q", bahasa, baris)
			}
		}
	}
}

func TestSusunDokumenAspirasi_DataLengkap_AngkaTerformat(t *testing.T) {
	teks := susunAspirasiUji(t, bahanAspirasiLengkap(t), BahasaIndonesia).Teks()

	for _, harap := range []string{"kebutuhan 1.240 siswa", "rasio 0,75", "3 paket dengan total HPS Rp4.500.000.000",
		"23 Oktober 2026", "label resmi", "Suku Dinas Pendidikan Kota Administrasi Jakarta Utara", "Anggota Uji (Fraksi Uji)",
		"Kekurangan kursi SMP negeri di Kecamatan Kelapa Gading", "tahun ajaran 2027"} {
		if !strings.Contains(teks, harap) {
			t.Errorf("teks tidak memuat %q:\n%s", harap, teks)
		}
	}
}

func TestSusunDokumenAspirasi_TanpaIndikatorDanAnggaran_TanpaAngkaData(t *testing.T) {
	bahan := bahanAspirasiLengkap(t)
	bahan.Indikator, bahan.Anggaran = nil, nil

	dokumen := susunAspirasiUji(t, bahan, BahasaIndonesia)

	baris := strings.Join(dokumen.Bagian[0].Baris, "\n")
	if strings.Contains(baris, "Sumber:") || !strings.Contains(baris, "belum tersedia") || !strings.Contains(baris, "Belum ada paket") {
		t.Errorf("baris data = %q", baris)
	}
}

func TestSusunDokumenAspirasi_TanpaAnggotaDiKomisi_KalimatPengganti(t *testing.T) {
	bahan := bahanAspirasiLengkap(t)
	bahan.AnggotaKomisi = nil

	baris := strings.Join(susunAspirasiUji(t, bahan, BahasaIndonesia).Bagian[2].Baris, "\n")

	if !strings.Contains(baris, "Tidak ada anggota Dapil 2 yang duduk di Komisi E") {
		t.Errorf("baris pengawas = %q", baris)
	}
}

func TestSusunDokumenAspirasi_SemuaPlaceholderTerisiDanKosakataNetral(t *testing.T) {
	terlarang := bacaKosakataTerlarang(t)
	for _, bahasa := range []string{BahasaIndonesia, BahasaInggris} {
		teks := susunAspirasiUji(t, bahanAspirasiLengkap(t), bahasa).Teks()

		if strings.ContainsAny(teks, "{}") {
			t.Errorf("%s: placeholder tersisa:\n%s", bahasa, teks)
		}
		if kata := terlarang.Temukan(teks); len(kata) > 0 {
			t.Errorf("%s: kosakata terlarang %v", bahasa, kata)
		}
	}
}

func TestFormatTanggal(t *testing.T) {
	tanggal := time.Date(2026, time.March, 2, 0, 0, 0, 0, time.UTC)
	if id, en := FormatTanggal(tanggal, BahasaIndonesia), FormatTanggal(tanggal, BahasaInggris); id != "2 Maret 2026" || en != "2 March 2026" {
		t.Errorf("id %q, en %q", id, en)
	}
}
