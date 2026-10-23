package ujiwilayah

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pengambilFixture struct {
	akar string
}

func MuatDKI(t *testing.T, tx pgx.Tx) {
	t.Helper()
	akar := akarBackend()
	var cfg wilayah.Konfigurasi
	if err := config.BacaYAML(filepath.Join(akar, "config/wilayah.yaml"), &cfg); err != nil {
		t.Fatal(err)
	}
	impor := usecase.ImporWilayah{
		Pengambil:   pengambilFixture{akar: akar},
		Penyimpan:   repository.WilayahRepository{DB: tx},
		Pencatat:    repository.IngestRunRepository{DB: tx},
		AturanAlias: cfg.AturanAlias(),
		Provinsi:    cfg.Provinsi,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func (p pengambilFixture) Ambil(context.Context) (usecase.DataSumberWilayah, error) {
	berkas, err := filepath.Glob(filepath.Join(p.akar, "internal/ingest/wilayah/testdata/*.sql"))
	if err != nil {
		return usecase.DataSumberWilayah{}, err
	}
	data := usecase.DataSumberWilayah{Provenance: domain.Provenance{
		ID: provenance.IDBaru(), Sumber: domain.SumberKemendagri, Label: domain.LabelResmi,
	}}
	for _, path := range berkas {
		isi, err := os.ReadFile(path)
		if err != nil {
			return usecase.DataSumberWilayah{}, err
		}
		daftar, err := wilayah.UraiSQL(isi)
		if err != nil {
			return usecase.DataSumberWilayah{}, err
		}
		data.Wilayah = append(data.Wilayah, daftar...)
	}
	return data, nil
}

func akarBackend() string {
	_, berkas, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(berkas), "../../..")
}
