package domain

import (
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	SumberLPSE       SumberData = "lpse"
	SumberSIRUP      SumberData = "sirup"
	SumberOpentender SumberData = "opentender"
)

const (
	awalanKodeTender    = "tender-"
	awalanKodeNonTender = "nontender-"
)

type JenisInstansi string

const (
	JenisInstansiDinas     JenisInstansi = "dinas"
	JenisInstansiSudin     JenisInstansi = "sudin"
	JenisInstansiBadan     JenisInstansi = "badan"
	JenisInstansiKelurahan JenisInstansi = "kelurahan"
	JenisInstansiLainnya   JenisInstansi = "lainnya"
)

type Instansi struct {
	Nama     string
	Jenis    JenisInstansi
	KodeKota string
}

type PaketDaftar struct {
	Kode           string
	KodeSPSE       string
	Nama           string
	NamaInstansi   string
	Tahap          string
	JenisPengadaan string
	Metode         string
	IsTender       bool
	TahunAnggaran  int
	HPS            *float64
	NilaiKontrak   *float64
}

type RUPTercatat struct {
	Kode       string
	Nama       string
	SumberDana string
}

type DetailPaket struct {
	Kode           string
	KodeSPSE       string
	Nama           string
	SatuanKerja    string
	IsBatal        bool
	RUP            *RUPTercatat
	TanggalDibuat  time.Time
	Tahap          string
	JenisPengadaan string
	Metode         string
	JenisKontrak   string
	TahunAnggaran  int
	Pagu           *float64
	HPS            *float64
	JumlahPeserta  *int
	Pemenang       *PemenangPaket
	TanggalKontrak *time.Time
}

type PemenangPaket struct {
	NamaPenyedia string
	Nilai        float64
	IsBerkontrak bool
}

type Penyedia struct {
	Nama       string
	NamaNormal string
}

type AturanJenisInstansi struct {
	Jenis  JenisInstansi `yaml:"jenis"`
	Awalan []string      `yaml:"awalan"`
}

type AturanInstansi struct {
	Jenis []AturanJenisInstansi `yaml:"jenis"`
	Kota  map[string][]string   `yaml:"kota"`
}

var (
	polaBentukUsaha = regexp.MustCompile(`\b(pt|cv|tbk|persero)\b`)
	polaTandaBaca   = regexp.MustCompile(`[^\p{L}\p{N}\s]`)
	polaNonHuruf    = regexp.MustCompile(`[^\p{L}\p{N}]+`)
)

func KodePaket(isTender bool, kodeSPSE string) string {
	if isTender {
		return awalanKodeTender + kodeSPSE
	}
	return awalanKodeNonTender + kodeSPSE
}

func NormalisasiNamaPenyedia(nama string) string {
	hasil := polaTandaBaca.ReplaceAllString(strings.ToLower(nama), " ")
	hasil = polaBentukUsaha.ReplaceAllString(hasil, " ")
	return strings.Join(strings.Fields(hasil), " ")
}

func (a AturanInstansi) Klasifikasi(nama string) Instansi {
	return Instansi{Nama: strings.TrimSpace(nama), Jenis: a.jenis(nama), KodeKota: a.kodeKota(nama)}
}

func (a AturanInstansi) jenis(nama string) JenisInstansi {
	kecil := strings.ToLower(strings.TrimSpace(nama))
	for _, aturan := range a.Jenis {
		if slices.ContainsFunc(aturan.Awalan, func(awalan string) bool { return strings.HasPrefix(kecil, awalan) }) {
			return aturan.Jenis
		}
	}
	return JenisInstansiLainnya
}

func (a AturanInstansi) kodeKota(nama string) string {
	kata := " " + strings.Join(strings.Fields(polaNonHuruf.ReplaceAllString(strings.ToLower(nama), " ")), " ") + " "
	kodeTerpilih := ""
	for _, kode := range slices.Sorted(maps.Keys(a.Kota)) {
		if slices.ContainsFunc(a.Kota[kode], func(sebutan string) bool { return strings.Contains(kata, " "+sebutan+" ") }) {
			if kodeTerpilih != "" {
				return ""
			}
			kodeTerpilih = kode
		}
	}
	return kodeTerpilih
}
