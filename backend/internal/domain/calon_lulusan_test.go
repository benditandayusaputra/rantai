package domain

import "testing"

var estimasiUji = AturanEstimasi{Metode: MetodeEstimasiTotalSD, Pembagi: 5}

func TestHitungCalonLulusan_MetodeSesuaiKetersediaanData(t *testing.T) {
	kasus := []struct {
		nama  string
		siswa SiswaKecamatan
		harap CalonLulusan
		ada   bool
	}{
		{"kelas 6 per sekolah tersedia", SiswaKecamatan{"31.72.06", 1250, 7400},
			CalonLulusan{"31.72.06", 1250, LabelResmi, MetodePerSekolah}, true},
		{"hanya total siswa SD", SiswaKecamatan{"31.72.06", 0, 7403},
			CalonLulusan{"31.72.06", 1481, LabelEstimasi, MetodeEstimasiTotalSD}, true},
		{"tanpa data", SiswaKecamatan{"31.72.06", 0, 0}, CalonLulusan{}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil, ada := HitungCalonLulusan(k.siswa, estimasiUji)

			if hasil != k.harap || ada != k.ada {
				t.Errorf("dapat %+v %v, harap %+v %v", hasil, ada, k.harap, k.ada)
			}
		})
	}
}

func TestHitungCalonLulusan_EstimasiSMP_MetodeMengikutiAturan(t *testing.T) {
	siswa := SiswaKecamatan{KodeKecamatan: "31.72.06", TotalSiswa: 3001}

	hasil, ada := HitungCalonLulusan(siswa, AturanEstimasi{Metode: MetodeEstimasiTotalSMP, Pembagi: 3})

	if !ada || hasil.Jumlah != 1000 || hasil.Metode != MetodeEstimasiTotalSMP || hasil.Label != LabelEstimasi {
		t.Errorf("hasil = %+v, %v", hasil, ada)
	}
}
