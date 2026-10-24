package http

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const (
	pathPDFAspirasi         = "/aspirasi/pdf"
	awalanBerkasAspirasi    = "aspirasi-"
	akhiranPDF              = ".pdf"
	formatTanggalVerifikasi = "2006-01-02"
)

type LayananAspirasi interface {
	Susun(ctx context.Context, p usecase.PermintaanAspirasi) (domain.Dokumen, error)
}

type LayananJalurAspirasi interface {
	Untuk(ctx context.Context, kodeWilayah string) (usecase.JalurWilayah, error)
}

type permintaanAspirasiJSON struct {
	KodeKecamatan string `json:"kode_kecamatan"`
	Tema          string `json:"tema"`
	Masalah       string `json:"masalah"`
	Bahasa        string `json:"bahasa"`
}

type dokumenJSON struct {
	Teks   string `json:"teks"`
	PDFURL string `json:"pdf_url"`
}

type jalurAspirasiJSON struct {
	Kode       string      `json:"kode"`
	Nama       domain.Teks `json:"nama"`
	Peran      domain.Teks `json:"peran"`
	Cara       domain.Teks `json:"cara"`
	DasarHukum string      `json:"dasar_hukum"`
	URL        string      `json:"url"`
}

type jalurWilayahJSON struct {
	Wilayah          ringkasanWilayahJSON `json:"wilayah"`
	DiverifikasiPada string               `json:"diverifikasi_pada"`
	Jalur            []jalurAspirasiJSON  `json:"jalur"`
}

type handlerAspirasi struct {
	layanan LayananAspirasi
}

func (h handlerAspirasi) susun(c fiber.Ctx) error {
	var masukan permintaanAspirasiJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi kode_kecamatan, tema, dan masalah")
	}
	permintaan := kePermintaanAspirasi(masukan)
	dokumen, err := h.layanan.Susun(c.Context(), permintaan)
	if err != nil {
		return err
	}
	return kirimData(c, dokumenJSON{Teks: dokumen.Teks(), PDFURL: urlPDFAspirasi(permintaan)})
}

func (h handlerAspirasi) pdf(c fiber.Ctx) error {
	permintaan := kePermintaanAspirasi(permintaanAspirasiJSON{
		KodeKecamatan: c.Query("kode_kecamatan"), Tema: c.Query("tema"), Masalah: c.Query("masalah"), Bahasa: c.Query("bahasa"),
	})
	dokumen, err := h.layanan.Susun(c.Context(), permintaan)
	if err != nil {
		return err
	}
	return kirimPDF(c, awalanBerkasAspirasi+permintaan.KodeKecamatan+akhiranPDF, dokumen)
}

func kePermintaanAspirasi(m permintaanAspirasiJSON) usecase.PermintaanAspirasi {
	if m.Bahasa == "" {
		m.Bahasa = domain.BahasaIndonesia
	}
	return usecase.PermintaanAspirasi(m)
}

func urlPDFAspirasi(p usecase.PermintaanAspirasi) string {
	kueri := url.Values{"kode_kecamatan": {p.KodeKecamatan}, "tema": {p.Tema}, "masalah": {p.Masalah}, "bahasa": {p.Bahasa}}
	return prefiksAPI + pathPDFAspirasi + "?" + kueri.Encode()
}

type handlerJalurAspirasi struct {
	layanan LayananJalurAspirasi
}

func (h handlerJalurAspirasi) untuk(c fiber.Ctx) error {
	hasil, err := h.layanan.Untuk(c.Context(), c.Query("kode_wilayah"))
	if err != nil {
		return err
	}
	jalur := make([]jalurAspirasiJSON, len(hasil.Jalur))
	for i, j := range hasil.Jalur {
		jalur[i] = jalurAspirasiJSON(j)
	}
	return kirimData(c, jalurWilayahJSON{
		Wilayah:          keRingkasanWilayahJSON(hasil.Wilayah.RingkasanWilayah),
		DiverifikasiPada: hasil.DiverifikasiPada.Format(formatTanggalVerifikasi),
		Jalur:            jalur,
	})
}

func daftarkanAspirasi(api fiber.Router, layanan LayananAspirasi) {
	h := handlerAspirasi{layanan: layanan}
	api.Post("/aspirasi", h.susun)
	api.Get(pathPDFAspirasi, h.pdf)
}
