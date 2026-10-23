package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type wilayahRingkasTiruan struct{}

func (wilayahRingkasTiruan) CariWilayah(context.Context, string, int) ([]domain.WilayahDitemukan, error) {
	return []domain.WilayahDitemukan{}, nil
}

func (wilayahRingkasTiruan) KelurahanDiTitik(context.Context, domain.Titik) (domain.WilayahDitemukan, bool, error) {
	return domain.WilayahDitemukan{}, false, nil
}

func (wilayahRingkasTiruan) AmbilWilayah(_ context.Context, kode string) (domain.DetailWilayah, bool, error) {
	tingkat, _ := domain.TingkatDariKode(kode)
	return domain.DetailWilayah{WilayahDitemukan: domain.WilayahDitemukan{RingkasanWilayah: domain.RingkasanWilayah{Kode: kode, Tingkat: tingkat}}}, true, nil
}

type sumberRingkasTiruan struct {
	anggaranTahun int
	komisi        []string
}

func (s *sumberRingkasTiruan) Daftar(context.Context) ([]domain.TemaTersimpan, error) {
	return []domain.TemaTersimpan{{Tema: domain.Tema{Kode: "pendidikan", KonfigurasiTema: domain.KonfigurasiTema{TanggungJawab: []domain.TanggungJawab{{Komisi: "E"}}}}}}, nil
}

func (s *sumberRingkasTiruan) IndikatorTerbaru(context.Context, string, string) (domain.IndikatorKecamatan, bool, error) {
	return domain.IndikatorKecamatan{}, false, nil
}

func (s *sumberRingkasTiruan) RingkasanAnggaran(_ context.Context, _ string, tahun int) (domain.RingkasanAnggaran, bool, error) {
	s.anggaranTahun = tahun
	return domain.RingkasanAnggaran{Tahun: tahun, JumlahPaket: 3}, true, nil
}

func (s *sumberRingkasTiruan) WakilKecamatanKomisi(_ context.Context, _ string, komisi []string) (domain.WakilKecamatan, error) {
	s.komisi = komisi
	return domain.WakilKecamatan{KomisiPengawas: komisi}, nil
}

func ringkasKecamatanUji(sumber *sumberRingkasTiruan) RingkasKecamatan {
	return RingkasKecamatan{
		Wilayah: wilayahRingkasTiruan{}, Tema: sumber, Indikator: sumber, Anggaran: sumber, Wakil: sumber,
		Sekarang: func() time.Time { return time.Date(2026, time.October, 23, 9, 0, 0, 0, time.UTC) },
	}
}

func TestRingkasKecamatan_KodeKelurahan_DiarahkanKeKecamatanDenganSorotan(t *testing.T) {
	sumber := &sumberRingkasTiruan{}

	hasil, err := ringkasKecamatanUji(sumber).Ringkas(t.Context(), "31.72.06.1003")

	if err != nil || hasil.Kecamatan.Kode != "31.72.06" || hasil.KelurahanDisorot == nil || hasil.KelurahanDisorot.Kode != "31.72.06.1003" {
		t.Fatalf("hasil = %+v, err = %v", hasil, err)
	}
	if sumber.anggaranTahun != 2026 || len(sumber.komisi) != 1 || sumber.komisi[0] != "E" {
		t.Errorf("tahun anggaran %d, komisi pengawas %v", sumber.anggaranTahun, sumber.komisi)
	}
	if len(hasil.Tema) != 1 || hasil.Tema[0].Indikator != nil {
		t.Errorf("harap tema tanpa indikator ditandai kosong, dapat %+v", hasil.Tema)
	}
}

func TestRingkasKecamatan_KodeKota_Ditolak(t *testing.T) {
	_, err := ringkasKecamatanUji(&sumberRingkasTiruan{}).Ringkas(t.Context(), "31.72")

	var kesalahan *domain.Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriInputTidakValid {
		t.Errorf("harap input tidak valid, dapat %v", err)
	}
}
