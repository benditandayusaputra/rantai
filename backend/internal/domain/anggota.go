package domain

import (
	"fmt"
	"slices"
	"strings"
)

type pencocokAnggota struct {
	dapilPerSlug   map[string][]BarisDapil
	dapilPerKunci  map[string][]BarisDapil
	fraksiPerSlug  map[string][]AnggotaFraksi
	fraksiPerKunci map[string][]AnggotaFraksi
}

func GabungkanAnggota(barisDapil []BarisDapil, anggotaFraksi []AnggotaFraksi) ([]Anggota, error) {
	pencocok := pencocokAnggota{
		dapilPerSlug:   kelompokkan(barisDapil, func(b BarisDapil) string { return b.Slug }),
		dapilPerKunci:  kelompokkan(barisDapil, func(b BarisDapil) string { return KunciNama(b.Nama) }),
		fraksiPerSlug:  kelompokkan(anggotaFraksi, func(f AnggotaFraksi) string { return f.Slug }),
		fraksiPerKunci: kelompokkan(anggotaFraksi, func(f AnggotaFraksi) string { return KunciNama(f.Nama) }),
	}
	hasil, tidakCocok := pencocok.cocokkanSemua(barisDapil)
	if len(tidakCocok) > 0 {
		return nil, fmt.Errorf("%s: %w", strings.Join(tidakCocok, "; "), ErrAnggotaTidakCocok)
	}
	return hasil, nil
}

func (p pencocokAnggota) cocokkanSemua(barisDapil []BarisDapil) ([]Anggota, []string) {
	hasil := []Anggota{}
	tidakCocok := []string{}
	dipakai := map[string]bool{}
	for _, b := range barisDapil {
		f, ada := p.cariFraksi(b)
		id := p.tentukanID(b, f)
		if !ada || id == "" || dipakai[id] {
			tidakCocok = append(tidakCocok, fmt.Sprintf("dapil %d %q", b.NomorDapil, b.Nama))
			continue
		}
		dipakai[id] = true
		hasil = append(hasil, buatAnggota(id, b, f))
	}
	return hasil, tidakCocok
}

func (p pencocokAnggota) tentukanID(b BarisDapil, f AnggotaFraksi) string {
	if f.Slug != "" {
		return f.Slug
	}
	if len(p.dapilPerSlug[b.Slug]) == 1 {
		return b.Slug
	}
	return ""
}

func (p pencocokAnggota) cariFraksi(b BarisDapil) (AnggotaFraksi, bool) {
	if b.Slug != "" && len(p.dapilPerSlug[b.Slug]) == 1 && len(p.fraksiPerSlug[b.Slug]) == 1 {
		return p.fraksiPerSlug[b.Slug][0], true
	}
	kunci := KunciNama(b.Nama)
	if len(p.dapilPerKunci[kunci]) == 1 && len(p.fraksiPerKunci[kunci]) == 1 {
		return p.fraksiPerKunci[kunci][0], true
	}
	return AnggotaFraksi{}, false
}

func buatAnggota(id string, b BarisDapil, f AnggotaFraksi) Anggota {
	namaLain := []string{NamaTanpaGelar(b.Nama), f.Nama, NamaTanpaGelar(f.Nama)}
	slices.Sort(namaLain)
	namaLain = slices.DeleteFunc(slices.Compact(namaLain), func(n string) bool { return n == b.Nama || n == "" })
	return Anggota{
		ID:         id,
		Nama:       b.Nama,
		NamaLain:   namaLain,
		Partai:     b.Partai,
		Fraksi:     f.Fraksi,
		NomorDapil: b.NomorDapil,
		FotoURL:    f.Foto,
	}
}

func kelompokkan[T any](daftar []T, kunci func(T) string) map[string][]T {
	hasil := map[string][]T{}
	for _, item := range daftar {
		k := kunci(item)
		hasil[k] = append(hasil[k], item)
	}
	return hasil
}
