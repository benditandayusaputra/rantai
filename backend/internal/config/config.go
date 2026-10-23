package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      Server `yaml:"server"`
	Peta        Peta   `yaml:"peta"`
	DatabaseURL string `yaml:"-"`
}

type Peta struct {
	ToleransiPenyederhanaan float64 `yaml:"toleransi_penyederhanaan"`
	DesimalKoordinat        int     `yaml:"desimal_koordinat"`
}

type Server struct {
	Port            string        `yaml:"port"`
	BatasWaktuBaca  time.Duration `yaml:"batas_waktu_baca"`
	BatasWaktuTulis time.Duration `yaml:"batas_waktu_tulis"`
	BatasWaktuHenti time.Duration `yaml:"batas_waktu_henti"`
	BatasLaju       BatasLaju     `yaml:"batas_laju"`
	CORS            CORS          `yaml:"cors"`
}

type BatasLaju struct {
	MaksPermintaan int           `yaml:"maks_permintaan"`
	Jendela        time.Duration `yaml:"jendela"`
}

type CORS struct {
	Origin []string `yaml:"origin"`
}

func Load(path string) (Config, error) {
	var cfg Config
	if err := BacaYAML(path, &cfg); err != nil {
		return Config{}, err
	}
	return terapkanEnv(cfg), nil
}

func BacaYAML(path string, tujuan any) error {
	isi, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("baca konfigurasi %s: %w", path, err)
	}
	if err := yaml.Unmarshal(isi, tujuan); err != nil {
		return fmt.Errorf("urai konfigurasi %s: %w", path, err)
	}
	return nil
}

func terapkanEnv(cfg Config) Config {
	if port, ada := os.LookupEnv("API_PORT"); ada {
		cfg.Server.Port = port
	}
	if origin, ada := os.LookupEnv("CORS_ORIGIN"); ada {
		cfg.Server.CORS.Origin = strings.Split(origin, ",")
	}
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	return cfg
}
