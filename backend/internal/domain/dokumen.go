package domain

import "strings"

type BagianDokumen struct {
	Judul string
	Baris []string
}

type Dokumen struct {
	Judul  string
	Bagian []BagianDokumen
}

func (d Dokumen) Teks() string {
	var b strings.Builder
	b.WriteString(d.Judul + "\n")
	for _, bagian := range d.Bagian {
		b.WriteString("\n" + bagian.Judul + "\n")
		for _, baris := range bagian.Baris {
			b.WriteString(baris + "\n")
		}
	}
	return b.String()
}
