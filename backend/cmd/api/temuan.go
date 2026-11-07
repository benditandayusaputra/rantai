package main

import (
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func denganTemuan(layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	var terlarang domain.KosakataTerlarang
	if err := config.BacaYAML(pathKosakataTerlarang, &terlarang); err != nil {
		return layanan, err
	}
	layanan.Temuan = usecase.KelolaTemuan{Repo: repository.TemuanRepository{DB: pool}, Terlarang: terlarang, Sekarang: time.Now}
	return layanan, nil
}
