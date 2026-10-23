package dayatampung

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func TestUraiCSV_LampiranNegeri_BarisDanNPSNKosong(t *testing.T) {
	isi, err := os.ReadFile("testdata/daya_tampung_negeri.csv")
	if err != nil {
		t.Fatal(err)
	}

	baris, err := UraiCSV(isi)

	if err != nil || len(baris) != 3 || domain.TotalKursi(baris) != 396 {
		t.Fatalf("baris %+v, err %v", baris, err)
	}
	if baris[1].NPSN != "" || baris[1].Nama != "SMP NEGERI 296 JAKARTA" {
		t.Errorf("baris tanpa NPSN = %+v", baris[1])
	}
}

func TestUraiCSV_KolomNPSNHilang_Ditolak(t *testing.T) {
	isi, err := os.ReadFile("testdata/tanpa_npsn.csv")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := UraiCSV(isi); !errors.Is(err, ErrKolomTidakLengkap) {
		t.Errorf("harap ErrKolomTidakLengkap, dapat %v", err)
	}
}

func TestUraiCSV_KursiBukanAngka_Ditolak(t *testing.T) {
	if _, err := UraiCSV([]byte("nama,npsn,kursi\nSMP A,1,sepuluh\n")); err == nil {
		t.Error("kursi bukan angka harus ditolak")
	}
}

func TestAmbil_BerkasNegeri_ProvenanceDisdikResmi(t *testing.T) {
	sumber := Sumber{
		Konfigurasi: Konfigurasi{TahunAjaran: "2026/2027", Berkas: []KonfigurasiBerkas{
			{Path: "testdata/daya_tampung_negeri.csv", Jalur: domain.JalurNegeri, URL: "https://contoh.test/sk258.pdf"},
		}},
		Sekarang: time.Now,
	}

	berkas, err := sumber.Ambil(t.Context())

	if err != nil || len(berkas) != 1 || berkas[0].TahunAjaran != "2026/2027" {
		t.Fatalf("berkas %+v, err %v", berkas, err)
	}
	p := berkas[0].Provenance
	if p.Sumber != domain.SumberDisdik || p.Label != domain.LabelResmi || p.HashKonten == "" {
		t.Errorf("provenance = %+v", p)
	}
}

func TestAmbil_JalurTidakDikenal_Ditolak(t *testing.T) {
	sumber := Sumber{Konfigurasi: Konfigurasi{Berkas: []KonfigurasiBerkas{{Path: "testdata/daya_tampung_negeri.csv", Jalur: "swasta"}}}, Sekarang: time.Now}

	if _, err := sumber.Ambil(t.Context()); err == nil {
		t.Error("jalur tidak dikenal harus ditolak")
	}
}
