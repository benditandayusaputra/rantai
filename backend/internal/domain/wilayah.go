package domain

import (
	"regexp"
	"slices"
	"strings"
)

type Tingkat string

const (
	TingkatProvinsi  Tingkat = "provinsi"
	TingkatKota      Tingkat = "kota"
	TingkatKecamatan Tingkat = "kecamatan"
	TingkatKelurahan Tingkat = "kelurahan"
)

const pemisahKode = "."

var polaKodeWilayah = regexp.MustCompile(`^\d{2}(\.\d{2}(\.\d{2}(\.\d{4})?)?)?$`)

var tingkatMenurutJumlahSegmen = map[int]Tingkat{
	1: TingkatProvinsi,
	2: TingkatKota,
	3: TingkatKecamatan,
	4: TingkatKelurahan,
}

type Koordinat [2]float64

type MultiPolygon [][][]Koordinat

type Wilayah struct {
	Kode      string
	Nama      string
	Tingkat   Tingkat
	IndukKode string
	Alias     []string
	Geometri  MultiPolygon
}

type AturanAlias struct {
	PrefiksPerTingkat map[Tingkat]map[string][]string
	Singkatan         map[string]string
}

func TingkatDariKode(kode string) (Tingkat, error) {
	if !polaKodeWilayah.MatchString(kode) {
		return "", InputTidakValid("kode wilayah tidak valid: " + kode)
	}
	return tingkatMenurutJumlahSegmen[strings.Count(kode, pemisahKode)+1], nil
}

func IndukKode(kode string) string {
	posisi := strings.LastIndex(kode, pemisahKode)
	if posisi < 0 {
		return ""
	}
	return kode[:posisi]
}

func (a AturanAlias) BuatAlias(nama string, tingkat Tingkat) []string {
	prefiksResmi, prefiksAlias := a.cocokkanPrefiksResmi(nama, tingkat)
	inti := strings.TrimPrefix(nama, prefiksResmi)
	alias := []string{}
	for _, dasar := range []string{inti, a.singkat(inti)} {
		for _, prefiks := range prefiksAlias {
			alias = append(alias, prefiks+dasar)
		}
	}
	slices.Sort(alias)
	alias = slices.Compact(alias)
	return slices.DeleteFunc(alias, func(s string) bool { return s == nama })
}

func (a AturanAlias) cocokkanPrefiksResmi(nama string, tingkat Tingkat) (string, []string) {
	terpilih := ""
	var prefiksAlias []string
	for prefiks, daftarAlias := range a.PrefiksPerTingkat[tingkat] {
		if strings.HasPrefix(nama, prefiks) && len(prefiks) >= len(terpilih) {
			terpilih, prefiksAlias = prefiks, daftarAlias
		}
	}
	return terpilih, prefiksAlias
}

func (a AturanAlias) singkat(nama string) string {
	kata := strings.Fields(nama)
	for i, k := range kata {
		if singkatan, ada := a.Singkatan[k]; ada {
			kata[i] = singkatan
		}
	}
	return strings.Join(kata, " ")
}

type Titik struct {
	Lat float64
	Lng float64
}

type RingkasanWilayah struct {
	Kode    string
	Nama    string
	Tingkat Tingkat
}

type WilayahDitemukan struct {
	RingkasanWilayah
	Hierarki     []RingkasanWilayah
	ProvenanceID string
}

type DetailWilayah struct {
	WilayahDitemukan
	Alias       []string
	TitikTengah Titik
}
