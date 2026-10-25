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

const sqlDasborUji = `INSERT INTO provenance (id, sumber, diambil_pada, label, catatan) VALUES
    ('prv_uji_dasbor_lpse', 'lpse', now(), 'resmi', NULL), ('prv_uji_dasbor_dprd', 'dprd', now(), 'resmi', NULL),
    ('prv_uji_dasbor_indikator', 'rantai', now(), 'proksi', 'rumus uji');
INSERT INTO instansi (nama, jenis) VALUES ('Instansi Uji Dasbor', 'lainnya'), ('Pemda Uji Dasbor', 'lainnya');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, is_batal, is_simulation, jenis_pengadaan, metode, tahun_anggaran,
    tahap, hps, tanggal_dibuat, lokasi_tingkat, kode_wilayah, provenance_id)
SELECT v.kode, v.kode, v.kode, (SELECT id FROM instansi WHERE nama = 'Instansi Uji Dasbor'), false, v.batal, v.simulasi,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', v.tahun, 'Paket Sudah Selesai', v.hps, v.tanggal::date,
    v.tingkat::tingkat_wilayah, w.kode, 'prv_uji_dasbor_lpse'
FROM (VALUES
    ('dasbor-1', false, false, 2025, 100000000, '2025-03-10', 'kelurahan', 'Kelapa Gading Barat'),
    ('dasbor-2', false, false, 2025, 50000000, '2025-03-20', 'kecamatan', 'Kelapa Gading'),
    ('dasbor-3', false, false, 2025, 30000000, NULL, 'kota', 'Kota Administrasi Jakarta Utara'),
    ('dasbor-4', true, false, 2025, 70000000, '2025-04-01', 'kelurahan', 'Kelapa Gading Timur'),
    ('dasbor-5', false, true, 2025, 90000000, '2025-04-01', 'kelurahan', 'Kelapa Gading Timur'),
    ('dasbor-6', false, false, 2026, 10000000, '2026-01-05', 'kelurahan', 'Kelapa Gading Timur')
) AS v (kode, batal, simulasi, tahun, hps, tanggal, tingkat, nama_wilayah)
JOIN wilayah w ON w.nama = v.nama_wilayah AND w.tingkat = v.tingkat::tingkat_wilayah;
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, hps, provenance_id)
SELECT 'dasbor-7', 'dasbor-7', 'dasbor-7', id, false, 'Jasa Lainnya', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 5000000000,
    'prv_uji_dasbor_lpse'
FROM instansi WHERE nama = 'Pemda Uji Dasbor';
INSERT INTO umpan_kegiatan (jenis_sumber, media_nama, media_domain, judul, url, terbit_pada, tag_isu, disembunyikan, provenance_id)
VALUES
    ('resmi_dprd', 'DPRD DKI', 'dprd-dkijakartaprov.go.id', 'a', 'https://uji.test/a', '2025-12-31T18:00:00Z', '{sekolah}', false, 'prv_uji_dasbor_dprd'),
    ('resmi_dprd', 'DPRD DKI', 'dprd-dkijakartaprov.go.id', 'b', 'https://uji.test/b', '2026-02-10T03:00:00Z', '{sekolah,guru}', false, 'prv_uji_dasbor_dprd'),
    ('resmi_dprd', 'DPRD DKI', 'dprd-dkijakartaprov.go.id', 'c', 'https://uji.test/c', '2026-02-11T03:00:00Z', '{}', true, 'prv_uji_dasbor_dprd');
INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, provenance_id)
SELECT 'pendidikan', w.kode, v.periode, v.kebutuhan, v.kapasitas, v.kebutuhan - v.kapasitas, NULL, v.kelas, v.label::label_kepastian,
    'prv_uji_dasbor_indikator'
FROM (VALUES
    ('Kelapa Gading', '2026/2027', 1200, 600, 'kurang', 'proksi'),
    ('Koja', '2026/2027', 500, 700, 'cukup', 'resmi'),
    ('Koja', '2025/2026', 400, 700, 'cukup', 'resmi')
) AS v (nama, periode, kebutuhan, kapasitas, kelas, label)
JOIN wilayah w ON w.nama = v.nama AND w.tingkat = 'kecamatan'`

func appUjiDasbor(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	muatTema(t, tx)
	if _, err := tx.Exec(t.Context(), sqlDasborUji); err != nil {
		t.Fatal(err)
	}
	u := usecase.BacaDasbor{
		Tema:                 usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}},
		Pembaca:              repository.DasborRepository{DB: tx},
		Sekarang:             func() time.Time { return time.Date(2026, time.October, 25, 0, 0, 0, 0, time.UTC) },
		InstansiTanpaRincian: "Pemda Uji Dasbor",
	}
	return NewApp(config.Server{}, Layanan{Dasbor: u})
}

func dataDasbor(t *testing.T, app *fiber.App, path string) map[string]any {
	t.Helper()
	hasil := panggil(t, app, path)
	if hasil.kode != fiber.StatusOK {
		t.Fatalf("kode %d, pesan %s", hasil.kode, hasil.respons.Message)
	}
	data, _ := hasil.respons.Data.(map[string]any)
	return data
}

func isiSeri(data map[string]any, nama string) []map[string]any {
	seri, _ := data[nama].(map[string]any)
	daftar, _ := seri["isi"].([]any)
	hasil := make([]map[string]any, len(daftar))
	for i, item := range daftar {
		hasil[i], _ = item.(map[string]any)
	}
	return hasil
}

func TestDasborUang_TanpaTahun_PilihTahunTerlengkapDanAbaikanBatalSertaSimulasi(t *testing.T) {
	data := dataDasbor(t, appUjiDasbor(t), "/api/v1/dasbor/uang")

	tahunan := isiSeri(data, "tahunan")
	if data["tahun"] != 2025.0 || len(tahunan) != 2 || tahunan[0]["jumlah"] != 3.0 || tahunan[0]["nilai"] != 180000000.0 {
		t.Fatalf("tahun %v, tahunan %v", data["tahun"], tahunan)
	}
	bulanan := isiSeri(data, "bulanan")
	if len(bulanan) != 12 || bulanan[2]["jumlah"] != 2.0 || bulanan[3]["jumlah"] != 0.0 {
		t.Errorf("bulanan %v", bulanan)
	}
	kecamatan := isiSeri(data, "kecamatan")
	if len(kecamatan) != 1 || kecamatan[0]["nama"] != "Kelapa Gading" || kecamatan[0]["jumlah"] != 2.0 {
		t.Errorf("kecamatan %v", kecamatan)
	}
	if seri, _ := data["kota"].(map[string]any); seri["label"] != "resmi" || seri["provenance_id"] != "prv_uji_dasbor_lpse" {
		t.Errorf("rujukan kota %v", seri)
	}
}

func TestDasborTanggungJawab_KegiatanDihitungMenurutWaktuJakarta(t *testing.T) {
	data := dataDasbor(t, appUjiDasbor(t), "/api/v1/dasbor/tanggung-jawab?tahun=2026")

	tahunan := isiSeri(data, "kegiatan_tahunan")
	if len(tahunan) != 1 || tahunan[0]["kunci"] != "2026" || tahunan[0]["jumlah"] != 2.0 {
		t.Fatalf("kegiatan tahunan %v", tahunan)
	}
	if bulanan := isiSeri(data, "kegiatan_bulanan"); bulanan[0]["jumlah"] != 1.0 || bulanan[1]["jumlah"] != 1.0 {
		t.Errorf("kegiatan bulanan %v", bulanan)
	}
	if topik := isiSeri(data, "topik"); len(topik) != 2 || topik[0]["kunci"] != "sekolah" || topik[0]["jumlah"] != 2.0 {
		t.Errorf("topik %v", topik)
	}
	if instansi := isiSeri(data, "instansi"); len(instansi) != 1 || instansi[0]["nama"] != "Instansi Uji Dasbor" {
		t.Errorf("instansi %v", instansi)
	}
}

func TestDasborLayanan_RekapPerPeriodeDanKotaDenganLabelTerlemah(t *testing.T) {
	data := dataDasbor(t, appUjiDasbor(t), "/api/v1/tema/pendidikan/dasbor")

	periode, _ := data["periode"].([]any)
	terbaru, _ := periode[len(periode)-1].(map[string]any)
	perKelas, _ := terbaru["per_kelas"].(map[string]any)
	if len(periode) != 2 || terbaru["kebutuhan"] != 1700.0 || perKelas["kurang"] != 1.0 || perKelas["cukup"] != 1.0 {
		t.Fatalf("periode %v", periode)
	}
	kota, _ := data["kota"].([]any)
	if len(kota) != 1 || kota[0].(map[string]any)["kapasitas"] != 1300.0 || data["label"] != "proksi" {
		t.Errorf("kota %v, label %v", kota, data["label"])
	}
}

func TestDasborLayanan_TemaTidakAda_404(t *testing.T) {
	hasil := panggil(t, appUjiDasbor(t), "/api/v1/tema/tidak-ada/dasbor")

	if hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode %d", hasil.kode)
	}
}

func TestDasborUang_TahunBukanAngka_400(t *testing.T) {
	hasil := panggil(t, appUji(), "/api/v1/dasbor/uang?tahun=dua")

	if hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode %d", hasil.kode)
	}
}
