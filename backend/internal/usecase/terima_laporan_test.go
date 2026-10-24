package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pemrosesFotoTiruan struct{}

func (pemrosesFotoTiruan) Proses(isi []byte) (domain.FotoDiproses, error) {
	if string(isi) == "rusak" {
		return domain.FotoDiproses{}, domain.ErrFormatFotoTidakSah
	}
	return domain.FotoDiproses{Utama: isi, Thumbnail: isi, Lebar: 10, Tinggi: 8}, nil
}

type penyimpanFotoTiruan struct{ kunci []string }

func (p *penyimpanFotoTiruan) SimpanJPEG(_ context.Context, kunci string, _ []byte) error {
	p.kunci = append(p.kunci, kunci)
	return nil
}

type repoLaporanTiruan struct {
	jumlahHariIni int
	jarakPaket    *float64
	diDKI         bool
	disimpan      *domain.LaporanBaru
}

func (r *repoLaporanTiruan) JumlahLaporanSejak(context.Context, int64, time.Time) (int, error) {
	return r.jumlahHariIni, nil
}

func (r *repoLaporanTiruan) JarakPaketMeter(_ context.Context, kode string, _ domain.Titik) (*float64, bool, error) {
	return r.jarakPaket, kode != "tidak-ada", nil
}

func (r *repoLaporanTiruan) WilayahDiTitik(context.Context, domain.Titik) (string, bool, error) {
	return "31.72.06.1003", r.diDKI, nil
}

func (r *repoLaporanTiruan) SimpanLaporan(_ context.Context, l domain.LaporanBaru) (int64, error) {
	r.disimpan = &l
	return 7, nil
}

func terimaLaporanUji(repo *repoLaporanTiruan) (TerimaLaporan, *penyimpanFotoTiruan) {
	penyimpan := &penyimpanFotoTiruan{}
	return TerimaLaporan{
		Pemroses: pemrosesFotoTiruan{}, Penyimpan: penyimpan, Repo: repo,
		Sekarang: time.Now, BuatKunci: func() string { return "acak" },
	}, penyimpan
}

func laporanMasukUji(foto ...string) domain.LaporanMasuk {
	isi := make([][]byte, len(foto))
	for i, f := range foto {
		isi[i] = []byte(f)
	}
	return domain.LaporanMasuk{
		PenggunaID: 3, PaketKode: "paket-uji", Titik: domain.Titik{Lat: -6.121437, Lng: 106.901288},
		AkurasiMeter: 20, Kondisi: domain.KondisiTerhenti, Foto: isi,
	}
}

func TestTerimaLaporan_Sah_DisimpanDenganLokasiPublikDanDuaObjekPerFoto(t *testing.T) {
	repo := &repoLaporanTiruan{diDKI: true, jarakPaket: new(120.0)}
	u, penyimpan := terimaLaporanUji(repo)

	hasil, err := u.Terima(context.Background(), laporanMasukUji("foto-1", "foto-2"))

	if err != nil || hasil.ID != 7 || hasil.Status != domain.StatusMenungguSaksi || len(penyimpan.kunci) != 4 {
		t.Fatalf("hasil %+v, err %v, kunci %v", hasil, err, penyimpan.kunci)
	}
	l := repo.disimpan
	if l.TitikPublik != domain.BulatkanLokasiPublik(l.Titik) || l.KodeWilayah != "31.72.06.1003" || l.Foto[0].ThumbnailKey != "laporan/acak-480.jpg" {
		t.Errorf("laporan = %+v", l)
	}
}

func TestTerimaLaporan_Penolakan(t *testing.T) {
	kasus := []struct {
		nama  string
		repo  repoLaporanTiruan
		ubah  func(*domain.LaporanMasuk)
		harap error
	}{
		{nama: "batas harian tercapai", repo: repoLaporanTiruan{diDKI: true, jumlahHariIni: domain.BatasLaporanPerHari}, harap: domain.ErrBatasLaporanHarian},
		{nama: "paket lebih dari 500 m", repo: repoLaporanTiruan{diDKI: true, jarakPaket: new(650.0)}, harap: domain.ErrPaketLaporanJauh},
		{nama: "titik di luar DKI", repo: repoLaporanTiruan{}, harap: domain.ErrLokasiDiLuarDKI},
		{nama: "foto rusak", repo: repoLaporanTiruan{diDKI: true}, ubah: func(l *domain.LaporanMasuk) { l.Foto = [][]byte{[]byte("foto-1"), []byte("rusak")} }, harap: domain.ErrFormatFotoTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			repo := k.repo
			u, penyimpan := terimaLaporanUji(&repo)
			masuk := laporanMasukUji("foto-1")
			if k.ubah != nil {
				k.ubah(&masuk)
			}
			if _, err := u.Terima(context.Background(), masuk); !errors.Is(err, k.harap) || repo.disimpan != nil || len(penyimpan.kunci) != 0 {
				t.Errorf("err %v, disimpan %v, kunci %v", err, repo.disimpan, penyimpan.kunci)
			}
		})
	}
}
