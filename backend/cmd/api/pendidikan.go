package main

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func denganLayananPendidikan(layanan deliveryhttp.Layanan, pool *pgxpool.Pool, cfg config.Config) deliveryhttp.Layanan {
	layanan.PendidikanKota = usecase.BacaPendidikanKota{
		Pembaca: repository.PendidikanKotaRepository{DB: pool},
		Wilayah: repository.WilayahRepository{DB: pool},
	}
	layanan.Indikator = usecase.BacaIndikator{
		Pembaca: repository.IndikatorBacaRepository{DB: pool, Peta: cfg.Peta},
		Tema:    repository.TemaRepository{DB: pool},
	}
	layanan.Sekolah = usecase.CariSekolahTerdekat{
		Pembaca:     repository.SekolahTerdekatRepository{DB: pool},
		Konfigurasi: usecase.KonfigurasiSekolahTerdekat(cfg.SekolahTerdekat),
	}
	return layanan
}
