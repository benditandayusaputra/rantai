package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dapodik"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dayatampung"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiDapodik = "config/dapodik.yaml"
	pathDayaTampung        = "config/daya_tampung.yaml"
)

func pilihJobPendidikan(nama string, pool *pgxpool.Pool) (job, error) {
	switch nama {
	case usecase.JobDapodik:
		return buatImporSekolah(pool)
	case usecase.JobDayaTampung:
		return buatImporDayaTampung(pool)
	default:
		return nil, fmt.Errorf("job tidak dikenal: %s", nama)
	}
}

func buatImporSekolah(pool *pgxpool.Pool) (job, error) {
	var cfg dapodik.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiDapodik, &cfg); err != nil {
		return nil, err
	}
	repo := repository.SekolahRepository{DB: pool}
	impor := usecase.ImporSekolah{
		Pengambil: dapodik.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Pengunduh),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Pemeta:    repo,
		Penyimpan: repo,
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Pelapor:   pelaporLog{},
	}
	return impor.Jalankan, nil
}

func buatImporDayaTampung(pool *pgxpool.Pool) (job, error) {
	var cfg dayatampung.Konfigurasi
	if err := config.BacaYAML(pathDayaTampung, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporDayaTampung{
		Pengambil: dayatampung.Sumber{Konfigurasi: cfg, Sekarang: time.Now},
		Penyimpan: repository.DayaTampungRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Pelapor:   pelaporLog{},
	}
	return impor.Jalankan, nil
}
