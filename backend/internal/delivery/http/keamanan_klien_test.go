package http

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

func TestKunciKlien(t *testing.T) {
	kasus := []struct {
		nama, rahasiaServer, rahasiaHeader, ipAsli, harap string
	}{
		{"rahasia cocok memakai IP pengunjung", "rahasia-uji", "rahasia-uji", "203.0.113.7", "203.0.113.7"},
		{"rahasia salah memakai IP koneksi", "rahasia-uji", "tebakan", "203.0.113.7", "0.0.0.0"},
		{"tanpa rahasia server mengabaikan header", "", "", "203.0.113.7", "0.0.0.0"},
		{"rahasia cocok tanpa IP memakai IP koneksi", "rahasia-uji", "rahasia-uji", "", "0.0.0.0"},
	}
	app := fiber.New()
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)
			ctx.Request().Header.Set(headerRahasiaProxy, k.rahasiaHeader)
			ctx.Request().Header.Set(headerIPKlienAsli, k.ipAsli)

			if kunci := kunciKlien(k.rahasiaServer)(ctx); kunci != k.harap {
				t.Errorf("kunci %s, harap %s", kunci, k.harap)
			}
		})
	}
}
