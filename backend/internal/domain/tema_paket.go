package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type MasukanTema struct {
	KodePaket    string
	NamaPaket    string
	NamaInstansi string
}

type PenandaTema struct {
	filter []filterTerkompilasi
}

type filterTerkompilasi struct {
	kodeTema     string
	instansiPola []string
	kataKunci    []*regexp.Regexp
}

func PenandaTemaBaru(daftar []Tema) (PenandaTema, error) {
	p := PenandaTema{}
	for _, t := range daftar {
		f := filterTerkompilasi{kodeTema: t.Kode}
		for _, pola := range t.FilterAnggaran.InstansiPola {
			f.instansiPola = append(f.instansiPola, strings.ToLower(pola))
		}
		for _, kata := range t.FilterAnggaran.KataKunci {
			r, err := regexp.Compile("(?i)" + kata)
			if err != nil {
				return PenandaTema{}, fmt.Errorf("filter anggaran tema %s: %w", t.Kode, err)
			}
			f.kataKunci = append(f.kataKunci, r)
		}
		p.filter = append(p.filter, f)
	}
	return p, nil
}

func (p PenandaTema) Tandai(m MasukanTema) []string {
	kode := []string{}
	for _, f := range p.filter {
		if f.cocok(m) {
			kode = append(kode, f.kodeTema)
		}
	}
	return kode
}

func (f filterTerkompilasi) cocok(m MasukanTema) bool {
	instansi := strings.ToLower(m.NamaInstansi)
	for _, pola := range f.instansiPola {
		if strings.Contains(instansi, pola) {
			return true
		}
	}
	for _, r := range f.kataKunci {
		if r.MatchString(m.NamaPaket) {
			return true
		}
	}
	return false
}
