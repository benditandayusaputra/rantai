package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

func appUjiSumberData(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	_, err := tx.Exec(t.Context(), `DELETE FROM ingest_run;
INSERT INTO ingest_run (sumber, mulai, selesai, status) VALUES
    ('uji_dapodik', '2026-10-01T01:00:00Z', '2026-10-01T02:00:00Z', 'berhasil'),
    ('uji_dapodik', '2026-10-20T01:00:00Z', '2026-10-20T01:05:00Z', 'gagal'),
    ('uji_kpu', '2026-10-22T03:00:00Z', NULL, 'berjalan')`)
	if err != nil {
		t.Fatal(err)
	}
	layanan := usecase.BacaSumberData{
		Pembaca:            repository.IngestRunRepository{DB: tx},
		Media:              []domain.MediaTerverifikasi{{Nama: "ANTARA", Domain: "antaranews.com"}},
		AtribusiOpentender: "Data Opentender oleh ICW, lisensi ODbL",
	}
	return NewApp(config.Server{}, Layanan{SumberData: layanan})
}

func TestSumberData_PengambilanTerakhir_TanggalBerhasilDariIngestRun(t *testing.T) {
	hasil := panggil(t, appUjiSumberData(t), "/api/v1/sumber-data")

	data, _ := hasil.respons.Data.(map[string]any)
	pengambilan, _ := data["pengambilan"].([]any)
	if hasil.kode != fiber.StatusOK || len(pengambilan) != 2 {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	dapodik, _ := pengambilan[0].(map[string]any)
	if dapodik["job"] != "uji_dapodik" || dapodik["status_terakhir"] != "gagal" || dapodik["berhasil_pada"] != "2026-10-01T02:00:00Z" {
		t.Errorf("pengambilan dapodik %v", dapodik)
	}
	kpu, _ := pengambilan[1].(map[string]any)
	if kpu["berhasil_pada"] != nil || kpu["status_terakhir"] != "berjalan" {
		t.Errorf("pengambilan kpu belum pernah berhasil: %v", kpu)
	}
}

func TestSumberData_MediaDanAtribusi_DariKonfigurasi(t *testing.T) {
	hasil := panggil(t, appUjiSumberData(t), "/api/v1/sumber-data")

	data, _ := hasil.respons.Data.(map[string]any)
	media, _ := data["media"].([]any)
	if len(media) != 1 || data["atribusi_opentender"] != "Data Opentender oleh ICW, lisensi ODbL" {
		t.Errorf("media %v, atribusi %v", media, data["atribusi_opentender"])
	}
}
