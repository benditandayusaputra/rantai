package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/llm"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	pathKonfigurasiTanya  = "config/tanya.yaml"
	pathKosakataTerlarang = "config/kosakata_terlarang.yaml"
	pathCacheTanyaDemo    = "../data/snapshot/tanya_demo.json"
)

type konfigurasiTanya struct {
	tanya.Aturan `yaml:",inline"`
	MaksToken    int64  `yaml:"maks_token"`
	Effort       string `yaml:"effort"`
}

func denganTanya(layanan deliveryhttp.Layanan) (deliveryhttp.Layanan, error) {
	penjawab, err := buatPenjawab(layanan)
	if err != nil {
		return layanan, err
	}
	if os.Getenv("RANTAI_MODE") == usecase.ModeDemo {
		cache, err := bacaCacheTanyaDemo()
		if err != nil {
			return layanan, err
		}
		layanan.Tanya = tanya.PenjawabTercache{Cache: cache, Cadangan: penjawab}
		return layanan, nil
	}
	if penjawab.Model == nil {
		log.Println("TANYA_MODEL kosong, Tanya Rantai tidak dipasang")
		return layanan, nil
	}
	layanan.Tanya = penjawab
	return layanan, nil
}

func buatPenjawab(layanan deliveryhttp.Layanan) (tanya.Penjawab, error) {
	var cfg konfigurasiTanya
	var terlarang domain.KosakataTerlarang
	if err := errors.Join(config.BacaYAML(pathKonfigurasiTanya, &cfg), config.BacaYAML(pathKosakataTerlarang, &terlarang)); err != nil {
		return tanya.Penjawab{}, err
	}
	penjawab := tanya.Penjawab{
		Aturan: cfg.Aturan, Validator: tanya.ValidatorJawaban{Terlarang: terlarang}, Laporkan: func(pesan string) { log.Println(pesan) },
	}
	klien, err := llm.Baru(llm.Konfigurasi{
		Model: os.Getenv("TANYA_MODEL"), MaksToken: cfg.MaksToken, Effort: anthropic.OutputConfigEffort(cfg.Effort),
	})
	if errors.Is(err, llm.ErrModelBelumDiatur) {
		return penjawab, nil
	}
	if err != nil {
		return tanya.Penjawab{}, fmt.Errorf("TANYA_MODEL: %w", err)
	}
	penjawab.Model = klien
	penjawab.Alat, err = deliveryhttp.RegistryTanya(layanan)
	return penjawab, err
}

func bacaCacheTanyaDemo() ([]tanya.EntriCache, error) {
	isi, err := os.ReadFile(pathCacheTanyaDemo)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%s belum ada; buat dengan make tanya-demo di mesin yang datanya lengkap", pathCacheTanyaDemo)
		return []tanya.EntriCache{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baca cache tanya demo: %w", err)
	}
	cache := []tanya.EntriCache{}
	if err := json.Unmarshal(isi, &cache); err != nil {
		return nil, fmt.Errorf("urai %s: %w", pathCacheTanyaDemo, err)
	}
	return cache, nil
}
