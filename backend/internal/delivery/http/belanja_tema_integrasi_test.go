package http

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const sqlBelanjaUji = `INSERT INTO provenance (id, sumber, diambil_pada, label, catatan) VALUES
    ('prv_uji_belanja', 'lpse', now(), 'resmi', NULL), ('prv_uji_indikator', 'rantai', now(), 'estimasi', 'rumus uji');
INSERT INTO instansi (nama, jenis) VALUES ('Instansi Uji Belanja', 'lainnya');
INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, provenance_id)
SELECT 'pendidikan', kode, '2026/2027', 1200, 600, 600, 0.5, 'kurang', 'estimasi', 'prv_uji_indikator'
FROM wilayah WHERE tingkat = 'kecamatan' AND nama = 'Kelapa Gading';
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, is_batal, jenis_pengadaan, metode, tahun_anggaran, tahap, hps,
    tema_kode, lokasi_tingkat, kode_wilayah, provenance_id)
SELECT v.kode, v.kode, v.kode, (SELECT id FROM instansi WHERE nama = 'Instansi Uji Belanja'), false, v.batal, 'Pekerjaan Konstruksi',
    'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', v.hps, v.tema::text[], v.tingkat::tingkat_wilayah, w.kode, 'prv_uji_belanja'
FROM (VALUES
    ('nontender-b1', false, 100000000, '{pendidikan}', 'kelurahan', 'Kelapa Gading Barat'),
    ('nontender-b2', false, 50000000, '{pendidikan}', 'kecamatan', 'Kelapa Gading'),
    ('nontender-b3', true, 70000000, '{pendidikan}', 'kelurahan', 'Kelapa Gading Timur'),
    ('nontender-b4', false, 900000000, '{}', 'kelurahan', 'Kelapa Gading Barat'),
    ('nontender-b5', false, 300000000, '{pendidikan}', 'kota', 'Kota Administrasi Jakarta Utara')
) AS v (kode, batal, hps, tema, tingkat, nama_wilayah)
JOIN wilayah w ON w.nama = v.nama_wilayah AND w.tingkat = v.tingkat::tingkat_wilayah`

func appUjiBelanja(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	muatTema(t, tx)
	if _, err := tx.Exec(t.Context(), sqlBelanjaUji); err != nil {
		t.Fatal(err)
	}
	u := usecase.BandingBelanja{
		Tema:     usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}},
		Pembaca:  repository.BelanjaTemaRepository{DB: tx},
		Sekarang: func() time.Time { return time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC) },
	}
	return NewApp(config.Server{}, Layanan{Belanja: u})
}

func cariKecamatan(t *testing.T, data map[string]any, nama string) map[string]any {
	t.Helper()
	daftar, _ := data["kecamatan"].([]any)
	for _, item := range daftar {
		k, _ := item.(map[string]any)
		if k["kecamatan"].(map[string]any)["nama"] == nama {
			return k
		}
	}
	t.Fatalf("kecamatan %s tidak ada di %d baris", nama, len(daftar))
	return nil
}

func TestBelanjaTema_PaketTingkatKota_DikeluarkanDariKecamatanDanDilaporkanTerpisah(t *testing.T) {
	hasil := panggil(t, appUjiBelanja(t), "/api/v1/tema/pendidikan/belanja")

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusOK {
		t.Fatalf("kode %d, pesan %s", hasil.kode, hasil.respons.Message)
	}
	kg := cariKecamatan(t, data, "Kelapa Gading")
	belanja, _ := kg["belanja"].(map[string]any)
	total, _ := belanja["total_hps"].(map[string]any)
	jumlah, _ := belanja["jumlah_paket"].(map[string]any)
	if total["nilai"] != 150000000.0 || jumlah["nilai"] != 2.0 {
		t.Errorf("belanja Kelapa Gading = %v", belanja)
	}
	rasio, _ := kg["indikator"].(map[string]any)["rasio"].(map[string]any)
	if rasio["nilai"] != 0.5 || rasio["label"] != "estimasi" {
		t.Errorf("rasio = %v", rasio)
	}
	luar, _ := data["luar_kecamatan"].([]any)
	if len(luar) != 1 || luar[0].(map[string]any)["belanja"].(map[string]any)["total_hps"].(map[string]any)["nilai"] != 300000000.0 {
		t.Errorf("luar kecamatan = %v", luar)
	}
	if len(data["kecamatan"].([]any)) != 44 || data["tahun"] != 2026.0 {
		t.Errorf("jumlah kecamatan %d, tahun %v", len(data["kecamatan"].([]any)), data["tahun"])
	}
}

func TestBelanjaTema_TemaTidakAda_404(t *testing.T) {
	if hasil := panggil(t, appUjiBelanja(t), "/api/v1/tema/tidak-ada/belanja"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode %d", hasil.kode)
	}
}
