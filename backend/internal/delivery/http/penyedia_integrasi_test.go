package http

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/opentender"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

type pengunduhOpentenderUji struct{}

func (pengunduhOpentenderUji) Ambil(context.Context, string) ([]byte, error) {
	return os.ReadFile("../../ingest/opentender/testdata/sampel_2023.jsonl.gz")
}

func imporOpentenderUji(t *testing.T, tx pgx.Tx) opentender.Konfigurasi {
	t.Helper()
	var cfg opentender.Konfigurasi
	if err := config.BacaYAML("../../../config/opentender.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	impor := usecase.ImporRiwayatPenyedia{
		Pengambil: opentender.Sumber{Konfigurasi: cfg, Klien: pengunduhOpentenderUji{}, DirRaw: t.TempDir(), Sekarang: time.Now},
		Penyimpan: repository.RiwayatPenyediaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		Tahun:     []int{2023},
	}
	for range 2 {
		if err := impor.Jalankan(t.Context()); err != nil {
			t.Fatal(err)
		}
	}
	return cfg
}

func TestProfilPenyedia_RiwayatOpentender_BeratribusiODbLTanpaSkor(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	cfg := imporOpentenderUji(t, tx)
	var id int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM penyedia WHERE nama_normal = 'nindya karya'`).Scan(&id); err != nil {
		t.Fatal(err)
	}
	var jumlahRiwayat int
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM riwayat_penyedia`).Scan(&jumlahRiwayat); err != nil || jumlahRiwayat != 4 {
		t.Fatalf("riwayat = %d setelah dua kali impor, err %v", jumlahRiwayat, err)
	}
	u := usecase.BacaPenyedia{Pembaca: repository.PenyediaRepository{DB: tx}, Batas: 100, Atribusi: cfg.Atribusi}
	app := NewApp(config.Server{}, Layanan{Penyedia: u})

	hasil := panggil(t, app, "/api/v1/penyedia/"+strconv.FormatInt(id, 10))

	data, _ := hasil.respons.Data.(map[string]any)
	mentah, _ := json.Marshal(data)
	if hasil.kode != fiber.StatusOK || !strings.Contains(data["atribusi"].(string), "ODbL") || polaKunciKinerja.Match(mentah) {
		t.Errorf("kode %d, data %s", hasil.kode, mentah)
	}
	riwayat, _ := data["riwayat_opentender"].([]any)
	if len(riwayat) != 1 || riwayat[0].(map[string]any)["pembeli"] != "DINAS SUMBER DAYA AIR" {
		t.Errorf("riwayat = %v", riwayat)
	}
	if salah := panggil(t, app, "/api/v1/penyedia/abc"); salah.kode != fiber.StatusBadRequest {
		t.Errorf("id bukan angka: %d", salah.kode)
	}
}
