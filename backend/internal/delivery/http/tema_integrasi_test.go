package http

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

func muatTema(t *testing.T, tx pgx.Tx) {
	t.Helper()
	impor := usecase.ImporTema{
		Pengambil: tema.Sumber{Dir: "../../../config/tema", URLRujukan: "https://contoh.test", Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func appUjiTema(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	muatTema(t, tx)
	return NewApp(config.Server{}, Layanan{Tema: usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}}})
}

func TestDaftarTema_TemaPendidikanAktif_TampilDenganProvenance(t *testing.T) {
	hasil := panggil(t, appUjiTema(t), "/api/v1/tema")

	daftar, _ := hasil.respons.Data.([]any)
	if hasil.kode != fiber.StatusOK || len(daftar) == 0 {
		t.Fatalf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
	pertama, _ := daftar[0].(map[string]any)
	if pertama["kode"] != "pendidikan" || pertama["provenance_id"] == "" || pertama["indikator"] == nil {
		t.Errorf("tema pertama = %v", pertama)
	}
}

func TestAmbilTema_Pendidikan_KontrakLengkap(t *testing.T) {
	hasil := panggil(t, appUjiTema(t), "/api/v1/tema/pendidikan")

	data, _ := hasil.respons.Data.(map[string]any)
	tanggungJawab, _ := data["tanggung_jawab"].([]any)
	if hasil.kode != fiber.StatusOK || len(tanggungJawab) != 1 || data["templat_aspirasi"] == nil {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestAmbilTema_KodeTidakAda_404(t *testing.T) {
	hasil := panggil(t, appUjiTema(t), "/api/v1/tema/antariksa")

	if hasil.kode != fiber.StatusNotFound || hasil.respons.Message != "tema tidak ditemukan" {
		t.Errorf("kode %d, pesan %q", hasil.kode, hasil.respons.Message)
	}
}
