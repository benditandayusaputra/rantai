package domain

import "strings"

type SumberFoto string

const (
	SumberFotoSekolahKita SumberFoto = "sekolah_kita"
	SumberFotoWikimedia   SumberFoto = "wikimedia"
	pemisahKredit                    = " · "
)

type Foto struct {
	URL     string
	Sumber  SumberFoto
	Halaman string
	Kredit  string
}

type FotoSekolah struct {
	NPSN string
	Foto
}

func SusunKredit(pembuat, lisensi string) string {
	bagian := []string{}
	for _, b := range []string{strings.TrimSpace(pembuat), strings.TrimSpace(lisensi)} {
		if b != "" {
			bagian = append(bagian, b)
		}
	}
	return strings.Join(bagian, pemisahKredit)
}
