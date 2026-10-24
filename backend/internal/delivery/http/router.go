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
	Aspirasi   LayananAspirasi
	Jalur      LayananJalurAspirasi
	BelumBuka  LayananDataBelumDibuka
	Masuk      LayananMasuk
	Akses      PemeriksaAkses
	Paripurna  LayananParipurna
	Laporan    LayananLaporan

	BacaLaporan       LayananBacaLaporan
	KonfirmasiLaporan LayananKonfirmasiLaporan
	Moderasi          LayananModerasi
	Pantau            LayananPantau
	WebPush           LayananWebPush
	Panel             LayananPanel
}

func NewApp(cfg config.Server, layanan Layanan) *fiber.App {
	app := fiber.New(konfigurasiFiber(cfg))
	pasangPengaman(app, cfg)
	api := app.Group(prefiksAPI)
	api.Get("/kesehatan", periksaKesehatan)
	daftarkanWilayah(api, handlerWilayah{layanan: layanan.Wilayah})
	api.Get("/provenance/:id", handlerProvenance{layanan: layanan.Provenance}.ambil)
	daftarkanTema(api, handlerTema{layanan: layanan.Tema})
	daftarkanWakil(api, handlerWakil{layanan: layanan.Wakil})
	daftarkanUmpan(api, handlerUmpan{layanan: layanan.Umpan})
	api.Get("/kecamatan/:kode/ringkasan", handlerRingkasanKecamatan{layanan: layanan.Ringkasan}.ringkas)
	daftarkanPaket(api, handlerPaket{layanan: layanan.Paket})
	api.Get("/tema/:kode/belanja", handlerBelanja{layanan: layanan.Belanja}.hitung)
	api.Get("/paket/:kode/pembanding", handlerPembanding{layanan: layanan.Pembanding}.bandingkan)
	api.Get("/penyedia/:id", handlerPenyedia{layanan: layanan.Penyedia}.ambil)
	daftarkanPendidikan(api, handlerPendidikan{layanan: layanan.Indikator})
	api.Get("/sumber-data", handlerSumberData{layanan: layanan.SumberData}.baca)
	api.Get("/sekolah/terdekat", handlerSekolahTerdekat{layanan: layanan.Sekolah}.cari)
	daftarkanTanya(api, layanan.Tanya, cfg.Tanya)
	daftarkanAspirasi(api, layanan.Aspirasi)
	api.Get("/jalur-aspirasi", handlerJalurAspirasi{layanan: layanan.Jalur}.untuk)
	daftarkanDataBelumDibuka(api, layanan.BelumBuka)
	daftarkanAuth(api, layanan, cfg.Auth)
	daftarkanParipurna(api, handlerParipurna{layanan: layanan.Paripurna})
	daftarkanAksiWarga(api, layanan)
	return app
}

func konfigurasiFiber(cfg config.Server) fiber.Config {
	return fiber.Config{
		ErrorHandler:     tanganiError,
		ReadTimeout:      cfg.BatasWaktuBaca,
		WriteTimeout:     cfg.BatasWaktuTulis,
		BodyLimit:        BatasBodyPermintaan,
		TrustProxy:       len(cfg.ProxyTepercaya) > 0,
		TrustProxyConfig: fiber.TrustProxyConfig{Proxies: cfg.ProxyTepercaya},
		ProxyHeader:      fiber.HeaderXForwardedFor,
	}
}

func daftarkanAksiWarga(api fiber.Router, layanan Layanan) {
	daftarkanLaporan(api, layanan)
	daftarkanModerasi(api, layanan)
	daftarkanPantau(api, layanan)
	daftarkanWebPush(api, layanan)
	daftarkanPanel(api, layanan)
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

func daftarkanUmpan(api fiber.Router, h handlerUmpan) {
	api.Get("/anggota/:id/umpan", h.umpanAnggota)
	api.Get("/anggota/:id/isu", h.isuAnggota)
}

func daftarkanPendidikan(api fiber.Router, h handlerPendidikan) {
	api.Get("/peta/tema/:kode", h.petaTema)
	api.Get("/tema/pendidikan/kecamatan/:kode", h.detailKecamatan)
	api.Get("/tema/pendidikan/kecamatan/:kode/sekolah", h.sekolahKecamatan)
}
