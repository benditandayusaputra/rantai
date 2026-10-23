package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const sqlSekolahTerdekatUji = `WITH prv AS (SELECT provenance_id AS id FROM wilayah WHERE kode = '31.72.06'),
data (npsn, bentuk, jenjang, status, lat, kursi) AS (VALUES
    ('90000301', 'SMP', 'menengah_pertama', 'negeri', -6.16452, 288),
    ('90000302', 'SMP', 'menengah_pertama', 'negeri', -6.17356, 252),
    ('90000303', 'SMP', 'menengah_pertama', 'negeri', -6.19617, 216),
    ('90000304', 'SD', 'dasar', 'negeri', -6.16100, NULL)
),
sekolah_baru AS (
    INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, lokasi, provenance_id)
    SELECT npsn, 'SEKOLAH UJI ' || npsn, bentuk, jenjang, status, false, '31.72.06',
        ST_SetSRID(ST_MakePoint(106.91, lat), 4326), prv.id
    FROM data, prv
    RETURNING npsn
),
kursi_baru AS (
    INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
    SELECT d.npsn, '9999/10000', 'negeri', d.kursi, prv.id FROM data d, prv, sekolah_baru s WHERE s.npsn = d.npsn AND d.kursi IS NOT NULL
)
INSERT INTO sekolah_rombel (npsn, tahun_ajaran, tingkat, jumlah_rombel, jumlah_siswa, provenance_id)
SELECT s.npsn, '9999/10000', 6, 4, 120, prv.id FROM sekolah_baru s, prv WHERE s.npsn = '90000304'`

func appUjiSekolahTerdekat(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	siapkanSekolahTerdekat(t, tx)
	var cfg config.Config
	if err := config.BacaYAML("../../../config/api.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{Sekolah: usecase.CariSekolahTerdekat{
		Pembaca:     repository.SekolahTerdekatRepository{DB: tx},
		Konfigurasi: usecase.KonfigurasiSekolahTerdekat(cfg.SekolahTerdekat),
	}})
}

func siapkanSekolahTerdekat(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	for _, sql := range []string{`DELETE FROM daya_tampung`, `DELETE FROM sekolah_rombel`, sqlSekolahTerdekatUji} {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSekolahTerdekat_TigaSMPJarakBerbeda_DuaDalamRadiusUrutTerdekat(t *testing.T) {
	hasil := panggil(t, appUjiSekolahTerdekat(t), "/api/v1/sekolah/terdekat?lat=-6.16&lng=106.91&radius_m=3000")

	data, _ := hasil.respons.Data.(map[string]any)
	negeri, _ := data["smp_negeri"].([]any)
	if hasil.kode != fiber.StatusOK || len(negeri) != 2 {
		t.Fatalf("kode %d, smp negeri %v", hasil.kode, data["smp_negeri"])
	}
	pertama, kedua := negeri[0].(map[string]any), negeri[1].(map[string]any)
	if pertama["npsn"] != "90000301" || kedua["npsn"] != "90000302" || pertama["jarak_m"].(float64) >= kedua["jarak_m"].(float64) {
		t.Errorf("urutan = %v, %v", pertama, kedua)
	}
}

func TestSekolahTerdekat_TekananRadius_SiswaKelas6DibagiKursiNegeri(t *testing.T) {
	hasil := panggil(t, appUjiSekolahTerdekat(t), "/api/v1/sekolah/terdekat?lat=-6.16&lng=106.91&radius_m=3000")

	data, _ := hasil.respons.Data.(map[string]any)
	tekanan, _ := data["tekanan_radius"].(map[string]any)
	siswa, _ := tekanan["siswa_kelas_6"].(map[string]any)
	kursi, _ := tekanan["kursi_negeri"].(map[string]any)
	rasio, _ := tekanan["rasio"].(map[string]any)
	if siswa["nilai"] != 120.0 || kursi["nilai"] != 540.0 || rasio["label"] != "resmi" {
		t.Errorf("tekanan = %v", tekanan)
	}
}

func TestSekolahTerdekat_RadiusMelebihiBatas_400(t *testing.T) {
	hasil := panggil(t, appUjiSekolahTerdekat(t), "/api/v1/sekolah/terdekat?lat=-6.16&lng=106.91&radius_m=50000")

	if hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode = %d", hasil.kode)
	}
}
