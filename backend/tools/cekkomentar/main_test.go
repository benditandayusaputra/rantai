package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCariKomentar_IsiBerkas_HanyaKomentarNonDirektifDilaporkan(t *testing.T) {
	kasus := []struct {
		nama   string
		isi    string
		jumlah int
	}{
		{"tanpa komentar", "package x\n\nvar s = \"http://contoh\"\n", 0},
		{"direktif go", "//go:build linux\n\npackage x\n", 0},
		{"direktif nolint", "package x\n\nvar a = 1 //nolint:gochecknoglobals\n", 0},
		{"komentar baris", "package x\n\n// penjelasan\nvar a = 1\n", 1},
		{"komentar blok", "package x\n\n/* blok */\nvar a = 1\n", 1},
		{"godoc paket", "// Package x ...\npackage x\n", 1},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			akar := t.TempDir()
			if err := os.WriteFile(filepath.Join(akar, "x.go"), []byte(k.isi), 0o600); err != nil {
				t.Fatal(err)
			}

			temuan, err := cariKomentar(akar)

			if err != nil {
				t.Fatal(err)
			}
			if len(temuan) != k.jumlah {
				t.Errorf("temuan = %v, harap %d", temuan, k.jumlah)
			}
		})
	}
}
