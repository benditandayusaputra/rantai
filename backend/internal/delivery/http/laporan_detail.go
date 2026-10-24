package http

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type LayananBacaLaporan interface {
	Ambil(ctx context.Context, id int64) (usecase.DetailLaporan, error)
}

type LayananKonfirmasiLaporan interface {
	Konfirmasi(ctx context.Context, masuk domain.KonfirmasiMasuk) (usecase.HasilKonfirmasi, error)
}

type saksiJSON struct {
	Jumlah     int `json:"jumlah"`
	Dibutuhkan int `json:"dibutuhkan"`
}

type paketLaporanJSON struct {
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

type fotoLaporanJSON struct {
	URL          string `json:"url"`
	URLThumbnail string `json:"url_thumbnail"`
	Lebar        int    `json:"lebar"`
	Tinggi       int    `json:"tinggi"`
}

type riwayatLaporanJSON struct {
	StatusLama *string   `json:"status_lama"`
	StatusBaru string    `json:"status_baru"`
	Alasan     string    `json:"alasan"`
	DibuatPada time.Time `json:"dibuat_pada"`
}

type detailLaporanJSON struct {
	ID           int64                `json:"id"`
	Status       string               `json:"status"`
	Kondisi      string               `json:"kondisi"`
	Catatan      string               `json:"catatan"`
	KodeWilayah  string               `json:"kode_wilayah"`
	NamaWilayah  string               `json:"nama_wilayah"`
	LokasiPublik titikJSON            `json:"lokasi_publik"`
	Paket        *paketLaporanJSON    `json:"paket"`
	IsSimulation bool                 `json:"is_simulation"`
	DibuatPada   time.Time            `json:"dibuat_pada"`
	Saksi        saksiJSON            `json:"saksi"`
	Foto         []fotoLaporanJSON    `json:"foto"`
	Riwayat      []riwayatLaporanJSON `json:"riwayat"`
}

type konfirmasiMasukJSON struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Akurasi float64 `json:"akurasi"`
	Kondisi string  `json:"kondisi"`
}

type hasilKonfirmasiJSON struct {
	Status string    `json:"status"`
	Saksi  saksiJSON `json:"saksi"`
}

type handlerDetailLaporan struct {
	baca       LayananBacaLaporan
	konfirmasi LayananKonfirmasiLaporan
}

func uraiIDLaporan(c fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return 0, domain.InputTidakValid("id laporan wajib berupa angka")
	}
	return id, nil
}

func (h handlerDetailLaporan) ambil(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	detail, err := h.baca.Ambil(c.Context(), id)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, keDetailLaporanJSON(detail))
}

func (h handlerDetailLaporan) konfirmasikan(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	var masukan konfirmasiMasukJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi lat, lng, akurasi, dan kondisi")
	}
	hasil, err := h.konfirmasi.Konfirmasi(c.Context(), domain.KonfirmasiMasuk{
		LaporanID: id, PenggunaID: klaimAkses(c).PenggunaID, Titik: domain.Titik{Lat: masukan.Lat, Lng: masukan.Lng},
		AkurasiMeter: masukan.Akurasi, Kondisi: domain.KondisiLaporan(masukan.Kondisi),
	})
	if err != nil {
		return err
	}
	return kirimData(c, hasilKonfirmasiJSON{Status: string(hasil.Status), Saksi: saksiDari(hasil.JumlahSaksi)})
}

func saksiDari(jumlah int) saksiJSON {
	return saksiJSON{Jumlah: jumlah, Dibutuhkan: domain.JumlahKonfirmasiTerverifikasi}
}

func keDetailLaporanJSON(d usecase.DetailLaporan) detailLaporanJSON {
	hasil := detailLaporanJSON{
		ID: d.ID, Status: string(d.Status), Kondisi: string(d.Kondisi), Catatan: d.Catatan, KodeWilayah: d.KodeWilayah,
		NamaWilayah: d.NamaWilayah, LokasiPublik: titikJSON{Lat: d.TitikPublik.Lat, Lng: d.TitikPublik.Lng},
		IsSimulation: d.IsSimulation, DibuatPada: d.DibuatPada, Saksi: saksiDari(d.JumlahSaksi),
		Foto: make([]fotoLaporanJSON, len(d.FotoBertanda)), Riwayat: make([]riwayatLaporanJSON, len(d.Riwayat)),
	}
	if d.PaketKode != "" {
		hasil.Paket = &paketLaporanJSON{Kode: d.PaketKode, Nama: d.PaketNama}
	}
	for i, f := range d.FotoBertanda {
		hasil.Foto[i] = fotoLaporanJSON{URL: f.URL, URLThumbnail: f.URLThumbnail, Lebar: f.Lebar, Tinggi: f.Tinggi}
	}
	for i, r := range d.Riwayat {
		hasil.Riwayat[i] = keRiwayatLaporanJSON(r)
	}
	return hasil
}

func keRiwayatLaporanJSON(r domain.RiwayatLaporan) riwayatLaporanJSON {
	hasil := riwayatLaporanJSON{StatusBaru: string(r.StatusBaru), Alasan: r.Alasan, DibuatPada: r.DibuatPada}
	if r.StatusLama != "" {
		lama := string(r.StatusLama)
		hasil.StatusLama = &lama
	}
	return hasil
}
