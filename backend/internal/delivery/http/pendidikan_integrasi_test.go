package http

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const batasUkuranPetaDKI = 300 * 1024

const rincianUji = `{"metode_calon_lulusan":"per_sekolah","calon_lulusan":{"nilai":1000,"label":"resmi"},
"kursi_negeri":{"nilai":600,"label":"resmi"},"kursi_spmb_bersama":{"nilai":20,"label":"resmi"},
"kursi_madrasah":{"nilai":null,"label":"proksi"},"kesenjangan_dengan_spmb_bersama":{"nilai":380,"label":"resmi"},
"kesenjangan_semua_jalur":{"nilai":null,"label":"proksi"}}`

func siapkanIndikatorUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	muatTema(t, tx)
	langkah := []string{
		`DELETE FROM indikator_kecamatan`,
		`INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, rincian, provenance_id)
SELECT 'pendidikan', kode, '2026/2027', 1000, 600, 400, 0.6, 'kurang', 'resmi', $1::jsonb, provenance_id
FROM wilayah WHERE tingkat = 'kecamatan'`,
		`INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, lokasi, provenance_id)
VALUES ('90000201', 'SMP UJI NEGERI', 'SMP', 'menengah_pertama', 'negeri', false, '31.72.06',
    ST_SetSRID(ST_MakePoint(106.91, -6.16), 4326), (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))
ON CONFLICT (npsn) DO NOTHING`,
		`INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
VALUES ('90000201', '9999/10000', 'negeri', 288, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))`,
	}
	for i, sql := range langkah {
		argumen := []any{}
		if i == 1 {
			argumen = append(argumen, rincianUji)
		}
		if _, err := tx.Exec(t.Context(), sql, argumen...); err != nil {
			t.Fatal(err)
		}
	}
}

func appUjiPendidikan(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	siapkanIndikatorUji(t, tx)
	var cfg config.Config
	if err := config.BacaYAML("../../../config/api.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{Indikator: usecase.BacaIndikator{
		Pembaca: repository.IndikatorBacaRepository{DB: tx, Peta: cfg.Peta},
		Tema:    repository.TemaRepository{DB: tx},
	}})
}

func TestPetaTema_Pendidikan_SeluruhKecamatanDiBawah300KB(t *testing.T) {
	app := appUjiPendidikan(t)
	resp, err := app.Test(httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/api/v1/peta/tema/pendidikan", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	isi, err := io.ReadAll(resp.Body)

	if err != nil || resp.StatusCode != fiber.StatusOK || len(isi) >= batasUkuranPetaDKI {
		t.Errorf("status %d, ukuran %d byte, err %v", resp.StatusCode, len(isi), err)
	}
	if resp.Header.Get(fiber.HeaderCacheControl) == "" {
		t.Error("respons peta harus bisa di-cache")
	}
}

func TestPetaTema_Pendidikan_FiturBerisiAngkaBerlabel(t *testing.T) {
	hasil := panggil(t, appUjiPendidikan(t), "/api/v1/peta/tema/pendidikan")

	data, _ := hasil.respons.Data.(map[string]any)
	fitur, _ := data["features"].([]any)
	if len(fitur) != 44 {
		t.Fatalf("jumlah fitur = %d", len(fitur))
	}
	properti := fitur[0].(map[string]any)["properties"].(map[string]any)
	rasio, _ := properti["rasio"].(map[string]any)
	if properti["kelas"] != "kurang" || rasio["label"] != "resmi" || rasio["provenance_id"] == "" {
		t.Errorf("properti = %v", properti)
	}
}

func TestPetaTema_TemaTidakAda_404(t *testing.T) {
	if hasil := panggil(t, appUjiPendidikan(t), "/api/v1/peta/tema/antariksa"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d", hasil.kode)
	}
}

func TestDetailPendidikan_KelapaGading_RincianKursiDanKeterbatasan(t *testing.T) {
	hasil := panggil(t, appUjiPendidikan(t), "/api/v1/tema/pendidikan/kecamatan/31.72.06")

	data, _ := hasil.respons.Data.(map[string]any)
	kursiNegeri, _ := data["kursi_negeri"].(map[string]any)
	keterbatasan, _ := data["keterbatasan"].([]any)
	if hasil.kode != fiber.StatusOK || kursiNegeri["nilai"] != 600.0 || data["kursi_madrasah"] != nil || len(keterbatasan) != 3 {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestDetailPendidikan_SatuPeriode_TrenKosong(t *testing.T) {
	hasil := panggil(t, appUjiPendidikan(t), "/api/v1/tema/pendidikan/kecamatan/31.72.06")

	data, _ := hasil.respons.Data.(map[string]any)
	tren, isDaftar := data["tren"].([]any)
	if hasil.kode != fiber.StatusOK || !isDaftar || len(tren) != 0 {
		t.Errorf("kode %d, tren %v", hasil.kode, data["tren"])
	}
}

func TestDetailPendidikan_DuaPeriode_TrenUrutPeriode(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanIndikatorUji(t, tx)
	_, err := tx.Exec(t.Context(), `INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, rincian, provenance_id)
SELECT 'pendidikan', kode, '2025/2026', 900, 600, 300, 1.5, 'kurang', 'resmi', '{}'::jsonb, provenance_id FROM wilayah WHERE kode = '31.72.06'`)
	if err != nil {
		t.Fatal(err)
	}
	app := NewApp(config.Server{}, Layanan{Indikator: usecase.BacaIndikator{
		Pembaca: repository.IndikatorBacaRepository{DB: tx},
		Tema:    repository.TemaRepository{DB: tx},
	}})

	hasil := panggil(t, app, "/api/v1/tema/pendidikan/kecamatan/31.72.06")

	data, _ := hasil.respons.Data.(map[string]any)
	tren, _ := data["tren"].([]any)
	if len(tren) != 2 {
		t.Fatalf("tren = %v", data["tren"])
	}
	pertama, _ := tren[0].(map[string]any)
	rasio, _ := pertama["rasio"].(map[string]any)
	if pertama["periode"] != "2025/2026" || rasio["nilai"] != 1.5 || rasio["provenance_id"] == "" {
		t.Errorf("titik pertama = %v", pertama)
	}
}

func TestDetailPendidikan_KodeKelurahan_400(t *testing.T) {
	if hasil := panggil(t, appUjiPendidikan(t), "/api/v1/tema/pendidikan/kecamatan/31.72.06.1001"); hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode = %d", hasil.kode)
	}
}

func TestSekolahKecamatan_KelapaGading_SMPDenganDayaTampungResmi(t *testing.T) {
	hasil := panggil(t, appUjiPendidikan(t), "/api/v1/tema/pendidikan/kecamatan/31.72.06/sekolah")

	daftar, _ := hasil.respons.Data.([]any)
	if len(daftar) == 0 {
		t.Fatalf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
	pertama := daftar[0].(map[string]any)
	dayaTampung, _ := pertama["daya_tampung"].(map[string]any)
	if pertama["npsn"] != "90000201" || dayaTampung["nilai"] != 288.0 || dayaTampung["label"] != "resmi" || pertama["lokasi"] == nil {
		t.Errorf("sekolah pertama = %v", pertama)
	}
}

const rincianSMASMKUji = `{"metode_calon_lulusan":"per_sekolah","calon_lulusan":{"nilai":1500,"label":"resmi"},
"kursi_negeri":{"nilai":540,"label":"resmi"},"kursi_spmb_bersama":{"nilai":30,"label":"resmi"},
"kursi_madrasah":{"nilai":null,"label":"proksi"},"kesenjangan_dengan_spmb_bersama":{"nilai":930,"label":"resmi"},
"kesenjangan_semua_jalur":{"nilai":null,"label":"proksi"},
"kursi_negeri_per_bentuk":[{"bentuk":"SMA","kursi":324},{"bentuk":"SMK","kursi":216}],
"kursi_spmb_bersama_per_bentuk":[{"bentuk":"SMA","kursi":0},{"bentuk":"SMK","kursi":30}],"batas_murid_per_rombel":36}`

func appUjiSMASMK(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	siapkanIndikatorUji(t, tx)
	langkah := []string{
		`INSERT INTO indikator_kecamatan (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, rincian, provenance_id)
SELECT 'pendidikan_sma_smk', kode, '2026/2027', 1500, 540, 960, 0.36, 'sangat_kurang', 'resmi', $1::jsonb, provenance_id
FROM wilayah WHERE kode = '31.72.06'`,
		`INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, provenance_id)
VALUES ('90000211', 'SMAN UJI', 'SMA', 'menengah_atas', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000212', 'SMKN UJI', 'SMK', 'menengah_atas', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))
ON CONFLICT (npsn) DO NOTHING`,
		`INSERT INTO provenance (id, sumber, url, diambil_pada, hash_konten, label, catatan)
VALUES ('prv_uji_smk', 'disdik', 'https://contoh.test/smk', now(), 'x', 'resmi', 'SMK uji')`,
		`INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
VALUES ('90000211', '9999/10000', 'negeri', 324, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000212', '9999/10000', 'negeri', 216, 'prv_uji_smk')`,
	}
	for i, sql := range langkah {
		argumen := []any{}
		if i == 0 {
			argumen = append(argumen, rincianSMASMKUji)
		}
		if _, err := tx.Exec(t.Context(), sql, argumen...); err != nil {
			t.Fatal(err)
		}
	}
	return NewApp(config.Server{}, Layanan{Indikator: usecase.BacaIndikator{
		Pembaca: repository.IndikatorBacaRepository{DB: tx},
		Tema:    repository.TemaRepository{DB: tx},
	}})
}

func TestDetailPendidikan_SMASMK_KursiPerBentukDenganProvenanceMasingMasing(t *testing.T) {
	hasil := panggil(t, appUjiSMASMK(t), "/api/v1/tema/pendidikan_sma_smk/kecamatan/31.72.06")

	data, _ := hasil.respons.Data.(map[string]any)
	perBentuk, _ := data["kursi_negeri_per_bentuk"].([]any)
	if hasil.kode != fiber.StatusOK || data["tema"] != "pendidikan_sma_smk" || len(perBentuk) != 2 {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	smk := perBentuk[1].(map[string]any)
	kursiSMK, _ := smk["kursi"].(map[string]any)
	if smk["bentuk"] != "SMK" || kursiSMK["nilai"] != 216.0 || kursiSMK["provenance_id"] != "prv_uji_smk" {
		t.Errorf("SMK = %v", smk)
	}
	kursiNegeri, _ := data["kursi_negeri"].(map[string]any)
	if kursiNegeri["provenance_id"] != data["provenance_id"] {
		t.Errorf("kursi negeri gabungan harus merujuk provenance rumus, dapat %v", kursiNegeri)
	}
}

func TestDetailPendidikan_SMASMK_RumusDanKeterbatasanJenjangAtas(t *testing.T) {
	hasil := panggil(t, appUjiSMASMK(t), "/api/v1/tema/pendidikan_sma_smk/kecamatan/31.72.06")

	data, _ := hasil.respons.Data.(map[string]any)
	rumus, _ := data["rumus"].(map[string]any)
	ambang, _ := rumus["ambang_rasio"].(map[string]any)
	if rumus["bentuk_asal"] != "SMP" || rumus["tingkat_asal"] != 9.0 || rumus["tingkat_tujuan"] != 10.0 ||
		rumus["batas_murid_per_rombel"] != 36.0 || ambang["sangat_kurang"] != 0.5 {
		t.Errorf("rumus = %v", rumus)
	}
	keterbatasan, _ := data["keterbatasan"].([]any)
	if len(keterbatasan) != 4 || keterbatasan[1] != "sma_smk_lintas_kecamatan" || keterbatasan[2] != "calon_lulusan_tanpa_mts" {
		t.Errorf("keterbatasan = %v", keterbatasan)
	}
}

func TestSekolahKecamatan_SMASMK_HanyaJenjangMenengahAtas(t *testing.T) {
	hasil := panggil(t, appUjiSMASMK(t), "/api/v1/tema/pendidikan_sma_smk/kecamatan/31.72.06/sekolah")

	daftar, _ := hasil.respons.Data.([]any)
	if len(daftar) != 2 {
		t.Fatalf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
	for _, butir := range daftar {
		if bentuk := butir.(map[string]any)["bentuk"]; bentuk != "SMA" && bentuk != "SMK" {
			t.Errorf("bentuk %v bukan jenjang menengah atas", bentuk)
		}
	}
}

func TestDetailPendidikan_TemaTanpaPeralihan_404(t *testing.T) {
	if hasil := panggil(t, appUjiPendidikan(t), "/api/v1/tema/banjir/kecamatan/31.72.06"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d", hasil.kode)
	}
}
