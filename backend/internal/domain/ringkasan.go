package domain

import (
	"errors"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

const BatasRingkasanUmpan = 160

var (
	ErrRingkasanKosong         = errors.New("ringkasan kosong")
	ErrRingkasanTerlaluPanjang = errors.New("ringkasan melebihi 160 karakter")
	ErrRingkasanLebihSatu      = errors.New("ringkasan lebih dari satu kalimat")
	ErrRingkasanTerlarang      = errors.New("ringkasan memuat kosakata terlarang")
	ErrRingkasanMenilai        = errors.New("ringkasan memuat kata penilaian")
	ErrRingkasanAngkaKarangan  = errors.New("ringkasan memuat angka yang tidak ada di sumber")
	polaAngka                  = regexp.MustCompile(`\d+(?:[.,]\d+)*`)
	polaBatasKalimat           = regexp.MustCompile(`[.!?]\s+\p{Lu}`)
)

type UmpanUntukRingkas struct {
	ID       int64
	Judul    string
	Cuplikan string
}

type AturanRingkasan struct {
	Terlarang     KosakataTerlarang
	KataPenilaian []string
}

func (a AturanRingkasan) Validasi(ringkasan string, sumber UmpanUntukRingkas) error {
	ringkasan = strings.TrimSpace(ringkasan)
	switch {
	case ringkasan == "":
		return ErrRingkasanKosong
	case utf8.RuneCountInString(ringkasan) > BatasRingkasanUmpan:
		return ErrRingkasanTerlaluPanjang
	case polaBatasKalimat.MatchString(ringkasan):
		return ErrRingkasanLebihSatu
	case len(a.Terlarang.Temukan(ringkasan)) > 0:
		return ErrRingkasanTerlarang
	case memuatSalahSatu(ringkasan, a.KataPenilaian):
		return ErrRingkasanMenilai
	case !IsAngkaBersumber(ringkasan, sumber.Judul+" "+sumber.Cuplikan):
		return ErrRingkasanAngkaKarangan
	default:
		return nil
	}
}

func IsAngkaBersumber(ringkasan, sumber string) bool {
	angkaSumber := angkaDalam(sumber)
	for _, angka := range angkaDalam(ringkasan) {
		if !slices.Contains(angkaSumber, angka) {
			return false
		}
	}
	return true
}

func angkaDalam(teks string) []string {
	angka := []string{}
	for _, a := range polaAngka.FindAllString(teks, -1) {
		angka = append(angka, strings.NewReplacer(".", "", ",", "").Replace(a))
	}
	return angka
}
