package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type AturanKategori struct {
	Kategori string   `yaml:"kategori"`
	Sektor   string   `yaml:"sektor"`
	Pola     []string `yaml:"pola"`
}

type Kategori struct {
	Kategori string
	Sektor   string
}

type PenggolongKategori struct {
	aturan []aturanTerkompilasi
}

type aturanTerkompilasi struct {
	Kategori
	pola []*regexp.Regexp
}

func PenggolongKategoriBaru(daftar []AturanKategori) (PenggolongKategori, error) {
	p := PenggolongKategori{}
	for _, a := range daftar {
		terkompilasi := aturanTerkompilasi{Kategori: Kategori{Kategori: a.Kategori, Sektor: a.Sektor}}
		for _, pola := range a.Pola {
			r, err := regexp.Compile(pola)
			if err != nil {
				return PenggolongKategori{}, fmt.Errorf("pola kategori %s %q: %w", a.Kategori, pola, err)
			}
			terkompilasi.pola = append(terkompilasi.pola, r)
		}
		p.aturan = append(p.aturan, terkompilasi)
	}
	return p, nil
}

func (p PenggolongKategori) Golongkan(namaPaket string) (Kategori, bool) {
	kecil := strings.ToLower(namaPaket)
	for _, a := range p.aturan {
		for _, pola := range a.pola {
			if pola.MatchString(kecil) {
				return a.Kategori, true
			}
		}
	}
	return Kategori{}, false
}
