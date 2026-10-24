package http

import (
	"context"
	"encoding/json"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananPantau interface {
	Langgani(ctx context.Context, penggunaID int64, kode string) (domain.Langganan, error)
	Daftar(ctx context.Context, penggunaID int64) ([]domain.Langganan, error)
	Berhenti(ctx context.Context, penggunaID, id int64) error
	Notifikasi(ctx context.Context, penggunaID int64) (domain.HalamanNotifikasi, error)
	TandaiDibaca(ctx context.Context, penggunaID, id int64) error
}

type langgananJSON struct {
	ID          int64     `json:"id"`
	KodeWilayah string    `json:"kode_wilayah"`
	NamaWilayah string    `json:"nama_wilayah"`
	DibuatPada  time.Time `json:"dibuat_pada"`
}

type notifikasiJSON struct {
	ID            int64     `json:"id"`
	Jenis         string    `json:"jenis"`
	KodeKecamatan string    `json:"kode_kecamatan"`
	NamaKecamatan string    `json:"nama_kecamatan"`
	Judul         string    `json:"judul"`
	Tautan        string    `json:"tautan"`
	TerjadiPada   time.Time `json:"terjadi_pada"`
	Dibaca        bool      `json:"dibaca"`
	IsSimulation  bool      `json:"is_simulation"`
}

type halamanNotifikasiJSON struct {
	Item        []notifikasiJSON `json:"item"`
	BelumDibaca int              `json:"belum_dibaca"`
}

type handlerPantau struct {
	layanan LayananPantau
}

func daftarkanPantau(api fiber.Router, layanan Layanan) {
	if layanan.Pantau == nil || layanan.Akses == nil {
		return
	}
	h := handlerPantau{layanan: layanan.Pantau}
	masuk := wajibMasuk(layanan.Akses)
	api.Get("/langganan", masuk, h.daftar)
	api.Post("/langganan", masuk, h.langgani)
	api.Delete("/langganan/:id", masuk, h.berhenti)
	api.Get("/notifikasi", masuk, h.notifikasi)
	api.Patch("/notifikasi/:id/dibaca", masuk, h.dibaca)
}

func keLanggananJSON(l domain.Langganan) langgananJSON {
	return langgananJSON{ID: l.ID, KodeWilayah: l.KodeWilayah, NamaWilayah: l.NamaWilayah, DibuatPada: l.DibuatPada}
}

func (h handlerPantau) daftar(c fiber.Ctx) error {
	daftar, err := h.layanan.Daftar(c.Context(), klaimAkses(c).PenggunaID)
	if err != nil {
		return err
	}
	hasil := make([]langgananJSON, len(daftar))
	for i, l := range daftar {
		hasil[i] = keLanggananJSON(l)
	}
	return kirimData(c, hasil)
}

func (h handlerPantau) langgani(c fiber.Ctx) error {
	var masukan struct {
		KodeWilayah string `json:"kode_wilayah"`
	}
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi kode_wilayah")
	}
	l, err := h.layanan.Langgani(c.Context(), klaimAkses(c).PenggunaID, masukan.KodeWilayah)
	if err != nil {
		return err
	}
	c.Status(fiber.StatusCreated)
	return kirimData(c, keLanggananJSON(l))
}

func (h handlerPantau) berhenti(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	if err := h.layanan.Berhenti(c.Context(), klaimAkses(c).PenggunaID, id); err != nil {
		return err
	}
	return kirimData(c, nil)
}

func (h handlerPantau) notifikasi(c fiber.Ctx) error {
	halaman, err := h.layanan.Notifikasi(c.Context(), klaimAkses(c).PenggunaID)
	if err != nil {
		return err
	}
	hasil := halamanNotifikasiJSON{Item: make([]notifikasiJSON, len(halaman.Item)), BelumDibaca: halaman.JumlahBelumDibaca}
	for i, n := range halaman.Item {
		hasil.Item[i] = notifikasiJSON{
			ID: n.ID, Jenis: string(n.Jenis), KodeKecamatan: n.KodeKecamatan, NamaKecamatan: n.NamaKecamatan, Judul: n.Judul,
			Tautan: n.Tautan, TerjadiPada: n.TerjadiPada, Dibaca: n.IsDibaca, IsSimulation: n.IsSimulation,
		}
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, hasil)
}

func (h handlerPantau) dibaca(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	if err := h.layanan.TandaiDibaca(c.Context(), klaimAkses(c).PenggunaID, id); err != nil {
		return err
	}
	return kirimData(c, nil)
}
