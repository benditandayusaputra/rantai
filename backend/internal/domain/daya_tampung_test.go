package domain

import (
	"slices"
	"testing"
)

func TestPisahkanTerpetakan_NPSNTakAdaAtauKosong_DikeluarkanDanDilaporkan(t *testing.T) {
	daftar := []BarisDayaTampung{
		{NPSN: "20100251", Nama: "SMP NEGERI 1 JAKARTA", Kursi: 252},
		{NPSN: "", Nama: "SMP NEGERI 296 JAKARTA", Kursi: 108},
		{NPSN: "99999999", Nama: "SMP TIDAK ADA", Kursi: 36},
	}

	terpetakan, tidak := PisahkanTerpetakan(daftar, map[string]bool{"20100251": true})

	if len(terpetakan) != 1 || terpetakan[0].NPSN != "20100251" {
		t.Errorf("terpetakan = %+v", terpetakan)
	}
	nama := []string{tidak[0].Nama, tidak[1].Nama}
	if !slices.Equal(nama, []string{"SMP NEGERI 296 JAKARTA", "SMP TIDAK ADA"}) || TotalKursi(tidak) != 144 {
		t.Errorf("tidak terpetakan = %+v", tidak)
	}
}

func TestPisahkanTerpetakan_DaftarKosong_KoleksiKosong(t *testing.T) {
	terpetakan, tidak := PisahkanTerpetakan(nil, map[string]bool{})

	if terpetakan == nil || tidak == nil || TotalKursi(nil) != 0 {
		t.Errorf("harap koleksi kosong, dapat %v %v", terpetakan, tidak)
	}
}
