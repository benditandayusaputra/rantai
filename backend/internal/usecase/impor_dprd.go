package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobKPUDPRD = "kpu_dprd"

type DataSumberDPRD struct {
	BarisDapil        []domain.BarisDapil
	AnggotaFraksi     []domain.AnggotaFraksi
	ProvenanceDapil   domain.Provenance
	ProvenanceAnggota domain.Provenance
}

type DataSimpanWakil struct {
	TahunPemilu       int
	URLLHKPN          string
	Dapil             []domain.Dapil
	Anggota           []domain.Anggota
	ProvenanceDapil   domain.Provenance
	ProvenanceAnggota domain.Provenance
}

type PengambilDPRD interface {
	Ambil(ctx context.Context) (DataSumberDPRD, error)
}

type PembacaKecamatan interface {
	DaftarKecamatan(ctx context.Context) (map[string]string, error)
}

type PenyimpanWakil interface {
	SimpanWakil(ctx context.Context, data DataSimpanWakil) (HasilImpor, error)
}

type ImporDPRD struct {
	Pengambil        PengambilDPRD
	PembacaKecamatan PembacaKecamatan
	Penyimpan        PenyimpanWakil
	Pencatat         PencatatIngest
	Dapil            []domain.Dapil
	TahunPemilu      int
	URLLHKPN         string
	PolaURLProfil    string
}

func (u ImporDPRD) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobKPUDPRD)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobKPUDPRD, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporDPRD) impor(ctx context.Context) (HasilImpor, error) {
	if err := u.validasiDapil(ctx); err != nil {
		return HasilImpor{}, err
	}
	sumber, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil halaman DPRD: %w", err)
	}
	anggota, err := u.susunAnggota(sumber)
	if err != nil {
		return HasilImpor{}, err
	}
	hasil, err := u.Penyimpan.SimpanWakil(ctx, DataSimpanWakil{
		TahunPemilu: u.TahunPemilu, URLLHKPN: u.URLLHKPN, Dapil: u.Dapil, Anggota: anggota,
		ProvenanceDapil: sumber.ProvenanceDapil, ProvenanceAnggota: sumber.ProvenanceAnggota,
	})
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan dapil dan anggota: %w", err)
	}
	return hasil, nil
}

func (u ImporDPRD) validasiDapil(ctx context.Context) error {
	kecamatan, err := u.PembacaKecamatan.DaftarKecamatan(ctx)
	if err != nil {
		return fmt.Errorf("baca kecamatan DKI: %w", err)
	}
	return domain.ValidasiCakupanDapil(u.Dapil, kecamatan)
}

func (u ImporDPRD) susunAnggota(sumber DataSumberDPRD) ([]domain.Anggota, error) {
	anggota, err := domain.GabungkanAnggota(sumber.BarisDapil, sumber.AnggotaFraksi)
	if err != nil {
		return nil, err
	}
	for i := range anggota {
		anggota[i].URLProfil = fmt.Sprintf(u.PolaURLProfil, anggota[i].ID)
	}
	return anggota, domain.ValidasiKursiDapil(u.Dapil, anggota)
}
