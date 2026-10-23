package ujiwakil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

var padananFixture = [][2]string{
	{"/fraksi/periode-2024-2029/", "fraksi_"},
	{"/akd/komisi/periode-2024-2029/", "komisi_"},
	{"/akd/dapil/", "dapil_"},
	{"/akd/", "akd_"},
}

type pengunduhFixture struct {
	akarURL string
	dir     string
}

func MuatDPRD(t *testing.T, tx pgx.Tx) []string {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	akar := akarBackend()
	var cfg dprd.Konfigurasi
	if err := config.BacaYAML(filepath.Join(akar, "config/dprd.yaml"), &cfg); err != nil {
		t.Fatal(err)
	}
	laporan := []string{}
	impor := usecase.ImporDPRD{
		Pengambil: dprd.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   pengunduhFixture{akarURL: cfg.Sumber.URLDasar, dir: filepath.Join(akar, "internal/ingest/dprd/testdata")},
			DirRaw:      t.TempDir(),
			Sekarang:    func() time.Time { return time.Date(2026, time.October, 23, 9, 0, 0, 0, time.UTC) },
		},
		PembacaKecamatan: repository.WakilRepository{DB: tx},
		Penyimpan:        repository.WakilRepository{DB: tx},
		Pencatat:         repository.IngestRunRepository{DB: tx},
		Dapil:            cfg.DapilKPU.Daftar,
		TahunPemilu:      cfg.DapilKPU.TahunPemilu,
		URLLHKPN:         cfg.URLLHKPN,
		PolaURLProfil:    cfg.Sumber.URLDasar + cfg.Sumber.HalamanProfil,
		Laporkan:         func(pesan string) { laporan = append(laporan, pesan) },
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return laporan
}

func (p pengunduhFixture) Ambil(_ context.Context, url string) ([]byte, error) {
	halaman := strings.TrimPrefix(url, p.akarURL)
	for _, padanan := range padananFixture {
		if nama, ada := strings.CutPrefix(halaman, padanan[0]); ada {
			return os.ReadFile(filepath.Join(p.dir, padanan[1]+strings.TrimSuffix(nama, "/")+".html"))
		}
	}
	return nil, fmt.Errorf("fixture untuk %s tidak ada", url)
}

func akarBackend() string {
	_, berkas, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(berkas), "../../..")
}
