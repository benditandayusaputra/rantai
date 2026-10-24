package tanya

import (
	"regexp"
	"strconv"
)

var polaRujukanProvenance = regexp.MustCompile(`\[(prv_[A-Za-z0-9_]+)\]`)

func SusunRujukan(teks string) (string, []Rujukan) {
	nomor := map[string]int{}
	rujukan := []Rujukan{}
	hasil := polaRujukanProvenance.ReplaceAllStringFunc(teks, func(cocok string) string {
		id := polaRujukanProvenance.FindStringSubmatch(cocok)[1]
		if _, ada := nomor[id]; !ada {
			nomor[id] = len(rujukan) + 1
			rujukan = append(rujukan, Rujukan{Nomor: nomor[id], ProvenanceID: id})
		}
		return "[" + strconv.Itoa(nomor[id]) + "]"
	})
	return hasil, rujukan
}
