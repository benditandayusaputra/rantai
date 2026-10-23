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
