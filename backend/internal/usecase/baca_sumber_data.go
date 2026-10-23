package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaPengambilan interface {
	DaftarPengambilanTerakhir(ctx context.Context) ([]domain.PengambilanSumber, error)
}

type BacaSumberData struct {
	Pembaca            PembacaPengambilan
	Media              []domain.MediaTerverifikasi
	AtribusiOpentender string
}

func (u BacaSumberData) Baca(ctx context.Context) (domain.InfoSumberData, error) {
	pengambilan, err := u.Pembaca.DaftarPengambilanTerakhir(ctx)
	if err != nil {
		return domain.InfoSumberData{}, fmt.Errorf("baca pengambilan sumber: %w", err)
	}
	media := u.Media
	if media == nil {
		media = []domain.MediaTerverifikasi{}
	}
	return domain.InfoSumberData{Pengambilan: pengambilan, Media: media, AtribusiOpentender: u.AtribusiOpentender}, nil
}
