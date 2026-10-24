package main

import (
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func denganPanel(layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	var pemeliharaan usecase.KonfigurasiPemeliharaan
	if err := config.BacaYAML(pathKonfigurasiPemeliharaan, &pemeliharaan); err != nil {
		return layanan, err
	}
	layanan.Panel = usecase.Panel{
		Pengguna:       repository.PenggunaRepository{DB: pool},
		Penyimpan:      repository.PanelRepository{DB: pool},
		Laporan:        repository.LaporanRepository{DB: pool},
		Pemeliharaan:   pemeliharaan,
		BuatProvenance: buatProvenanceInstansi,
		Sekarang:       time.Now,
	}
	return layanan, nil
}

func buatProvenanceInstansi(catatan usecase.CatatanProvenance) domain.Provenance {
	return provenance.Baru(provenance.Catatan(catatan), []byte(catatan.Catatan), time.Now())
}
