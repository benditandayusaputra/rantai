package domain

import (
	"fmt"
	"strconv"
)

const JumlahBulanSetahun = 12

type Porsi struct {
	Kunci  string
	Nama   string
	Jumlah int
	Nilai  float64
}

type Seri struct {
	Label        Label
	ProvenanceID string
	Isi          []Porsi
}

type RekapKursi struct {
	Kunci     string
	Nama      string
	Kebutuhan float64
	Kapasitas float64
	PerKelas  map[KelasRasio]int
}

type DasborLayanan struct {
	Tema         string
	Label        Label
	ProvenanceID string
	Periode      []RekapKursi
	Kota         []RekapKursi
}

type DasborUang struct {
	Tahun     int
	Tahunan   Seri
	Bulanan   Seri
	Jenis     Seri
	Metode    Seri
	Kota      Seri
	Kecamatan Seri
}

type DasborTanggungJawab struct {
	Tahun           int
	KegiatanTahunan Seri
	KegiatanBulanan Seri
	Topik           Seri
	Reses           Seri
	Instansi        Seri
}

func KunciBulan(tahun, bulan int) string {
	return fmt.Sprintf("%04d-%02d", tahun, bulan)
}

func LengkapiBulan(tahun int, isi []Porsi) []Porsi {
	menurutKunci := make(map[string]Porsi, len(isi))
	for _, p := range isi {
		menurutKunci[p.Kunci] = p
	}
	lengkap := make([]Porsi, JumlahBulanSetahun)
	for i := range lengkap {
		kunci := KunciBulan(tahun, i+1)
		lengkap[i] = Porsi{Kunci: kunci}
		if p, ada := menurutKunci[kunci]; ada {
			lengkap[i] = p
		}
	}
	return lengkap
}

func TahunTerlengkap(tahunan []Porsi) int {
	terpilih, jumlahTerbanyak := 0, -1
	for _, p := range tahunan {
		tahun, err := strconv.Atoi(p.Kunci)
		if err != nil {
			continue
		}
		if p.Jumlah > jumlahTerbanyak || (p.Jumlah == jumlahTerbanyak && tahun > terpilih) {
			terpilih, jumlahTerbanyak = tahun, p.Jumlah
		}
	}
	return terpilih
}
