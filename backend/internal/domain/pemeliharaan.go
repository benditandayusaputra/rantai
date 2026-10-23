package domain

import (
	"math"
	"time"
)

const jamPerHari = 24 * time.Hour

type StatusPemeliharaan string

const (
	StatusPelaksanaan          StatusPemeliharaan = "pelaksanaan"
	StatusPemeliharaanBerjalan StatusPemeliharaan = "pemeliharaan"
	StatusPemeliharaanBerakhir StatusPemeliharaan = "pemeliharaan_berakhir"
	StatusTidakDapatDiestimasi StatusPemeliharaan = "tidak_dapat_diestimasi"
)

type MasaKategori struct {
	Jenis            string `yaml:"jenis"`
	DurasiBawaanHari int    `yaml:"durasi_bawaan_hari"`
}

type AturanPemeliharaan struct {
	JenisMasaHari map[string]int          `yaml:"jenis_masa_hari"`
	Kategori      map[string]MasaKategori `yaml:"kategori"`
	Bawaan        MasaKategori            `yaml:"bawaan"`
}

type JadwalPaket struct {
	Kategori           string
	TanggalSerahTerima *time.Time
	JadwalSelesai      *time.Time
	TanggalKontrak     *time.Time
}

type EstimasiPemeliharaan struct {
	Status            StatusPemeliharaan
	AkhirPelaksanaan  *time.Time
	AkhirPemeliharaan *time.Time
	SisaHari          int
	Label             Label
}

func (a AturanPemeliharaan) Estimasi(j JadwalPaket, hariIni time.Time) EstimasiPemeliharaan {
	masa := a.masaUntuk(j.Kategori)
	label := LabelEstimasi
	if j.TanggalSerahTerima != nil {
		label = LabelResmi
	}
	akhirPelaksanaan, ada := akhirPelaksanaan(j, masa)
	if !ada {
		return EstimasiPemeliharaan{Status: StatusTidakDapatDiestimasi, Label: LabelEstimasi}
	}
	akhirPemeliharaan := tambahHari(akhirPelaksanaan, a.JenisMasaHari[masa.Jenis])
	e := EstimasiPemeliharaan{AkhirPelaksanaan: &akhirPelaksanaan, AkhirPemeliharaan: &akhirPemeliharaan, Label: label}
	switch {
	case hariIni.Before(akhirPelaksanaan):
		e.Status = StatusPelaksanaan
	case hariIni.Before(akhirPemeliharaan):
		e.Status = StatusPemeliharaanBerjalan
		e.SisaHari = int(math.Ceil(akhirPemeliharaan.Sub(hariIni).Hours() / jamPerHari.Hours()))
	default:
		e.Status = StatusPemeliharaanBerakhir
	}
	return e
}

func (a AturanPemeliharaan) masaUntuk(kategori string) MasaKategori {
	if masa, ada := a.Kategori[kategori]; ada {
		return masa
	}
	return a.Bawaan
}

func akhirPelaksanaan(j JadwalPaket, masa MasaKategori) (time.Time, bool) {
	switch {
	case j.TanggalSerahTerima != nil:
		return *j.TanggalSerahTerima, true
	case j.JadwalSelesai != nil:
		return *j.JadwalSelesai, true
	case j.TanggalKontrak != nil:
		return tambahHari(*j.TanggalKontrak, masa.DurasiBawaanHari), true
	default:
		return time.Time{}, false
	}
}

func tambahHari(t time.Time, hari int) time.Time {
	return t.AddDate(0, 0, hari)
}
