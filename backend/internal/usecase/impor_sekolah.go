package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobDapodik = "dapodik"

type DataSumberSekolah struct {
	Sekolah    []domain.DataSekolah
	Provenance domain.Provenance
}

type SekolahTerpetakan struct {
	domain.DataSekolah
	domain.PenempatanSekolah
}

type DataSimpanSekolah struct {
	Sekolah    []SekolahTerpetakan
	Provenance domain.Provenance
}

type LaporanImporSekolah struct {
	JumlahSumber          int
	JumlahTersimpan       int
	TanpaKoordinat        []string
	KoordinatDiLuarBatas  []string
	TidakTerpetakan       []string
	KecamatanTanpaSekolah []string
}

type PengambilSekolah interface {
	Ambil(ctx context.Context) (DataSumberSekolah, error)
}

type PemetaWilayahSekolah interface {
	KandidatWilayah(ctx context.Context, lokasi *domain.Titik, namaKecamatanNormal string) (domain.KandidatWilayah, error)
}

type PenyimpanSekolah interface {
	SimpanSekolah(ctx context.Context, data DataSimpanSekolah) (HasilImpor, error)
	KecamatanTanpaSekolah(ctx context.Context) ([]string, error)
}

type PelaporImpor interface {
	LaporkanSekolah(laporan LaporanImporSekolah)
}

type ImporSekolah struct {
	Pengambil PengambilSekolah
	Pemeta    PemetaWilayahSekolah
	Penyimpan PenyimpanSekolah
	Pencatat  PencatatIngest
	Pelapor   PelaporImpor
}

func (u ImporSekolah) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobDapodik)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobDapodik, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporSekolah) impor(ctx context.Context) (HasilImpor, error) {
	sumber, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil sumber sekolah: %w", err)
	}
	terpetakan, laporan, err := u.petakan(ctx, sumber.Sekolah)
	if err != nil {
		return HasilImpor{}, err
	}
	hasil, err := u.Penyimpan.SimpanSekolah(ctx, DataSimpanSekolah{Sekolah: terpetakan, Provenance: sumber.Provenance})
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan sekolah: %w", err)
	}
	if laporan.KecamatanTanpaSekolah, err = u.Penyimpan.KecamatanTanpaSekolah(ctx); err != nil {
		return HasilImpor{}, fmt.Errorf("periksa kecamatan tanpa sekolah: %w", err)
	}
	u.Pelapor.LaporkanSekolah(laporan)
	return hasil, nil
}

func (u ImporSekolah) petakan(ctx context.Context, daftar []domain.DataSekolah) ([]SekolahTerpetakan, LaporanImporSekolah, error) {
	terpetakan := []SekolahTerpetakan{}
	laporan := LaporanImporSekolah{JumlahSumber: len(daftar), TanpaKoordinat: []string{}, KoordinatDiLuarBatas: []string{}, TidakTerpetakan: []string{}}
	for _, s := range daftar {
		kandidat, err := u.Pemeta.KandidatWilayah(ctx, s.Lokasi, domain.NormalkanNamaKecamatan(s.NamaKecamatanSumber))
		if err != nil {
			return nil, LaporanImporSekolah{}, fmt.Errorf("petakan sekolah %s: %w", s.NPSN, err)
		}
		penempatan, ada := domain.TempatkanSekolah(kandidat)
		laporan.catat(s, penempatan, ada)
		if ada {
			terpetakan = append(terpetakan, SekolahTerpetakan{DataSekolah: s, PenempatanSekolah: penempatan})
		}
	}
	laporan.JumlahTersimpan = len(terpetakan)
	return terpetakan, laporan, nil
}

func (l *LaporanImporSekolah) catat(s domain.DataSekolah, penempatan domain.PenempatanSekolah, ada bool) {
	identitas := fmt.Sprintf("%s %s (%s)", s.NPSN, s.Nama, s.NamaKecamatanSumber)
	switch {
	case s.Lokasi == nil:
		l.TanpaKoordinat = append(l.TanpaKoordinat, identitas)
	case penempatan.Metode != domain.MetodeLokasiTitik:
		l.KoordinatDiLuarBatas = append(l.KoordinatDiLuarBatas, identitas)
	}
	if !ada {
		l.TidakTerpetakan = append(l.TidakTerpetakan, identitas)
	}
}
