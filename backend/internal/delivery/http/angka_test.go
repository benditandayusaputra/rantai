package http

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func TestAngkaMarshalJSON_LabelDanRujukan_HanyaLengkapYangLolos(t *testing.T) {
	kasus := []struct {
		nama  string
		angka Angka
		lolos bool
	}{
		{"lengkap", Angka{Nilai: 7360, Satuan: "kursi", Label: domain.LabelProksi, ProvenanceID: "prv_1"}, true},
		{"tanpa label", Angka{Nilai: 7360, Satuan: "kursi", ProvenanceID: "prv_1"}, false},
		{"label tak dikenal", Angka{Nilai: 7360, Satuan: "kursi", Label: "tebakan", ProvenanceID: "prv_1"}, false},
		{"tanpa provenance", Angka{Nilai: 7360, Satuan: "kursi", Label: domain.LabelResmi}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			_, err := json.Marshal(struct{ Kursi Angka }{k.angka})

			if lolos := err == nil; lolos != k.lolos {
				t.Errorf("lolos = %v, harap %v (err %v)", lolos, k.lolos, err)
			}
			if !k.lolos && !errors.Is(err, ErrAngkaTanpaRujukan) {
				t.Errorf("err = %v, harap ErrAngkaTanpaRujukan", err)
			}
		})
	}
}

func TestAngkaMarshalJSON_Lengkap_BentukSesuaiSpesifikasi(t *testing.T) {
	isi, err := json.Marshal(Angka{Nilai: 7360, Satuan: "kursi", Label: domain.LabelProksi, ProvenanceID: "prv_1"})

	harap := `{"nilai":7360,"satuan":"kursi","label":"proksi","provenance_id":"prv_1"}`
	if err != nil || string(isi) != harap {
		t.Errorf("dapat %s (err %v), harap %s", isi, err, harap)
	}
}
