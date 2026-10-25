package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobIndikator = "indikator"

var ErrPenghitungBelumTersedia = errors.New("penghitung indikator belum tersedia di kode")

type IndikatorTercatat struct {
	domain.IndikatorKesenjangan
	ProvenanceID string
}

type DataIndikator struct {
	TemaKode   string
	Daftar     []IndikatorTercatat
	Provenance []domain.Provenance
}

type PenghitungIndikator interface {
	Hitung(ctx context.Context, tema domain.Tema) (DataIndikator, error)
}

type PembacaTemaAktif interface {
	DaftarTemaAktif(ctx context.Context) ([]domain.TemaTersimpan, error)
}

type PenyimpanIndikator interface {
	SimpanIndikator(ctx context.Context, data DataIndikator) (HasilImpor, error)
}

type HitungIndikator struct {
	Tema       PembacaTemaAktif
	Penghitung map[string]PenghitungIndikator
	Penyimpan  PenyimpanIndikator
	Pencatat   PencatatIngest
}

func (u HitungIndikator) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobIndikator)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobIndikator, err)
	}
	hasil, err := u.hitungSemua(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u HitungIndikator) hitungSemua(ctx context.Context) (HasilImpor, error) {
	daftarTema, err := u.Tema.DaftarTemaAktif(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca tema aktif: %w", err)
	}
	total := HasilImpor{}
	for _, tema := range daftarTema {
		hasil, err := u.hitungTema(ctx, tema.Tema)
		if err != nil {
			return HasilImpor{}, fmt.Errorf("indikator tema %s: %w", tema.Kode, err)
		}
		total.JumlahBaru += hasil.JumlahBaru
		total.JumlahBerubah += hasil.JumlahBerubah
	}
	return total, nil
}

func (u HitungIndikator) hitungTema(ctx context.Context, tema domain.Tema) (HasilImpor, error) {
	penghitung, ada := u.Penghitung[tema.Indikator.Penghitung]
	if !ada {
		return HasilImpor{}, fmt.Errorf("%q: %w", tema.Indikator.Penghitung, ErrPenghitungBelumTersedia)
	}
	data, err := penghitung.Hitung(ctx, tema)
	if err != nil {
		return HasilImpor{}, err
	}
	data.TemaKode = tema.Kode
	return u.Penyimpan.SimpanIndikator(ctx, data)
}
