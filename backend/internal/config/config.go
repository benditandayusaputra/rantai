package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      Server `yaml:"server"`
	DatabaseURL string `yaml:"-"`
}

type Server struct {
	Port            string        `yaml:"port"`
	BatasWaktuBaca  time.Duration `yaml:"batas_waktu_baca"`
	BatasWaktuTulis time.Duration `yaml:"batas_waktu_tulis"`
	BatasWaktuHenti time.Duration `yaml:"batas_waktu_henti"`
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
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	return cfg
}
