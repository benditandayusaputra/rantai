package domain

import (
	"math"
)

type MetodeCalonLulusan string

const (
	MetodePerSekolah       MetodeCalonLulusan = "per_sekolah"
	MetodeAgregatResmi     MetodeCalonLulusan = "agregat_resmi"
	MetodeEstimasiTotalSD  MetodeCalonLulusan = "estimasi_total_sd"
	MetodeEstimasiTotalSMP MetodeCalonLulusan = "estimasi_total_smp"
)

type SiswaKecamatan struct {
	KodeKecamatan     string
	SiswaTingkatAkhir int
	TotalSiswa        int
}

type AturanEstimasi struct {
	Metode  MetodeCalonLulusan
	Pembagi int
}

type CalonLulusan struct {
	KodeKecamatan string
	Jumlah        int
	Label         Label
	Metode        MetodeCalonLulusan
}

func HitungCalonLulusan(siswa SiswaKecamatan, estimasi AturanEstimasi) (CalonLulusan, bool) {
	hasil := CalonLulusan{KodeKecamatan: siswa.KodeKecamatan}
	switch {
	case siswa.SiswaTingkatAkhir > 0:
		hasil.Jumlah, hasil.Label, hasil.Metode = siswa.SiswaTingkatAkhir, LabelResmi, MetodePerSekolah
	case siswa.TotalSiswa > 0 && estimasi.Pembagi > 0:
		perkiraan := int(math.Round(float64(siswa.TotalSiswa) / float64(estimasi.Pembagi)))
		hasil.Jumlah, hasil.Label, hasil.Metode = perkiraan, LabelEstimasi, estimasi.Metode
	default:
		return CalonLulusan{}, false
	}
	return hasil, true
}
