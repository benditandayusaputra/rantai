package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/lpse"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/sirup"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiWilayah = "config/wilayah.yaml"
	dirKonfigurasiTema     = "config/tema"
	pathKonfigurasiDPRD    = "config/dprd.yaml"
	pathKonfigurasiLPSE    = "config/lpse.yaml"
	pathKonfigurasiSIRUP   = "config/sirup.yaml"
	urlRujukanTema         = "https://github.com/benditandayusaputra/rantai/tree/main/backend/config/tema"
	dirRawBawaan           = "../data/raw"
	batasWaktuUnduh        = time.Minute
)

type job func(ctx context.Context) error

var jobTersedia = []string{usecase.JobWilayah, usecase.JobTema, usecase.JobKPUDPRD, usecase.JobLPSEDaftar, usecase.JobLPSEDetail, usecase.JobSIRUP, usecase.JobUmpanDPRD}

func main() {
	if err := jalankan(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func jalankan(argumen []string) error {
	if len(argumen) != 1 {
		return errors.New("pemakaian: worker <job>; job tersedia: " + strings.Join(jobTersedia, ", "))
	}
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("sambung basis data: %w", err)
	}
	defer pool.Close()
	jalankanJob, err := pilihJob(argumen[0], pool)
	if err != nil {
		return err
	}
	return jalankanJob(ctx)
}

func pilihJob(nama string, pool *pgxpool.Pool) (job, error) {
	switch nama {
	case usecase.JobWilayah:
		return buatImporWilayah(pool)
	case usecase.JobTema:
		return buatImporTema(pool), nil
	case usecase.JobKPUDPRD:
		return buatImporDPRD(pool)
	case usecase.JobLPSEDaftar, usecase.JobLPSEDetail:
		return buatImporLPSE(nama, pool)
	case usecase.JobSIRUP:
		return buatImporRUP(pool)
	case usecase.JobUmpanDPRD:
		return buatImporUmpanDPRD(pool)
	default:
		return nil, fmt.Errorf("job tidak dikenal: %s", nama)
	}
}

func buatImporWilayah(pool *pgxpool.Pool) (job, error) {
	var cfg wilayah.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiWilayah, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporWilayah{
		Pengambil: wilayah.Sumber{
			Konfigurasi: cfg.Sumber,
			Klien:       &http.Client{Timeout: batasWaktuUnduh},
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Penyimpan:   repository.WilayahRepository{DB: pool},
		Pencatat:    repository.IngestRunRepository{DB: pool},
		AturanAlias: cfg.AturanAlias(),
		Provinsi:    cfg.Provinsi,
	}
	return impor.Jalankan, nil
}

func buatImporTema(pool *pgxpool.Pool) job {
	impor := usecase.ImporTema{
		Pengambil: tema.Sumber{Dir: dirKonfigurasiTema, URLRujukan: urlRujukanTema, Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
	}
	return impor.Jalankan
}

func buatImporDPRD(pool *pgxpool.Pool) (job, error) {
	var cfg dprd.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiDPRD, &cfg); err != nil {
		return nil, err
	}
	wakil := repository.WakilRepository{DB: pool}
	impor := usecase.ImporDPRD{
		Pengambil: dprd.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Pengunduh),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		PembacaKecamatan: wakil,
		Penyimpan:        wakil,
		Pencatat:         repository.IngestRunRepository{DB: pool},
		Dapil:            cfg.DapilKPU.Daftar,
		TahunPemilu:      cfg.DapilKPU.TahunPemilu,
		URLLHKPN:         cfg.URLLHKPN,
		PolaURLProfil:    cfg.Sumber.URLDasar + cfg.Sumber.HalamanProfil,
		Laporkan:         func(pesan string) { log.Println(pesan) },
	}
	return impor.Jalankan, nil
}

func buatImporLPSE(nama string, pool *pgxpool.Pool) (job, error) {
	var cfg lpse.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiLPSE, &cfg); err != nil {
		return nil, err
	}
	toples, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	sumber := lpse.Sumber{
		Klien:    klien.Baru(&http.Client{Timeout: batasWaktuUnduh, Jar: toples}, cfg.Klien),
		URLDasar: cfg.URLDasar,
		Ukuran:   cfg.UkuranHalaman,
		DirRaw:   dirRaw(),
		Sekarang: time.Now,
	}
	paket := repository.PaketRepository{DB: pool}
	pencatat := repository.IngestRunRepository{DB: pool}
	if nama == usecase.JobLPSEDaftar {
		return usecase.ImporDaftarPaket{Pengambil: sumber, Penyimpan: paket, Pencatat: pencatat, Tahun: cfg.TahunAnggaran}.Jalankan, nil
	}
	return usecase.ImporDetailPaket{
		Pengambil: sumber, Penyimpan: paket, Pencatat: pencatat,
		AturanInstansi: cfg.Instansi, JenisPengadaan: cfg.JenisPengadaanDetail,
	}.Jalankan, nil
}

func buatImporUmpanDPRD(pool *pgxpool.Pool) (job, error) {
	var cfg dprd.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiDPRD, &cfg); err != nil {
		return nil, err
	}
	umpan := repository.UmpanRepository{DB: pool}
	impor := usecase.ImporUmpan{
		Job:           usecase.JobUmpanDPRD,
		JenisSumber:   domain.SumberUmpanResmiDPRD,
		MulaiBawaan:   cfg.Rilis.Mulai,
		TumpangTindih: cfg.Rilis.Tumpang,
		Pengambil: dprd.SumberRilis{
			Konfigurasi: cfg.Rilis,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Pengunduh),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Anggota:   umpan,
		Tema:      usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}},
		Penyimpan: umpan,
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Pautkan:   domain.PencariNama.PautkanRilisResmi,
	}
	return impor.Jalankan, nil
}

func buatImporRUP(pool *pgxpool.Pool) (job, error) {
	var cfg sirup.Konfigurasi
	var cfgLPSE lpse.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiSIRUP, &cfg); err != nil {
		return nil, err
	}
	if err := config.BacaYAML(pathKonfigurasiLPSE, &cfgLPSE); err != nil {
		return nil, err
	}
	impor := usecase.ImporRUP{
		Pengambil: sirup.Sumber{
			Klien:     klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Klien),
			URLDetail: cfg.URLDetail,
			DirRaw:    dirRaw(),
			Sekarang:  time.Now,
		},
		Penyimpan:      repository.RUPRepository{DB: pool},
		Pencatat:       repository.IngestRunRepository{DB: pool},
		AturanInstansi: cfgLPSE.Instansi,
	}
	return impor.Jalankan, nil
}

func dirRaw() string {
	if dir, ada := os.LookupEnv("RANTAI_DATA_RAW"); ada {
		return dir
	}
	return dirRawBawaan
}
