package http

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananParipurna interface {
	SikapFraksi(ctx context.Context, fraksi string) ([]domain.SikapFraksi, error)
	AspirasiReses(ctx context.Context, filter domain.FilterAspirasiReses) ([]domain.AspirasiReses, error)
	PaketTerkait(ctx context.Context, idAspirasi int64) (domain.HasilPaketTerkait, error)
}

type sikapFraksiJSON struct {
	Fraksi       string `json:"fraksi"`
	Agenda       string `json:"agenda"`
	JudulAgenda  string `json:"judul_agenda"`
	Tanggal      string `json:"tanggal"`
	Ringkasan    string `json:"ringkasan"`
	URL          string `json:"url"`
	ProvenanceID string `json:"provenance_id"`
}

type aspirasiResesJSON struct {
	ID               int64                 `json:"id"`
	MasaReses        string                `json:"masa_reses"`
	TanggalParipurna string                `json:"tanggal_paripurna"`
	Bidang           string                `json:"bidang"`
	Uraian           string                `json:"uraian"`
	Wilayah          *ringkasanWilayahJSON `json:"wilayah"`
	URL              string                `json:"url"`
	ProvenanceID     string                `json:"provenance_id"`
}

type paketTerkaitJSON struct {
	paketRingkasJSON
	DasarCocok string `json:"dasar_cocok"`
}

type hasilPaketTerkaitJSON struct {
	Aspirasi      aspirasiResesJSON  `json:"aspirasi"`
	TahunAnggaran int                `json:"tahun_anggaran"`
	Status        string             `json:"status"`
	Paket         []paketTerkaitJSON `json:"paket"`
}

type handlerParipurna struct {
	layanan LayananParipurna
}

func daftarkanParipurna(api fiber.Router, h handlerParipurna) {
	api.Get("/fraksi/:nama/sikap", h.sikapFraksi)
	api.Get("/aspirasi-reses", h.aspirasiReses)
	api.Get("/aspirasi-reses/:id/paket-terkait", h.paketTerkait)
}

func (h handlerParipurna) sikapFraksi(c fiber.Ctx) error {
	fraksi, err := url.PathUnescape(c.Params("nama"))
	if err != nil {
		return domain.InputTidakValid("nama fraksi tidak sah")
	}
	daftar, err := h.layanan.SikapFraksi(c.Context(), fraksi)
	if err != nil {
		return err
	}
	hasil := make([]sikapFraksiJSON, len(daftar))
	for i, s := range daftar {
		hasil[i] = sikapFraksiJSON{
			Fraksi: s.Fraksi, Agenda: string(s.Agenda), JudulAgenda: s.JudulAgenda, Tanggal: s.Tanggal.Format(time.DateOnly),
			Ringkasan: s.Ringkasan, URL: s.URL, ProvenanceID: s.ProvenanceID,
		}
	}
	return kirimData(c, hasil)
}

func (h handlerParipurna) aspirasiReses(c fiber.Ctx) error {
	dapil, err := uraiDapilOpsional(c.Query("dapil"))
	if err != nil {
		return err
	}
	daftar, err := h.layanan.AspirasiReses(c.Context(), domain.FilterAspirasiReses{DapilID: dapil, Bidang: c.Query("bidang")})
	if err != nil {
		return err
	}
	hasil := make([]aspirasiResesJSON, len(daftar))
	for i, a := range daftar {
		hasil[i] = keAspirasiResesJSON(a)
	}
	return kirimData(c, hasil)
}

func uraiDapilOpsional(nilai string) (int, error) {
	if nilai == "" {
		return 0, nil
	}
	dapil, err := strconv.Atoi(nilai)
	if err != nil || dapil < 1 {
		return 0, domain.InputTidakValid("dapil harus id dapil berupa angka 1 atau lebih")
	}
	return dapil, nil
}

func keAspirasiResesJSON(a domain.AspirasiReses) aspirasiResesJSON {
	hasil := aspirasiResesJSON{
		ID: a.ID, MasaReses: a.MasaReses, TanggalParipurna: a.TanggalParipurna.Format(time.DateOnly), Bidang: a.Bidang,
		Uraian: a.Uraian, URL: a.URL, ProvenanceID: a.ProvenanceID,
	}
	if a.KodeWilayah != "" {
		wilayah := keRingkasanWilayahJSON(domain.RingkasanWilayah{Kode: a.KodeWilayah, Nama: a.NamaWilayah, Tingkat: a.TingkatWilayah})
		hasil.Wilayah = &wilayah
	}
	return hasil
}

func (h handlerParipurna) paketTerkait(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return domain.InputTidakValid("id aspirasi reses wajib berupa angka")
	}
	hasil, err := h.layanan.PaketTerkait(c.Context(), id)
	if err != nil {
		return err
	}
	paket := make([]paketTerkaitJSON, len(hasil.Paket))
	for i, p := range hasil.Paket {
		paket[i] = paketTerkaitJSON{paketRingkasJSON: kePaketRingkasJSON(p.PaketRingkas), DasarCocok: dasarCocokPaket(p)}
	}
	return kirimData(c, hasilPaketTerkaitJSON{
		Aspirasi: keAspirasiResesJSON(hasil.Aspirasi), TahunAnggaran: hasil.TahunAnggaran, Status: string(hasil.Status()), Paket: paket,
	})
}

func dasarCocokPaket(p domain.PaketTerkait) string {
	if p.IsCocokKategori {
		return "kategori"
	}
	return "kata_kunci"
}
