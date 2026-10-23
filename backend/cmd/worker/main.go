package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiWilayah = "config/wilayah.yaml"
	dirKonfigurasiTema     = "config/tema"
	pathKonfigurasiDPRD    = "config/dprd.yaml"
	urlRujukanTema         = "https://github.com/benditandayusaputra/rantai/tree/main/backend/config/tema"
	dirRawBawaan           = "../data/raw"
	batasWaktuUnduh        = time.Minute
)

type job func(ctx context.Context) error

var jobTersedia = []string{usecase.JobWilayah, usecase.JobTema, usecase.JobKPUDPRD}

func main() {
	if err := jalankan(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func jalankan(argumen []string) error {
	if len(argumen) != 1 {
		return errors.New("pemakaian: worker <job>; job tersedia: " + strings.Join(jobTersedia, ", "))
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
	case usecase.JobTema:
		return buatImporTema(pool), nil
	case usecase.JobKPUDPRD:
		return buatImporDPRD(pool)
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

func buatImporTema(pool *pgxpool.Pool) job {
	impor := usecase.ImporTema{
		Pengambil: tema.Sumber{Dir: dirKonfigurasiTema, URLRujukan: urlRujukanTema, Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
	}
	return impor.Jalankan
}

func buatImporDPRD(pool *pgxpool.Pool) (job, error) {
	var cfg dprd.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiDPRD, &cfg); err != nil {
		return nil, err
	}
	wakil := repository.WakilRepository{DB: pool}
	impor := usecase.ImporDPRD{
		Pengambil: dprd.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Pengunduh),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		PembacaKecamatan: wakil,
		Penyimpan:        wakil,
		Pencatat:         repository.IngestRunRepository{DB: pool},
		Dapil:            cfg.DapilKPU.Daftar,
		TahunPemilu:      cfg.DapilKPU.TahunPemilu,
		URLLHKPN:         cfg.URLLHKPN,
		PolaURLProfil:    cfg.Sumber.URLDasar + cfg.Sumber.HalamanProfil,
	}
	return impor.Jalankan, nil
}

func dirRaw() string {
	if dir, ada := os.LookupEnv("RANTAI_DATA_RAW"); ada {
		return dir
	}
	return dirRawBawaan
}
