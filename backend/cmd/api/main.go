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
	"github.com/gofiber/fiber/v3"
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
	app := deliveryhttp.NewApp(cfg.Server)
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	go hentikanSaatSinyal(ctx, app, cfg.Server)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		return fmt.Errorf("jalankan server: %w", err)
	}
	return nil
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
