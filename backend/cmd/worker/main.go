package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiWilayah = "config/wilayah.yaml"
	dirRawBawaan           = "../data/raw"
	batasWaktuUnduh        = time.Minute
)

type job func(ctx context.Context) error

func main() {
	if err := jalankan(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func jalankan(argumen []string) error {
	if len(argumen) != 1 {
		return errors.New("pemakaian: worker <job>; job tersedia: " + usecase.JobWilayah)
	}
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("sambung basis data: %w", err)
	}
	defer pool.Close()
	jalankanJob, err := pilihJob(argumen[0], pool)
	if err != nil {
		return err
	}
	return jalankanJob(ctx)
}

func pilihJob(nama string, pool *pgxpool.Pool) (job, error) {
	switch nama {
	case usecase.JobWilayah:
		return buatImporWilayah(pool)
	default:
		return nil, fmt.Errorf("job tidak dikenal: %s", nama)
	}
}

func buatImporWilayah(pool *pgxpool.Pool) (job, error) {
	var cfg wilayah.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiWilayah, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporWilayah{
		Pengambil: wilayah.Sumber{
			Konfigurasi: cfg.Sumber,
			Klien:       &http.Client{Timeout: batasWaktuUnduh},
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Penyimpan:   repository.WilayahRepository{DB: pool},
		Pencatat:    repository.IngestRunRepository{DB: pool},
		AturanAlias: cfg.AturanAlias(),
		Provinsi:    cfg.Provinsi,
	}
	return impor.Jalankan, nil
}

func dirRaw() string {
	if dir, ada := os.LookupEnv("RANTAI_DATA_RAW"); ada {
		return dir
	}
	return dirRawBawaan
}
