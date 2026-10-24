package main

import (
	"errors"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/paripurna"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathKonfigurasiParipurna = "config/paripurna.yaml"

func buatImporParipurna(pool *pgxpool.Pool) (job, error) {
	var cfg paripurna.Konfigurasi
	var terlarang domain.KosakataTerlarang
	if err := errors.Join(config.BacaYAML(pathKonfigurasiParipurna, &cfg), config.BacaYAML(pathKosakataTerlarang, &terlarang)); err != nil {
		return nil, err
	}
	impor := usecase.ImporParipurna{
		Pengambil: paripurna.Sumber{Konfigurasi: cfg, Sekarang: time.Now},
		Penyimpan: repository.ParipurnaRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Terlarang: terlarang,
	}
	return impor.Jalankan, nil
}
