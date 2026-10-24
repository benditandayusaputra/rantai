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
	IsuAnggota(ctx context.Context, anggotaID string) ([]domain.IsuDisuarakan, error)
}

type isuJSON struct {
	Tag              string      `json:"tag"`
	JumlahSumber     int         `json:"jumlah_sumber"`
	JumlahRilisResmi int         `json:"jumlah_rilis_resmi"`
	JumlahMedia      int         `json:"jumlah_media"`
	Item             []umpanJSON `json:"item"`
}

type umpanJSON struct {
	ID           int64          `json:"id"`
	JenisSumber  string         `json:"jenis_sumber"`
	MediaNama    string         `json:"media_nama"`
	MediaDomain  string         `json:"media_domain"`
	Judul        string         `json:"judul"`
	Ringkasan    *string        `json:"ringkasan"`
	URL          string         `json:"url"`
	TerbitPada   time.Time      `json:"terbit_pada"`
	TagIsu       []string       `json:"tag_isu"`
	DasarCocok   string         `json:"dasar_cocok"`
	ProvenanceID string         `json:"provenance_id"`
	HakJawab     []hakJawabJSON `json:"hak_jawab"`
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
	return kirimHalaman(c, keDaftarUmpanJSON(hasil.Item), buatMeta(halaman, hasil))
}

func (h handlerUmpan) isuAnggota(c fiber.Ctx) error {
	daftar, err := h.layanan.IsuAnggota(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	hasil := make([]isuJSON, len(daftar))
	for i, isu := range daftar {
		hasil[i] = isuJSON{
			Tag: isu.Tag, JumlahSumber: isu.JumlahSumber(), JumlahRilisResmi: isu.JumlahResmi,
			JumlahMedia: isu.JumlahMedia, Item: keDaftarUmpanJSON(isu.Item),
		}
	}
	return kirimData(c, hasil)
}

func keDaftarUmpanJSON(daftar []domain.UmpanTampil) []umpanJSON {
	hasil := make([]umpanJSON, len(daftar))
	for i, u := range daftar {
		hasil[i] = keUmpanJSON(u)
	}
	return hasil
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
		DasarCocok: string(u.DasarCocok), ProvenanceID: u.ProvenanceID, HakJawab: keDaftarHakJawabJSON(u.HakJawab),
	}
}
