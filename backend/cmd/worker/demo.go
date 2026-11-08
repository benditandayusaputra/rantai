package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/foto"
	"github.com/benditandayusaputra/rantai/backend/internal/objek"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathSkenarioDemo = "../data/snapshot/skenario_demo.yaml"

func pilihJobDemo(nama string, pool *pgxpool.Pool) (job, error) {
	switch nama {
	case usecase.JobSkenarioDemo:
		return buatMuatSkenarioDemo(pool)
	case usecase.JobParipurna:
		return buatImporParipurna(pool)
	default:
		return pilihJobAnggaran(nama, pool)
	}
}

func buatMuatSkenarioDemo(pool *pgxpool.Pool) (job, error) {
	var skenario domain.SkenarioDemo
	if err := config.BacaYAML(pathSkenarioDemo, &skenario); err != nil {
		return nil, err
	}
	penyimpanFoto, err := objek.BaruDariEnv(context.Background())
	if err != nil {
		return nil, err
	}
	muat := usecase.MuatSkenarioDemo{
		Skenario:       skenario,
		Penyimpan:      repository.SkenarioDemoRepository{DB: pool},
		Pelapor:        pelaporLog{},
		BuatProvenance: buatProvenanceTurunan,
		Pemroses:       foto.Pemroses{},
		PenyimpanFoto:  penyimpanFoto,
		BacaFoto: func(nama string) ([]byte, error) {
			return os.ReadFile(filepath.Join(filepath.Dir(pathSkenarioDemo), nama))
		},
	}
	return muat.Jalankan, nil
}
