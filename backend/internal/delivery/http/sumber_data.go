package http

import (
	"context"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananSumberData interface {
	Baca(ctx context.Context) (domain.InfoSumberData, error)
}

type pengambilanJSON struct {
	Job            string     `json:"job"`
	StatusTerakhir string     `json:"status_terakhir"`
	DijalankanPada time.Time  `json:"dijalankan_pada"`
	BerhasilPada   *time.Time `json:"berhasil_pada"`
}

type mediaJSON struct {
	Nama   string `json:"nama"`
	Domain string `json:"domain"`
}

type sumberDataJSON struct {
	Pengambilan        []pengambilanJSON `json:"pengambilan"`
	Media              []mediaJSON       `json:"media"`
	AtribusiOpentender string            `json:"atribusi_opentender"`
}

type handlerSumberData struct {
	layanan LayananSumberData
}

func (h handlerSumberData) baca(c fiber.Ctx) error {
	data, err := h.layanan.Baca(c.Context())
	if err != nil {
		return err
	}
	hasil := sumberDataJSON{
		Pengambilan:        make([]pengambilanJSON, len(data.Pengambilan)),
		Media:              make([]mediaJSON, len(data.Media)),
		AtribusiOpentender: data.AtribusiOpentender,
	}
	for i, p := range data.Pengambilan {
		hasil.Pengambilan[i] = pengambilanJSON{Job: p.Job, StatusTerakhir: p.StatusTerakhir, DijalankanPada: p.DijalankanPada.UTC(), BerhasilPada: keUTC(p.BerhasilPada)}
	}
	for i, m := range data.Media {
		hasil.Media[i] = mediaJSON(m)
	}
	return kirimData(c, hasil)
}

func keUTC(waktu *time.Time) *time.Time {
	if waktu == nil {
		return nil
	}
	utc := waktu.UTC()
	return &utc
}
