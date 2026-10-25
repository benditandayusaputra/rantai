package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type repoRingkasanTiruan struct {
	RepositoriRingkasanEmail
	penerima []domain.PenerimaRingkasan
	kabar    map[int64][]domain.Notifikasi
	ditandai map[int64]time.Time
	batas    time.Time
}

func (r *repoRingkasanTiruan) PenerimaRingkasan(_ context.Context, terkirimSebelum time.Time) ([]domain.PenerimaRingkasan, error) {
	r.batas = terkirimSebelum
	return r.penerima, nil
}

func (r *repoRingkasanTiruan) KabarRingkasan(_ context.Context, penggunaID int64, _ time.Time, _ int) ([]domain.Notifikasi, error) {
	return r.kabar[penggunaID], nil
}

func (r *repoRingkasanTiruan) TandaiRingkasanTerkirim(_ context.Context, penggunaID int64, pada time.Time) error {
	r.ditandai[penggunaID] = pada
	return nil
}

type pengirimEmailTiruan struct{ tujuan []string }

func (p *pengirimEmailTiruan) Kirim(_ context.Context, email, _, _ string) error {
	p.tujuan = append(p.tujuan, email)
	return nil
}

func TestKirimRingkasanMingguan_HanyaYangPunyaKabar_SemuaDitandai(t *testing.T) {
	sekarang := time.Date(2026, time.October, 26, 7, 0, 0, 0, time.UTC)
	repo := &repoRingkasanTiruan{
		penerima: []domain.PenerimaRingkasan{{PenggunaID: 1, Email: "a@uji.test"}, {PenggunaID: 2, Email: "b@uji.test"}},
		kabar:    map[int64][]domain.Notifikasi{1: {{Jenis: domain.KejadianPaketBaru, Judul: "Paket", Tautan: "/paket/x"}}},
		ditandai: map[int64]time.Time{},
	}
	pengirim := &pengirimEmailTiruan{}
	u := KirimRingkasanMingguan{Repo: repo, Pengirim: pengirim, AlamatSitus: "https://rantai.uji", Sekarang: func() time.Time { return sekarang }}

	terkirim, err := u.Jalankan(t.Context())

	if err != nil || terkirim != 1 || len(pengirim.tujuan) != 1 || pengirim.tujuan[0] != "a@uji.test" {
		t.Fatalf("terkirim %d, tujuan %v, galat %v", terkirim, pengirim.tujuan, err)
	}
	if len(repo.ditandai) != 2 || !repo.ditandai[2].Equal(sekarang) {
		t.Errorf("ditandai = %v", repo.ditandai)
	}
	if !repo.batas.Equal(sekarang.Add(-domain.JendelaKejadian)) {
		t.Errorf("batas penerima = %v", repo.batas)
	}
}
