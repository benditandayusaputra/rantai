package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobUmpanDPRD = "umpan_dprd"

type DataSumberUmpan struct {
	Item       []domain.ItemUmpan
	Provenance domain.Provenance
}

type DataSimpanUmpan struct {
	Umpan      []domain.UmpanTerpaut
	Provenance domain.Provenance
}

type PengambilUmpan interface {
	Ambil(ctx context.Context, sejak time.Time) (DataSumberUmpan, error)
}

type PembacaAnggotaUmpan interface {
	DaftarAnggota(ctx context.Context) ([]domain.Anggota, error)
}

type PembacaDaftarTema interface {
	Daftar(ctx context.Context) ([]domain.TemaTersimpan, error)
}

type PenyimpanUmpan interface {
	TerbitTerakhir(ctx context.Context, jenis domain.JenisSumberUmpan) (time.Time, bool, error)
	SimpanUmpan(ctx context.Context, data DataSimpanUmpan) (HasilImpor, error)
}

type ImporUmpan struct {
	Job           string
	JenisSumber   domain.JenisSumberUmpan
	MulaiBawaan   time.Time
	TumpangTindih time.Duration
	Pengambil     PengambilUmpan
	Anggota       PembacaAnggotaUmpan
	Tema          PembacaDaftarTema
	Penyimpan     PenyimpanUmpan
	Pencatat      PencatatIngest
	Pautkan       func(domain.PencariNama, domain.ItemUmpan) []domain.PautanUmpan
}

func (u ImporUmpan) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, u.Job)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", u.Job, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporUmpan) impor(ctx context.Context) (HasilImpor, error) {
	sejak, err := u.sejak(ctx)
	if err != nil {
		return HasilImpor{}, err
	}
	sumber, err := u.Pengambil.Ambil(ctx, sejak)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil umpan %s: %w", u.JenisSumber, err)
	}
	umpan, err := u.pautkanSemua(ctx, sumber.Item)
	if err != nil {
		return HasilImpor{}, err
	}
	hasil, err := u.Penyimpan.SimpanUmpan(ctx, DataSimpanUmpan{Umpan: umpan, Provenance: sumber.Provenance})
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan umpan %s: %w", u.JenisSumber, err)
	}
	return hasil, nil
}

func (u ImporUmpan) sejak(ctx context.Context) (time.Time, error) {
	terakhir, ada, err := u.Penyimpan.TerbitTerakhir(ctx, u.JenisSumber)
	if err != nil {
		return time.Time{}, fmt.Errorf("baca umpan terakhir: %w", err)
	}
	if !ada {
		return u.MulaiBawaan, nil
	}
	return terakhir.Add(-u.TumpangTindih), nil
}

func (u ImporUmpan) pautkanSemua(ctx context.Context, daftar []domain.ItemUmpan) ([]domain.UmpanTerpaut, error) {
	anggota, err := u.Anggota.DaftarAnggota(ctx)
	if err != nil {
		return nil, fmt.Errorf("baca anggota: %w", err)
	}
	tagIsu, err := u.semuaTagIsu(ctx)
	if err != nil {
		return nil, err
	}
	pencari := domain.BuatPencariNama(anggota)
	umpan := []domain.UmpanTerpaut{}
	for _, item := range daftar {
		pautan := u.Pautkan(pencari, item)
		if len(pautan) == 0 {
			continue
		}
		tag := domain.TandaiTagIsu(item.Judul+" "+item.Cuplikan, tagIsu)
		item.Isi = ""
		umpan = append(umpan, domain.UmpanTerpaut{ItemUmpan: item, TagIsu: tag, Pautan: pautan})
	}
	return umpan, nil
}

func (u ImporUmpan) semuaTagIsu(ctx context.Context) ([]string, error) {
	daftar, err := u.Tema.Daftar(ctx)
	if err != nil {
		return nil, fmt.Errorf("baca tag isu tema: %w", err)
	}
	tag := []string{}
	for _, t := range daftar {
		tag = append(tag, t.TagIsu...)
	}
	slices.Sort(tag)
	return slices.Compact(tag), nil
}
