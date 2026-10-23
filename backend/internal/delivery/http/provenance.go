package http

import (
	"context"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananProvenance interface {
	Ambil(ctx context.Context, id string) (domain.Provenance, error)
}

type provenanceJSON struct {
	ID          string    `json:"id"`
	Sumber      string    `json:"sumber"`
	URL         string    `json:"url"`
	DiambilPada time.Time `json:"diambil_pada"`
	Label       string    `json:"label"`
	HashKonten  string    `json:"hash_konten"`
	Catatan     string    `json:"catatan"`
}

type handlerProvenance struct {
	layanan LayananProvenance
}

func (h handlerProvenance) ambil(c fiber.Ctx) error {
	p, err := h.layanan.Ambil(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	return kirimData(c, provenanceJSON{
		ID:          p.ID,
		Sumber:      string(p.Sumber),
		URL:         p.URL,
		DiambilPada: p.DiambilPada.UTC(),
		Label:       string(p.Label),
		HashKonten:  p.HashKonten,
		Catatan:     p.Catatan,
	})
}
