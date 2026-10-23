package main

import (
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathSkenarioDemo = "../data/snapshot/skenario_demo.yaml"

func pilihJobDemo(nama string, pool *pgxpool.Pool) (job, error) {
	if nama != usecase.JobSkenarioDemo {
		return nil, fmt.Errorf("job tidak dikenal: %s", nama)
	}
	return buatMuatSkenarioDemo(pool)
}

func buatMuatSkenarioDemo(pool *pgxpool.Pool) (job, error) {
	var skenario domain.SkenarioDemo
	if err := config.BacaYAML(pathSkenarioDemo, &skenario); err != nil {
		return nil, err
	}
	muat := usecase.MuatSkenarioDemo{
		Skenario:       skenario,
		Penyimpan:      repository.SkenarioDemoRepository{DB: pool},
		Pelapor:        pelaporLog{},
		BuatProvenance: buatProvenanceTurunan,
	}
	return muat.Jalankan, nil
}
