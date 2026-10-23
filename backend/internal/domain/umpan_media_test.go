package domain

import (
	"slices"
	"testing"
)

func aturanMediaUji() AturanMedia {
	return AturanMedia{
		DomainSah:   []string{"contohberita.id"},
		KonteksDPRD: []string{"DPRD DKI", "DPRD Jakarta"},
		AliasFraksi: map[string][]string{"Fraksi Partai Pelangi": {"Partai Pelangi", "PPL"}},
	}
}

func pencariMediaUji() PencariNama {
	return BuatPencariNama([]Anggota{{ID: "budi-santosa", Nama: "Budi Santosa, S.H.", Fraksi: "Fraksi Partai Pelangi"}})
}

func TestAturanMedia_PautanArtikel_SesuaiKonteks(t *testing.T) {
	kasus := []struct {
		nama  string
		item  ItemUmpan
		harap []PautanUmpan
	}{
		{"nama sama tanpa DPRD DKI dibuang", ItemUmpan{URL: "https://www.contohberita.id/a", Judul: "Budi Santosa juara lomba lari maraton"}, []PautanUmpan{}},
		{"nama dengan DPRD DKI dipautkan", ItemUmpan{URL: "https://contohberita.id/b", Judul: "Anggota DPRD DKI Budi Santosa tinjau banjir"},
			[]PautanUmpan{{AnggotaID: "budi-santosa", DasarCocok: DasarNamaDanDPRD}}},
		{"nama dengan fraksi di cuplikan dipautkan", ItemUmpan{URL: "https://contohberita.id/c", Judul: "Budi Santosa soroti tarif bus", Cuplikan: "Politikus PPL itu meminta kajian ulang."},
			[]PautanUmpan{{AnggotaID: "budi-santosa", DasarCocok: DasarNamaDanFraksi}}},
		{"domain di luar daftar ditolak", ItemUmpan{URL: "https://contohberita.id.palsu.net/d", Judul: "Anggota DPRD DKI Budi Santosa tinjau banjir"}, []PautanUmpan{}},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if pautan := aturanMediaUji().Pautkan(pencariMediaUji(), k.item); !slices.Equal(pautan, k.harap) {
				t.Errorf("pautan = %v, harap %v", pautan, k.harap)
			}
		})
	}
}

func TestIsDomainSah_SubdomainDanSkema_Diperiksa(t *testing.T) {
	aturan := aturanMediaUji()
	kasus := map[string]bool{
		"https://megapolitan.contohberita.id/x": true,
		"https://contohberita.id/x":             true,
		"https://bukancontohberita.id/x":        false,
		"ftp://contohberita.id/x":               false,
		"bukan url":                             false,
	}
	for alamat, harap := range kasus {
		if aturan.IsDomainSah(alamat) != harap {
			t.Errorf("IsDomainSah(%q) harap %v", alamat, harap)
		}
	}
}
