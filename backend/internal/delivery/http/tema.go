package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananTema interface {
	Daftar(ctx context.Context) ([]domain.TemaTersimpan, error)
	Ambil(ctx context.Context, kode string) (domain.TemaTersimpan, error)
}

type temaJSON struct {
	domain.Tema
	ProvenanceID string `json:"provenance_id"`
}

type handlerTema struct {
	layanan LayananTema
}

func (h handlerTema) daftar(c fiber.Ctx) error {
	hasil, err := h.layanan.Daftar(c.Context())
	if err != nil {
		return err
	}
	daftar := make([]temaJSON, len(hasil))
	for i, t := range hasil {
		daftar[i] = keTemaJSON(t)
	}
	return kirimData(c, daftar)
}

func (h handlerTema) ambil(c fiber.Ctx) error {
	tema, err := h.layanan.Ambil(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, keTemaJSON(tema))
}

func keTemaJSON(t domain.TemaTersimpan) temaJSON {
	return temaJSON{Tema: t.Tema, ProvenanceID: t.ProvenanceID}
}
