package domain

import (
	"math"
	"slices"
	"strconv"
	"strings"
)

const (
	kuartilBawah     = 0.25
	kuartilTengah    = 0.5
	kuartilAtas      = 0.75
	digitPerKelompok = 3
)

type Statistik struct {
	Jumlah int     `json:"jumlah"`
	Min    float64 `json:"min"`
	Q1     float64 `json:"q1"`
	Median float64 `json:"median"`
	Q3     float64 `json:"q3"`
	Maks   float64 `json:"maks"`
}

func HitungStatistik(nilai []float64) Statistik {
	if len(nilai) == 0 {
		return Statistik{}
	}
	urut := slices.Sorted(slices.Values(nilai))
	return Statistik{
		Jumlah: len(urut),
		Min:    urut[0],
		Q1:     kuantil(urut, kuartilBawah),
		Median: kuantil(urut, kuartilTengah),
		Q3:     kuantil(urut, kuartilAtas),
		Maks:   urut[len(urut)-1],
	}
}

func (s Statistik) BatasAtas(pengaliIQR float64) float64 {
	return s.Q3 + pengaliIQR*(s.Q3-s.Q1)
}

func kuantil(urut []float64, p float64) float64 {
	posisi := p * float64(len(urut)-1)
	bawah := int(math.Floor(posisi))
	atas := int(math.Ceil(posisi))
	return urut[bawah] + (posisi-float64(bawah))*(urut[atas]-urut[bawah])
}

func FormatRupiah(nilai float64) string {
	return "Rp" + FormatRibuan(nilai)
}

func FormatRibuan(nilai float64) string {
	digit := strconv.FormatInt(int64(math.Round(nilai)), 10)
	kelompok := []string{}
	for len(digit) > digitPerKelompok {
		kelompok = append([]string{digit[len(digit)-digitPerKelompok:]}, kelompok...)
		digit = digit[:len(digit)-digitPerKelompok]
	}
	kelompok = append([]string{digit}, kelompok...)
	return strings.Join(kelompok, ".")
}

type NilaiPembanding struct {
	Kode         string
	HPS          float64
	ProvenanceID string
}

type PaketDibandingkan struct {
	Kode           string
	Kategori       string
	JenisPengadaan string
	TahunAnggaran  int
	HPS            *float64
	ProvenanceID   string
}

type PembandingPaket struct {
	Paket             PaketDibandingkan
	IsTersedia        bool
	Statistik         Statistik
	Nilai             []NilaiPembanding
	ProvenanceTerbaru string
}

func BandingkanPaket(p PaketDibandingkan, pembanding []NilaiPembanding, provenanceTerbaru string, minPembanding int) PembandingPaket {
	hasil := PembandingPaket{Paket: p, Nilai: []NilaiPembanding{}, ProvenanceTerbaru: provenanceTerbaru}
	if p.Kategori == "" || p.HPS == nil || len(pembanding) < minPembanding {
		return hasil
	}
	nilai := make([]float64, len(pembanding))
	for i, n := range pembanding {
		nilai[i] = n.HPS
	}
	hasil.IsTersedia, hasil.Statistik, hasil.Nilai = true, HitungStatistik(nilai), pembanding
	return hasil
}
