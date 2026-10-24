package usecase

import (
	"context"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type repoPushTiruan struct {
	RepositoriPush
	tujuan   []domain.TujuanPush
	dihapus  []int64
	didorong []int64
}

func (r *repoPushTiruan) KejadianBelumDidorong(context.Context) ([]domain.KejadianUntukPush, error) {
	return []domain.KejadianUntukPush{{ID: 9, Jenis: domain.KejadianPaketBaru, KodeKecamatan: "31.72.06", Judul: "Rehab", Tautan: "/paket/x"}}, nil
}

func (r *repoPushTiruan) TujuanPush(context.Context, string) ([]domain.TujuanPush, error) {
	return r.tujuan, nil
}

func (r *repoPushTiruan) HapusLanggananPushByID(_ context.Context, id int64) error {
	r.dihapus = append(r.dihapus, id)
	return nil
}

func (r *repoPushTiruan) TandaiDidorong(_ context.Context, id int64) error {
	r.didorong = append(r.didorong, id)
	return nil
}

type pengirimPushTiruan struct{ pesan []domain.PesanPush }

func (p *pengirimPushTiruan) Kirim(_ context.Context, l domain.LanggananPush, pesan domain.PesanPush) error {
	if l.Endpoint == "https://kedaluwarsa.test" {
		return ErrLanggananPushKedaluwarsa
	}
	p.pesan = append(p.pesan, pesan)
	return nil
}

func TestDorongPush_LanggananKedaluwarsaDihapusKejadianDitandai(t *testing.T) {
	repo := &repoPushTiruan{tujuan: []domain.TujuanPush{
		{LanggananID: 1, Langganan: domain.LanggananPush{Endpoint: "https://aktif.test"}},
		{LanggananID: 2, Langganan: domain.LanggananPush{Endpoint: "https://kedaluwarsa.test"}},
	}}
	pengirim := &pengirimPushTiruan{}

	terkirim, err := DorongPush{Repo: repo, Pengirim: pengirim}.Dorong(context.Background())

	if err != nil || terkirim != 1 || len(repo.dihapus) != 1 || repo.dihapus[0] != 2 || len(repo.didorong) != 1 {
		t.Fatalf("terkirim %d, err %v, dihapus %v, didorong %v", terkirim, err, repo.dihapus, repo.didorong)
	}
	if pengirim.pesan[0].Judul != "Paket baru" || pengirim.pesan[0].Tautan != "/paket/x" {
		t.Errorf("pesan = %+v", pengirim.pesan[0])
	}
}
