package domain

import (
	"net/url"
	"slices"
	"strings"
)

const SumberMedia SumberData = "media"

type AturanMedia struct {
	DomainSah   []string
	KonteksDPRD []string
	AliasFraksi map[string][]string
}

func (a AturanMedia) Pautkan(pencari PencariNama, item ItemUmpan) []PautanUmpan {
	pautan := []PautanUmpan{}
	if !a.IsDomainSah(item.URL) {
		return pautan
	}
	teks := item.Judul + " " + item.Cuplikan
	adaKonteksDPRD := memuatSalahSatu(teks, a.KonteksDPRD)
	for _, id := range pencari.AnggotaDisebut(teks) {
		switch {
		case adaKonteksDPRD:
			pautan = append(pautan, PautanUmpan{AnggotaID: id, DasarCocok: DasarNamaDanDPRD})
		case memuatSalahSatu(teks, a.AliasFraksi[pencari.fraksiPerAnggota[id]]):
			pautan = append(pautan, PautanUmpan{AnggotaID: id, DasarCocok: DasarNamaDanFraksi})
		}
	}
	return pautan
}

func (a AturanMedia) IsDomainSah(alamat string) bool {
	u, err := url.Parse(alamat)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return false
	}
	host := strings.ToLower(u.Hostname())
	return slices.ContainsFunc(a.DomainSah, func(domain string) bool {
		return host == domain || strings.HasSuffix(host, "."+domain)
	})
}

func memuatSalahSatu(teks string, frasa []string) bool {
	normal := " " + NormalkanTeks(teks) + " "
	return slices.ContainsFunc(frasa, func(f string) bool { return strings.Contains(normal, " "+NormalkanTeks(f)+" ") })
}
