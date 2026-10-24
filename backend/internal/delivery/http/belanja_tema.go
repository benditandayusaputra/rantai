package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananBelanja interface {
	Hitung(ctx context.Context, kodeTema string, tahun int) (domain.BelanjaTema, error)
}

type belanjaKecamatanJSON struct {
	Kecamatan ringkasanWilayahJSON `json:"kecamatan"`
	KodeKota  string               `json:"kode_kota"`
	Indikator *indikatorJSON       `json:"indikator"`
	Belanja   *anggaranJSON        `json:"belanja"`
}

type belanjaLuarKecamatanJSON struct {
	Wilayah ringkasanWilayahJSON `json:"wilayah"`
	Belanja *anggaranJSON        `json:"belanja"`
}

type belanjaTemaJSON struct {
	Tema          string                     `json:"tema"`
	Tahun         int                        `json:"tahun"`
	Kecamatan     []belanjaKecamatanJSON     `json:"kecamatan"`
	LuarKecamatan []belanjaLuarKecamatanJSON `json:"luar_kecamatan"`
}

type handlerBelanja struct {
	layanan LayananBelanja
}

func (h handlerBelanja) hitung(c fiber.Ctx) error {
	tahun, err := bilanganOpsional(c.Query("tahun"))
	if err != nil {
		return domain.InputTidakValid("tahun wajib berupa angka")
	}
	b, err := h.layanan.Hitung(c.Context(), c.Params("kode"), tahun)
	if err != nil {
		return err
	}
	return kirimData(c, keBelanjaTemaJSON(b))
}

func keBelanjaTemaJSON(b domain.BelanjaTema) belanjaTemaJSON {
	hasil := belanjaTemaJSON{
		Tema: b.Tema.Kode, Tahun: b.Tahun,
		Kecamatan:     make([]belanjaKecamatanJSON, len(b.Kecamatan)),
		LuarKecamatan: make([]belanjaLuarKecamatanJSON, len(b.LuarKecamatan)),
	}
	for i, k := range b.Kecamatan {
		tema := keTemaKecamatanJSON(domain.TemaKecamatan{Tema: b.Tema, Indikator: k.Indikator})
		hasil.Kecamatan[i] = belanjaKecamatanJSON{
			Kecamatan: keRingkasanWilayahJSON(k.Kecamatan), KodeKota: k.KodeKota, Indikator: tema.Indikator, Belanja: keAnggaranJSON(k.Anggaran),
		}
	}
	for i, l := range b.LuarKecamatan {
		hasil.LuarKecamatan[i] = belanjaLuarKecamatanJSON{Wilayah: keRingkasanWilayahJSON(l.Wilayah), Belanja: keAnggaranJSON(&l.Anggaran)}
	}
	return hasil
}
