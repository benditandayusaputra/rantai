package domain

import (
	"slices"
	"strings"
)

type KosakataTerlarang []string

func (k KosakataTerlarang) Temukan(teks string) []string {
	normal := " " + NormalkanTeks(teks) + " "
	ditemukan := []string{}
	for _, kata := range k {
		if strings.Contains(normal, " "+NormalkanTeks(kata)+" ") {
			ditemukan = append(ditemukan, kata)
		}
	}
	slices.Sort(ditemukan)
	return ditemukan
}
