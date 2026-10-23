package http

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const batasWaktuResponsCari = 150 * time.Millisecond

func appUjiDB(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	repo := repository.WilayahRepository{DB: tx}
	return NewApp(config.Server{}, Layanan{Wilayah: usecase.CariWilayah{Pembaca: repo}})
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
	provenanceID, _ := data["provenance_id"].(string)
	if data["nama"] != "Kelapa Gading Barat" || !strings.HasPrefix(provenanceID, "prv_") {
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
