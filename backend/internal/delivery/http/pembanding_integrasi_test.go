package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const sqlPembandingUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_banding', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis) VALUES ('Instansi Uji Banding', 'lainnya');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, hps, kategori, provenance_id)
SELECT 'nontender-' || k || n, k || n, 'Saluran ' || n, (SELECT id FROM instansi WHERE nama = 'Instansi Uji Banding'), false,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 100000000 * n, k, 'prv_uji_banding'
FROM (VALUES ('saluran', 9), ('jalan', 8)) AS v (k, jumlah), generate_series(1, v.jumlah) AS n`

func appUjiPembanding(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	if _, err := tx.Exec(t.Context(), sqlPembandingUji); err != nil {
		t.Fatal(err)
	}
	u := usecase.BandingPaket{Pembaca: repository.PaketRepository{DB: tx}, MinPembanding: 8}
	return NewApp(config.Server{}, Layanan{Pembanding: u})
}

func TestPembanding_DelapanPaketSejenis_StatistikTampil(t *testing.T) {
	hasil := panggil(t, appUjiPembanding(t), "/api/v1/paket/nontender-saluran9/pembanding")

	data, _ := hasil.respons.Data.(map[string]any)
	statistik, _ := data["statistik"].(map[string]any)
	median, _ := statistik["median"].(map[string]any)
	if hasil.kode != fiber.StatusOK || data["tersedia"] != true || median["nilai"] != 450000000.0 || len(data["nilai"].([]any)) != 8 {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestPembanding_TujuhPaketSejenis_TidakTampil(t *testing.T) {
	hasil := panggil(t, appUjiPembanding(t), "/api/v1/paket/nontender-jalan1/pembanding")

	data, _ := hasil.respons.Data.(map[string]any)
	if data["tersedia"] != false || data["statistik"] != nil || len(data["nilai"].([]any)) != 0 {
		t.Errorf("data %v", data)
	}
}
