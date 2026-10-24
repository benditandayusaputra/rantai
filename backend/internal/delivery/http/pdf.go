package http

import (
	"bytes"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/go-pdf/fpdf"
	"github.com/gofiber/fiber/v3"
)

const (
	tipePDF          = "application/pdf"
	hurufPDF         = "Helvetica"
	gayaTebal        = "B"
	gayaBiasa        = ""
	ukuranJudul      = 14
	ukuranSubjudul   = 11
	ukuranIsi        = 10
	tinggiBarisJudul = 7
	tinggiBarisSub   = 6
	tinggiBarisIsi   = 5
	jarakAntarBagian = 4
	marginHalamanMM  = 20
	lebarPenuh       = 0
	rataKiri         = "L"
	tanpaBingkai     = ""
	encodingBawaan   = ""
)

func tulisPDF(d domain.Dokumen) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginHalamanMM, marginHalamanMM, marginHalamanMM)
	pdf.SetAutoPageBreak(true, marginHalamanMM)
	pdf.SetTitle(d.Judul, true)
	pdf.AddPage()
	terjemah := pdf.UnicodeTranslatorFromDescriptor(encodingBawaan)
	tulis := func(gaya string, ukuran, tinggi float64, teks string) {
		pdf.SetFont(hurufPDF, gaya, ukuran)
		pdf.MultiCell(lebarPenuh, tinggi, terjemah(teks), tanpaBingkai, rataKiri, false)
	}
	tulis(gayaTebal, ukuranJudul, tinggiBarisJudul, d.Judul)
	for _, bagian := range d.Bagian {
		pdf.Ln(jarakAntarBagian)
		tulis(gayaTebal, ukuranSubjudul, tinggiBarisSub, bagian.Judul)
		for _, baris := range bagian.Baris {
			tulis(gayaBiasa, ukuranIsi, tinggiBarisIsi, baris)
		}
	}
	var keluaran bytes.Buffer
	if err := pdf.Output(&keluaran); err != nil {
		return nil, fmt.Errorf("tulis pdf %q: %w", d.Judul, err)
	}
	return keluaran.Bytes(), nil
}

func kirimPDF(c fiber.Ctx, namaBerkas string, d domain.Dokumen) error {
	isi, err := tulisPDF(d)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, tipePDF)
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+namaBerkas+`"`)
	return c.Send(isi)
}
