package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobParipurna = "paripurna"

var (
	ErrFraksiTidakDikenal         = errors.New("fraksi tidak ada di data anggota")
	ErrKosakataTerlarangParipurna = errors.New("teks paripurna memuat kosakata terlarang")
)

type DataParipurna struct {
	Sikap      []domain.SikapFraksi
	Aspirasi   []domain.AspirasiReses
	Provenance domain.Provenance
}

type PengambilParipurna interface {
	Ambil(ctx context.Context) (DataParipurna, error)
}

type PenyimpanParipurna interface {
	DaftarFraksi(ctx context.Context) ([]string, error)
	SimpanParipurna(ctx context.Context, data DataParipurna) (HasilImpor, error)
}

type ImporParipurna struct {
	Pengambil PengambilParipurna
	Penyimpan PenyimpanParipurna
	Pencatat  PencatatIngest
	Terlarang domain.KosakataTerlarang
}

func (u ImporParipurna) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobParipurna)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobParipurna, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporParipurna) impor(ctx context.Context) (HasilImpor, error) {
	data, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil data paripurna: %w", err)
	}
	fraksi, err := u.Penyimpan.DaftarFraksi(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca daftar fraksi: %w", err)
	}
	if err := errors.Join(u.validasiSikap(data.Sikap, fraksi), u.validasiAspirasi(data.Aspirasi)); err != nil {
		return HasilImpor{}, err
	}
	return u.Penyimpan.SimpanParipurna(ctx, data)
}

func (u ImporParipurna) validasiSikap(daftar []domain.SikapFraksi, fraksi []string) error {
	pelanggaran := []error{}
	for i, s := range daftar {
		err := errors.Join(s.Validasi(), u.periksaKosakata(s.Ringkasan))
		if !slices.Contains(fraksi, s.Fraksi) {
			err = errors.Join(err, fmt.Errorf("%q: %w", s.Fraksi, ErrFraksiTidakDikenal))
		}
		if err != nil {
			pelanggaran = append(pelanggaran, fmt.Errorf("sikap fraksi baris %d: %w", i+1, err))
		}
	}
	return errors.Join(pelanggaran...)
}

func (u ImporParipurna) validasiAspirasi(daftar []domain.AspirasiReses) error {
	pelanggaran := []error{}
	for i, a := range daftar {
		if err := errors.Join(a.Validasi(), u.periksaKosakata(a.Uraian)); err != nil {
			pelanggaran = append(pelanggaran, fmt.Errorf("aspirasi reses baris %d: %w", i+1, err))
		}
	}
	return errors.Join(pelanggaran...)
}

func (u ImporParipurna) periksaKosakata(teks string) error {
	if ditemukan := u.Terlarang.Temukan(teks); len(ditemukan) > 0 {
		return fmt.Errorf("%s: %w", strings.Join(ditemukan, ", "), ErrKosakataTerlarangParipurna)
	}
	return nil
}
