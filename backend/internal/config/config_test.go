package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tulisKonfigurasi(t *testing.T, isi string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "api.yaml")
	if err := os.WriteFile(path, []byte(isi), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoad_BerkasYAML_TerbacaLengkap(t *testing.T) {
	path := tulisKonfigurasi(t, "server:\n  port: \"9000\"\n  batas_waktu_baca: 5s\n")

	cfg, err := Load(path)

	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "9000" || cfg.Server.BatasWaktuBaca != 5*time.Second {
		t.Errorf("konfigurasi tidak sesuai: %+v", cfg.Server)
	}
}

func TestLoad_EnvTersedia_MenimpaYAML(t *testing.T) {
	path := tulisKonfigurasi(t, "server:\n  port: \"9000\"\n")
	t.Setenv("API_PORT", "7000")
	t.Setenv("DATABASE_URL", "postgres://uji")

	cfg, err := Load(path)

	if err != nil {
		t.Fatal(err)
	}
	if cfg.Server.Port != "7000" || cfg.DatabaseURL != "postgres://uji" {
		t.Errorf("env tidak diterapkan: %+v", cfg)
	}
}

func TestLoad_BerkasTidakAda_MengembalikanError(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "tidak-ada.yaml"))

	if err == nil {
		t.Error("diharapkan error untuk berkas yang tidak ada")
	}
}
