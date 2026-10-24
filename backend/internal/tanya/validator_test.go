package tanya

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func TestValidatorJawaban_Periksa(t *testing.T) {
	hasilAlat := []json.RawMessage{
		json.RawMessage(`{"kursi":{"nilai":7360,"satuan":"kursi","label":"resmi","provenance_id":"prv_01M5KPYF"}}`),
		json.RawMessage(`{"rasio":{"nilai":0.85,"provenance_id":"prv_rantai_x"},"hps":1500000000}`),
	}
	kasus := []struct {
		nama  string
		teks  string
		harap error
	}{
		{"angka persis bersumber", "Ada 7360 kursi [prv_01M5KPYF].", nil},
		{"pemisah ribuan dan desimal Indonesia", "Ada 7.360 kursi [prv_01M5KPYF], rasio 0,85 [prv_rantai_x], HPS Rp1.500.000.000.", nil},
		{"tanpa angka", "Data kursi madrasah belum tersedia.", nil},
		{"angka karangan", "Ada 8000 kursi [prv_01M5KPYF].", ErrJawabanAngkaKarangan},
		{"angka hasil penjumlahan", "Totalnya 7361 kursi.", ErrJawabanAngkaKarangan},
		{"rujukan tak dikenal", "Ada 7360 kursi [prv_karangan].", ErrJawabanRujukanTakDikenal},
		{"kosakata terlarang", "Paket ini terindikasi korupsi.", ErrJawabanTerlarang},
		{"hanya rujukan", " [prv_01M5KPYF] ", ErrJawabanKosong},
	}
	validator := ValidatorJawaban{Terlarang: domain.KosakataTerlarang{"korupsi", "mark up"}}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := validator.Periksa(k.teks, hasilAlat); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
