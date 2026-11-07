package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	deliveryhttp "github.com/benditandayusaputra/rantai/backend/internal/delivery/http"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/kartu"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiBawaan       = "config/api.yaml"
	pathKonfigurasiArusAnggaran = "config/arus_anggaran.yaml"
	pathKonfigurasiPemeliharaan = "config/pemeliharaan.yaml"
	pathKonfigurasiSinyal       = "config/sinyal.yaml"
	pathKonfigurasiOpentender   = "config/opentender.yaml"
	pathKonfigurasiMedia        = "config/media.yaml"
	pathKonfigurasiKartu        = "config/kartu.yaml"
	pathKonfigurasiSitus        = "config/situs.yaml"
)

func main() {
	if err := jalankan(); err != nil {
		log.Fatal(err)
	}
}

func jalankan() error {
	cfg, err := config.Load(pathKonfigurasi())
	if err != nil {
		return err
	}
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("sambung basis data: %w", err)
	}
	defer pool.Close()
	layanan, err := susunLayanan(pool)
	if err != nil {
		return err
	}
	layanan = denganLayananPendidikan(layanan, pool, cfg)
	if layanan, err = denganSumberDasborDanKartu(layanan, pool); err != nil {
		return err
	}
	if layanan, err = denganTanyaAuthLaporan(ctx, layanan, pool, cfg.Server.Auth); err != nil {
		return err
	}
	app := deliveryhttp.NewApp(cfg.Server, layanan)
	go hentikanSaatSinyal(ctx, app, cfg.Server)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		return fmt.Errorf("jalankan server: %w", err)
	}
	return nil
}

func denganTanyaAuthLaporan(ctx context.Context, layanan deliveryhttp.Layanan, pool *pgxpool.Pool, cfg config.ServerAuth) (deliveryhttp.Layanan, error) {
	layanan, err := denganTanya(layanan)
	if err != nil {
		return layanan, err
	}
	if layanan, err = denganAuth(layanan, pool, cfg); err != nil {
		return layanan, err
	}
	if layanan, err = denganTemuan(layanan, pool); err != nil {
		return layanan, err
	}
	return denganLaporan(ctx, layanan, pool)
}

func susunLayanan(pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}}
	wilayah := repository.WilayahRepository{DB: pool}
	wakil := usecase.BacaWakil{Pembaca: repository.WakilRepository{DB: pool}, Tema: tema, Umpan: repository.UmpanRepository{DB: pool}}
	ringkasan := repository.RingkasanKecamatanRepository{DB: pool}
	paket, err := buatBacaPaket(pool)
	if err != nil {
		return deliveryhttp.Layanan{}, err
	}
	var sinyal domain.AturanSinyal
	if err := config.BacaYAML(pathKonfigurasiSinyal, &sinyal); err != nil {
		return deliveryhttp.Layanan{}, err
	}
	penyedia, err := buatBacaPenyedia(pool, paket.Aturan.BatasDaftarPenyedia)
	if err != nil {
		return deliveryhttp.Layanan{}, err
	}
	return denganLayananAksi(deliveryhttp.Layanan{
		Wilayah:    usecase.CariWilayah{Pembaca: wilayah},
		Provenance: usecase.AmbilProvenance{Pembaca: repository.ProvenanceRepository{DB: pool}},
		Tema:       tema,
		Wakil:      wakil,
		Umpan:      usecase.BacaUmpan{Pembaca: repository.UmpanRepository{DB: pool}},
		Ringkasan: usecase.RingkasKecamatan{
			Wilayah: wilayah, Tema: tema, Indikator: ringkasan, Anggaran: ringkasan, Wakil: wakil, Sekarang: time.Now,
		},
		Paket:    paket,
		Penyedia: penyedia, Jaringan: usecase.BacaJaringan{Pembaca: repository.JaringanRepository{DB: pool}},
		Pembanding: usecase.BandingPaket{Pembaca: repository.PaketRepository{DB: pool}, MinPembanding: sinyal.MinPembanding},
		Belanja:    usecase.BandingBelanja{Tema: tema, Pembaca: repository.BelanjaTemaRepository{DB: pool}, Sekarang: time.Now},
	}, pool)
}

func buatBacaDasbor(pool *pgxpool.Pool) (usecase.BacaDasbor, error) {
	var cfg struct {
		InstansiTanpaRincian string `yaml:"instansi_tanpa_rincian"`
	}
	if err := config.BacaYAML(pathKonfigurasiArusAnggaran, &cfg); err != nil {
		return usecase.BacaDasbor{}, err
	}
	return usecase.BacaDasbor{
		Tema:    usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}},
		Pembaca: repository.DasborRepository{DB: pool}, Sekarang: time.Now, InstansiTanpaRincian: cfg.InstansiTanpaRincian,
	}, nil
}

func bacaAtribusiOpentender() (string, error) {
	var cfg struct {
		Atribusi string `yaml:"atribusi"`
	}
	err := config.BacaYAML(pathKonfigurasiOpentender, &cfg)
	return cfg.Atribusi, err
}

func buatBacaPenyedia(pool *pgxpool.Pool, batas int) (usecase.BacaPenyedia, error) {
	atribusi, err := bacaAtribusiOpentender()
	if err != nil {
		return usecase.BacaPenyedia{}, err
	}
	return usecase.BacaPenyedia{Pembaca: repository.PenyediaRepository{DB: pool}, Batas: batas, Atribusi: atribusi}, nil
}

func denganSumberDasborDanKartu(layanan deliveryhttp.Layanan, pool *pgxpool.Pool) (deliveryhttp.Layanan, error) {
	var err error
	if layanan.SumberData, err = buatBacaSumberData(pool); err != nil {
		return layanan, err
	}
	if layanan.Dasbor, err = buatBacaDasbor(pool); err != nil {
		return layanan, err
	}
	layanan.Kartu, err = buatPembuatKartu()
	return layanan, err
}

func buatPembuatKartu() (kartu.Pembuat, error) {
	var k kartu.Konfigurasi
	var situs domain.Situs
	if err := config.BacaYAML(pathKonfigurasiKartu, &k); err != nil {
		return kartu.Pembuat{}, err
	}
	if err := config.BacaYAML(pathKonfigurasiSitus, &situs); err != nil {
		return kartu.Pembuat{}, err
	}
	k.AlamatSitus = situs.Alamat
	return kartu.Baru(k)
}

func buatBacaSumberData(pool *pgxpool.Pool) (usecase.BacaSumberData, error) {
	var cfg struct {
		Media []domain.MediaTerverifikasi `yaml:"media"`
	}
	if err := config.BacaYAML(pathKonfigurasiMedia, &cfg); err != nil {
		return usecase.BacaSumberData{}, err
	}
	atribusi, err := bacaAtribusiOpentender()
	if err != nil {
		return usecase.BacaSumberData{}, err
	}
	return usecase.BacaSumberData{Pembaca: repository.IngestRunRepository{DB: pool}, Media: cfg.Media, AtribusiOpentender: atribusi}, nil
}

func buatBacaPaket(pool *pgxpool.Pool) (usecase.BacaPaket, error) {
	u := usecase.BacaPaket{Pembaca: repository.PaketRepository{DB: pool}, Sekarang: time.Now}
	if err := config.BacaYAML(pathKonfigurasiArusAnggaran, &u.Aturan); err != nil {
		return usecase.BacaPaket{}, err
	}
	if err := config.BacaYAML(pathKonfigurasiPemeliharaan, &u.Pemeliharaan); err != nil {
		return usecase.BacaPaket{}, err
	}
	return u, nil
}

func pathKonfigurasi() string {
	if path, ada := os.LookupEnv("RANTAI_CONFIG"); ada {
		return path
	}
	return pathKonfigurasiBawaan
}

func hentikanSaatSinyal(ctx context.Context, app *fiber.App, cfg config.Server) {
	<-ctx.Done()
	ctxHenti, batal := context.WithTimeout(context.Background(), cfg.BatasWaktuHenti)
	defer batal()
	if err := app.ShutdownWithContext(ctxHenti); err != nil {
		log.Printf("hentikan server: %v", err)
	}
}
