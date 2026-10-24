package http

import (
	"encoding/json"
	"net/url"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const awalanBerkasPermohonan = "permohonan-"

type LayananDataBelumDibuka interface {
	Daftar() domain.KonfigurasiDataBelumDibuka
	Permohonan(kode, bahasa string) (domain.Dokumen, error)
}

type ppidJSON struct {
	Nama   string `json:"nama"`
	Alamat string `json:"alamat"`
	URL    string `json:"url"`
}

type dataBelumDibukaJSON struct {
	Kode           string      `json:"kode"`
	Nama           domain.Teks `json:"nama"`
	Alasan         domain.Teks `json:"alasan"`
	Pengganti      domain.Teks `json:"pengganti"`
	FiturPengganti string      `json:"fitur_pengganti"`
}

type daftarDataBelumDibukaJSON struct {
	PPID ppidJSON              `json:"ppid"`
	Data []dataBelumDibukaJSON `json:"data"`
}

type permintaanPermohonanJSON struct {
	Bahasa string `json:"bahasa"`
}

type handlerDataBelumDibuka struct {
	layanan LayananDataBelumDibuka
}

func daftarkanDataBelumDibuka(api fiber.Router, layanan LayananDataBelumDibuka) {
	h := handlerDataBelumDibuka{layanan: layanan}
	api.Get("/data-belum-dibuka", h.daftar)
	api.Post("/data-belum-dibuka/:kode/permohonan", h.permohonan)
	api.Get("/data-belum-dibuka/:kode/permohonan/pdf", h.pdf)
}

func (h handlerDataBelumDibuka) daftar(c fiber.Ctx) error {
	k := h.layanan.Daftar()
	data := make([]dataBelumDibukaJSON, len(k.Data))
	for i, d := range k.Data {
		data[i] = dataBelumDibukaJSON{Kode: d.Kode, Nama: d.Nama, Alasan: d.Alasan, Pengganti: d.Pengganti, FiturPengganti: d.FiturPengganti}
	}
	return kirimData(c, daftarDataBelumDibukaJSON{PPID: ppidJSON(k.PPID), Data: data})
}

func (h handlerDataBelumDibuka) permohonan(c fiber.Ctx) error {
	var masukan permintaanPermohonanJSON
	if badan := c.Body(); len(badan) > 0 {
		if err := json.Unmarshal(badan, &masukan); err != nil {
			return domain.InputTidakValid("badan permintaan wajib JSON berisi bahasa")
		}
	}
	bahasa := bahasaAtauBawaan(masukan.Bahasa)
	dokumen, err := h.layanan.Permohonan(c.Params("kode"), bahasa)
	if err != nil {
		return err
	}
	pdfURL := prefiksAPI + "/data-belum-dibuka/" + url.PathEscape(c.Params("kode")) + "/permohonan/pdf?" + url.Values{"bahasa": {bahasa}}.Encode()
	return kirimData(c, dokumenJSON{Teks: dokumen.Teks(), PDFURL: pdfURL})
}

func (h handlerDataBelumDibuka) pdf(c fiber.Ctx) error {
	dokumen, err := h.layanan.Permohonan(c.Params("kode"), bahasaAtauBawaan(c.Query("bahasa")))
	if err != nil {
		return err
	}
	return kirimPDF(c, awalanBerkasPermohonan+c.Params("kode")+akhiranPDF, dokumen)
}

func bahasaAtauBawaan(bahasa string) string {
	if bahasa == "" {
		return domain.BahasaIndonesia
	}
	return bahasa
}
