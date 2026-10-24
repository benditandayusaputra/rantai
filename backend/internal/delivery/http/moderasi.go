package http

import (
	"context"
	"encoding/json"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type LayananModerasi interface {
	Antrean(ctx context.Context) (usecase.AntreanModerasi, error)
	ModerasiLaporan(ctx context.Context, p domain.PermintaanModerasiLaporan) error
	SembunyikanUmpan(ctx context.Context, p usecase.PermintaanModerasiUmpan) error
}

type fotoModerasiJSON struct {
	ID            int64  `json:"id"`
	URL           string `json:"url"`
	URLThumbnail  string `json:"url_thumbnail"`
	Disembunyikan bool   `json:"disembunyikan"`
}

type laporanModerasiJSON struct {
	ID           int64              `json:"id"`
	Kondisi      string             `json:"kondisi"`
	Catatan      string             `json:"catatan"`
	NamaWilayah  string             `json:"nama_wilayah"`
	PaketNama    string             `json:"paket_nama"`
	IsSimulation bool               `json:"is_simulation"`
	DibuatPada   time.Time          `json:"dibuat_pada"`
	Foto         []fotoModerasiJSON `json:"foto"`
}

type umpanModerasiJSON struct {
	ID          int64     `json:"id"`
	JenisSumber string    `json:"jenis_sumber"`
	MediaNama   string    `json:"media_nama"`
	Judul       string    `json:"judul"`
	URL         string    `json:"url"`
	TerbitPada  time.Time `json:"terbit_pada"`
	Anggota     []string  `json:"anggota"`
}

type antreanModerasiJSON struct {
	Laporan     []laporanModerasiJSON `json:"laporan"`
	Umpan       []umpanModerasiJSON   `json:"umpan"`
	AlasanTolak []string              `json:"alasan_tolak"`
}

type aksiModerasiJSON struct {
	Aksi   string `json:"aksi"`
	Alasan string `json:"alasan"`
	Pesan  string `json:"pesan"`
	FotoID int64  `json:"foto_id"`
}

type handlerModerasi struct {
	layanan LayananModerasi
}

func daftarkanModerasi(api fiber.Router, layanan Layanan) {
	if layanan.Moderasi == nil || layanan.Akses == nil {
		return
	}
	h := handlerModerasi{layanan: layanan.Moderasi}
	moderasi := api.Group("/moderasi", wajibMasuk(layanan.Akses), wajibModerator)
	moderasi.Get("/antrean", h.antrean)
	moderasi.Patch("/laporan/:id", h.laporan)
	moderasi.Patch("/umpan/:id", h.umpan)
}

func wajibModerator(c fiber.Ctx) error {
	if !domain.IsModerator(klaimAkses(c).Peran) {
		return domain.ErrBukanModerator
	}
	return c.Next()
}

func (h handlerModerasi) antrean(c fiber.Ctx) error {
	antrean, err := h.layanan.Antrean(c.Context())
	if err != nil {
		return err
	}
	hasil := antreanModerasiJSON{
		Laporan: make([]laporanModerasiJSON, len(antrean.Laporan)), Umpan: make([]umpanModerasiJSON, len(antrean.Umpan)),
		AlasanTolak: domain.AlasanTolakSah,
	}
	for i, l := range antrean.Laporan {
		hasil.Laporan[i] = keLaporanModerasiJSON(l)
	}
	for i, u := range antrean.Umpan {
		hasil.Umpan[i] = umpanModerasiJSON{
			ID: u.ID, JenisSumber: string(u.JenisSumber), MediaNama: u.MediaNama, Judul: u.Judul, URL: u.URL,
			TerbitPada: u.TerbitPada, Anggota: u.Anggota,
		}
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, hasil)
}

func keLaporanModerasiJSON(l usecase.LaporanAntrean) laporanModerasiJSON {
	foto := make([]fotoModerasiJSON, len(l.FotoBertanda))
	for i, f := range l.FotoBertanda {
		foto[i] = fotoModerasiJSON{ID: f.ID, URL: f.URL, URLThumbnail: f.URLThumbnail, Disembunyikan: f.Disembunyikan}
	}
	return laporanModerasiJSON{
		ID: l.ID, Kondisi: string(l.Kondisi), Catatan: l.Catatan, NamaWilayah: l.NamaWilayah, PaketNama: l.PaketNama,
		IsSimulation: l.IsSimulation, DibuatPada: l.DibuatPada, Foto: foto,
	}
}

func (h handlerModerasi) laporan(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	var aksi aksiModerasiJSON
	if err := json.Unmarshal(c.Body(), &aksi); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi aksi")
	}
	err = h.layanan.ModerasiLaporan(c.Context(), domain.PermintaanModerasiLaporan{
		LaporanID: id, ModeratorID: klaimAkses(c).PenggunaID, Aksi: domain.AksiModerasi(aksi.Aksi),
		Alasan: aksi.Alasan, Pesan: aksi.Pesan, FotoID: aksi.FotoID,
	})
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}

func (h handlerModerasi) umpan(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	var aksi aksiModerasiJSON
	if err := json.Unmarshal(c.Body(), &aksi); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi alasan")
	}
	err = h.layanan.SembunyikanUmpan(c.Context(), usecase.PermintaanModerasiUmpan{
		UmpanID: id, ModeratorID: klaimAkses(c).PenggunaID, Alasan: aksi.Alasan,
	})
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}
