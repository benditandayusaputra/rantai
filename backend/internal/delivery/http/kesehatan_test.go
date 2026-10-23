package http

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestPeriksaKesehatan_ServerHidup_MengembalikanStatusSehat(t *testing.T) {
	hasil := panggil(t, appUji(), "/api/v1/kesehatan")

	if hasil.kode != fiber.StatusOK || hasil.respons.Code != fiber.StatusOK {
		t.Fatalf("kode = %d, amplop = %d", hasil.kode, hasil.respons.Code)
	}
	data, _ := hasil.respons.Data.(map[string]any)
	if data["status"] != statusSehat || hasil.respons.Message != pesanOK {
		t.Errorf("respons tidak sesuai: %+v", hasil.respons)
	}
}
