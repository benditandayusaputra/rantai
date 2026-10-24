package main

import (
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiAspirasi      = "config/aspirasi.yaml"
	pathKonfigurasiJalurAspirasi = "config/jalur_aspirasi.yaml"
)

func denganAspirasi(layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	var kosakata domain.KosakataAspirasi
	if err := config.BacaYAML(pathKonfigurasiAspirasi, &kosakata); err != nil {
		return layanan, err
	}
	var jalur domain.KonfigurasiJalurAspirasi
	if err := config.BacaYAML(pathKonfigurasiJalurAspirasi, &jalur); err != nil {
		return layanan, err
	}
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}}
	wilayah := usecase.CariWilayah{Pembaca: repository.WilayahRepository{DB: pool}}
	layanan.Aspirasi = usecase.SusunAspirasi{
		Tema:       tema,
		Wilayah:    wilayah,
		Belanja:    usecase.BandingBelanja{Tema: tema, Pembaca: repository.BelanjaTemaRepository{DB: pool}, Sekarang: time.Now},
		Wakil:      usecase.BacaWakil{Pembaca: repository.WakilRepository{DB: pool}, Tema: tema},
		Provenance: usecase.AmbilProvenance{Pembaca: repository.ProvenanceRepository{DB: pool}},
		Jalur:      jalur,
		Kosakata:   kosakata,
	}
	layanan.Jalur = usecase.BacaJalurAspirasi{Wilayah: wilayah, Jalur: jalur}
	return layanan, nil
}
