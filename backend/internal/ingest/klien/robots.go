package klien

import (
	"slices"
	"strings"
)

const agenSemua = "*"

type aturanPath struct {
	prefiks   string
	isIzinkan bool
}

type aturanRobots struct {
	daftar []aturanPath
}

type grupRobots struct {
	agen   []string
	aturan []aturanPath
}

func uraiRobots(isi, userAgent string) aturanRobots {
	grup := uraiGrup(isi)
	agenKita := tokenAgen(userAgent)
	for _, g := range grup {
		if slices.Contains(g.agen, agenKita) {
			return aturanRobots{daftar: g.aturan}
		}
	}
	for _, g := range grup {
		if slices.Contains(g.agen, agenSemua) {
			return aturanRobots{daftar: g.aturan}
		}
	}
	return aturanRobots{}
}

func (a aturanRobots) Mengizinkan(path string) bool {
	cocok := aturanPath{isIzinkan: true}
	for _, aturan := range a.daftar {
		if !strings.HasPrefix(path, aturan.prefiks) {
			continue
		}
		lebihPanjang := len(aturan.prefiks) > len(cocok.prefiks)
		samaPanjangIzinkan := len(aturan.prefiks) == len(cocok.prefiks) && aturan.isIzinkan
		if lebihPanjang || samaPanjangIzinkan {
			cocok = aturan
		}
	}
	return cocok.isIzinkan
}

func uraiGrup(isi string) []grupRobots {
	daftar := []grupRobots{}
	sedangAgen := false
	for _, baris := range strings.Split(isi, "\n") {
		kunci, nilai, ada := pecahBaris(baris)
		if !ada {
			continue
		}
		if kunci == "user-agent" {
			if !sedangAgen {
				daftar = append(daftar, grupRobots{})
			}
			daftar[len(daftar)-1].agen = append(daftar[len(daftar)-1].agen, strings.ToLower(nilai))
			sedangAgen = true
			continue
		}
		sedangAgen = false
		if aturan, sah := buatAturan(kunci, nilai); sah && len(daftar) > 0 {
			daftar[len(daftar)-1].aturan = append(daftar[len(daftar)-1].aturan, aturan)
		}
	}
	return daftar
}

func pecahBaris(baris string) (string, string, bool) {
	if posisi := strings.Index(baris, "#"); posisi >= 0 {
		baris = baris[:posisi]
	}
	kunci, nilai, ada := strings.Cut(baris, ":")
	return strings.ToLower(strings.TrimSpace(kunci)), strings.TrimSpace(nilai), ada
}

func buatAturan(kunci, nilai string) (aturanPath, bool) {
	switch {
	case kunci == "allow" && nilai != "":
		return aturanPath{prefiks: nilai, isIzinkan: true}, true
	case kunci == "disallow" && nilai != "":
		return aturanPath{prefiks: nilai}, true
	default:
		return aturanPath{}, false
	}
}

func tokenAgen(userAgent string) string {
	token, _, _ := strings.Cut(userAgent, "/")
	return strings.ToLower(strings.TrimSpace(token))
}
