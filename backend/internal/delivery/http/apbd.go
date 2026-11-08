package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananAPBD interface {
	RingkasBelanja(ctx context.Context, tahun int) (domain.RingkasanBelanjaAPBD, error)
}

type akunAPBDJSON struct {
	Nama      string `json:"nama"`
	Anggaran  Angka  `json:"anggaran"`
	Realisasi Angka  `json:"realisasi"`
}

type belanjaAPBDJSON struct {
	Tahun          int            `json:"tahun"`
	BulanRealisasi int            `json:"bulan_realisasi"`
	Belanja        akunAPBDJSON   `json:"belanja"`
	Jenis          []akunAPBDJSON `json:"jenis"`
}

type handlerAPBD struct {
	layanan LayananAPBD
}

func (h handlerAPBD) ringkas(c fiber.Ctx) error {
	tahun, err := angkaKueriOpsional(c.Query("tahun"))
	if err != nil {
		return domain.InputTidakValid("tahun wajib berupa angka")
	}
	ringkasan, err := h.layanan.RingkasBelanja(c.Context(), int(tahun))
	if err != nil {
		return err
	}
	return kirimData(c, keBelanjaAPBDJSON(ringkasan))
}

func keBelanjaAPBDJSON(r domain.RingkasanBelanjaAPBD) belanjaAPBDJSON {
	akunJSON := func(a domain.AkunAPBD) akunAPBDJSON {
		return akunAPBDJSON{
			Nama:      a.Nama,
			Anggaran:  Angka{Nilai: a.Anggaran, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: r.ProvenanceID},
			Realisasi: Angka{Nilai: a.Realisasi, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: r.ProvenanceID},
		}
	}
	jenis := make([]akunAPBDJSON, len(r.Jenis))
	for i, a := range r.Jenis {
		jenis[i] = akunJSON(a)
	}
	return belanjaAPBDJSON{Tahun: r.Tahun, BulanRealisasi: r.BulanRealisasi, Belanja: akunJSON(r.Belanja), Jenis: jenis}
}
