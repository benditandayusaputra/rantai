package http

import (
	"encoding/csv"
	"net/http/httptest"
	"strconv"
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

const sqlPaketAPIUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES
    ('prv_uji_paket', 'lpse', now(), 'resmi'), ('prv_uji_rup', 'sirup', now(), 'resmi');
INSERT INTO instansi (nama, jenis, kode_kota) VALUES ('SUKU DINAS PENDIDIKAN WILAYAH I JAKARTA UTARA', 'sudin', '31.72');
INSERT INTO penyedia (nama, nama_normal) VALUES ('CV Uji Bangun', 'uji bangun');
INSERT INTO rup (kode, nama, sumber_dana, pagu, jadwal_pelaksanaan_selesai, provenance_id)
VALUES ('rup-api', 'Rehab Ruang Kelas', 'APBD', 500000000, '2026-03-31', 'prv_uji_rup');
INSERT INTO paket (kode, kode_spse, nama, kode_rup, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap,
    hps, nilai_kontrak, penyedia_id, tema_kode, lokasi_tingkat, kode_wilayah, lokasi, kategori, provenance_id)
SELECT 'nontender-api', 'api', 'Rehab Ruang Kelas SDN Kelapa Gading Barat 01', 'rup-api', (SELECT id FROM instansi LIMIT 1), false,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 390000000, 385000000, (SELECT id FROM penyedia LIMIT 1),
    '{pendidikan}', 'kelurahan', w.kode, ST_SetSRID(ST_MakePoint(ST_X(w.titik_tengah), ST_Y(w.titik_tengah)), 4326), 'ruang_kelas', 'prv_uji_paket'
FROM wilayah w WHERE w.nama = 'Kelapa Gading Barat';
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap,
    hps, lokasi_tingkat, kode_wilayah, provenance_id)
SELECT 'tender-api', 'api', 'Perbaikan Saluran Kecamatan Kelapa Gading', (SELECT id FROM instansi LIMIT 1), true,
    'Pekerjaan Konstruksi', 'Tender', 2026, 'Tender Sedang Berjalan', 2500000000, 'kecamatan', w.kode, 'prv_uji_paket'
FROM wilayah w WHERE w.nama = 'Kelapa Gading' AND w.tingkat = 'kecamatan'`

func appUjiPaket(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	if _, err := tx.Exec(t.Context(), sqlPaketAPIUji); err != nil {
		t.Fatal(err)
	}
	u := usecase.BacaPaket{
		Pembaca:  repository.PaketRepository{DB: tx},
		Sekarang: func() time.Time { return time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC) },
	}
	if err := config.BacaYAML("../../../config/arus_anggaran.yaml", &u.Aturan); err != nil {
		t.Fatal(err)
	}
	if err := config.BacaYAML("../../../config/pemeliharaan.yaml", &u.Pemeliharaan); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{Paket: u})
}

func statusRantai(data map[string]any) map[string]string {
	hasil := map[string]string{}
	daftar, _ := data["rantai"].([]any)
	for _, item := range daftar {
		m, _ := item.(map[string]any)
		kode, _ := m["kode"].(string)
		hasil[kode], _ = m["status"].(string)
	}
	return hasil
}

func TestAmbilPaket_NonTenderLengkap_RantaiDanAngkaBersumber(t *testing.T) {
	hasil := panggil(t, appUjiPaket(t), "/api/v1/paket/nontender-api")

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusOK {
		t.Fatalf("kode %d, pesan %s", hasil.kode, hasil.respons.Message)
	}
	harap := map[string]string{"rencana": "ada", "pemilihan": "ada", "pemenang": "ada", "pelaksanaan": "estimasi",
		"lokasi": "ada", "serah_terima": "estimasi", "masa_pemeliharaan": "estimasi", "bukti_warga": "menunggu"}
	if status := statusRantai(data); len(status) != len(harap) {
		t.Errorf("rantai = %v", status)
	} else {
		for mata, s := range harap {
			if status[mata] != s {
				t.Errorf("%s = %s, harap %s", mata, status[mata], s)
			}
		}
	}
	hps, _ := data["hps"].(map[string]any)
	pemeliharaan, _ := data["pemeliharaan"].(map[string]any)
	if hps["provenance_id"] != "prv_uji_paket" || pemeliharaan["status"] != "pemeliharaan" || pemeliharaan["sisa_hari"] == nil {
		t.Errorf("hps %v, pemeliharaan %v", hps, pemeliharaan)
	}
}

func TestAmbilPaket_TenderTanpaRUPDanPemenang_PutusDanMenunggu(t *testing.T) {
	hasil := panggil(t, appUjiPaket(t), "/api/v1/paket/tender-api")

	status := statusRantai(hasil.respons.Data.(map[string]any))
	if status["rencana"] != "putus" || status["pemenang"] != "menunggu" || status["pelaksanaan"] != "putus" {
		t.Errorf("rantai = %v", status)
	}
}

func TestAmbilPaket_KodeTidakAda_404(t *testing.T) {
	if hasil := panggil(t, appUjiPaket(t), "/api/v1/paket/nontender-tidak-ada"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode %d", hasil.kode)
	}
}

func TestDaftarPaketKecamatan_FilterTema_HanyaPaketBertema(t *testing.T) {
	app := appUjiPaket(t)

	semua := panggil(t, app, "/api/v1/kecamatan/31.72.06/paket")
	pendidikan := panggil(t, app, "/api/v1/kecamatan/31.72.06/paket?tema=pendidikan&tahun=2026")

	daftarSemua, _ := semua.respons.Data.([]any)
	daftarPendidikan, _ := pendidikan.respons.Data.([]any)
	if len(daftarSemua) != 2 || len(daftarPendidikan) != 1 || pendidikan.respons.Meta == nil || pendidikan.respons.Meta.Total != 1 {
		t.Errorf("semua %d, pendidikan %d, meta %+v", len(daftarSemua), len(daftarPendidikan), pendidikan.respons.Meta)
	}
}

func TestPetaPaket_TitikKelurahanDanAreaKecamatan(t *testing.T) {
	app := appUjiPaket(t)
	hasil := panggil(t, app, "/api/v1/peta/paket?bbox=106.85,-6.2,106.95,-6.1")

	data, _ := hasil.respons.Data.(map[string]any)
	fitur, _ := data["features"].([]any)
	if hasil.kode != fiber.StatusOK || len(fitur) != 2 {
		t.Errorf("kode %d, fitur %d, pesan %s", hasil.kode, len(fitur), hasil.respons.Message)
	}
	if salah := panggil(t, app, "/api/v1/peta/paket?bbox=1,2,3"); salah.kode != fiber.StatusBadRequest {
		t.Errorf("bbox rusak diterima: %d", salah.kode)
	}
}

func TestPaketTerdekat_DalamRadius_BerjarakEstimasi(t *testing.T) {
	app := appUjiPaket(t)
	detail := panggil(t, app, "/api/v1/paket/nontender-api").respons.Data.(map[string]any)
	titik := detail["lokasi"].(map[string]any)["titik"].(map[string]any)

	hasil := panggil(t, app, "/api/v1/paket/terdekat?lat="+fmtAngka(titik["lat"])+"&lng="+fmtAngka(titik["lng"])+"&radius_m=200")

	daftar, _ := hasil.respons.Data.([]any)
	if len(daftar) != 1 {
		t.Fatalf("terdekat = %v", hasil.respons)
	}
	jarak := daftar[0].(map[string]any)["jarak"].(map[string]any)
	if jarak["label"] != "estimasi" || jarak["provenance_id"] != "prv_rantai_lokasi_paket" {
		t.Errorf("jarak = %v", jarak)
	}
}

func fmtAngka(nilai any) string {
	angka, _ := nilai.(float64)
	return strconv.FormatFloat(angka, 'f', -1, 64)
}

func TestUnduhPaketCSV_Kecamatan_BarisDenganSumber(t *testing.T) {
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/api/v1/unduh/paket.csv?kode_wilayah=31.72.06&tahun=2026", nil)
	resp, err := appUjiPaket(t).Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	baris, err := csv.NewReader(resp.Body).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK || !strings.HasPrefix(resp.Header.Get(fiber.HeaderContentType), "text/csv") ||
		!strings.Contains(resp.Header.Get(fiber.HeaderContentDisposition), "rantai-paket-31.72.06-2026.csv") {
		t.Fatalf("kode %d, header %v", resp.StatusCode, resp.Header)
	}
	if len(baris) != 3 || baris[0][0] != "kode" {
		t.Fatalf("baris = %v", baris)
	}
	kolom := map[string]int{}
	for i, nama := range baris[0] {
		kolom[nama] = i
	}
	tender := baris[1]
	if tender[kolom["kode"]] != "tender-api" || !strings.HasSuffix(tender[kolom["url_sumber"]], "/pengumumanlelang") ||
		tender[kolom["label"]] != "resmi" || tender[kolom["diambil_pada"]] == "" || tender[kolom["data_simulasi"]] != "false" {
		t.Errorf("baris tender = %v", tender)
	}
}

func TestUnduhPaketCSV_KodeWilayahTidakValid_400(t *testing.T) {
	if hasil := panggil(t, appUjiPaket(t), "/api/v1/unduh/paket.csv?kode_wilayah=jakarta"); hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode %d", hasil.kode)
	}
}
