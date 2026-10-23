package domain

import (
	"slices"
	"testing"
)

func penandaUji(t *testing.T) PenandaTema {
	t.Helper()
	p, err := PenandaTemaBaru([]Tema{
		{Kode: "pendidikan", KonfigurasiTema: KonfigurasiTema{FilterAnggaran: FilterAnggaran{
			InstansiPola: []string{"dinas pendidikan", "suku dinas pendidikan"},
			KataKunci:    []string{"ruang kelas", "rehab sekolah", `\b(sdn|smpn)\b`},
		}}},
		{Kode: "banjir", KonfigurasiTema: KonfigurasiTema{FilterAnggaran: FilterAnggaran{KataKunci: []string{"pintu air", "saluran"}}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestTandai_FilterAnggaranTema(t *testing.T) {
	kasus := []struct {
		nama     string
		instansi string
		harap    []string
	}{
		{"Pengadaan Alat Tulis", "SUKU DINAS PENDIDIKAN WILAYAH I KOTA ADMINISTRASI JAKARTA UTARA", []string{"pendidikan"}},
		{"Pembangunan Ruang Kelas Baru SDN 01", "Provinsi DKI Jakarta", []string{"pendidikan"}},
		{"Rehab SMPN 12 dan Saluran Halaman", "Provinsi DKI Jakarta", []string{"pendidikan", "banjir"}},
		{"Pembangunan Pintu Air Kel. Semper Barat", "SUKU DINAS SUMBER DAYA AIR", []string{"banjir"}},
		{"Perbaikan Gedung SDNegeri", "Provinsi DKI Jakarta", []string{}},
	}
	for _, k := range kasus {
		hasil := penandaUji(t).Tandai(MasukanTema{NamaPaket: k.nama, NamaInstansi: k.instansi})
		if !slices.Equal(hasil, k.harap) {
			t.Errorf("%q = %v, harap %v", k.nama, hasil, k.harap)
		}
	}
}
