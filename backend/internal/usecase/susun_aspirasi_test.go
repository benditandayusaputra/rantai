package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const kodeKecamatanUji = "31.72.06"

type temaAspirasiTiruan struct{}

func (temaAspirasiTiruan) Ambil(context.Context, string) (domain.TemaTersimpan, error) {
	return domain.TemaTersimpan{Tema: domain.Tema{
		Kode: "pendidikan", Nama: domain.Teks{ID: "Pendidikan"},
		KonfigurasiTema: domain.KonfigurasiTema{
			TanggungJawab:   []domain.TanggungJawab{{Masalah: "kursi_smp_kurang", Komisi: "E", UnitWilayah: "Suku Dinas {kota}"}},
			TemplatAspirasi: domain.TemplatAspirasi{Judul: domain.Teks{ID: "Kursi di {kecamatan}"}},
		},
	}}, nil
}

type wilayahAspirasiTiruan struct{}

func (wilayahAspirasiTiruan) Ambil(_ context.Context, kode string) (domain.DetailWilayah, error) {
	return domain.DetailWilayah{WilayahDitemukan: domain.WilayahDitemukan{
		RingkasanWilayah: domain.RingkasanWilayah{Kode: kode, Nama: "Kelapa Gading", Tingkat: domain.TingkatKecamatan},
		Hierarki:         []domain.RingkasanWilayah{{Nama: "DKI Jakarta", Tingkat: domain.TingkatProvinsi}, {Nama: "Jakarta Utara", Tingkat: domain.TingkatKota}},
	}}, nil
}

type belanjaAspirasiTiruan struct {
	kecamatan []domain.BelanjaKecamatan
}

func (b belanjaAspirasiTiruan) Hitung(context.Context, string, int) (domain.BelanjaTema, error) {
	return domain.BelanjaTema{Tahun: 2026, Kecamatan: b.kecamatan}, nil
}

type wakilAspirasiTiruan struct{}

func (wakilAspirasiTiruan) WakilKecamatanKomisi(_ context.Context, _ string, komisi []string) (domain.WakilKecamatan, error) {
	anggota := []domain.Wakil{
		{Nama: "Anggota Komisi E", KodeKomisi: "E"},
		{Nama: "Anggota Komisi A", KodeKomisi: "A"},
	}
	return domain.WakilKecamatan{DetailDapil: domain.DetailDapil{
		RingkasanDapil: domain.RingkasanDapil{Nama: "Dapil Uji"}, Anggota: domain.UrutkanWakil(anggota, komisi),
	}}, nil
}

type sumberAspirasiTiruan struct{}

func (sumberAspirasiTiruan) Ambil(_ context.Context, id string) (domain.Provenance, error) {
	return domain.Provenance{ID: id, Sumber: "lpse", DiambilPada: time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC)}, nil
}

func susunAspirasiUji(belanja []domain.BelanjaKecamatan) SusunAspirasi {
	return SusunAspirasi{
		Tema: temaAspirasiTiruan{}, Wilayah: wilayahAspirasiTiruan{}, Belanja: belanjaAspirasiTiruan{kecamatan: belanja},
		Wakil: wakilAspirasiTiruan{}, Provenance: sumberAspirasiTiruan{},
		Kosakata: domain.KosakataAspirasi{Kalimat: domain.KalimatAspirasi{
			Anggaran:       domain.Teks{ID: "{jumlah} paket"},
			AnggaranKosong: domain.Teks{ID: "belum ada paket"},
			Sumber:         domain.Teks{ID: "diambil {tanggal}"},
			Instansi:       domain.Teks{ID: "{unit}"},
		}},
	}
}

func permintaanAspirasiUji() PermintaanAspirasi {
	return PermintaanAspirasi{KodeKecamatan: kodeKecamatanUji, Tema: "pendidikan", Masalah: "kursi_smp_kurang", Bahasa: domain.BahasaIndonesia}
}

func TestSusunAspirasi_PermintaanTidakValid(t *testing.T) {
	kasus := map[string]func(p *PermintaanAspirasi){
		"bahasa asing":        func(p *PermintaanAspirasi) { p.Bahasa = "fr" },
		"masalah tak dikenal": func(p *PermintaanAspirasi) { p.Masalah = "banjir" },
		"kode kota":           func(p *PermintaanAspirasi) { p.KodeKecamatan = "31.72" },
	}
	for nama, ubah := range kasus {
		t.Run(nama, func(t *testing.T) {
			p := permintaanAspirasiUji()
			ubah(&p)
			_, err := susunAspirasiUji(nil).Susun(t.Context(), p)
			var kesalahan *domain.Kesalahan
			if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriInputTidakValid {
				t.Errorf("err = %v, harap input tidak valid", err)
			}
		})
	}
}

func TestSusunAspirasi_DataKecamatan_BersumberDanHanyaAnggotaKomisiPengawas(t *testing.T) {
	belanja := []domain.BelanjaKecamatan{{
		Kecamatan: domain.RingkasanWilayah{Kode: kodeKecamatanUji},
		Anggaran:  &domain.RingkasanAnggaran{JumlahPaket: 3, ProvenanceID: "prv_lpse"},
	}}

	dokumen, err := susunAspirasiUji(belanja).Susun(t.Context(), permintaanAspirasiUji())

	teks := dokumen.Teks()
	switch {
	case err != nil:
		t.Fatal(err)
	case dokumen.Judul != "Kursi di Kelapa Gading":
		t.Errorf("judul = %q", dokumen.Judul)
	case !strings.Contains(teks, "3 paket diambil 23 Oktober 2026"):
		t.Errorf("anggaran tanpa sumber bertanggal:\n%s", teks)
	case !strings.Contains(teks, "Suku Dinas Jakarta Utara"):
		t.Errorf("unit wilayah tanpa kota:\n%s", teks)
	case !strings.Contains(teks, "Anggota Komisi E") || strings.Contains(teks, "Anggota Komisi A"):
		t.Errorf("anggota pengawas salah:\n%s", teks)
	}
}

func TestSusunAspirasi_KecamatanTanpaBelanja_KalimatKosong(t *testing.T) {
	dokumen, err := susunAspirasiUji(nil).Susun(t.Context(), permintaanAspirasiUji())

	if err != nil || !strings.Contains(dokumen.Teks(), "belum ada paket") {
		t.Errorf("err %v, teks:\n%s", err, dokumen.Teks())
	}
}

func TestBacaJalurAspirasi_KodeKelurahan_MemakaiKecamatanDanKota(t *testing.T) {
	jalur := BacaJalurAspirasi{Wilayah: wilayahAspirasiTiruan{}, Jalur: domain.KonfigurasiJalurAspirasi{
		Jalur: []domain.JalurAspirasi{{Nama: domain.Teks{ID: "Dewan {kota}"}, Cara: domain.Teks{ID: "Kantor {kecamatan}"}}},
	}}

	hasil, err := jalur.Untuk(t.Context(), kodeKecamatanUji+".1001")

	if err != nil || hasil.Wilayah.Kode != kodeKecamatanUji || hasil.Jalur[0].Nama.ID != "Dewan Jakarta Utara" || hasil.Jalur[0].Cara.ID != "Kantor Kelapa Gading" {
		t.Errorf("hasil %+v, err %v", hasil, err)
	}
}
