package domain

import (
	"errors"
	"slices"
	"testing"
	"time"
)

func konfigurasiAnggaranUji() KonfigurasiAspirasiAnggaran {
	return KonfigurasiAspirasiAnggaran{
		Bidang: map[string]AturanBidangAnggaran{
			"infrastruktur": {Kategori: []string{"jalan", "saluran"}, KataKunci: []string{"penerangan jalan"}},
			"kesehatan":     {KataKunci: []string{"posyandu"}},
		},
		BatasPaket: 10,
	}
}

func TestKueriUntuk_ResesTahunIni_TahunAnggaranBerikutnyaDanPolaKataUtuh(t *testing.T) {
	aspirasi := AspirasiReses{Bidang: "kesehatan", KodeWilayah: "31.75.02.1006", TanggalParipurna: time.Date(2026, time.August, 26, 0, 0, 0, 0, time.UTC)}

	kueri := konfigurasiAnggaranUji().KueriUntuk(aspirasi)

	if kueri.TahunAnggaran != 2027 || kueri.KodeWilayah != "31.75.02.1006" || kueri.Batas != 10 ||
		!slices.Equal(kueri.PolaKataKunci, []string{`\mposyandu\M`}) || kueri.Kategori == nil || len(kueri.Kategori) != 0 {
		t.Errorf("kueri = %+v", kueri)
	}
}

func TestKueriUntuk_BidangTanpaAturan_KategoriDanPolaKosong(t *testing.T) {
	kueri := konfigurasiAnggaranUji().KueriUntuk(AspirasiReses{Bidang: "budaya", KodeWilayah: "31.71.03"})

	if len(kueri.Kategori) != 0 || len(kueri.PolaKataKunci) != 0 {
		t.Errorf("kueri = %+v", kueri)
	}
}

func TestStatusPaketTerkait(t *testing.T) {
	kasus := []struct {
		nama  string
		hasil HasilPaketTerkait
		harap StatusPaketTerkait
	}{
		{nama: "aspirasi tanpa wilayah", hasil: HasilPaketTerkait{}, harap: StatusTanpaWilayah},
		{nama: "berwilayah tanpa paket cocok", hasil: HasilPaketTerkait{Aspirasi: AspirasiReses{KodeWilayah: "31.71.03"}}, harap: StatusBelumDitemukan},
		{
			nama:  "ada paket cocok",
			hasil: HasilPaketTerkait{Aspirasi: AspirasiReses{KodeWilayah: "31.71.03"}, Paket: []PaketTerkait{{}}},
			harap: StatusKemungkinanTerkait,
		},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if status := k.hasil.Status(); status != k.harap {
				t.Errorf("status = %s, harap %s", status, k.harap)
			}
		})
	}
}

func TestKonfigurasiAspirasiAnggaranValidasi(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*KonfigurasiAspirasiAnggaran)
		harap error
	}{
		{nama: "sah", ubah: func(*KonfigurasiAspirasiAnggaran) {}},
		{nama: "batas nol", ubah: func(k *KonfigurasiAspirasiAnggaran) { k.BatasPaket = 0 }, harap: ErrBatasPaketTerkaitTidakSah},
		{
			nama: "kata kunci huruf besar",
			ubah: func(k *KonfigurasiAspirasiAnggaran) {
				k.Bidang["ekonomi"] = AturanBidangAnggaran{KataKunci: []string{"Pasar"}}
			},
			harap: ErrKataKunciAnggaranTidakSah,
		},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			konfigurasi := konfigurasiAnggaranUji()
			k.ubah(&konfigurasi)
			if err := konfigurasi.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
