package http

import (
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestKelurahanDiTitik_KoordinatBukanAngka_400(t *testing.T) {
	hasil := panggil(t, appUji(), "/api/v1/wilayah/titik?lat=utara&lng=106.9")

	if hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode = %d, harap 400", hasil.kode)
	}
}
