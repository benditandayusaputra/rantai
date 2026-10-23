package wilayah

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const userAgentUji = "RantaiBot/1.0 (+uji)"

var berkasDKI = []string{
	"wilayah_boundaries_kab_31.sql",
	"wilayah_boundaries_kec_31.sql",
	"wilayah_boundaries_kel_31.01.sql",
	"wilayah_boundaries_kel_31.71.sql",
	"wilayah_boundaries_kel_31.72.sql",
	"wilayah_boundaries_kel_31.73.sql",
	"wilayah_boundaries_kel_31.74.sql",
	"wilayah_boundaries_kel_31.75.sql",
}

func sumberUji(t *testing.T) Sumber {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() != userAgentUji {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		http.ServeFile(w, r, filepath.Join("testdata", filepath.Base(r.URL.Path)))
	}))
	t.Cleanup(server.Close)
	return Sumber{
		Konfigurasi: KonfigurasiSumber{URLDasar: server.URL, Berkas: berkasDKI, UserAgent: userAgentUji},
		Klien:       server.Client(),
		DirRaw:      t.TempDir(),
		Sekarang:    func() time.Time { return time.Date(2026, 10, 23, 0, 0, 0, 0, time.UTC) },
	}
}

func hitungPerTingkat(daftar []domain.Wilayah) map[domain.Tingkat]int {
	jumlah := map[domain.Tingkat]int{}
	for _, w := range daftar {
		jumlah[w.Tingkat]++
	}
	return jumlah
}

func TestAmbil_SeluruhBerkasDKI_JumlahWilayahSesuai(t *testing.T) {
	data, err := sumberUji(t).Ambil(t.Context())

	if err != nil {
		t.Fatal(err)
	}
	jumlah := hitungPerTingkat(data.Wilayah)
	harap := map[domain.Tingkat]int{domain.TingkatKota: 6, domain.TingkatKecamatan: 44, domain.TingkatKelurahan: 267}
	for tingkat, n := range harap {
		if jumlah[tingkat] != n {
			t.Errorf("jumlah %s = %d, harap %d", tingkat, jumlah[tingkat], n)
		}
	}
	if data.Provenance.Label != domain.LabelResmi || data.Provenance.HashKonten == "" {
		t.Errorf("provenance tidak lengkap: %+v", data.Provenance)
	}
}

func TestAmbil_BerkasDiunduh_DataMentahTersimpan(t *testing.T) {
	sumber := sumberUji(t)

	if _, err := sumber.Ambil(t.Context()); err != nil {
		t.Fatal(err)
	}

	berkas := filepath.Join(sumber.DirRaw, "wilayah", "2026-10-23", berkasDKI[0])
	if _, err := os.Stat(berkas); err != nil {
		t.Errorf("data mentah tidak tersimpan: %v", err)
	}
}

func TestUraiSQL_PathLintangBujur_DitukarKeBujurLintang(t *testing.T) {
	isi := []byte(`('31.72.06.1003','Kelapa Gading Barat',-6.15,106.9,'[[[-6.1,106.8],[-6.2,106.9],[-6.1,106.9],[-6.1,106.8]]]')`)

	daftar, err := UraiSQL(isi)

	if err != nil || len(daftar) != 1 {
		t.Fatalf("daftar = %v, err = %v", daftar, err)
	}
	titikPertama := daftar[0].Geometri[0][0][0]
	if titikPertama != (domain.Koordinat{106.8, -6.1}) || daftar[0].IndukKode != "31.72.06" {
		t.Errorf("wilayah = %+v", daftar[0])
	}
}
