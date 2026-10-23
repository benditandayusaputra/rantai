package lpse

import (
	"bytes"
	"slices"
	"strings"

	"golang.org/x/net/html"
)

const kelasLabel = "bgwarning"

type tabelLabel map[string]*html.Node

func bacaTabelLabel(isi []byte) (tabelLabel, error) {
	dokumen, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, err
	}
	tabel := tabelLabel{}
	for _, th := range cariSemua(dokumen, isLabel) {
		if nilai := elemenBerikut(th); nilai != nil {
			tabel[teks(th)] = nilai
		}
	}
	return tabel, nil
}

func (t tabelLabel) teks(label string) string {
	if simpul, ada := t[label]; ada {
		return teks(simpul)
	}
	return ""
}

func isLabel(n *html.Node) bool {
	return isElemen(n, "th") && slices.Contains(strings.Fields(atribut(n, "class")), kelasLabel)
}

func isElemen(n *html.Node, tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func isLencana(n *html.Node) bool {
	return isElemen(n, "span") && slices.Contains(strings.Fields(atribut(n, "class")), "badge")
}

func atribut(n *html.Node, kunci string) string {
	for _, a := range n.Attr {
		if a.Key == kunci {
			return a.Val
		}
	}
	return ""
}

func cariSemua(akar *html.Node, cocok func(*html.Node) bool) []*html.Node {
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

func teks(n *html.Node) string {
	return rapikan(kumpulkanTeks(n, func(*html.Node) bool { return false }))
}

func teksTanpaLencana(n *html.Node) string {
	return rapikan(kumpulkanTeks(n, isLencana))
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

func rapikan(s string) string {
	return strings.Join(strings.Fields(strings.ReplaceAll(s, " ", " ")), " ")
}

func selBaris(baris *html.Node) []string {
	sel := []string{}
	for anak := baris.FirstChild; anak != nil; anak = anak.NextSibling {
		if isElemen(anak, "td") {
			sel = append(sel, teks(anak))
		}
	}
	return sel
}
