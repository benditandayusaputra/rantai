package main

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"time"

	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/foto"
	"github.com/benditandayusaputra/rantai/backend/internal/objek"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

func denganLaporan(ctx context.Context, layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	penyimpan, err := objek.BaruDariEnv(ctx)
	if errors.Is(err, objek.ErrKredensialObjekKosong) {
		log.Println("kredensial MinIO kosong, Lapor Bukti tidak dipasang")
		return layanan, nil
	}
	if err != nil {
		return layanan, err
	}
	layanan.Laporan = usecase.TerimaLaporan{
		Pemroses: foto.Pemroses{}, Penyimpan: penyimpan, Repo: repository.LaporanRepository{DB: pool},
		Sekarang: time.Now, BuatKunci: rand.Text,
	}
	return layanan, nil
}
