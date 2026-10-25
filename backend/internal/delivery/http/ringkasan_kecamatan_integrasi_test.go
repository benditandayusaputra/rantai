package http

import (
	"image/png"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/kartu"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const sqlDataRingkasanUji = `INSERT INTO provenance (id, sumber, diambil_pada, label, catatan) VALUES
    ('prv_uji_indikator', 'rantai', now(), 'estimasi', 'uji'),
    ('prv_uji_paket', 'lpse', now(), 'resmi', NULL);
INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, provenance_id)
VALUES ('pendidikan', '31.72.06', '2026/2027', 1200, 900, 300, 0.75, 'kurang', 'estimasi', 'prv_uji_indikator');
INSERT INTO instansi (nama, jenis, kode_kota) VALUES ('INSTANSI UJI RINGKASAN', 'sudin', '31.72');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, hps,
    lokasi_tingkat, kode_wilayah, is_batal, provenance_id)
SELECT k, k, k, (SELECT id FROM instansi WHERE nama = 'INSTANSI UJI RINGKASAN'), false, 'Pekerjaan Konstruksi', 'Pengadaan Langsung',
    2026, 'Selesai', hps, tingkat::tingkat_wilayah, wilayah, batal, 'prv_uji_paket'
FROM (VALUES
    ('uji-ringkasan-1', 100000000, 'kelurahan', '31.72.06.1003', false),
    ('uji-ringkasan-2', 50000000, 'kecamatan', '31.72.06', false),
    ('uji-ringkasan-3', 70000000, 'kecamatan', '31.72.06', true),
    ('uji-ringkasan-4', 90000000, 'kota', '31.72', false)) AS v (k, hps, tingkat, wilayah, batal)`

var kunciPengenal = []string{"id", "nomor", "tahun", "tahun_pemilu"}

func appUjiRingkasan(t *testing.T) *fiber.App {
	t.Helper()
	return appUjiRingkasanDengan(t, nil)
}

func appUjiRingkasanDengan(t *testing.T, pembuatKartu PembuatKartu) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	ujiwakil.MuatTema(t, tx)
	if _, err := tx.Exec(t.Context(), sqlDataRingkasanUji); err != nil {
		t.Fatal(err)
	}
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}}
	ringkasan := repository.RingkasanKecamatanRepository{DB: tx}
	return NewApp(config.Server{}, Layanan{Kartu: pembuatKartu, Ringkasan: usecase.RingkasKecamatan{
		Wilayah: repository.WilayahRepository{DB: tx}, Tema: tema, Indikator: ringkasan, Anggaran: ringkasan,
		Wakil:    usecase.BacaWakil{Pembaca: repository.WakilRepository{DB: tx}, Tema: tema},
		Sekarang: func() time.Time { return time.Date(2026, time.October, 23, 9, 0, 0, 0, time.UTC) },
	}})
}

func TestRingkasanKecamatan_KodeKelurahan_DiarahkanDenganSorotan(t *testing.T) {
	hasil := panggil(t, appUjiRingkasan(t), "/api/v1/kecamatan/31.72.06.1003/ringkasan")

	data, _ := hasil.respons.Data.(map[string]any)
	kecamatan, _ := data["kecamatan"].(map[string]any)
	disorot, _ := data["kelurahan_disorot"].(map[string]any)
	if hasil.kode != fiber.StatusOK || kecamatan["kode"] != "31.72.06" || disorot["kode"] != "31.72.06.1003" {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
}

func TestRingkasanKecamatan_TigaBagian_AnggaranTanpaBatalDanTingkatKota(t *testing.T) {
	hasil := panggil(t, appUjiRingkasan(t), "/api/v1/kecamatan/31.72.06/ringkasan")

	data, _ := hasil.respons.Data.(map[string]any)
	anggaran, _ := data["anggaran"].(map[string]any)
	jumlah, _ := anggaran["jumlah_paket"].(map[string]any)
	total, _ := anggaran["total_hps"].(map[string]any)
	if jumlah["nilai"] != float64(2) || total["nilai"] != float64(150000000) || jumlah["label"] != "proksi" {
		t.Errorf("anggaran = %v", anggaran)
	}
	tema, _ := data["tema"].([]any)
	pertama, _ := tema[0].(map[string]any)
	indikator, _ := pertama["indikator"].(map[string]any)
	kesenjangan, _ := indikator["kesenjangan"].(map[string]any)
	if kesenjangan["nilai"] != float64(300) || kesenjangan["label"] != "estimasi" || kesenjangan["satuan"] != "siswa" {
		t.Errorf("indikator = %v", indikator)
	}
	wakilku, _ := data["wakilku"].(map[string]any)
	if anggota := daftarAnggota(t, wakilku["dapil"]); len(anggota) != 9 {
		t.Errorf("harap 9 anggota dapil 2, dapat %d", len(anggota))
	}
}

func TestRingkasanKecamatan_SetiapAngka_BerbentukAngkaBerlabel(t *testing.T) {
	hasil := panggil(t, appUjiRingkasan(t), "/api/v1/kecamatan/31.72.06/ringkasan")

	for _, pelanggaran := range angkaTelanjang(hasil.respons.Data, "data") {
		t.Errorf("angka tanpa label dan provenance: %s", pelanggaran)
	}
}

func angkaTelanjang(nilai any, jalur string) []string {
	switch v := nilai.(type) {
	case map[string]any:
		if _, ok := v["provenance_id"]; ok && v["label"] != nil && v["nilai"] != nil {
			return []string{}
		}
		temuan := []string{}
		for kunci, isi := range v {
			if _, isAngka := isi.(float64); isAngka && slices.Contains(kunciPengenal, kunci) {
				continue
			}
			temuan = append(temuan, angkaTelanjang(isi, jalur+"."+kunci)...)
		}
		return temuan
	case []any:
		temuan := []string{}
		for _, isi := range v {
			temuan = append(temuan, angkaTelanjang(isi, jalur+"[]")...)
		}
		return temuan
	case float64:
		return []string{jalur}
	default:
		return []string{}
	}
}

func TestAngkaTelanjang_NilaiTanpaRujukan_Terdeteksi(t *testing.T) {
	data := map[string]any{"id": float64(1), "jumlah": float64(5), "kursi": map[string]any{"nilai": float64(9), "label": "resmi", "provenance_id": "prv_x"}}

	if temuan := angkaTelanjang(data, "data"); !slices.Equal(temuan, []string{"data.jumlah"}) {
		t.Errorf("temuan = %v", temuan)
	}
}

func TestKartuKecamatan_RingkasanAda_GambarPNG(t *testing.T) {
	var k kartu.Konfigurasi
	if err := config.BacaYAML("../../../config/kartu.yaml", &k); err != nil {
		t.Fatal(err)
	}
	pembuat, err := kartu.Baru(k)
	if err != nil {
		t.Fatal(err)
	}
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/api/v1/kecamatan/31.72.06/kartu.png?bahasa=en", nil)
	resp, err := appUjiRingkasanDengan(t, pembuat).Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	gambar, err := png.Decode(resp.Body)
	if resp.StatusCode != fiber.StatusOK || resp.Header.Get(fiber.HeaderContentType) != "image/png" || err != nil {
		t.Fatalf("kode %d, tipe %s, galat %v", resp.StatusCode, resp.Header.Get(fiber.HeaderContentType), err)
	}
	if ukuran := gambar.Bounds().Size(); ukuran.X != 1080 || ukuran.Y != 1350 {
		t.Errorf("ukuran = %v", ukuran)
	}
}
