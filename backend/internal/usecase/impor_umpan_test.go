package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilUmpanTiruan struct {
	item  []domain.ItemUmpan
	sejak time.Time
}

func (p *pengambilUmpanTiruan) Ambil(_ context.Context, sejak time.Time) (DataSumberUmpan, error) {
	p.sejak = sejak
	return DataSumberUmpan{Item: p.item}, nil
}

type umpanTiruan struct {
	terakhir *time.Time
	diterima DataSimpanUmpan
}

func (u *umpanTiruan) DaftarAnggota(context.Context) ([]domain.Anggota, error) {
	return []domain.Anggota{{ID: "yuke", Nama: "Yuke Yurike"}}, nil
}

func (u *umpanTiruan) TerbitTerakhir(context.Context, domain.JenisSumberUmpan) (time.Time, bool, error) {
	if u.terakhir == nil {
		return time.Time{}, false, nil
	}
	return *u.terakhir, true, nil
}

func (u *umpanTiruan) SimpanUmpan(_ context.Context, data DataSimpanUmpan) (HasilImpor, error) {
	u.diterima = data
	return HasilImpor{JumlahBaru: len(data.Umpan)}, nil
}

type temaTiruan struct{}

func (temaTiruan) Daftar(context.Context) ([]domain.TemaTersimpan, error) {
	return []domain.TemaTersimpan{{Tema: domain.Tema{KonfigurasiTema: domain.KonfigurasiTema{TagIsu: []string{"sekolah"}}}}}, nil
}

var mulaiUmpanUji = time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC)

func imporUmpanUji(pengambil *pengambilUmpanTiruan, umpan *umpanTiruan) ImporUmpan {
	return ImporUmpan{
		Job: JobUmpanDPRD, JenisSumber: domain.SumberUmpanResmiDPRD, MulaiBawaan: mulaiUmpanUji, TumpangTindih: time.Hour,
		Pengambil: pengambil, Anggota: umpan, Tema: temaTiruan{}, Penyimpan: umpan, Pencatat: &pencatatTiruan{},
		Pautkan: domain.PencariNama.PautkanRilisResmi,
	}
}

func TestImporUmpan_BelumAdaUmpan_MulaiDariTanggalBawaan(t *testing.T) {
	pengambil := &pengambilUmpanTiruan{}

	if err := imporUmpanUji(pengambil, &umpanTiruan{}).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if !pengambil.sejak.Equal(mulaiUmpanUji) {
		t.Errorf("sejak = %v, harap %v", pengambil.sejak, mulaiUmpanUji)
	}
}

func TestImporUmpan_SudahAdaUmpan_MulaiDariTerakhirDikurangiTumpangTindih(t *testing.T) {
	terakhir := time.Date(2026, time.October, 20, 10, 0, 0, 0, time.UTC)
	pengambil := &pengambilUmpanTiruan{}

	if err := imporUmpanUji(pengambil, &umpanTiruan{terakhir: &terakhir}).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if harap := terakhir.Add(-time.Hour); !pengambil.sejak.Equal(harap) {
		t.Errorf("sejak = %v, harap %v", pengambil.sejak, harap)
	}
}

func TestImporUmpan_ItemTanpaAnggota_DibuangDanIsiTidakDiteruskan(t *testing.T) {
	pengambil := &pengambilUmpanTiruan{item: []domain.ItemUmpan{
		{Judul: "Rapat sekolah", Isi: "Yuke Yurike meninjau sekolah", URL: "https://a"},
		{Judul: "Kegiatan lain", Isi: "tanpa nama anggota", URL: "https://b"},
	}}
	umpan := &umpanTiruan{}

	if err := imporUmpanUji(pengambil, umpan).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if len(umpan.diterima.Umpan) != 1 {
		t.Fatalf("harap satu umpan, dapat %d", len(umpan.diterima.Umpan))
	}
	u := umpan.diterima.Umpan[0]
	if u.URL != "https://a" || u.Isi != "" || len(u.TagIsu) != 1 || u.TagIsu[0] != "sekolah" {
		t.Errorf("umpan = %+v", u)
	}
}
