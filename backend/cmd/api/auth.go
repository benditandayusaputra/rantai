package main

import (
	"log"
	"os"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/surel"
	"github.com/benditandayusaputra/rantai/backend/internal/token"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/benditandayusaputra/rantai/backend/internal/webpush"
	"github.com/jackc/pgx/v5/pgxpool"
)

func denganAuth(layanan deliveryhttp.Layanan, pool *pgxpool.Pool, cfg config.ServerAuth) (deliveryhttp.Layanan, error) {
	rahasia := os.Getenv("AUTH_RAHASIA")
	if rahasia == "" {
		log.Println("AUTH_RAHASIA kosong, masuk dengan OTP tidak dipasang")
		return layanan, nil
	}
	akses, err := token.Baru(rahasia, cfg.UmurAkses)
	if err != nil {
		return layanan, err
	}
	pengirim, ada := pengirimOTP(cfg)
	if !ada {
		log.Println("SMTP_HOST kosong di luar mode demo, masuk dengan OTP tidak dipasang")
		return layanan, nil
	}
	akun := repository.PenggunaRepository{DB: pool}
	layanan.Akses = akses
	layanan.Pantau = usecase.Pantau{Repo: repository.PantauRepository{DB: pool}}
	if pengirim, ada := webpush.DariEnv(nil); ada {
		layanan.WebPush = usecase.LayananPush{Repo: repository.PantauRepository{DB: pool}, KunciPublik: pengirim.KunciPublik}
	}
	layanan.Masuk = usecase.Masuk{
		OTP: akun, Pengguna: akun, Sesi: akun, Pengirim: pengirim, Akses: akses,
		Aturan:   usecase.AturanMasuk{OTP: cfg.OTP, UmurRefresh: cfg.UmurRefresh},
		Rahasia:  []byte(rahasia),
		Sekarang: time.Now,
	}
	return denganPanel(layanan, pool)
}

func pengirimOTP(cfg config.ServerAuth) (usecase.PengirimOTP, bool) {
	if host := os.Getenv("SMTP_HOST"); host != "" {
		return surel.SMTP{
			Host: host, Port: os.Getenv("SMTP_PORT"), Pengguna: os.Getenv("SMTP_PENGGUNA"),
			Sandi: os.Getenv("SMTP_SANDI"), Dari: os.Getenv("SMTP_DARI"), UmurKode: cfg.OTP.Umur,
		}, true
	}
	if os.Getenv("RANTAI_MODE") == usecase.ModeDemo {
		return surel.PengirimLog{}, true
	}
	return nil, false
}
