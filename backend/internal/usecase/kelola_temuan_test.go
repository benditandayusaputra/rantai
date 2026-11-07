package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

type repoTemuanPalsu struct {
	isPaketAda      bool
	jumlahHariIni   int
	pengirim        int64
	tersimpan       []domain.TemuanMasuk
	tanggapan       []domain.TanggapanMasuk
	catatan         []domain.CatatanModerasi
	isDisembunyikan bool
}

func (r *repoTemuanPalsu) IsPaketAda(context.Context, string) (bool, error) { return r.isPaketAda, nil }
func (r *repoTemuanPalsu) HitungTemuanSejak(context.Context, int64, time.Time) (int, error) {
	return r.jumlahHariIni, nil
}
func (r *repoTemuanPalsu) SimpanTemuan(_ context.Context, m domain.TemuanMasuk) (int64, error) {
	r.tersimpan = append(r.tersimpan, m)
	return int64(len(r.tersimpan)), nil
}
func (r *repoTemuanPalsu) AmbilPengirimTemuan(context.Context, int64) (int64, bool, error) {
	return r.pengirim, r.pengirim != 0, nil
}
func (r *repoTemuanPalsu) SimpanTanggapan(_ context.Context, m domain.TanggapanMasuk) error {
	r.tanggapan = append(r.tanggapan, m)
	return nil
}
func (r *repoTemuanPalsu) DaftarTemuan(context.Context, string, int) ([]domain.Temuan, error) {
	return []domain.Temuan{}, nil
}
func (r *repoTemuanPalsu) SembunyikanTemuan(context.Context, int64) (bool, error) {
	r.isDisembunyikan = true
	return true, nil
}
func (r *repoTemuanPalsu) CatatModerasi(_ context.Context, c domain.CatatanModerasi) error {
	r.catatan = append(r.catatan, c)
	return nil
}

func kelolaUji(repo *repoTemuanPalsu) usecase.KelolaTemuan {
	return usecase.KelolaTemuan{Repo: repo, Terlarang: domain.KosakataTerlarang{"korupsi"}, Sekarang: time.Now}
}

var temuanSah = domain.TemuanMasuk{
	PenggunaID: 1, PaketKode: "P1", Isi: "  Papan proyek mencantumkan nilai kontrak berbeda.  ", TautanSumber: " https://contoh.id/foto ",
}

func TestKirimTemuan_Sah_DisimpanTerpangkas(t *testing.T) {
	repo := &repoTemuanPalsu{isPaketAda: true}
	if _, err := kelolaUji(repo).Kirim(t.Context(), temuanSah); err != nil {
		t.Fatal(err)
	}
	if got := repo.tersimpan[0]; got.Isi != "Papan proyek mencantumkan nilai kontrak berbeda." || got.TautanSumber != "https://contoh.id/foto" {
		t.Errorf("tersimpan %+v", got)
	}
}

func TestKirimTemuan_BatasHarianTercapai_Ditolak(t *testing.T) {
	repo := &repoTemuanPalsu{isPaketAda: true, jumlahHariIni: domain.BatasTemuanPerHari}
	if _, err := kelolaUji(repo).Kirim(t.Context(), temuanSah); !errors.Is(err, domain.ErrBatasTemuanHarian) {
		t.Errorf("err %v", err)
	}
}

func TestKirimTemuan_PaketTidakAda_TidakDitemukan(t *testing.T) {
	_, err := kelolaUji(&repoTemuanPalsu{}).Kirim(t.Context(), temuanSah)
	var k *domain.Kesalahan
	if !errors.As(err, &k) || k.Kategori != domain.KategoriTidakDitemukan {
		t.Errorf("err %v", err)
	}
}

func TestTanggapiTemuan_PengirimSendiri_Ditolak(t *testing.T) {
	repo := &repoTemuanPalsu{pengirim: 7}
	err := kelolaUji(repo).Tanggapi(t.Context(), domain.TanggapanMasuk{TemuanID: 1, PenggunaID: 7, Sikap: domain.SikapMendukung})
	if !errors.Is(err, domain.ErrTanggapiTemuanSendiri) || len(repo.tanggapan) != 0 {
		t.Errorf("err %v tanggapan %d", err, len(repo.tanggapan))
	}
}

func TestTanggapiTemuan_PenggunaLain_Disimpan(t *testing.T) {
	repo := &repoTemuanPalsu{pengirim: 7}
	err := kelolaUji(repo).Tanggapi(t.Context(), domain.TanggapanMasuk{TemuanID: 1, PenggunaID: 8, Sikap: domain.SikapMembantah})
	if err != nil || len(repo.tanggapan) != 1 {
		t.Errorf("err %v tanggapan %d", err, len(repo.tanggapan))
	}
}

func TestSembunyikanTemuan_TanpaAlasan_Ditolak(t *testing.T) {
	repo := &repoTemuanPalsu{}
	err := kelolaUji(repo).Sembunyikan(t.Context(), usecase.PermintaanModerasiTemuan{TemuanID: 1, ModeratorID: 2})
	if !errors.Is(err, domain.ErrAlasanModerasiKosong) || repo.isDisembunyikan {
		t.Errorf("err %v", err)
	}
}

func TestSembunyikanTemuan_DenganAlasan_Dicatat(t *testing.T) {
	repo := &repoTemuanPalsu{}
	err := kelolaUji(repo).Sembunyikan(t.Context(), usecase.PermintaanModerasiTemuan{TemuanID: 1, ModeratorID: 2, Alasan: "spam"})
	if err != nil || len(repo.catatan) != 1 || repo.catatan[0].TargetJenis != domain.TargetTemuan {
		t.Errorf("err %v catatan %+v", err, repo.catatan)
	}
}
