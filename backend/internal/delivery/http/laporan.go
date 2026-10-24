package http

import (
	"context"
	"io"
	"math"
	"mime/multipart"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const (
	cadanganBodyMultipart = 1 << 20
	BatasBodyPermintaan   = domain.BatasFotoLaporan*domain.BatasUkuranFotoByte + cadanganBodyMultipart
)

type LayananLaporan interface {
	Terima(ctx context.Context, masuk domain.LaporanMasuk) (usecase.LaporanDiterima, error)
}

type laporanDiterimaJSON struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type handlerLaporan struct {
	layanan LayananLaporan
}

func daftarkanLaporan(api fiber.Router, layanan Layanan) {
	if layanan.Laporan == nil || layanan.Akses == nil {
		return
	}
	api.Post("/laporan", wajibMasuk(layanan.Akses), handlerLaporan{layanan: layanan.Laporan}.terima)
}

func (h handlerLaporan) terima(c fiber.Ctx) error {
	masuk, err := bacaLaporanMasuk(c)
	if err != nil {
		return err
	}
	masuk.PenggunaID = klaimAkses(c).PenggunaID
	hasil, err := h.layanan.Terima(c.Context(), masuk)
	if err != nil {
		return err
	}
	c.Status(fiber.StatusCreated)
	return kirimData(c, laporanDiterimaJSON{ID: hasil.ID, Status: string(hasil.Status)})
}

func bacaLaporanMasuk(c fiber.Ctx) (domain.LaporanMasuk, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return domain.LaporanMasuk{}, domain.InputTidakValid("kirim laporan sebagai multipart/form-data")
	}
	nilai := func(kunci string) string {
		if isi := form.Value[kunci]; len(isi) > 0 {
			return isi[0]
		}
		return ""
	}
	lat, errLat := uraiDesimalHingga(nilai("lat"))
	lng, errLng := uraiDesimalHingga(nilai("lng"))
	akurasi, errAkurasi := uraiDesimalHingga(nilai("akurasi"))
	if errLat != nil || errLng != nil || errAkurasi != nil {
		return domain.LaporanMasuk{}, domain.InputTidakValid("lat, lng, dan akurasi wajib berupa angka")
	}
	foto, err := bacaFotoLaporan(form.File["foto"])
	return domain.LaporanMasuk{
		PaketKode: nilai("paket_kode"), Titik: domain.Titik{Lat: lat, Lng: lng}, AkurasiMeter: akurasi,
		Kondisi: domain.KondisiLaporan(nilai("kondisi")), Catatan: nilai("catatan"), Foto: foto,
	}, err
}

func uraiDesimalHingga(teks string) (float64, error) {
	nilai, err := strconv.ParseFloat(teks, 64)
	if err == nil && (math.IsNaN(nilai) || math.IsInf(nilai, 0)) {
		return 0, strconv.ErrSyntax
	}
	return nilai, err
}

func bacaFotoLaporan(berkas []*multipart.FileHeader) ([][]byte, error) {
	if len(berkas) == 0 || len(berkas) > domain.BatasFotoLaporan {
		return nil, domain.ErrJumlahFotoTidakSah
	}
	semua := make([][]byte, len(berkas))
	for i, fh := range berkas {
		if fh.Size > domain.BatasUkuranFotoByte {
			return nil, domain.ErrFotoTerlaluBesar
		}
		isi, err := bacaBerkas(fh)
		if err != nil {
			return nil, err
		}
		semua[i] = isi
	}
	return semua, nil
}

func bacaBerkas(fh *multipart.FileHeader) ([]byte, error) {
	berkas, err := fh.Open()
	if err != nil {
		return nil, domain.InputTidakValid("foto tidak dapat dibaca")
	}
	defer func() { _ = berkas.Close() }()
	return io.ReadAll(io.LimitReader(berkas, domain.BatasUkuranFotoByte+1))
}
