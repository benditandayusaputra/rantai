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
	Ringkasan  LayananRingkasanKecamatan
	Paket      LayananPaket
	Belanja    LayananBelanja
	Pembanding LayananPembanding
	Penyedia   LayananPenyedia
	Indikator  LayananIndikator
	Sekolah    LayananSekolahTerdekat
	SumberData LayananSumberData
	Tanya      LayananTanya
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
	api.Get("/kecamatan/:kode/ringkasan", handlerRingkasanKecamatan{layanan: layanan.Ringkasan}.ringkas)
	daftarkanPaket(api, handlerPaket{layanan: layanan.Paket})
	api.Get("/tema/:kode/belanja", handlerBelanja{layanan: layanan.Belanja}.hitung)
	api.Get("/paket/:kode/pembanding", handlerPembanding{layanan: layanan.Pembanding}.bandingkan)
	api.Get("/penyedia/:id", handlerPenyedia{layanan: layanan.Penyedia}.ambil)
	daftarkanPendidikan(api, handlerPendidikan{layanan: layanan.Indikator})
	api.Get("/sumber-data", handlerSumberData{layanan: layanan.SumberData}.baca)
	api.Get("/sekolah/terdekat", handlerSekolahTerdekat{layanan: layanan.Sekolah}.cari)
	daftarkanTanya(api, layanan.Tanya, cfg.Tanya)
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

func daftarkanPendidikan(api fiber.Router, h handlerPendidikan) {
	api.Get("/peta/tema/:kode", h.petaTema)
	api.Get("/tema/pendidikan/kecamatan/:kode", h.detailKecamatan)
	api.Get("/tema/pendidikan/kecamatan/:kode/sekolah", h.sekolahKecamatan)
}
