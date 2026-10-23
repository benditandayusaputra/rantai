package domain

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type BentukSekolah string

const (
	BentukSD  BentukSekolah = "SD"
	BentukMI  BentukSekolah = "MI"
	BentukSMP BentukSekolah = "SMP"
	BentukMTs BentukSekolah = "MTs"
)

type Jenjang string

const (
	JenjangDasar           Jenjang = "dasar"
	JenjangMenengahPertama Jenjang = "menengah_pertama"
)

type StatusSekolah string

const (
	StatusNegeri StatusSekolah = "negeri"
	StatusSwasta StatusSekolah = "swasta"
)

const (
	TingkatSDPertama   = 1
	TingkatSDTerakhir  = 6
	TingkatSMPPertama  = 7
	TingkatSMPTerakhir = 9
	panjangSemesterID  = 5
	digitTahunSemester = 4
)

var (
	jenjangBentuk = map[BentukSekolah]Jenjang{
		BentukSD: JenjangDasar, BentukMI: JenjangDasar,
		BentukSMP: JenjangMenengahPertama, BentukMTs: JenjangMenengahPertama,
	}
	prefiksKecamatanSumber = regexp.MustCompile(`(?i)^\s*kec(amatan)?\.?\s+`)
	bukanHurufKecil        = regexp.MustCompile(`[^a-z]`)
)

type Sekolah struct {
	NPSN                string
	Nama                string
	Bentuk              BentukSekolah
	Status              StatusSekolah
	NamaKecamatanSumber string
	Lokasi              *Titik
}

type RombelTingkat struct {
	Tingkat      int
	JumlahRombel int
	JumlahSiswa  int
}

type DataSekolah struct {
	Sekolah
	TahunAjaran string
	Rombel      []RombelTingkat
}

func (b BentukSekolah) IsValid() bool {
	_, ada := jenjangBentuk[b]
	return ada
}

func (b BentukSekolah) Jenjang() Jenjang {
	return jenjangBentuk[b]
}

func (b BentukSekolah) IsMadrasah() bool {
	return b == BentukMI || b == BentukMTs
}

func (b BentukSekolah) RentangTingkat() (pertama, terakhir int) {
	if b.Jenjang() == JenjangDasar {
		return TingkatSDPertama, TingkatSDTerakhir
	}
	return TingkatSMPPertama, TingkatSMPTerakhir
}

func StatusDariSumber(status string) (StatusSekolah, error) {
	hasil := StatusSekolah(strings.ToLower(strings.TrimSpace(status)))
	if hasil != StatusNegeri && hasil != StatusSwasta {
		return "", InputTidakValid("status sekolah tidak dikenal: " + status)
	}
	return hasil, nil
}

func TahunAjaranDariSemester(semesterID string) (string, error) {
	if len(semesterID) != panjangSemesterID {
		return "", InputTidakValid("semester_id tidak valid: " + semesterID)
	}
	tahun, err := strconv.Atoi(semesterID[:digitTahunSemester])
	if err != nil {
		return "", InputTidakValid("semester_id tidak valid: " + semesterID)
	}
	return fmt.Sprintf("%d/%d", tahun, tahun+1), nil
}

func NormalkanNamaKecamatan(nama string) string {
	tanpaPrefiks := prefiksKecamatanSumber.ReplaceAllString(nama, "")
	return bukanHurufKecil.ReplaceAllString(strings.ToLower(tanpaPrefiks), "")
}

type MetodeLokasi string

const (
	MetodeLokasiTitik MetodeLokasi = "titik"
	MetodeLokasiNama  MetodeLokasi = "nama"
)

type KandidatWilayah struct {
	KecamatanDariTitik string
	KelurahanDariTitik string
	KecamatanDariNama  string
}

type PenempatanSekolah struct {
	KodeKecamatan string
	KodeKelurahan string
	Metode        MetodeLokasi
}

func TempatkanSekolah(k KandidatWilayah) (PenempatanSekolah, bool) {
	if k.KecamatanDariTitik != "" {
		return PenempatanSekolah{KodeKecamatan: k.KecamatanDariTitik, KodeKelurahan: k.KelurahanDariTitik, Metode: MetodeLokasiTitik}, true
	}
	if k.KecamatanDariNama != "" {
		return PenempatanSekolah{KodeKecamatan: k.KecamatanDariNama, Metode: MetodeLokasiNama}, true
	}
	return PenempatanSekolah{}, false
}
