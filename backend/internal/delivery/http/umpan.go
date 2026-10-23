package http

import (
	"context"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const halamanBawaan = 1

type LayananUmpan interface {
	UmpanAnggota(ctx context.Context, filter domain.FilterUmpan) (domain.HalamanUmpan, error)
}

type umpanJSON struct {
	ID           int64     `json:"id"`
	JenisSumber  string    `json:"jenis_sumber"`
	MediaNama    string    `json:"media_nama"`
	MediaDomain  string    `json:"media_domain"`
	Judul        string    `json:"judul"`
	Ringkasan    *string   `json:"ringkasan"`
	URL          string    `json:"url"`
	TerbitPada   time.Time `json:"terbit_pada"`
	TagIsu       []string  `json:"tag_isu"`
	DasarCocok   string    `json:"dasar_cocok"`
	ProvenanceID string    `json:"provenance_id"`
}

type handlerUmpan struct {
	layanan LayananUmpan
}

func (h handlerUmpan) umpanAnggota(c fiber.Ctx) error {
	halaman, err := uraiHalaman(c.Query("halaman"))
	if err != nil {
		return err
	}
	filter := domain.FilterUmpan{
		AnggotaID: c.Params("id"), Jenis: domain.JenisSumberUmpan(c.Query("jenis")), Tag: c.Query("tag"), Halaman: halaman,
	}
	hasil, err := h.layanan.UmpanAnggota(c.Context(), filter)
	if err != nil {
		return err
	}
	daftar := make([]umpanJSON, len(hasil.Item))
	for i, u := range hasil.Item {
		daftar[i] = keUmpanJSON(u)
	}
	return kirimHalaman(c, daftar, buatMeta(halaman, hasil))
}

func uraiHalaman(nilai string) (int, error) {
	if nilai == "" {
		return halamanBawaan, nil
	}
	halaman, err := strconv.Atoi(nilai)
	if err != nil {
		return 0, domain.InputTidakValid("halaman harus angka 1 atau lebih")
	}
	return halaman, nil
}

func buatMeta(halaman int, hasil domain.HalamanUmpan) Meta {
	return Meta{Page: halaman, PerPage: hasil.PerHalaman, Total: hasil.Total, TotalPages: (hasil.Total + hasil.PerHalaman - 1) / hasil.PerHalaman}
}

func keUmpanJSON(u domain.UmpanTampil) umpanJSON {
	var ringkasan *string
	if u.Ringkasan != "" {
		ringkasan = &u.Ringkasan
	}
	return umpanJSON{
		ID: u.ID, JenisSumber: string(u.JenisSumber), MediaNama: u.MediaNama, MediaDomain: u.MediaDomain,
		Judul: u.Judul, Ringkasan: ringkasan, URL: u.URL, TerbitPada: u.TerbitPada, TagIsu: u.TagIsu,
		DasarCocok: string(u.DasarCocok), ProvenanceID: u.ProvenanceID,
	}
}
