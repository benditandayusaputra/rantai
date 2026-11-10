package domain

import (
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

type JenisSumberUmpan string

const (
	SumberUmpanResmiDPRD JenisSumberUmpan = "resmi_dprd"
	SumberUmpanMedia     JenisSumberUmpan = "media"
)

type DasarCocok string

const (
	DasarSitusResmi    DasarCocok = "situs_resmi"
	DasarNamaDanDPRD   DasarCocok = "nama_dan_dprd"
	DasarNamaDanFraksi DasarCocok = "nama_dan_fraksi"
)

const (
	BatasCuplikanUmpan    = 300
	jumlahKataNamaMinimal = 2
)

var (
	polaBukanHurufAngka = regexp.MustCompile(`[^\p{L}\p{N}]+`)
	konteksNamaTunggal  = []string{"dprd dki jakarta", "dprd dki"}
)

type ItemUmpan struct {
	JenisSumber JenisSumberUmpan
	MediaNama   string
	MediaDomain string
	Judul       string
	Cuplikan    string
	Isi         string
	URL         string
	TerbitPada  time.Time
}

type PautanUmpan struct {
	AnggotaID  string
	DasarCocok DasarCocok
}

type UmpanTerpaut struct {
	ItemUmpan
	TagIsu []string
	Pautan []PautanUmpan
}

type PencariNama struct {
	polaPerAnggota   map[string][]string
	fraksiPerAnggota map[string]string
}

func BuatPencariNama(anggota []Anggota) PencariNama {
	pemilik := map[string][]string{}
	fraksi := map[string]string{}
	for _, a := range anggota {
		fraksi[a.ID] = a.Fraksi
		for _, nama := range namaPencarian(a) {
			pemilik[nama] = append(pemilik[nama], a.ID)
		}
	}
	polaPerAnggota := map[string][]string{}
	for nama, id := range pemilik {
		if len(id) == 1 {
			polaPerAnggota[id[0]] = append(polaPerAnggota[id[0]], polaSebutan(nama)...)
		}
	}
	return PencariNama{polaPerAnggota: polaPerAnggota, fraksiPerAnggota: fraksi}
}

func namaPencarian(a Anggota) []string {
	nama := []string{}
	for _, n := range append([]string{a.Nama}, a.NamaLain...) {
		kata := strings.Fields(NormalkanTeks(NamaTanpaGelar(n)))
		if len(kata) > 1 && utf8.RuneCountInString(kata[len(kata)-1]) == 1 {
			kata = kata[:len(kata)-1]
		}
		if len(kata) > 0 {
			nama = append(nama, strings.Join(kata, " "))
		}
	}
	slices.Sort(nama)
	return slices.Compact(nama)
}

func polaSebutan(nama string) []string {
	if len(strings.Fields(nama)) >= jumlahKataNamaMinimal {
		return []string{nama}
	}
	pola := []string{}
	for _, konteks := range konteksNamaTunggal {
		pola = append(pola, konteks+" "+nama)
	}
	return pola
}

func (p PencariNama) AnggotaDisebut(teks string) []string {
	normal := " " + NormalkanTeks(teks) + " "
	disebut := []string{}
	for id, daftarPola := range p.polaPerAnggota {
		if slices.ContainsFunc(daftarPola, func(pola string) bool { return strings.Contains(normal, " "+pola+" ") }) {
			disebut = append(disebut, id)
		}
	}
	slices.Sort(disebut)
	return disebut
}

func (p PencariNama) PautkanRilisResmi(item ItemUmpan) []PautanUmpan {
	pautan := []PautanUmpan{}
	for _, id := range p.AnggotaDisebut(item.Judul + " " + item.Cuplikan + " " + item.Isi) {
		pautan = append(pautan, PautanUmpan{AnggotaID: id, DasarCocok: DasarSitusResmi})
	}
	return pautan
}

func TandaiTagIsu(teks string, tagIsu []string) []string {
	normal := " " + NormalkanTeks(teks) + " "
	tag := []string{}
	for _, t := range tagIsu {
		if strings.Contains(normal, " "+NormalkanTeks(t)+" ") {
			tag = append(tag, t)
		}
	}
	slices.Sort(tag)
	return slices.Compact(tag)
}

func NormalkanTeks(teks string) string {
	return strings.TrimSpace(polaBukanHurufAngka.ReplaceAllString(strings.ToLower(teks), " "))
}

func PotongCuplikan(teks string) string {
	teks = strings.Join(strings.Fields(teks), " ")
	if utf8.RuneCountInString(teks) <= BatasCuplikanUmpan {
		return teks
	}
	return string([]rune(teks)[:BatasCuplikanUmpan])
}

type UmpanTampil struct {
	ID           int64
	JenisSumber  JenisSumberUmpan
	MediaNama    string
	MediaDomain  string
	Judul        string
	Ringkasan    string
	Cuplikan     string
	URL          string
	TerbitPada   time.Time
	TagIsu       []string
	DasarCocok   DasarCocok
	ProvenanceID string
	HakJawab     []HakJawab
}

const tandaCuplikanTerpotong = "…"

func (u UmpanTampil) CuplikanRilisResmi() string {
	if u.JenisSumber != SumberUmpanResmiDPRD {
		return ""
	}
	if utf8.RuneCountInString(u.Cuplikan) < BatasCuplikanUmpan {
		return u.Cuplikan
	}
	utuh := u.Cuplikan
	if akhirKata := strings.LastIndexByte(utuh, ' '); akhirKata > 0 {
		utuh = utuh[:akhirKata]
	}
	return utuh + tandaCuplikanTerpotong
}

type FilterUmpan struct {
	AnggotaID  string
	Jenis      JenisSumberUmpan
	Tag        string
	Halaman    int
	PerHalaman int
}

type HalamanUmpan struct {
	Item       []UmpanTampil
	Total      int
	PerHalaman int
}
