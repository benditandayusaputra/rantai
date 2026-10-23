package domain

import "testing"

func penggolongUji(t *testing.T) PenggolongKategori {
	t.Helper()
	p, err := PenggolongKategoriBaru([]AturanKategori{
		{Kategori: "ruang_kelas", Sektor: "pendidikan", Pola: []string{"(rehab|pembangunan).*(ruang kelas|gedung sekolah)", `\b(sdn|smpn)\b.*(rehab|pembangunan)`}},
		{Kategori: "pintu_air", Sektor: "air", Pola: []string{"pintu air"}},
		{Kategori: "saluran", Sektor: "air", Pola: []string{"saluran", "drainase"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestGolongkan_KecocokanPertamaMenang(t *testing.T) {
	kasus := []struct {
		nama  string
		harap string
	}{
		{"Rehab Total Gedung Sekolah SDN 01 Cakung", "ruang_kelas"},
		{"SMPN 123 Pembangunan Lanjutan", "ruang_kelas"},
		{"Pembangunan Pintu Air dan Saluran Kel. Semper Barat", "pintu_air"},
		{"Perbaikan Saluran Penghubung RW 05", "saluran"},
	}
	for _, k := range kasus {
		hasil, ada := penggolongUji(t).Golongkan(k.nama)
		if !ada || hasil.Kategori != k.harap {
			t.Errorf("%q = %+v, harap %s", k.nama, hasil, k.harap)
		}
	}
}

func TestGolongkan_TidakAdaPolaCocok_TidakBerkategori(t *testing.T) {
	if hasil, ada := penggolongUji(t).Golongkan("Pengadaan Alat Tulis Kantor"); ada {
		t.Errorf("hasil = %+v", hasil)
	}
}

func TestPenggolongKategoriBaru_RegexRusak_Ditolak(t *testing.T) {
	if _, err := PenggolongKategoriBaru([]AturanKategori{{Kategori: "x", Pola: []string{"(rusak"}}}); err == nil {
		t.Error("regex rusak diterima")
	}
}
