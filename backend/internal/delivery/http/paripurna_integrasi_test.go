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
	baca := usecase.BacaParipurna{Pembaca: repository.ParipurnaRepository{DB: tx}}
	if err := config.BacaYAML("../../../config/aspirasi_anggaran.yaml", &baca.Anggaran); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{Paripurna: baca}), tx
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

const sqlPaketTerkaitUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_terkait', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis, kode_kota) VALUES ('SUKU DINAS KESEHATAN JAKARTA TIMUR', 'sudin', '31.75');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, hps,
    lokasi_tingkat, kode_wilayah, kategori, is_simulation, provenance_id)
SELECT x.kode, x.kode, x.nama, (SELECT id FROM instansi WHERE nama = 'SUKU DINAS KESEHATAN JAKARTA TIMUR'), false,
    'Pengadaan Barang', 'Pengadaan Langsung', x.tahun, 'Paket Sudah Selesai', 150000000, 'kelurahan', '31.75.02.1006', x.kategori,
    x.simulasi, 'prv_uji_terkait'
FROM (VALUES
    ('uji-posyandu-2027', 'Pengadaan Alat Kesehatan Posyandu Kelurahan Kayu Putih', 2027, NULL, false),
    ('uji-posyandu-2026', 'Pengadaan Alat Kesehatan Posyandu Kelurahan Kayu Putih', 2026, NULL, false),
    ('uji-jalan-2027', 'Perbaikan Jalan Lingkungan Kayu Putih', 2027, 'jalan', false),
    ('simulasi-posyandu-2027', 'Pengadaan Alat Posyandu Simulasi', 2027, NULL, true)
) AS x (kode, nama, tahun, kategori, simulasi)`

func idAspirasi(t *testing.T, tx pgx.Tx, kondisi string) int64 {
	t.Helper()
	var id int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM aspirasi_reses WHERE `+kondisi+` ORDER BY id LIMIT 1`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestPaketTerkait_KataKunciBidangTahunBerikutnya_KemungkinanTerkait(t *testing.T) {
	app, tx := appUjiParipurna(t)
	if _, err := tx.Exec(t.Context(), sqlPaketTerkaitUji); err != nil {
		t.Fatal(err)
	}
	id := idAspirasi(t, tx, "kode_wilayah = '31.75.02.1006'")

	hasil := panggil(t, app, fmt.Sprintf("/api/v1/aspirasi-reses/%d/paket-terkait", id))

	data, _ := hasil.respons.Data.(map[string]any)
	paket, _ := data["paket"].([]any)
	if hasil.kode != fiber.StatusOK || data["status"] != "kemungkinan_terkait" || data["tahun_anggaran"] != float64(2027) || len(paket) != 1 {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	if p, _ := paket[0].(map[string]any); p["kode"] != "uji-posyandu-2027" || p["dasar_cocok"] != "kata_kunci" {
		t.Errorf("paket = %v", p)
	}
}

func TestPaketTerkait_TanpaPaketCocok_BelumDitemukan(t *testing.T) {
	app, tx := appUjiParipurna(t)
	id := idAspirasi(t, tx, "kode_wilayah = '31.75.02.1006'")

	data, _ := panggil(t, app, fmt.Sprintf("/api/v1/aspirasi-reses/%d/paket-terkait", id)).respons.Data.(map[string]any)

	if paket, _ := data["paket"].([]any); data["status"] != "belum_ditemukan" || paket == nil || len(paket) != 0 {
		t.Errorf("data = %v", data)
	}
}

func TestPaketTerkait_AspirasiTanpaWilayah_StatusTanpaWilayah(t *testing.T) {
	app, tx := appUjiParipurna(t)
	id := idAspirasi(t, tx, "kode_wilayah IS NULL")

	if data, _ := panggil(t, app, fmt.Sprintf("/api/v1/aspirasi-reses/%d/paket-terkait", id)).respons.Data.(map[string]any); data["status"] != "tanpa_wilayah" {
		t.Errorf("data = %v", data)
	}
}

func TestPaketTerkait_IDTidakSah_KodeError(t *testing.T) {
	app, _ := appUjiParipurna(t)
	kasus := map[string]int{
		"/api/v1/aspirasi-reses/abc/paket-terkait":       fiber.StatusBadRequest,
		"/api/v1/aspirasi-reses/999999999/paket-terkait": fiber.StatusNotFound,
	}
	for path, harap := range kasus {
		if hasil := panggil(t, app, path); hasil.kode != harap {
			t.Errorf("%s: kode %d, harap %d", path, hasil.kode, harap)
		}
	}
}
