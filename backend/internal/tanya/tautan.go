package tanya

import (
	"cmp"
	"encoding/json"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	TautanAnggota = "anggota"
	TautanDapil   = "dapil"
	TautanRilis   = "rilis"

	panjangTeksTautanMinimal = 4
	awalanURLAman            = "https://"
)

type Tautan struct {
	Teks   string
	Jenis  string
	Tujuan string
}

func SusunTautan(teks string, hasilAlat []json.RawMessage) []Tautan {
	temuan := map[string]Tautan{}
	for _, hasil := range hasilAlat {
		var isi any
		if err := json.Unmarshal(hasil, &isi); err == nil {
			kumpulkanTautan(isi, temuan)
		}
	}
	daftar := []Tautan{}
	for _, t := range temuan {
		if strings.Contains(teks, t.Teks) {
			daftar = append(daftar, t)
		}
	}
	slices.SortFunc(daftar, func(a, b Tautan) int {
		return cmp.Or(cmp.Compare(len(b.Teks), len(a.Teks)), cmp.Compare(a.Teks, b.Teks))
	})
	return daftar
}

func kumpulkanTautan(isi any, temuan map[string]Tautan) {
	switch nilai := isi.(type) {
	case []any:
		for _, anak := range nilai {
			kumpulkanTautan(anak, temuan)
		}
	case map[string]any:
		for _, t := range tautanObjek(nilai) {
			if utf8.RuneCountInString(t.Teks) >= panjangTeksTautanMinimal && t.Tujuan != "" {
				temuan[t.Teks] = t
			}
		}
		for _, anak := range nilai {
			kumpulkanTautan(anak, temuan)
		}
	}
}

func tautanObjek(objek map[string]any) []Tautan {
	daftar := []Tautan{}
	if _, isAnggota := objek["fraksi"]; isAnggota {
		daftar = append(daftar, Tautan{Teks: teksObjek(objek, "nama"), Jenis: TautanAnggota, Tujuan: teksObjek(objek, "id")})
	}
	if url := teksObjek(objek, "url"); strings.HasPrefix(url, awalanURLAman) {
		daftar = append(daftar, Tautan{Teks: teksObjek(objek, "judul"), Jenis: TautanRilis, Tujuan: url})
	}
	dapil, isBerdapil := objek["dapil"].(map[string]any)
	kecamatan, isBerkecamatan := objek["kecamatan"].(map[string]any)
	if isBerdapil && isBerkecamatan {
		daftar = append(daftar, Tautan{Teks: teksObjek(dapil, "nama"), Jenis: TautanDapil, Tujuan: teksObjek(kecamatan, "kode")})
	}
	return daftar
}

func teksObjek(objek map[string]any, kunci string) string {
	teks, _ := objek[kunci].(string)
	return teks
}
