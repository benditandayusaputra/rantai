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

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/lpse"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/media"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/opentender"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/sirup"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/wilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/llm"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pathKonfigurasiWilayah    = "config/wilayah.yaml"
	dirKonfigurasiTema        = "config/tema"
	pathKonfigurasiDPRD       = "config/dprd.yaml"
	pathKonfigurasiMedia      = "config/media.yaml"
	pathKonfigurasiRingkas    = "config/ringkasan.yaml"
	pathKosakataTerlarang     = "config/kosakata_terlarang.yaml"
	pathKonfigurasiLPSE       = "config/lpse.yaml"
	pathKonfigurasiSIRUP      = "config/sirup.yaml"
	pathKonfigurasiOpentender = "config/opentender.yaml"
	batasWaktuUnduhBesar      = 10 * time.Minute
	pathKonfigurasiGeocode    = "config/geocode.yaml"
	pathKonfigurasiKategori   = "config/kategori.yaml"
	pathKonfigurasiSinyal     = "config/sinyal.yaml"
	jobLaporanGeocode         = "laporan_geocode"
	jenisPengadaanLaporan     = "Pekerjaan Konstruksi"
	persenPenuh               = 100
	urlRujukanTema            = "https://github.com/benditandayusaputra/rantai/tree/main/backend/config/tema"
	dirRawBawaan              = "../data/raw"
	batasWaktuUnduh           = time.Minute
	perintahPeran             = "peran"
)

type job func(ctx context.Context) error

var jobTersedia = []string{usecase.JobWilayah, usecase.JobTema, usecase.JobKPUDPRD, usecase.JobLPSEDaftar, usecase.JobLPSEDetail, usecase.JobSIRUP, usecase.JobOpentender, usecase.JobTurunan, jobLaporanGeocode, usecase.JobUmpanDPRD, usecase.JobUmpanMedia, usecase.JobRingkasUmpan,
	usecase.JobDapodik, usecase.JobDayaTampung, usecase.JobCalonLulusan, usecase.JobIndikator, usecase.JobSkenarioDemo, usecase.JobParipurna}

func main() {
	if err := jalankan(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func jalankan(argumen []string) error {
	if !isArgumenSah(argumen) {
		return errors.New("pemakaian: worker <job> atau worker peran <email> moderator|admin; job tersedia: " + strings.Join(jobTersedia, ", "))
	}
	if err := usecase.IzinkanJob(os.Getenv("RANTAI_MODE"), argumen[0]); err != nil {
		return err
	}
	ctx, berhenti := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer berhenti()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("sambung basis data: %w", err)
	}
	defer pool.Close()
	jalankanJob, err := pilihPerintah(argumen, pool)
	if err != nil {
		return err
	}
	return jalankanJob(ctx)
}

func isArgumenSah(argumen []string) bool {
	return len(argumen) == 1 || (len(argumen) == 3 && argumen[0] == perintahPeran)
}

func pilihPerintah(argumen []string, pool *pgxpool.Pool) (job, error) {
	if argumen[0] != perintahPeran {
		return pilihJob(argumen[0], pool)
	}
	beri := usecase.BeriPeran{Penyimpan: repository.PenggunaRepository{DB: pool}}
	return func(ctx context.Context) error { return beri.Jalankan(ctx, argumen[1], domain.Peran(argumen[2])) }, nil
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
	case usecase.JobOpentender:
		return buatImporOpentender(pool)
	case usecase.JobTurunan:
		return buatTurunan(pool)
	case jobLaporanGeocode:
		return buatLaporanGeocode(pool), nil
	case usecase.JobUmpanDPRD:
		return buatImporUmpanDPRD(pool)
	case usecase.JobUmpanMedia:
		return buatImporUmpanMedia(pool)
	case usecase.JobRingkasUmpan:
		return buatRingkasUmpan(pool)
	default:
		return pilihJobPendidikan(nama, pool)
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

func buatImporUmpanMedia(pool *pgxpool.Pool) (job, error) {
	var cfg media.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiMedia, &cfg); err != nil {
		return nil, err
	}
	umpan := repository.UmpanRepository{DB: pool}
	impor := usecase.ImporUmpan{
		Job:           usecase.JobUmpanMedia,
		JenisSumber:   domain.SumberUmpanMedia,
		MulaiBawaan:   cfg.Mulai,
		TumpangTindih: cfg.TumpangTindih,
		Pengambil: media.Sumber{
			Konfigurasi: cfg,
			Pengunduh:   klien.Baru(&http.Client{Timeout: batasWaktuUnduh}, cfg.Pengunduh),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Anggota:   umpan,
		Tema:      usecase.BacaTema{Pembaca: repository.TemaRepository{DB: pool}},
		Penyimpan: umpan,
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Pautkan:   cfg.Aturan().Pautkan,
	}
	return impor.Jalankan, nil
}

func buatTurunan(pool *pgxpool.Pool) (job, error) {
	var aturan domain.AturanLokasi
	if err := config.BacaYAML(pathKonfigurasiGeocode, &aturan); err != nil {
		return nil, err
	}
	anggaran, err := bacaAturanTurunanAnggaran()
	if err != nil {
		return nil, err
	}
	lokasi := repository.LokasiPaketRepository{DB: pool}
	turunan := usecase.TurunanPaket{
		Gazetir:        lokasi,
		Lokasi:         lokasi,
		Tema:           repository.TemaRepository{DB: pool},
		TemaPaket:      lokasi,
		KategoriSinyal: lokasi,
		Pencatat:       repository.IngestRunRepository{DB: pool},
		AturanLokasi:   aturan,
		Anggaran:       anggaran,
		Sekarang:       time.Now,
	}
	return turunan.Jalankan, nil
}

func bacaAturanTurunanAnggaran() (usecase.AturanTurunanAnggaran, error) {
	var aturan usecase.AturanTurunanAnggaran
	if err := config.BacaYAML(pathKonfigurasiKategori, &aturan.Kategori); err != nil {
		return aturan, err
	}
	return aturan, config.BacaYAML(pathKonfigurasiSinyal, &aturan.Sinyal)
}

func buatLaporanGeocode(pool *pgxpool.Pool) job {
	return func(ctx context.Context) error {
		cakupan, err := repository.LokasiPaketRepository{DB: pool}.CakupanLokasi(ctx, jenisPengadaanLaporan)
		if err != nil {
			return err
		}
		cetakCakupan(cakupan)
		return nil
	}
}

func cetakCakupan(cakupan map[string]int) {
	total := 0
	for _, jumlah := range cakupan {
		total += jumlah
	}
	rinci := cakupan[string(domain.TingkatKelurahan)] + cakupan[string(domain.TingkatKecamatan)]
	fmt.Printf("Cakupan lokasi paket %s (%d paket)\n", jenisPengadaanLaporan, total)
	for _, tingkat := range []domain.Tingkat{domain.TingkatKelurahan, domain.TingkatKecamatan, domain.TingkatKota, domain.TingkatProvinsi} {
		fmt.Printf("  %-10s %6d\n", tingkat, cakupan[string(tingkat)])
	}
	if total > 0 {
		fmt.Printf("Kecamatan atau lebih rinci: %.1f%%\n", float64(rinci)*persenPenuh/float64(total))
	}
}

type konfigurasiRingkasan struct {
	BatasPerJalan int      `yaml:"batas_per_jalan"`
	MaksToken     int64    `yaml:"maks_token"`
	Effort        string   `yaml:"effort"`
	PromptSistem  string   `yaml:"prompt_sistem"`
	KataPenilaian []string `yaml:"kata_penilaian"`
}

func buatRingkasUmpan(pool *pgxpool.Pool) (job, error) {
	var cfg konfigurasiRingkasan
	var terlarang domain.KosakataTerlarang
	if err := errors.Join(config.BacaYAML(pathKonfigurasiRingkas, &cfg), config.BacaYAML(pathKosakataTerlarang, &terlarang)); err != nil {
		return nil, err
	}
	klienLLM, err := llm.Baru(llm.Konfigurasi{
		Model: os.Getenv("RINGKAS_MODEL"), MaksToken: cfg.MaksToken, Effort: anthropic.OutputConfigEffort(cfg.Effort),
	})
	if err != nil {
		return nil, fmt.Errorf("RINGKAS_MODEL: %w", err)
	}
	umpan := repository.UmpanRepository{DB: pool}
	ringkas := usecase.RingkasUmpan{
		Umpan:    umpan,
		Penulis:  llm.PeringkasUmpan{Penulis: klienLLM, PromptSistem: cfg.PromptSistem},
		Pencatat: repository.IngestRunRepository{DB: pool},
		Aturan:   domain.AturanRingkasan{Terlarang: terlarang, KataPenilaian: cfg.KataPenilaian},
		Batas:    cfg.BatasPerJalan,
		Laporkan: func(pesan string) { log.Println(pesan) },
	}
	return ringkas.Jalankan, nil
}

func buatImporOpentender(pool *pgxpool.Pool) (job, error) {
	var cfg opentender.Konfigurasi
	if err := config.BacaYAML(pathKonfigurasiOpentender, &cfg); err != nil {
		return nil, err
	}
	impor := usecase.ImporRiwayatPenyedia{
		Pengambil: opentender.Sumber{
			Konfigurasi: cfg,
			Klien:       klien.Baru(&http.Client{Timeout: batasWaktuUnduhBesar}, cfg.Klien),
			DirRaw:      dirRaw(),
			Sekarang:    time.Now,
		},
		Penyimpan: repository.RiwayatPenyediaRepository{DB: pool},
		Pencatat:  repository.IngestRunRepository{DB: pool},
		Tahun:     cfg.Tahun,
	}
	return impor.Jalankan, nil
}

func dirRaw() string {
	if dir, ada := os.LookupEnv("RANTAI_DATA_RAW"); ada {
		return dir
	}
	return dirRawBawaan
}
