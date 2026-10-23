package halaman

import (
	"bytes"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

type TabelLabel map[string]*html.Node

func BacaTabelLabel(isi []byte, tagLabel, kelasLabel string) (TabelLabel, error) {
	dokumen, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, err
	}
	tabel := TabelLabel{}
	isLabel := func(n *html.Node) bool { return IsElemen(n, tagLabel) && BerkelasCSS(n, kelasLabel) }
	for _, label := range CariSemua(dokumen, isLabel) {
		if _, sudahAda := tabel[Teks(label)]; sudahAda {
			continue
		}
		if nilai := elemenBerikut(label); nilai != nil {
			tabel[Teks(label)] = nilai
		}
	}
	return tabel, nil
}

func (t TabelLabel) Teks(label string) string {
	if simpul, ada := t[label]; ada {
		return Teks(simpul)
	}
	return ""
}

func BerkelasCSS(n *html.Node, kelas string) bool {
	return slices.Contains(strings.Fields(atribut(n, "class")), kelas)
}

func IsElemen(n *html.Node, tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func IsLencana(n *html.Node) bool {
	return IsElemen(n, "span") && BerkelasCSS(n, "badge")
}

func atribut(n *html.Node, kunci string) string {
	for _, a := range n.Attr {
		if a.Key == kunci {
			return a.Val
		}
	}
	return ""
}

func CariSemua(akar *html.Node, cocok func(*html.Node) bool) []*html.Node {
	hasil := []*html.Node{}
	for n := range akar.Descendants() {
		if cocok(n) {
			hasil = append(hasil, n)
		}
	}
	return hasil
}

func elemenBerikut(n *html.Node) *html.Node {
	for s := n.NextSibling; s != nil; s = s.NextSibling {
		if s.Type == html.ElementNode {
			return s
		}
	}
	return nil
}

func Teks(n *html.Node) string {
	return Rapikan(kumpulkanTeks(n, func(*html.Node) bool { return false }))
}

func TeksTanpaLencana(n *html.Node) string {
	return Rapikan(kumpulkanTeks(n, IsLencana))
}

func kumpulkanTeks(n *html.Node, lewati func(*html.Node) bool) string {
	if n == nil {
		return ""
	}
	var b strings.Builder
	var telusuri func(*html.Node)
	telusuri = func(simpul *html.Node) {
		if lewati(simpul) {
			return
		}
		if simpul.Type == html.TextNode {
			b.WriteString(simpul.Data)
			b.WriteString(" ")
		}
		for anak := simpul.FirstChild; anak != nil; anak = anak.NextSibling {
			telusuri(anak)
		}
	}
	telusuri(n)
	return b.String()
}

func Rapikan(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, " ", " ")), " ")
}

func SelBaris(baris *html.Node) []string {
	sel := []string{}
	for anak := baris.FirstChild; anak != nil; anak = anak.NextSibling {
		if IsElemen(anak, "td") {
			sel = append(sel, Teks(anak))
		}
	}
	return sel
}
