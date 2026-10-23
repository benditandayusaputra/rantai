package http

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

const originFrontendUji = "http://localhost:5173"

func appUjiPengaman(maksPermintaan int) *fiber.App {
	return NewApp(config.Server{
		BatasLaju: config.BatasLaju{MaksPermintaan: maksPermintaan, Jendela: time.Minute},
		CORS:      config.CORS{Origin: []string{originFrontendUji}},
	}, Layanan{})
}

func TestBatasLaju_MelebihiMaksimum_429DenganRetryAfter(t *testing.T) {
	app := appUjiPengaman(2)

	panggil(t, app, "/api/v1/kesehatan")
	panggil(t, app, "/api/v1/kesehatan")
	hasil := panggil(t, app, "/api/v1/kesehatan")

	if hasil.kode != fiber.StatusTooManyRequests || hasil.respons.Code != fiber.StatusTooManyRequests {
		t.Fatalf("kode = %d, amplop = %d", hasil.kode, hasil.respons.Code)
	}
	if hasil.header(fiber.HeaderRetryAfter) == "" {
		t.Error("header Retry-After tidak ada")
	}
}

func TestHeaderKeamanan_SetiapRespons_HeaderWajibAda(t *testing.T) {
	hasil := panggil(t, appUjiPengaman(10), "/api/v1/kesehatan")

	harap := map[string]string{
		fiber.HeaderXContentTypeOptions:   "nosniff",
		fiber.HeaderReferrerPolicy:        kebijakanPerujuk,
		fiber.HeaderContentSecurityPolicy: kebijakanKontenAPI,
		fiber.HeaderPermissionsPolicy:     kebijakanIzinAPI,
	}
	for nama, nilai := range harap {
		if hasil.header(nama) != nilai {
			t.Errorf("header %s = %q, harap %q", nama, hasil.header(nama), nilai)
		}
	}
}

func TestCORS_Origin_HanyaOriginTerdaftarDiizinkan(t *testing.T) {
	kasus := map[string]string{originFrontendUji: originFrontendUji, "https://situs-lain.test": ""}
	for origin, harap := range kasus {
		t.Run(origin, func(t *testing.T) {
			permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/api/v1/kesehatan", nil)
			permintaan.Header.Set(fiber.HeaderOrigin, origin)

			hasil := kirim(t, appUjiPengaman(10), permintaan)

			if izin := hasil.header(fiber.HeaderAccessControlAllowOrigin); izin != harap {
				t.Errorf("Access-Control-Allow-Origin = %q, harap %q", izin, harap)
			}
		})
	}
}
