package dprd

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func waktuUji() time.Time {
	return time.Date(2026, time.October, 23, 9, 0, 0, 0, time.UTC)
}

type pengunduhFixture struct {
	akarURL string
}

func (p pengunduhFixture) Ambil(_ context.Context, url string) ([]byte, error) {
	return os.ReadFile("testdata/" + namaFixture(url[len(p.akarURL):]))
}

func namaFixture(halaman string) string {
	var nomor int
	if _, err := fmt.Sscanf(halaman, "/akd/dapil/dki-jakarta-%d/", &nomor); err == nil {
		return fmt.Sprintf("dapil_dki-jakarta-%d.html", nomor)
	}
	var fraksi string
	if _, err := fmt.Sscanf(halaman, "/fraksi/periode-2024-2029/%s", &fraksi); err == nil {
		return "fraksi_" + fraksi[:len(fraksi)-1] + ".html"
	}
	return halaman
}

func konfigurasiUji(t *testing.T) Konfigurasi {
	t.Helper()
	var cfg Konfigurasi
	if err := config.BacaYAML("../../../config/dprd.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func jalankanImporDPRDUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	cfg := konfigurasiUji(t)
	impor := usecase.ImporDPRD{
		Pengambil: Sumber{
			Konfigurasi: cfg,
			Pengunduh:   pengunduhFixture{akarURL: cfg.Sumber.URLDasar},
			DirRaw:      t.TempDir(),
			Sekarang:    waktuUji,
		},
		PembacaKecamatan: repository.WakilRepository{DB: tx},
		Penyimpan:        repository.WakilRepository{DB: tx},
		Pencatat:         repository.IngestRunRepository{DB: tx},
		Dapil:            cfg.DapilKPU.Daftar,
		TahunPemilu:      cfg.DapilKPU.TahunPemilu,
		URLLHKPN:         cfg.URLLHKPN,
		PolaURLProfil:    cfg.Sumber.URLDasar + cfg.Sumber.HalamanProfil,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func hitung(t *testing.T, tx pgx.Tx, sql string) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(t.Context(), sql).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestImporDPRD_FixtureSitusResmi_SepuluhDapilDan106Anggota(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporDPRDUji(t, tx)

	if n := hitung(t, tx, `SELECT count(*) FROM dapil WHERE tahun_pemilu = 2024`); n != 10 {
		t.Errorf("jumlah dapil = %d, harap 10", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM anggota`); n != 106 {
		t.Errorf("jumlah anggota = %d, harap 106", n)
	}
}

func TestImporDPRD_FixtureSitusResmi_SetiapKecamatanTepatSatuDapil(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporDPRDUji(t, tx)

	tanpaDapil := hitung(t, tx, `SELECT count(*) FROM wilayah w
WHERE w.tingkat = 'kecamatan' AND (SELECT count(*) FROM dapil_wilayah dw WHERE dw.kode_kecamatan = w.kode) <> 1`)
	if tanpaDapil != 0 {
		t.Errorf("%d kecamatan tidak masuk tepat satu dapil", tanpaDapil)
	}
}

func TestImporDPRD_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporDPRDUji(t, tx)
	jalankanImporDPRDUji(t, tx)

	if n := hitung(t, tx, `SELECT count(*) FROM anggota`); n != 106 {
		t.Errorf("jumlah anggota setelah dua kali impor = %d, harap 106", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM dapil_wilayah`); n != 44 {
		t.Errorf("jumlah kecamatan dapil = %d, harap 44", n)
	}
}

func TestImporDPRD_AnggotaPPPdiFraksiNasDem_FraksiDariHalamanFraksi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporDPRDUji(t, tx)

	var partai, fraksi, url string
	err := tx.QueryRow(t.Context(), `SELECT partai, fraksi, url_profil FROM anggota WHERE id = 'h-matnoor-tindoan-sh'`).Scan(&partai, &fraksi, &url)
	if err != nil {
		t.Fatal(err)
	}
	if partai != "Partai Persatuan Pembangunan" || fraksi != "Fraksi Partai NasDem" || url != "https://dprd-dkijakartaprov.go.id/profile/h-matnoor-tindoan-sh/" {
		t.Errorf("dapat partai %q fraksi %q url %q", partai, fraksi, url)
	}
}
