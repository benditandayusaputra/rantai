package dapodik

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
)

func serverFixture(t *testing.T, jumlahPanggilan *atomic.Int32) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jumlahPanggilan.Add(1)
		_, _ = w.Write(bacaFixture(t, fixtureUntuk(t, r)))
	}))
	t.Cleanup(server.Close)
	return server
}

func fixtureUntuk(t *testing.T, r *http.Request) string {
	t.Helper()
	switch {
	case r.URL.Path == pathCari:
		badan := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(badan)
		if strings.Contains(string(badan), `"page":0`) {
			return "cari_sekolah.json"
		}
		return "cari_sekolah_kosong.json"
	case strings.HasPrefix(r.URL.Path, pathDetail):
		return "full_detail.json"
	case strings.HasPrefix(r.URL.Path, pathRombel):
		return "rombongan_belajar.json"
	default:
		return "peserta_didik.json"
	}
}

func sumberUji(url, dirRaw string) Sumber {
	return Sumber{
		Konfigurasi: Konfigurasi{
			URLAPI: url, UkuranHalaman: 10, Kota: []string{"Kab. Adm. Kep. Seribu"},
			Bentuk: []domain.BentukSekolah{domain.BentukSD}, BentukDenganRombel: []domain.BentukSekolah{domain.BentukSD},
		},
		Pengunduh: klien.Baru(http.DefaultClient, klien.Konfigurasi{UserAgent: "RantaiBot/1.0 (+uji)"}),
		DirRaw:    dirRaw,
		Sekarang:  func() time.Time { return time.Date(2026, 10, 23, 2, 0, 0, 0, time.UTC) },
	}
}

func TestAmbil_SatuKotaSD_SekolahDenganRombelPerTingkat(t *testing.T) {
	var panggilan atomic.Int32
	server := serverFixture(t, &panggilan)

	data, err := sumberUji(server.URL, t.TempDir()).Ambil(t.Context())

	if err != nil || len(data.Sekolah) != 4 {
		t.Fatalf("sekolah %d, err %v", len(data.Sekolah), err)
	}
	s := data.Sekolah[0]
	if s.TahunAjaran != "2026/2027" || len(s.Rombel) != 6 || s.Rombel[5] != (domain.RombelTingkat{Tingkat: 6, JumlahRombel: 2, JumlahSiswa: 63}) {
		t.Errorf("tahun ajaran %q, rombel %+v", s.TahunAjaran, s.Rombel)
	}
	if data.Provenance.Sumber != domain.SumberDapodik || data.Provenance.Label != domain.LabelResmi || data.Provenance.HashKonten == "" {
		t.Errorf("provenance = %+v", data.Provenance)
	}
}

func TestAmbil_DijalankanUlangHariSama_MemakaiCacheTanpaPermintaan(t *testing.T) {
	var panggilan atomic.Int32
	server := serverFixture(t, &panggilan)
	sumber := sumberUji(server.URL, t.TempDir())
	if _, err := sumber.Ambil(t.Context()); err != nil {
		t.Fatal(err)
	}
	sebelum := panggilan.Load()

	if _, err := sumber.Ambil(t.Context()); err != nil || panggilan.Load() != sebelum {
		t.Errorf("permintaan baru %d, err %v", panggilan.Load()-sebelum, err)
	}
}

func TestAmbil_Status403_BerhentiDenganErrDitolak(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	_, err := sumberUji(server.URL, t.TempDir()).Ambil(t.Context())

	if !errors.Is(err, klien.ErrDitolakServer) {
		t.Errorf("harap ErrDitolakServer, dapat %v", err)
	}
}
