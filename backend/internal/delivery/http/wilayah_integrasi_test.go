package http

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const batasWaktuResponsCari = 150 * time.Millisecond

type pengambilFixture struct{}

func (p pengambilFixture) Ambil(context.Context) (usecase.DataSumberWilayah, error) {
	berkas, err := filepath.Glob("../../ingest/wilayah/testdata/*.sql")
	if err != nil {
		return usecase.DataSumberWilayah{}, err
	}
	data := usecase.DataSumberWilayah{}
	for _, path := range berkas {
		isi, err := os.ReadFile(path)
		if err != nil {
			return usecase.DataSumberWilayah{}, err
		}
		daftar, err := wilayah.UraiSQL(isi)
		if err != nil {
			return usecase.DataSumberWilayah{}, err
		}
		data.Wilayah = append(data.Wilayah, daftar...)
	}
	data.Provenance.ID, data.Provenance.Sumber, data.Provenance.Label = "prv_uji_wilayah", "kemendagri", "resmi"
	return data, nil
}

func appUjiDB(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	muatWilayahDKI(t, tx)
	repo := repository.WilayahRepository{DB: tx}
	return NewApp(config.Server{}, Layanan{Wilayah: usecase.CariWilayah{Pembaca: repo}})
}

func muatWilayahDKI(t *testing.T, tx pgx.Tx) {
	t.Helper()
	var cfg wilayah.Konfigurasi
	if err := config.BacaYAML("../../../config/wilayah.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	impor := usecase.ImporWilayah{
		Pengambil:   pengambilFixture{},
		Penyimpan:   repository.WilayahRepository{DB: tx},
		Pencatat:    repository.IngestRunRepository{DB: tx},
		AturanAlias: cfg.AturanAlias(),
		Provinsi:    cfg.Provinsi,
	}
	if _, err := tx.Exec(t.Context(), "DELETE FROM wilayah"); err != nil {
		t.Fatal(err)
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func kodeHasil(t *testing.T, data any) []string {
	t.Helper()
	daftar, _ := data.([]any)
	kode := []string{}
	for _, item := range daftar {
		w, _ := item.(map[string]any)
		kode = append(kode, w["kode"].(string))
	}
	return kode
}

func TestCariWilayah_KelapaGading_KecamatanDanKetigaKelurahannya(t *testing.T) {
	app := appUjiDB(t)

	mulai := time.Now()
	hasil := panggil(t, app, "/api/v1/wilayah/cari?q=kelapa%20gading")
	durasi := time.Since(mulai)

	kode := kodeHasil(t, hasil.respons.Data)
	if len(kode) == 0 || kode[0] != "31.72.06" {
		t.Fatalf("hasil pertama harus Kecamatan Kelapa Gading, dapat %v", kode)
	}
	for _, kelurahan := range []string{"31.72.06.1001", "31.72.06.1002", "31.72.06.1003"} {
		if !slices.Contains(kode, kelurahan) {
			t.Errorf("kelurahan %s tidak ada di hasil %v", kelurahan, kode)
		}
	}
	if durasi > batasWaktuResponsCari {
		t.Errorf("respons %v melebihi %v", durasi, batasWaktuResponsCari)
	}
}

func TestKelurahanDiTitik_MallKelapaGading_KelapaGadingTimur(t *testing.T) {
	hasil := panggil(t, appUjiDB(t), "/api/v1/wilayah/titik?lat=-6.1576&lng=106.9081")

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusOK || data["kode"] != "31.72.06.1001" {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	if hierarki, _ := data["hierarki"].([]any); len(hierarki) != 3 {
		t.Errorf("hierarki = %v, harap provinsi, kota, kecamatan", data["hierarki"])
	}
}

func TestKelurahanDiTitik_Surabaya_TidakDitemukan(t *testing.T) {
	hasil := panggil(t, appUjiDB(t), "/api/v1/wilayah/titik?lat=-7.2575&lng=112.7521")

	if hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}

func TestAmbilWilayah_KodeKelurahan_DetailLengkap(t *testing.T) {
	hasil := panggil(t, appUjiDB(t), "/api/v1/wilayah/31.72.06.1003")

	data, _ := hasil.respons.Data.(map[string]any)
	if data["nama"] != "Kelapa Gading Barat" || data["provenance_id"] != "prv_uji_wilayah" {
		t.Fatalf("data = %v", data)
	}
	titik, _ := data["titik_tengah"].(map[string]any)
	if titik["lat"] == nil || titik["lng"] == nil {
		t.Errorf("titik tengah kosong: %v", data["titik_tengah"])
	}
}

func TestAmbilWilayah_KodeTidakAda_404(t *testing.T) {
	hasil := panggil(t, appUjiDB(t), "/api/v1/wilayah/31.72.99")

	if hasil.kode != fiber.StatusNotFound || hasil.respons.Message != "wilayah tidak ditemukan" {
		t.Errorf("kode %d, pesan %q", hasil.kode, hasil.respons.Message)
	}
}
