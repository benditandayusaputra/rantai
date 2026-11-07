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

type LayananTemuan interface {
	Daftar(ctx context.Context, paketKode string) ([]domain.Temuan, error)
	Kirim(ctx context.Context, masuk domain.TemuanMasuk) (int64, error)
	Tanggapi(ctx context.Context, masuk domain.TanggapanMasuk) error
	Sembunyikan(ctx context.Context, p usecase.PermintaanModerasiTemuan) error
}

type temuanJSON struct {
	ID              int64               `json:"id"`
	Isi             string              `json:"isi"`
	TautanSumber    string              `json:"tautan_sumber"`
	NamaKontributor string              `json:"nama_kontributor"`
	DibuatPada      time.Time           `json:"dibuat_pada"`
	JumlahMendukung int                 `json:"jumlah_mendukung"`
	JumlahMembantah int                 `json:"jumlah_membantah"`
	Status          domain.StatusTemuan `json:"status"`
}

type temuanMasukJSON struct {
	Isi          string `json:"isi"`
	TautanSumber string `json:"tautan_sumber"`
	Sikap        string `json:"sikap"`
	Alasan       string `json:"alasan"`
}

type temuanDiterimaJSON struct {
	ID int64 `json:"id"`
}

type handlerTemuan struct {
	layanan LayananTemuan
}

func daftarkanTemuan(api fiber.Router, layanan Layanan) {
	if layanan.Temuan == nil {
		return
	}
	h := handlerTemuan{layanan: layanan.Temuan}
	api.Get("/paket/:kode/temuan", h.daftar)
	if layanan.Akses == nil {
		return
	}
	api.Post("/paket/:kode/temuan", wajibMasuk(layanan.Akses), h.kirim)
	api.Post("/temuan/:id/tanggapan", wajibMasuk(layanan.Akses), h.tanggapi)
	api.Patch("/moderasi/temuan/:id", wajibMasuk(layanan.Akses), wajibModerator, h.sembunyikan)
}

func (h handlerTemuan) daftar(c fiber.Ctx) error {
	daftar, err := h.layanan.Daftar(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	hasil := make([]temuanJSON, len(daftar))
	for i, t := range daftar {
		hasil[i] = temuanJSON{
			ID: t.ID, Isi: t.Isi, TautanSumber: t.TautanSumber, NamaKontributor: t.NamaKontributor, DibuatPada: t.DibuatPada.UTC(),
			JumlahMendukung: t.JumlahMendukung, JumlahMembantah: t.JumlahMembantah, Status: t.Status(),
		}
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, hasil)
}

func (h handlerTemuan) kirim(c fiber.Ctx) error {
	badan, err := bacaTemuanMasuk(c)
	if err != nil {
		return err
	}
	id, err := h.layanan.Kirim(c.Context(), domain.TemuanMasuk{
		PenggunaID: klaimAkses(c).PenggunaID, PaketKode: c.Params("kode"), Isi: badan.Isi, TautanSumber: badan.TautanSumber,
	})
	if err != nil {
		return err
	}
	c.Status(fiber.StatusCreated)
	return kirimData(c, temuanDiterimaJSON{ID: id})
}

func (h handlerTemuan) tanggapi(c fiber.Ctx) error {
	id, errID := uraiIDTemuan(c)
	badan, errBadan := bacaTemuanMasuk(c)
	if errID != nil || errBadan != nil {
		return domain.InputTidakValid("id temuan wajib angka dan badan permintaan wajib JSON")
	}
	err := h.layanan.Tanggapi(c.Context(), domain.TanggapanMasuk{
		TemuanID: id, PenggunaID: klaimAkses(c).PenggunaID, Sikap: domain.SikapTemuan(badan.Sikap), TautanSumber: badan.TautanSumber,
	})
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}

func (h handlerTemuan) sembunyikan(c fiber.Ctx) error {
	id, errID := uraiIDTemuan(c)
	badan, errBadan := bacaTemuanMasuk(c)
	if errID != nil || errBadan != nil {
		return domain.InputTidakValid("id temuan wajib angka dan badan permintaan wajib JSON berisi alasan")
	}
	err := h.layanan.Sembunyikan(c.Context(), usecase.PermintaanModerasiTemuan{
		TemuanID: id, ModeratorID: klaimAkses(c).PenggunaID, Alasan: badan.Alasan,
	})
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}

func bacaTemuanMasuk(c fiber.Ctx) (temuanMasukJSON, error) {
	var badan temuanMasukJSON
	if err := json.Unmarshal(c.Body(), &badan); err != nil {
		return badan, domain.InputTidakValid("badan permintaan wajib JSON")
	}
	return badan, nil
}

func uraiIDTemuan(c fiber.Ctx) (int64, error) {
	return strconv.ParseInt(c.Params("id"), 10, 64)
}
