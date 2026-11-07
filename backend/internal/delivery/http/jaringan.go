package http

import (
	"context"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananJaringan interface {
	Susun(ctx context.Context, f domain.FilterJaringan) (domain.JaringanRelasi, error)
}

type simpulJaringanJSON struct {
	ID          string             `json:"id"`
	Jenis       domain.JenisSimpul `json:"jenis"`
	Kunci       string             `json:"kunci"`
	Label       string             `json:"label"`
	Nilai       Angka              `json:"nilai"`
	JumlahPaket int                `json:"jumlah_paket"`
}

type sisiJaringanJSON struct {
	Dari string `json:"dari"`
	Ke   string `json:"ke"`
}

type jaringanJSON struct {
	Simpul      []simpulJaringanJSON `json:"simpul"`
	Sisi        []sisiJaringanJSON   `json:"sisi"`
	IsTerpotong bool                 `json:"is_terpotong"`
	BatasPaket  int                  `json:"batas_paket"`
}

type handlerJaringan struct {
	layanan LayananJaringan
}

func (h handlerJaringan) susun(c fiber.Ctx) error {
	filter, err := bacaFilterJaringan(c)
	if err != nil {
		return err
	}
	jaringan, err := h.layanan.Susun(c.Context(), filter)
	if err != nil {
		return err
	}
	return kirimData(c, keJaringanJSON(jaringan))
}

func bacaFilterJaringan(c fiber.Ctx) (domain.FilterJaringan, error) {
	penyedia, errPenyedia := angkaKueriOpsional(c.Query("penyedia"))
	tahun, errTahun := angkaKueriOpsional(c.Query("tahun"))
	if errPenyedia != nil || errTahun != nil {
		return domain.FilterJaringan{}, domain.InputTidakValid("penyedia dan tahun wajib berupa angka")
	}
	return domain.FilterJaringan{KodeWilayah: c.Query("wilayah"), PenyediaID: penyedia, Tahun: int(tahun)}, nil
}

func angkaKueriOpsional(teks string) (int64, error) {
	if teks == "" {
		return 0, nil
	}
	return strconv.ParseInt(teks, 10, 64)
}

func keJaringanJSON(j domain.JaringanRelasi) jaringanJSON {
	hasil := jaringanJSON{
		Simpul: make([]simpulJaringanJSON, len(j.Simpul)), Sisi: make([]sisiJaringanJSON, len(j.Sisi)),
		IsTerpotong: j.IsTerpotong, BatasPaket: domain.BatasPaketJaringan,
	}
	for i, s := range j.Simpul {
		hasil.Simpul[i] = simpulJaringanJSON{
			ID: s.ID, Jenis: s.Jenis, Kunci: s.Kunci, Label: s.Label, JumlahPaket: s.JumlahPaket,
			Nilai: Angka{Nilai: s.Nilai, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: s.ProvenanceID},
		}
	}
	for i, s := range j.Sisi {
		hasil.Sisi[i] = sisiJaringanJSON(s)
	}
	return hasil
}
