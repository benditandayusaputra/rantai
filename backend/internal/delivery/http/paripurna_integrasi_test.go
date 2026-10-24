package http

import (
	"fmt"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/paripurna"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const dirFixtureParipurna = "../../ingest/paripurna/testdata/"

func appUjiParipurna(t *testing.T) (*fiber.App, pgx.Tx) {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	impor := usecase.ImporParipurna{
		Pengambil: paripurna.Sumber{Konfigurasi: paripurna.Konfigurasi{
			PathSikapFraksi: dirFixtureParipurna + "sikap_fraksi.csv", PathAspirasiReses: dirFixtureParipurna + "aspirasi_reses.csv",
			URLSumber: "https://dprd-dkijakartaprov.go.id/",
		}, Sekarang: time.Now},
		Penyimpan: repository.ParipurnaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		Terlarang: domain.KosakataTerlarang{},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{Paripurna: usecase.BacaParipurna{Pembaca: repository.ParipurnaRepository{DB: tx}}}), tx
}

func dapilKecamatan(t *testing.T, tx pgx.Tx, kode string) int {
	t.Helper()
	var id int
	if err := tx.QueryRow(t.Context(), `SELECT dapil_id FROM dapil_wilayah WHERE kode_kecamatan = $1`, kode).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestSikapFraksi_NamaTerenkode_DaftarDenganTautan(t *testing.T) {
	app, _ := appUjiParipurna(t)

	hasil := panggil(t, app, "/api/v1/fraksi/Fraksi%20Partai%20Golkar/sikap")

	daftar, _ := hasil.respons.Data.([]any)
	if hasil.kode != fiber.StatusOK || len(daftar) != 1 {
		t.Fatalf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
	sikap, _ := daftar[0].(map[string]any)
	if sikap["agenda"] != "apbd" || sikap["tanggal"] != "2026-10-19" || sikap["url"] == "" || sikap["provenance_id"] == "" {
		t.Errorf("sikap = %v", sikap)
	}
}

func TestSikapFraksi_FraksiTidakAda_404(t *testing.T) {
	app, _ := appUjiParipurna(t)

	if hasil := panggil(t, app, "/api/v1/fraksi/Fraksi%20Contoh/sikap"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}

func TestAspirasiReses_FilterDapil_EntriWilayahDapilDanSeluruhProvinsi(t *testing.T) {
	app, tx := appUjiParipurna(t)
	dapilPuloGadung := dapilKecamatan(t, tx, "31.75.02")
	dapilLain := dapilKecamatan(t, tx, "31.73.08")

	dalam, _ := panggil(t, app, fmt.Sprintf("/api/v1/aspirasi-reses?dapil=%d", dapilPuloGadung)).respons.Data.([]any)
	luar, _ := panggil(t, app, fmt.Sprintf("/api/v1/aspirasi-reses?dapil=%d&bidang=Kesehatan", dapilLain)).respons.Data.([]any)

	if len(dalam) != 2 || len(luar) != 1 {
		t.Fatalf("dapil Pulo Gadung %d entri, dapil lain %d entri", len(dalam), len(luar))
	}
	kayuPutih, _ := dalam[0].(map[string]any)
	wilayah, _ := kayuPutih["wilayah"].(map[string]any)
	if wilayah["kode"] != "31.75.02.1006" || wilayah["tingkat"] != "kelurahan" {
		t.Errorf("entri pertama = %v", kayuPutih)
	}
	if provinsi, _ := luar[0].(map[string]any); provinsi["wilayah"] != nil {
		t.Errorf("entri dapil lain harus tanpa wilayah, dapat %v", provinsi)
	}
}

func TestAspirasiReses_DapilTidakSah_KodeError(t *testing.T) {
	app, _ := appUjiParipurna(t)
	kasus := map[string]int{
		"/api/v1/aspirasi-reses?dapil=satu": fiber.StatusBadRequest,
		"/api/v1/aspirasi-reses?dapil=0":    fiber.StatusBadRequest,
		"/api/v1/aspirasi-reses?dapil=9999": fiber.StatusNotFound,
	}
	for path, harap := range kasus {
		if hasil := panggil(t, app, path); hasil.kode != harap {
			t.Errorf("%s: kode %d, harap %d", path, hasil.kode, harap)
		}
	}
}
