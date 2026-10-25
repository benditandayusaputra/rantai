package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/surel"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pathKonfigurasiSitus = "config/situs.yaml"

var errTanpaPengirimEmail = errors.New("SMTP_HOST kosong di luar mode demo, ringkasan mingguan tidak bisa dikirim")

func pilihJobPantau(nama string, pool *pgxpool.Pool) (job, error) {
	if nama != usecase.JobRingkasanMingguan {
		return nil, fmt.Errorf("job tidak dikenal: %s", nama)
	}
	var situs domain.Situs
	if err := config.BacaYAML(pathKonfigurasiSitus, &situs); err != nil {
		return nil, err
	}
	pengirim, err := pengirimEmail()
	if err != nil {
		return nil, err
	}
	kirim := usecase.KirimRingkasanMingguan{
		Repo: repository.PantauRepository{DB: pool}, Pengirim: pengirim, AlamatSitus: situs.Alamat, Sekarang: time.Now,
	}
	return func(ctx context.Context) error {
		terkirim, err := kirim.Jalankan(ctx)
		log.Printf("ringkasan mingguan terkirim ke %d pengguna", terkirim)
		return err
	}, nil
}

func pengirimEmail() (usecase.PengirimEmail, error) {
	if smtp, ada := surel.DariEnv(0); ada {
		return smtp, nil
	}
	if os.Getenv("RANTAI_MODE") == usecase.ModeDemo {
		return surel.PengirimLog{}, nil
	}
	return nil, errTanpaPengirimEmail
}
