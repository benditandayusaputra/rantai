package main

import (
	"net/http"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dapodik"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dayatampung"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wikimedia"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiDapodik = "config/dapodik.yaml"
	pathDayaTampung        = "config/daya_tampung.yaml"
	pathPendidikan         = "config/pendidikan.yaml"
	pathFotoWikimedia      = "config/foto_wikimedia.yaml"
)

func pilihJobPendidikan(nama string, pool *pgxpool.Pool) (job, error) {
	switch nama {
	case usecase.JobDapodik:
		return buatImporSekolah(pool)
	case usecase.JobDayaTampung:
		return buatImporDayaTampung(pool)
	case usecase.JobCalonLulusan:
		return buatHitungCalonLulusan(pool)
	case usecase.JobIndikator:
		return buatHitungIndikator(pool)
	case usecase.JobFotoWikimedia:
		return buatImporFotoWikimedia(pool)
	default:
		return pilihJobDemo(nama, pool)
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

func buatImporFotoWikimedia(pool *pgxpool.Pool) (job, error) {
	var cfg wikimedia.Konfigurasi
	if err := config.BacaYAML(pathFotoWikimedia, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporFotoWikimedia{
		Pengambil: wikimedia.Sumber{Konfigurasi: cfg},
		Penyimpan: repository.FotoSekolahRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
	}
	return impor.Jalankan, nil
}

func buatHitungCalonLulusan(pool *pgxpool.Pool) (job, error) {
	var konstanta usecase.KonstantaPendidikan
	if err := config.BacaYAML(pathPendidikan, &konstanta); err != nil {
		return nil, err
	}
	repo := repository.CalonLulusanRepository{DB: pool}
	hitung := usecase.HitungCalonLulusan{
		Pembaca: repo, Penyimpan: repo, Pencatat: repository.IngestRunRepository{DB: pool},
		Pelapor: pelaporLog{}, Konstanta: konstanta, BuatProvenance: buatProvenanceTurunan,
	}
	return hitung.Jalankan, nil
}

func buatHitungIndikator(pool *pgxpool.Pool) (job, error) {
	var konstanta usecase.KonstantaPendidikan
	if err := config.BacaYAML(pathPendidikan, &konstanta); err != nil {
		return nil, err
	}
	repo := repository.IndikatorRepository{DB: pool}
	hitung := usecase.HitungIndikator{
		Tema: repository.TemaRepository{DB: pool},
		Penghitung: map[string]usecase.PenghitungIndikator{
			domain.PenghitungKesenjanganKursiSMP: usecase.KesenjanganKursiSMP{Pembaca: repo, Konstanta: konstanta, BuatProvenance: buatProvenanceTurunan},
		},
		Penyimpan: repo,
		Pencatat:  repository.IngestRunRepository{DB: pool},
	}
	return hitung.Jalankan, nil
}

func buatProvenanceTurunan(catatan usecase.CatatanProvenance) domain.Provenance {
	return provenance.Baru(provenance.Catatan(catatan), []byte(catatan.Catatan), time.Now())
}
