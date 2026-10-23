package domain

import (
	"fmt"
	"math"
)

type MetodeCalonLulusan string

const (
	MetodePerSekolah      MetodeCalonLulusan = "per_sekolah"
	MetodeAgregatResmi    MetodeCalonLulusan = "agregat_resmi"
	MetodeEstimasiTotalSD MetodeCalonLulusan = "estimasi_total_sd"
)

type SiswaSDKecamatan struct {
	KodeKecamatan string
	SiswaKelas6   int
	TotalSiswa    int
}

type CalonLulusan struct {
	KodeKecamatan string
	Jumlah        int
	Label         Label
	Metode        MetodeCalonLulusan
}

func HitungCalonLulusan(siswa SiswaSDKecamatan, pembagiEstimasi int) (CalonLulusan, bool) {
	hasil := CalonLulusan{KodeKecamatan: siswa.KodeKecamatan}
	switch {
	case siswa.SiswaKelas6 > 0:
		hasil.Jumlah, hasil.Label, hasil.Metode = siswa.SiswaKelas6, LabelResmi, MetodePerSekolah
	case siswa.TotalSiswa > 0 && pembagiEstimasi > 0:
		perkiraan := int(math.Round(float64(siswa.TotalSiswa) / float64(pembagiEstimasi)))
		hasil.Jumlah, hasil.Label, hasil.Metode = perkiraan, LabelEstimasi, MetodeEstimasiTotalSD
	default:
		return CalonLulusan{}, false
	}
	return hasil, true
}

func RumusEstimasiTotalSD(pembagiEstimasi int) string {
	return fmt.Sprintf("calon lulusan = total siswa SD di kecamatan ÷ %d, dibulatkan; dipakai bila jumlah siswa kelas 6 per sekolah tidak tersedia", pembagiEstimasi)
}
