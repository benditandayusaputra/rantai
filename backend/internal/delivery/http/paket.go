package http

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const (
	satuanPeserta  = "peserta"
	satuanHari     = "hari"
	satuanMeter    = "meter"
	jumlahSisiBBox = 4
)

type LayananPaket interface {
	Ambil(ctx context.Context, kode string) (usecase.DetailPaket, error)
	DaftarKecamatan(ctx context.Context, f usecase.FilterPaketKecamatan) (usecase.HalamanPaket, error)
	Peta(ctx context.Context, f usecase.FilterPetaPaket) (usecase.PetaPaket, error)
	Terdekat(ctx context.Context, titik domain.Titik, radiusMeter float64) ([]domain.TitikPaket, error)
}

type handlerPaket struct {
	layanan LayananPaket
}

func daftarkanPaket(api fiber.Router, h handlerPaket) {
	api.Get("/paket/terdekat", h.terdekat)
	api.Get("/paket/:kode", h.ambil)
	api.Get("/kecamatan/:kode/paket", h.daftarKecamatan)
	api.Get("/peta/paket", h.peta)
}

func (h handlerPaket) ambil(c fiber.Ctx) error {
	detail, err := h.layanan.Ambil(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, keDetailPaketJSON(detail))
}

func (h handlerPaket) daftarKecamatan(c fiber.Ctx) error {
	tahun, errTahun := bilanganOpsional(c.Query("tahun"))
	halaman, errHalaman := bilanganOpsional(c.Query("halaman"))
	if errTahun != nil || errHalaman != nil {
		return domain.InputTidakValid("tahun dan halaman wajib berupa angka")
	}
	hasil, err := h.layanan.DaftarKecamatan(c.Context(), usecase.FilterPaketKecamatan{
		KodeKecamatan: c.Params("kode"), Tema: c.Query("tema"), Tahun: tahun, Halaman: halaman,
	})
	if err != nil {
		return err
	}
	return c.JSON(Respons{Code: fiber.StatusOK, Message: pesanOK, Data: keDaftarPaketJSON(hasil.Paket), Meta: metaHalaman(hasil)})
}

func (h handlerPaket) peta(c fiber.Ctx) error {
	kotak, err := uraiKotak(c.Query("bbox"))
	if err != nil {
		return err
	}
	tahun, err := bilanganOpsional(c.Query("tahun"))
	if err != nil {
		return domain.InputTidakValid("tahun wajib berupa angka")
	}
	peta, err := h.layanan.Peta(c.Context(), usecase.FilterPetaPaket{Kotak: kotak, Tema: c.Query("tema"), Tahun: tahun})
	if err != nil {
		return err
	}
	return kirimData(c, kePetaJSON(peta))
}

func (h handlerPaket) terdekat(c fiber.Ctx) error {
	titik, err := uraiTitik(c.Query("lat"), c.Query("lng"))
	if err != nil {
		return err
	}
	radius, err := pecahanOpsional(c.Query("radius_m"))
	if err != nil {
		return domain.InputTidakValid("radius_m wajib berupa angka")
	}
	daftar, err := h.layanan.Terdekat(c.Context(), titik, radius)
	if err != nil {
		return err
	}
	hasil := make([]paketTerdekatJSON, len(daftar))
	for i, t := range daftar {
		hasil[i] = kePaketTerdekatJSON(t)
	}
	return kirimData(c, hasil)
}

func bilanganOpsional(teks string) (int, error) {
	if teks == "" {
		return 0, nil
	}
	return strconv.Atoi(teks)
}

func pecahanOpsional(teks string) (float64, error) {
	if teks == "" {
		return 0, nil
	}
	return strconv.ParseFloat(teks, ukuranBitFloat)
}

func uraiKotak(teks string) (domain.Kotak, error) {
	bagian := strings.Split(teks, ",")
	if len(bagian) != jumlahSisiBBox {
		return domain.Kotak{}, domain.InputTidakValid("bbox wajib berupa barat,selatan,timur,utara")
	}
	nilai := make([]float64, jumlahSisiBBox)
	for i, b := range bagian {
		n, err := strconv.ParseFloat(strings.TrimSpace(b), ukuranBitFloat)
		if err != nil {
			return domain.Kotak{}, domain.InputTidakValid("bbox wajib berupa angka")
		}
		nilai[i] = n
	}
	return domain.Kotak{Barat: nilai[0], Selatan: nilai[1], Timur: nilai[2], Utara: nilai[3]}, nil
}

func metaHalaman(h usecase.HalamanPaket) *Meta {
	totalHalaman := 0
	if h.PerHalaman > 0 {
		totalHalaman = (h.Total + h.PerHalaman - 1) / h.PerHalaman
	}
	return &Meta{Page: h.Halaman, PerPage: h.PerHalaman, Total: h.Total, TotalPages: totalHalaman}
}

func angkaOpsional(nilai *float64, satuan string, label domain.Label, provenanceID string) *Angka {
	if nilai == nil {
		return nil
	}
	return &Angka{Nilai: *nilai, Satuan: satuan, Label: label, ProvenanceID: provenanceID}
}

func labelPaket(isSimulation bool) domain.Label {
	if isSimulation {
		return domain.LabelSimulasi
	}
	return domain.LabelResmi
}

func tanggalJSON(t *time.Time) *string {
	if t == nil {
		return nil
	}
	teks := t.Format(time.DateOnly)
	return &teks
}

func keDaftarPaketJSON(daftar []domain.PaketRingkas) []paketRingkasJSON {
	hasil := make([]paketRingkasJSON, len(daftar))
	for i, p := range daftar {
		hasil[i] = kePaketRingkasJSON(p)
	}
	return hasil
}
