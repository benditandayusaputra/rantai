package domain

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func aturanRingkasanUji() AturanRingkasan {
	return AturanRingkasan{
		Terlarang:     KosakataTerlarang{"korupsi", "mark up", "memperjuangkan"},
		KataPenilaian: []string{"hebat", "luar biasa", "gagal total"},
	}
}

var sumberRingkasanUji = UmpanUntukRingkas{
	Judul:    "Kasus DBD Jakarta Selatan Tembus 1.314",
	Cuplikan: "Ketua Komisi D meminta Sudinkes memperketat pengawasan DBD hingga September 2026.",
}

func TestValidasiRingkasan_KalimatNetralBersumber_Lolos(t *testing.T) {
	ringkasan := "Ketua Komisi D meminta pengawasan DBD di Jakarta Selatan diperketat karena kasus mencapai 1.314."
	if err := aturanRingkasanUji().Validasi(ringkasan, sumberRingkasanUji); err != nil {
		t.Errorf("harap lolos, dapat %v", err)
	}
}

func TestValidasiRingkasan_Pelanggaran_Ditolak(t *testing.T) {
	kasus := []struct {
		nama      string
		ringkasan string
		harap     error
	}{
		{"kosong", "  ", ErrRingkasanKosong},
		{"terlalu panjang", strings.Repeat("kata ", 40), ErrRingkasanTerlaluPanjang},
		{"dua kalimat", "Komisi D meminta pengawasan. Sudinkes diminta aktif.", ErrRingkasanLebihSatu},
		{"kosakata terlarang", "Komisi D menduga ada korupsi dalam penanganan DBD.", ErrRingkasanTerlarang},
		{"memperjuangkan", "Anggota memperjuangkan pengawasan DBD.", ErrRingkasanTerlarang},
		{"kata penilaian", "Langkah hebat Komisi D soal DBD.", ErrRingkasanMenilai},
		{"angka karangan", "Komisi D menyebut kasus DBD mencapai 2.500.", ErrRingkasanAngkaKarangan},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := aturanRingkasanUji().Validasi(k.ringkasan, sumberRingkasanUji); !errors.Is(err, k.harap) {
				t.Errorf("harap %v, dapat %v", k.harap, err)
			}
		})
	}
}

func TestKosakataTerlarang_FrasaUtuh_DitemukanTanpaPotonganKata(t *testing.T) {
	kosakata := KosakataTerlarang{"mark up", "fiktif"}

	if hasil := kosakata.Temukan("Ada MARK-UP harga dan proyek fiktif"); !slices.Equal(hasil, []string{"fiktif", "mark up"}) {
		t.Errorf("hasil = %v", hasil)
	}
	if hasil := kosakata.Temukan("markup bahasa dan fiktifnya"); len(hasil) != 0 {
		t.Errorf("potongan kata tidak boleh cocok, dapat %v", hasil)
	}
}
