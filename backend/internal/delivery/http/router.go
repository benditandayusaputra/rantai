package http

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

const prefiksAPI = "/api/v1"

type Layanan struct {
	Wilayah    LayananWilayah
	Provenance LayananProvenance
	Tema       LayananTema
	Wakil      LayananWakil
	Umpan      LayananUmpan
}

func NewApp(cfg config.Server, layanan Layanan) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: tanganiError,
		ReadTimeout:  cfg.BatasWaktuBaca,
		WriteTimeout: cfg.BatasWaktuTulis,
	})
	pasangPengaman(app, cfg)
	api := app.Group(prefiksAPI)
	api.Get("/kesehatan", periksaKesehatan)
	daftarkanWilayah(api, handlerWilayah{layanan: layanan.Wilayah})
	api.Get("/provenance/:id", handlerProvenance{layanan: layanan.Provenance}.ambil)
	daftarkanTema(api, handlerTema{layanan: layanan.Tema})
	daftarkanWakil(api, handlerWakil{layanan: layanan.Wakil})
	api.Get("/anggota/:id/umpan", handlerUmpan{layanan: layanan.Umpan}.umpanAnggota)
	return app
}

func daftarkanWilayah(api fiber.Router, h handlerWilayah) {
	api.Get("/wilayah/cari", h.cari)
	api.Get("/wilayah/titik", h.diTitik)
	api.Get("/wilayah/:kode", h.ambil)
}

func daftarkanTema(api fiber.Router, h handlerTema) {
	api.Get("/tema", h.daftar)
	api.Get("/tema/:kode", h.ambil)
}

func daftarkanWakil(api fiber.Router, h handlerWakil) {
	api.Get("/kecamatan/:kode/wakil", h.wakilKecamatan)
	api.Get("/dapil/:id", h.dapil)
	api.Get("/anggota/:id", h.anggota)
}
