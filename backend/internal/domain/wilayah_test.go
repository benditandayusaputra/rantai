package domain

import (
	"errors"
	"slices"
	"testing"
)

func aturanAliasUji() AturanAlias {
	return AturanAlias{
		PrefiksPerTingkat: map[Tingkat]map[string][]string{
			TingkatKota: {
				"Kota Administrasi ":      {"", "Kota ", "Kota Adm. "},
				"Kabupaten Administrasi ": {"", "Kab. Adm. "},
			},
			TingkatKelurahan: {"": {"", "Kel. ", "Kelurahan "}},
		},
		Singkatan: map[string]string{"Kelapa": "Klp.", "Tanjung": "Tj."},
	}
}

func TestTingkatDariKode_FormatKemendagri_TingkatSesuai(t *testing.T) {
	kasus := []struct {
		kode    string
		tingkat Tingkat
	}{
		{"31", TingkatProvinsi},
		{"31.72", TingkatKota},
		{"31.72.06", TingkatKecamatan},
		{"31.72.06.1003", TingkatKelurahan},
	}
	for _, k := range kasus {
		t.Run(k.kode, func(t *testing.T) {
			tingkat, err := TingkatDariKode(k.kode)

			if err != nil || tingkat != k.tingkat {
				t.Errorf("dapat (%q, %v), harap %q", tingkat, err, k.tingkat)
			}
		})
	}
}

func TestTingkatDariKode_KodeTidakValid_ErrorInputTidakValid(t *testing.T) {
	for _, kode := range []string{"31.72.06.1003.9", "31..06", ""} {
		t.Run(kode, func(t *testing.T) {
			_, err := TingkatDariKode(kode)

			var kesalahan *Kesalahan
			if !errors.As(err, &kesalahan) || kesalahan.Kategori != KategoriInputTidakValid {
				t.Errorf("harap error input tidak valid, dapat %v", err)
			}
		})
	}
}

func TestIndukKode_TiapTingkat_KodeIndukLangsung(t *testing.T) {
	kasus := map[string]string{"31.72.06.1003": "31.72.06", "31.72": "31", "31": ""}
	for kode, harap := range kasus {
		if induk := IndukKode(kode); induk != harap {
			t.Errorf("IndukKode(%q) = %q, harap %q", kode, induk, harap)
		}
	}
}

func TestBuatAlias_NamaDanTingkat_AliasSesuaiAturan(t *testing.T) {
	kasus := []struct {
		nama    string
		tingkat Tingkat
		harap   []string
	}{
		{"Kelapa Gading Barat", TingkatKelurahan, []string{
			"Kel. Kelapa Gading Barat", "Kel. Klp. Gading Barat", "Kelurahan Kelapa Gading Barat",
			"Kelurahan Klp. Gading Barat", "Klp. Gading Barat",
		}},
		{"Pluit", TingkatKelurahan, []string{"Kel. Pluit", "Kelurahan Pluit"}},
		{"Kota Administrasi Jakarta Utara", TingkatKota, []string{
			"Jakarta Utara", "Kota Adm. Jakarta Utara", "Kota Jakarta Utara",
		}},
		{"Kabupaten Administrasi Kepulauan Seribu", TingkatKota, []string{
			"Kab. Adm. Kepulauan Seribu", "Kepulauan Seribu",
		}},
		{"Gambir", TingkatKecamatan, []string{}},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			alias := aturanAliasUji().BuatAlias(k.nama, k.tingkat)

			if !slices.Equal(alias, k.harap) {
				t.Errorf("dapat %q, harap %q", alias, k.harap)
			}
		})
	}
}
