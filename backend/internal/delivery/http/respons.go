package http

import (
	"errors"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const pesanOK = "ok"

type Respons struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

var kodeHTTPKategori = map[domain.KategoriKesalahan]int{
	domain.KategoriTidakDitemukan:  fiber.StatusNotFound,
	domain.KategoriInputTidakValid: fiber.StatusBadRequest,
	domain.KategoriTidakBerwenang:  fiber.StatusUnauthorized,
	domain.KategoriTerlaluBanyak:   fiber.StatusTooManyRequests,
}

var pesanStatusUmum = map[int]string{
	fiber.StatusBadRequest:            "permintaan tidak valid",
	fiber.StatusUnauthorized:          "silakan masuk lebih dulu",
	fiber.StatusNotFound:              "halaman tidak ditemukan",
	fiber.StatusMethodNotAllowed:      "metode tidak diizinkan",
	fiber.StatusRequestEntityTooLarge: "permintaan terlalu besar",
	fiber.StatusTooManyRequests:       "terlalu banyak permintaan, coba lagi sebentar lagi",
	fiber.StatusInternalServerError:   "terjadi kesalahan di server",
}

func kirimData(c fiber.Ctx, data any) error {
	return c.JSON(Respons{Code: fiber.StatusOK, Message: pesanOK, Data: data})
}

func kirimHalaman(c fiber.Ctx, data any, meta Meta) error {
	return c.JSON(Respons{Code: fiber.StatusOK, Message: pesanOK, Data: data, Meta: &meta})
}

func tanganiError(c fiber.Ctx, err error) error {
	kode, pesan := petakanError(err)
	return c.Status(kode).JSON(Respons{Code: kode, Message: pesan})
}

func petakanError(err error) (int, string) {
	var kesalahanDomain *domain.Kesalahan
	if errors.As(err, &kesalahanDomain) {
		return kodeHTTPKategori[kesalahanDomain.Kategori], kesalahanDomain.Pesan
	}
	var kesalahanFiber *fiber.Error
	if errors.As(err, &kesalahanFiber) {
		return kesalahanFiber.Code, pesanStatus(kesalahanFiber.Code)
	}
	return fiber.StatusInternalServerError, pesanStatus(fiber.StatusInternalServerError)
}

func pesanStatus(kode int) string {
	if pesan, ada := pesanStatusUmum[kode]; ada {
		return pesan
	}
	return pesanStatusUmum[fiber.StatusBadRequest]
}
