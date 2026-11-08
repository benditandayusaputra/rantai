package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	perintahJadwal    = "jadwal"
	jamJalanHarian    = 2
	selisihWIBDariUTC = 7 * time.Hour
)

var (
	zonaWIB      = time.FixedZone("WIB", int(selisihWIBDariUTC.Seconds()))
	urutanHarian = []string{usecase.JobLPSEDaftar, usecase.JobLPSEDetail, usecase.JobSIRUP, usecase.JobUmpanDPRD, usecase.JobUmpanMedia, usecase.JobRingkasUmpan, usecase.JobAPBD, usecase.JobTurunan}
)

func buatJadwalHarian(pool *pgxpool.Pool) job {
	return func(ctx context.Context) error {
		for {
			jalan := waktuJalanBerikutnya(time.Now())
			log.Printf("jadwal: pembaruan berikutnya %s", jalan.Format(time.RFC3339))
			if err := tungguSampai(ctx, jalan); err != nil {
				return err
			}
			jalankanUrutanHarian(ctx, pool)
		}
	}
}

func waktuJalanBerikutnya(sekarang time.Time) time.Time {
	lokal := sekarang.In(zonaWIB)
	jalan := time.Date(lokal.Year(), lokal.Month(), lokal.Day(), jamJalanHarian, 0, 0, 0, zonaWIB)
	if !jalan.After(lokal) {
		jalan = jalan.AddDate(0, 0, 1)
	}
	return jalan
}

func tungguSampai(ctx context.Context, waktu time.Time) error {
	pewaktu := time.NewTimer(time.Until(waktu))
	defer pewaktu.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-pewaktu.C:
		return nil
	}
}

func jalankanUrutanHarian(ctx context.Context, pool *pgxpool.Pool) {
	for _, nama := range urutanHarian {
		if ctx.Err() != nil {
			return
		}
		mulai := time.Now()
		if err := jalankanSatuJob(ctx, nama, pool); err != nil {
			log.Printf("jadwal: %s gagal setelah %s: %v", nama, time.Since(mulai).Round(time.Second), err)
			continue
		}
		log.Printf("jadwal: %s selesai dalam %s", nama, time.Since(mulai).Round(time.Second))
	}
}

func jalankanSatuJob(ctx context.Context, nama string, pool *pgxpool.Pool) error {
	if err := usecase.IzinkanJob(os.Getenv("RANTAI_MODE"), nama); err != nil {
		return err
	}
	jalankanJob, err := pilihJob(nama, pool)
	if err != nil {
		return err
	}
	return jalankanJob(ctx)
}
