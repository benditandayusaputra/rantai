package http

import (
	"context"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const ukuranBitFloat = 64

type LayananWilayah interface {
	Cari(ctx context.Context, kueri string) ([]domain.WilayahDitemukan, error)
	DiTitik(ctx context.Context, titik domain.Titik) (domain.WilayahDitemukan, error)
	Ambil(ctx context.Context, kode string) (domain.DetailWilayah, error)
}

type ringkasanWilayahJSON struct {
	Kode    string `json:"kode"`
	Nama    string `json:"nama"`
	Tingkat string `json:"tingkat"`
}

type wilayahDitemukanJSON struct {
	ringkasanWilayahJSON
	Hierarki []ringkasanWilayahJSON `json:"hierarki"`
}

type titikJSON struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type detailWilayahJSON struct {
	wilayahDitemukanJSON
	Alias        []string  `json:"alias"`
	TitikTengah  titikJSON `json:"titik_tengah"`
	ProvenanceID string    `json:"provenance_id"`
}

type handlerWilayah struct {
	layanan LayananWilayah
}

func (h handlerWilayah) cari(c fiber.Ctx) error {
	hasil, err := h.layanan.Cari(c.Context(), c.Query("q"))
	if err != nil {
		return err
	}
	daftar := make([]wilayahDitemukanJSON, len(hasil))
	for i, w := range hasil {
		daftar[i] = keWilayahDitemukanJSON(w)
	}
	return kirimData(c, daftar)
}

func (h handlerWilayah) diTitik(c fiber.Ctx) error {
	titik, err := uraiTitik(c.Query("lat"), c.Query("lng"))
	if err != nil {
		return err
	}
	wilayah, err := h.layanan.DiTitik(c.Context(), titik)
	if err != nil {
		return err
	}
	return kirimData(c, keWilayahDitemukanJSON(wilayah))
}

func (h handlerWilayah) ambil(c fiber.Ctx) error {
	wilayah, err := h.layanan.Ambil(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, detailWilayahJSON{
		wilayahDitemukanJSON: keWilayahDitemukanJSON(wilayah.WilayahDitemukan),
		Alias:                wilayah.Alias,
		TitikTengah:          titikJSON(wilayah.TitikTengah),
		ProvenanceID:         wilayah.ProvenanceID,
	})
}

func uraiTitik(lat, lng string) (domain.Titik, error) {
	lintang, errLintang := strconv.ParseFloat(lat, ukuranBitFloat)
	bujur, errBujur := strconv.ParseFloat(lng, ukuranBitFloat)
	if errLintang != nil || errBujur != nil {
		return domain.Titik{}, domain.InputTidakValid("lat dan lng wajib berupa angka")
	}
	return domain.Titik{Lat: lintang, Lng: bujur}, nil
}

func keWilayahDitemukanJSON(w domain.WilayahDitemukan) wilayahDitemukanJSON {
	hierarki := make([]ringkasanWilayahJSON, len(w.Hierarki))
	for i, leluhur := range w.Hierarki {
		hierarki[i] = keRingkasanWilayahJSON(leluhur)
	}
	return wilayahDitemukanJSON{ringkasanWilayahJSON: keRingkasanWilayahJSON(w.RingkasanWilayah), Hierarki: hierarki}
}

func keRingkasanWilayahJSON(w domain.RingkasanWilayah) ringkasanWilayahJSON {
	return ringkasanWilayahJSON{Kode: w.Kode, Nama: w.Nama, Tingkat: string(w.Tingkat)}
}
