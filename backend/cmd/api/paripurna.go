package main

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathAspirasiAnggaran = "config/aspirasi_anggaran.yaml"

func buatBacaParipurna(pool *pgxpool.Pool) (usecase.BacaParipurna, error) {
	var anggaran domain.KonfigurasiAspirasiAnggaran
	if err := config.BacaYAML(pathAspirasiAnggaran, &anggaran); err != nil {
		return usecase.BacaParipurna{}, err
	}
	if err := anggaran.Validasi(); err != nil {
		return usecase.BacaParipurna{}, err
	}
	return usecase.BacaParipurna{Pembaca: repository.ParipurnaRepository{DB: pool}, Anggaran: anggaran}, nil
}
