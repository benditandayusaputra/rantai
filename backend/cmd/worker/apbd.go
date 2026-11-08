package main

import (
	"net/http"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/djpk"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathKonfigurasiDJPK = "config/djpk.yaml"

func pilihJobAnggaran(nama string, pool *pgxpool.Pool) (job, error) {
	if nama != usecase.JobAPBD {
		return pilihJobPantau(nama, pool)
	}
	var cfg djpk.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiDJPK, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporAPBD{
		Pengambil: djpk.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Klien),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Penyimpan: repository.APBDRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Tahun:     cfg.Tahun,
		Sekarang:  time.Now,
	}
	return impor.Jalankan, nil
}
