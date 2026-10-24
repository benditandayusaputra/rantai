package http

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
	"github.com/gofiber/fiber/v3"
)

type tanyaTiruan struct {
	errPeriksa error
}

func (t tanyaTiruan) Periksa(tanya.Permintaan) error {
	return t.errPeriksa
}

func (tanyaTiruan) Jawab(_ context.Context, _ tanya.Permintaan, kabar func(tanya.Kabar)) (tanya.Jawaban, error) {
	kabar(tanya.Kabar{Alat: "indikator_tema"})
	return tanya.Jawaban{Teks: "Kursi 288 [1].", Rujukan: []tanya.Rujukan{{Nomor: 1, ProvenanceID: "prv_1"}}}, nil
}

func appUjiTanya(layanan LayananTanya, maks int) *fiber.App {
	cfg := config.Server{Tanya: config.ServerTanya{BatasLaju: config.BatasLaju{MaksPermintaan: maks, Jendela: time.Hour}, BatasWaktu: time.Second}}
	cfg.BatasLaju = config.BatasLaju{MaksPermintaan: 100, Jendela: time.Minute}
	return NewApp(cfg, Layanan{Tanya: layanan})
}

func kirimTanya(t *testing.T, app *fiber.App, badan string) (int, string, string) {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPost, "/api/v1/tanya", strings.NewReader(badan))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	isi, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode, resp.Header.Get(fiber.HeaderContentType), string(isi)
}

func TestTanya_PertanyaanSah_AliranKabarAlatLaluJawaban(t *testing.T) {
	kode, tipe, isi := kirimTanya(t, appUjiTanya(tanyaTiruan{}, 20), `{"pertanyaan":"Kursi SMP cukup?","bahasa":"id"}`)

	harap := "event: alat\ndata: {\"alat\":\"indikator_tema\"}\n\n" +
		"event: jawaban\ndata: {\"teks\":\"Kursi 288 [1].\",\"rujukan\":[{\"nomor\":1,\"provenance_id\":\"prv_1\"}]}\n\n"
	if kode != fiber.StatusOK || tipe != tipeAliranPeristiwa || isi != harap {
		t.Errorf("kode %d, tipe %q, isi %q", kode, tipe, isi)
	}
}

func TestTanya_PermintaanDitolakPemeriksaan_400TanpaAliran(t *testing.T) {
	layanan := tanyaTiruan{errPeriksa: domain.InputTidakValid("pertanyaan wajib diisi")}

	kode, _, isi := kirimTanya(t, appUjiTanya(layanan, 20), `{"pertanyaan":""}`)

	if kode != fiber.StatusBadRequest || strings.Contains(isi, "event:") {
		t.Errorf("kode %d, isi %q", kode, isi)
	}
}

func TestTanya_MelewatiBatasLajuTanya_429(t *testing.T) {
	app := appUjiTanya(tanyaTiruan{}, 1)
	kirimTanya(t, app, `{"pertanyaan":"Pertama"}`)

	kode, _, _ := kirimTanya(t, app, `{"pertanyaan":"Kedua"}`)

	if kode != fiber.StatusTooManyRequests {
		t.Errorf("kode %d, harap 429", kode)
	}
}

func TestTanya_TanpaLayananTanya_RuteTidakDipasang(t *testing.T) {
	kode, _, _ := kirimTanya(t, appUji(), `{"pertanyaan":"Halo"}`)

	if kode != fiber.StatusNotFound && kode != fiber.StatusMethodNotAllowed {
		t.Errorf("kode %d, harap rute tidak ada", kode)
	}
}
