package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

type hasilPanggilan struct {
	kode    int
	header  func(string) string
	respons Respons
}

func appUji() *fiber.App {
	return NewApp(config.Server{}, Layanan{})
}

func panggil(t *testing.T, app *fiber.App, path string) hasilPanggilan {
	t.Helper()
	return kirim(t, app, httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, path, nil))
}

func kirim(t *testing.T, app *fiber.App, permintaan *http.Request) hasilPanggilan {
	t.Helper()
	resp, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	var respons Respons
	if err := json.NewDecoder(resp.Body).Decode(&respons); err != nil {
		t.Fatalf("respons bukan amplop JSON: %v", err)
	}
	return hasilPanggilan{kode: resp.StatusCode, header: resp.Header.Get, respons: respons}
}
