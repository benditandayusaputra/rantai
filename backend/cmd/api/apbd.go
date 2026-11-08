package main

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathKonfigurasiDJPK = "config/djpk.yaml"

func denganAPBD(layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	var cfg struct {
		Aturan domain.AturanAPBD `yaml:"aturan"`
	}
	if err := config.BacaYAML(pathKonfigurasiDJPK, &cfg); err != nil {
		return layanan, err
	}
	layanan.APBD = usecase.BacaAPBD{Pembaca: repository.APBDRepository{DB: pool}, Aturan: cfg.Aturan}
	return layanan, nil
}
