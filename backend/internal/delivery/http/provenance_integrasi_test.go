package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

func appUjiProvenance(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	_, err := tx.Exec(t.Context(), `INSERT INTO provenance (id, sumber, url, diambil_pada, label, catatan)
VALUES ('prv_uji_kursi', 'rantai', 'https://contoh.test', '2026-10-23T02:00:00Z', 'proksi', 'rombel kelas 7 dikali 32')`)
	if err != nil {
		t.Fatal(err)
	}
	layanan := usecase.AmbilProvenance{Pembaca: repository.ProvenanceRepository{DB: tx}}
	return NewApp(config.Server{}, Layanan{Provenance: layanan})
}

func TestAmbilProvenance_IDAda_SumberWaktuLabelDanRumus(t *testing.T) {
	hasil := panggil(t, appUjiProvenance(t), "/api/v1/provenance/prv_uji_kursi")

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusOK || data["label"] != "proksi" || data["sumber"] != "rantai" {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	if data["diambil_pada"] != "2026-10-23T02:00:00Z" || data["catatan"] != "rombel kelas 7 dikali 32" {
		t.Errorf("waktu atau rumus tidak sesuai: %v", data)
	}
}

func TestAmbilProvenance_IDTidakAda_404(t *testing.T) {
	hasil := panggil(t, appUjiProvenance(t), "/api/v1/provenance/prv_tidak_ada")

	if hasil.kode != fiber.StatusNotFound || hasil.respons.Message != "sumber data tidak ditemukan" {
		t.Errorf("kode %d, pesan %q", hasil.kode, hasil.respons.Message)
	}
}
