package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/llm"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

const pathKonfigurasiTanya = "config/tanya.yaml"

type konfigurasiTanya struct {
	tanya.Aturan `yaml:",inline"`
	MaksToken    int64  `yaml:"maks_token"`
	Effort       string `yaml:"effort"`
}

func denganTanya(layanan deliveryhttp.Layanan) (deliveryhttp.Layanan, error) {
	var cfg konfigurasiTanya
	if err := config.BacaYAML(pathKonfigurasiTanya, &cfg); err != nil {
		return layanan, err
	}
	klien, err := llm.Baru(llm.Konfigurasi{
		Model: os.Getenv("TANYA_MODEL"), MaksToken: cfg.MaksToken, Effort: anthropic.OutputConfigEffort(cfg.Effort),
	})
	if errors.Is(err, llm.ErrModelBelumDiatur) {
		log.Println("TANYA_MODEL kosong, Tanya Rantai tidak dipasang")
		return layanan, nil
	}
	if err != nil {
		return layanan, fmt.Errorf("TANYA_MODEL: %w", err)
	}
	registry, err := deliveryhttp.RegistryTanya(layanan)
	if err != nil {
		return layanan, err
	}
	layanan.Tanya = tanya.Penjawab{Model: klien, Alat: registry, Aturan: cfg.Aturan}
	return layanan, nil
}
