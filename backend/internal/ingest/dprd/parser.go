package dprd

import (
	"bytes"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	awalanBidang      = "Meliputi"
	kataPenutupDaftar = "serta"
)

var polaSlugProfil = regexp.MustCompile(`/profile/([^/?#]+)/?$`)

func UraiHalamanDapil(isi []byte, nomorDapil int) ([]domain.BarisDapil, error) {
	akar, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, fmt.Errorf("urai halaman dapil %d: %w", nomorDapil, err)
	}
	daftar := []domain.BarisDapil{}
	for n := range akar.Descendants() {
		slug := slugProfil(n)
		if slug == "" || adaAnak(n, atom.Img) {
			continue
		}
		daftar = append(daftar, domain.BarisDapil{
			NomorDapil: nomorDapil,
			Nama:       teks(n),
			Partai:     teks(saudaraBerikut(n.Parent)),
			Slug:       slug,
		})
	}
	return daftar, nil
}

func UraiHalamanStruktur(isi []byte) ([]domain.BarisStruktur, error) {
	akar, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, fmt.Errorf("urai halaman struktur: %w", err)
	}
	daftar := []domain.BarisStruktur{}
	for n := range akar.Descendants() {
		jabatan := saudaraBerikut(n)
		if n.DataAtom != atom.H3 || jabatan == nil || jabatan.DataAtom != atom.H4 {
			continue
		}
		daftar = append(daftar, domain.BarisStruktur{
			Nama:    teks(n),
			Jabatan: teks(jabatan),
			Slug:    slugDalam(saudaraSebelum(n)),
		})
	}
	return daftar, nil
}

func slugProfil(n *html.Node) string {
	if n.DataAtom != atom.A {
		return ""
	}
	for _, atribut := range n.Attr {
		if cocok := polaSlugProfil.FindStringSubmatch(atribut.Val); atribut.Key == "href" && cocok != nil {
			return cocok[1]
		}
	}
	return ""
}

func slugDalam(n *html.Node) string {
	if n == nil {
		return ""
	}
	for keturunan := range n.Descendants() {
		if slug := slugProfil(keturunan); slug != "" {
			return slug
		}
	}
	return ""
}

func adaAnak(n *html.Node, jenis atom.Atom) bool {
	for keturunan := range n.Descendants() {
		if keturunan.DataAtom == jenis {
			return true
		}
	}
	return false
}

func saudaraBerikut(n *html.Node) *html.Node {
	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if s.Type == html.ElementNode {
			return s
		}
	}
	return nil
}

func saudaraSebelum(n *html.Node) *html.Node {
	for s := n.PrevSibling; s != nil; s = s.PrevSibling {
		if s.Type == html.ElementNode {
			return s
		}
	}
	return nil
}

func teks(n *html.Node) string {
	if n == nil {
		return ""
	}
	var isi strings.Builder
	for keturunan := range n.Descendants() {
		if keturunan.Type == html.TextNode {
			isi.WriteString(keturunan.Data + " ")
		}
	}
	return strings.Join(strings.Fields(isi.String()), " ")
}

func UraiBidangKomisi(isi []byte, kataAwalUnit []string) ([]string, error) {
	akar, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, fmt.Errorf("urai halaman komisi: %w", err)
	}
	for n := range akar.Descendants() {
		if n.DataAtom != atom.P {
			continue
		}
		if daftar, ada := strings.CutPrefix(teks(n), awalanBidang); ada {
			return pisahkanUnit(strings.TrimLeft(daftar, " :"), kataAwalUnit), nil
		}
	}
	return []string{}, nil
}

func pisahkanUnit(daftar string, kataAwalUnit []string) []string {
	unit := []string{}
	for _, potongan := range strings.Split(strings.TrimSuffix(daftar, "."), ",") {
		potongan = strings.TrimSpace(potongan)
		kataPertama, _, _ := strings.Cut(potongan, " ")
		switch {
		case kataPertama == kataPenutupDaftar:
			continue
		case len(unit) == 0 || slices.Contains(kataAwalUnit, kataPertama):
			unit = append(unit, potongan)
		default:
			unit[len(unit)-1] += ", " + potongan
		}
	}
	return unit
}
