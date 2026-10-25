package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const (
	tipeCSV          = "text/csv; charset=utf-8"
	pemisahTema      = ";"
	formatTanggalCSV = time.RFC3339
)

var kepalaCSVPaket = []string{
	"kode", "nama", "tahun_anggaran", "tahap", "jenis_pengadaan", "metode", "instansi", "penyedia",
	"pagu", "hps", "nilai_kontrak", "sektor", "tema", "tingkat_lokasi", "kode_wilayah", "nama_wilayah",
	"kode_rup", "data_simulasi", "label", "url_sumber", "diambil_pada", "provenance_id",
}

func (h handlerPaket) unduhCSV(c fiber.Ctx) error {
	tahun, err := bilanganOpsional(c.Query("tahun"))
	if err != nil {
		return domain.InputTidakValid("tahun wajib berupa angka")
	}
	kodeWilayah := c.Query("kode_wilayah")
	baris, err := h.layanan.Unduh(c.Context(), usecase.FilterUnduhPaket{KodeWilayah: kodeWilayah, Tahun: tahun})
	if err != nil {
		return err
	}
	isi, err := tulisCSVPaket(baris)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, tipeCSV)
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, namaBerkasCSV(kodeWilayah, tahun)))
	return c.Send(isi)
}

func namaBerkasCSV(kodeWilayah string, tahun int) string {
	nama := "rantai-paket-" + kodeWilayah
	if tahun > 0 {
		nama += "-" + strconv.Itoa(tahun)
	}
	return nama + ".csv"
}

func tulisCSVPaket(daftar []domain.BarisUnduhPaket) ([]byte, error) {
	var isi bytes.Buffer
	penulis := csv.NewWriter(&isi)
	if err := penulis.Write(kepalaCSVPaket); err != nil {
		return nil, fmt.Errorf("tulis kepala csv: %w", err)
	}
	for _, p := range daftar {
		if err := penulis.Write(barisCSVPaket(p)); err != nil {
			return nil, fmt.Errorf("tulis baris csv %s: %w", p.Kode, err)
		}
	}
	penulis.Flush()
	return isi.Bytes(), penulis.Error()
}

func barisCSVPaket(p domain.BarisUnduhPaket) []string {
	return []string{
		p.Kode, p.Nama, strconv.Itoa(p.TahunAnggaran), p.Tahap, p.JenisPengadaan, p.Metode, p.Instansi, p.NamaPenyedia,
		angkaCSV(p.Pagu), angkaCSV(p.HPS), angkaCSV(p.NilaiKontrak), p.Sektor, strings.Join(p.TemaKode, pemisahTema),
		string(p.Lokasi.Tingkat), p.Lokasi.KodeWilayah, p.Lokasi.NamaWilayah,
		p.KodeRUP, strconv.FormatBool(p.IsSimulation), string(p.Label), p.URLSumber,
		p.DiambilPada.UTC().Format(formatTanggalCSV), p.ProvenanceID,
	}
}

func angkaCSV(nilai *float64) string {
	if nilai == nil {
		return ""
	}
	return strconv.FormatFloat(*nilai, 'f', -1, 64)
}
