package domain

import (
	"strings"
	"testing"
)

const pembagiUji = 5

func TestHitungCalonLulusan_MetodeSesuaiKetersediaanData(t *testing.T) {
	kasus := []struct {
		nama  string
		siswa SiswaSDKecamatan
		harap CalonLulusan
		ada   bool
	}{
		{"kelas 6 per sekolah tersedia", SiswaSDKecamatan{"31.72.06", 1250, 7400},
			CalonLulusan{"31.72.06", 1250, LabelResmi, MetodePerSekolah}, true},
		{"hanya total siswa SD", SiswaSDKecamatan{"31.72.06", 0, 7403},
			CalonLulusan{"31.72.06", 1481, LabelEstimasi, MetodeEstimasiTotalSD}, true},
		{"tanpa data", SiswaSDKecamatan{"31.72.06", 0, 0}, CalonLulusan{}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil, ada := HitungCalonLulusan(k.siswa, pembagiUji)

			if hasil != k.harap || ada != k.ada {
				t.Errorf("dapat %+v %v, harap %+v %v", hasil, ada, k.harap, k.ada)
			}
		})
	}
}

func TestRumusEstimasiTotalSD_MemuatPembagi(t *testing.T) {
	if rumus := RumusEstimasiTotalSD(pembagiUji); !strings.Contains(rumus, "÷ 5") {
		t.Errorf("rumus = %q", rumus)
	}
}
