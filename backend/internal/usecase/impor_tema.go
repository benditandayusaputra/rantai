package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobTema = "tema"

var PenghitungTerdaftar = []string{domain.PenghitungKesenjanganKursiSMP, domain.PenghitungKesenjanganKursiSMASMK}

var ErrKodeTemaGanda = errors.New("kode tema dipakai lebih dari satu berkas")

type DataSumberTema struct {
	Tema       []domain.Tema
	Provenance domain.Provenance
}

type PengambilTema interface {
	Ambil(ctx context.Context) (DataSumberTema, error)
}

type PenyimpanTema interface {
	SimpanTema(ctx context.Context, data DataSumberTema) (HasilImpor, error)
}

type ImporTema struct {
	Pengambil PengambilTema
	Penyimpan PenyimpanTema
	Pencatat  PencatatIngest
}

func (u ImporTema) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobTema)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobTema, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporTema) impor(ctx context.Context) (HasilImpor, error) {
	sumber, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil konfigurasi tema: %w", err)
	}
	if err := validasiSemuaTema(sumber.Tema); err != nil {
		return HasilImpor{}, err
	}
	hasil, err := u.Penyimpan.SimpanTema(ctx, sumber)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan tema: %w", err)
	}
	return hasil, nil
}

func validasiSemuaTema(daftar []domain.Tema) error {
	pelanggaran := []error{}
	kodeTerpakai := map[string]bool{}
	for _, tema := range daftar {
		if kodeTerpakai[tema.Kode] {
			pelanggaran = append(pelanggaran, fmt.Errorf("%q: %w", tema.Kode, ErrKodeTemaGanda))
		}
		kodeTerpakai[tema.Kode] = true
		pelanggaran = append(pelanggaran, tema.Validasi(PenghitungTerdaftar))
	}
	return errors.Join(pelanggaran...)
}
