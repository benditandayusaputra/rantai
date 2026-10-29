package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

func appUjiKota(t *testing.T) (*fiber.App, pgx.Tx) {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	siapkanIndikatorUji(t, tx)
	return NewApp(config.Server{}, Layanan{PendidikanKota: usecase.BacaPendidikanKota{
		Pembaca: repository.PendidikanKotaRepository{DB: tx}, Wilayah: repository.WilayahRepository{DB: tx},
	}}), tx
}

func TestPendidikanKota_HanyaKecamatanDalamKotaDanSatuPeriode(t *testing.T) {
	app, _ := appUjiKota(t)
	hasil := panggil(t, app, "/api/v1/tema/pendidikan/kota/31.72")
	data := hasil.respons.Data.(map[string]any)
	r := data["ringkasan"].(map[string]any)
	if hasil.kode != 200 || len(data["kecamatan"].([]any)) != 6 || len(data["kabupaten_kota"].([]any)) != 6 || data["periode"] != "2026/2027" {
		t.Fatalf("respons salah: %+v", data)
	}
	if r["kebutuhan"].(map[string]any)["nilai"] != 6000.0 || r["kapasitas"].(map[string]any)["nilai"] != 3600.0 {
		t.Fatalf("rekap kota salah: %+v", r)
	}
	if r["rasio"].(map[string]any)["nilai"] != 0.6 || r["kecamatan_kurang"] != 6.0 {
		t.Fatalf("rasio kota salah: %+v", r)
	}
}

func TestPendidikanKota_KecamatanTanpaIndikatorTerbaruTetapDitampilkan(t *testing.T) {
	app, tx := appUjiKota(t)
	_, err := tx.Exec(t.Context(), `UPDATE indikator_kecamatan SET periode = '2025/2026' WHERE kode_kecamatan = '31.72.06'`)
	if err != nil {
		t.Fatal(err)
	}
	data := panggil(t, app, "/api/v1/tema/pendidikan/kota/31.72").respons.Data.(map[string]any)
	r := data["ringkasan"].(map[string]any)
	if r["kebutuhan"] != nil || r["rasio"] != nil || r["kecamatan_dengan_data"] != 5.0 {
		t.Fatalf("data lama tidak boleh dicampur: %+v", r)
	}
	for _, baris := range data["kecamatan"].([]any) {
		k := baris.(map[string]any)
		if k["kode"] == "31.72.06" && k["indikator"] != nil {
			t.Fatal("indikator lama harus kosong")
		}
	}
}

func TestPendidikanKota_SekolahTidakGandaKarenaJalurAtauTahunKursi(t *testing.T) {
	app, tx := appUjiKota(t)
	_, err := tx.Exec(t.Context(), `INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
SELECT npsn, '2026/2027', 'spmb_bersama', 50, provenance_id FROM sekolah WHERE npsn = '90000201'`)
	if err != nil {
		t.Fatal(err)
	}
	data := panggil(t, app, "/api/v1/tema/pendidikan/kota/31.72").respons.Data.(map[string]any)
	sekolah := data["sekolah"].([]any)
	if len(sekolah) != 1 || sekolah[0].(map[string]any)["negeri"].(map[string]any)["nilai"] != 1.0 {
		t.Fatalf("sekolah ganda: %+v", sekolah)
	}
	lain := panggil(t, app, "/api/v1/tema/pendidikan/kota/31.01").respons.Data.(map[string]any)
	if len(lain["sekolah"].([]any)) != 0 {
		t.Fatal("sekolah dari kota lain ikut dihitung")
	}
}

func TestPendidikanKota_ValidasiWilayahDanTema(t *testing.T) {
	app, _ := appUjiKota(t)
	for _, uji := range []struct {
		path   string
		status int
	}{
		{"/api/v1/tema/pendidikan/kota/31.72.06", 400},
		{"/api/v1/tema/pendidikan/kota/31.99", 404},
		{"/api/v1/tema/tidak-ada/kota/31.72", 404},
	} {
		if h := panggil(t, app, uji.path); h.kode != uji.status {
			t.Errorf("%s: status %d, ingin %d", uji.path, h.kode, uji.status)
		}
	}
}

func TestPendidikanKota_DaftarKotaTetapAdaTanpaIndikatorJenjang(t *testing.T) {
	app, _ := appUjiKota(t)
	hasil := panggil(t, app, "/api/v1/tema/pendidikan_sma_smk/kota/31.01")
	data := hasil.respons.Data.(map[string]any)
	if hasil.kode != 200 || len(data["kabupaten_kota"].([]any)) != 6 || data["ringkasan"].(map[string]any)["rasio"] != nil {
		t.Fatalf("daftar kota harus tersedia tanpa indikator: %+v", data)
	}
}
