package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathKonfigurasiBawaan = "config/api.yaml"

func main() {
	if err := jalankan(); err != nil {
		log.Fatal(err)
	}
}

func jalankan() error {
	cfg, err := config.Load(pathKonfigurasi())
	if err != nil {
		return err
	}
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("sambung basis data: %w", err)
	}
	defer pool.Close()
	app := deliveryhttp.NewApp(cfg.Server, susunLayanan(pool))
	go hentikanSaatSinyal(ctx, app, cfg.Server)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		return fmt.Errorf("jalankan server: %w", err)
	}
	return nil
}

func susunLayanan(pool *pgxpool.Pool) deliveryhttp.Layanan {
	return deliveryhttp.Layanan{
		Wilayah:    usecase.CariWilayah{Pembaca: repository.WilayahRepository{DB: pool}},
		Provenance: usecase.AmbilProvenance{Pembaca: repository.ProvenanceRepository{DB: pool}},
		Tema:       usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}},
	}
}

func pathKonfigurasi() string {
	if path, ada := os.LookupEnv("RANTAI_CONFIG"); ada {
		return path
	}
	return pathKonfigurasiBawaan
}

func hentikanSaatSinyal(ctx context.Context, app *fiber.App, cfg config.Server) {
	<-ctx.Done()
	ctxHenti, batal := context.WithTimeout(context.Background(), cfg.BatasWaktuHenti)
	defer batal()
	if err := app.ShutdownWithContext(ctxHenti); err != nil {
		log.Printf("hentikan server: %v", err)
	}
}
