package http

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

const prefiksAPI = "/api/v1"

type Layanan struct {
	Wilayah    LayananWilayah
	Provenance LayananProvenance
}

func NewApp(cfg config.Server, layanan Layanan) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: tanganiError,
		ReadTimeout:  cfg.BatasWaktuBaca,
		WriteTimeout: cfg.BatasWaktuTulis,
	})
	api := app.Group(prefiksAPI)
	api.Get("/kesehatan", periksaKesehatan)
	daftarkanWilayah(api, handlerWilayah{layanan: layanan.Wilayah})
	api.Get("/provenance/:id", handlerProvenance{layanan: layanan.Provenance}.ambil)
	return app
}

func daftarkanWilayah(api fiber.Router, h handlerWilayah) {
	api.Get("/wilayah/cari", h.cari)
	api.Get("/wilayah/titik", h.diTitik)
	api.Get("/wilayah/:kode", h.ambil)
}
