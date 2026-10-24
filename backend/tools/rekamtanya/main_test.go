package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

func TestBacaAliran_AlatLaluJawaban_TerekamLengkap(t *testing.T) {
	aliran := "event: alat\ndata: {\"alat\":\"indikator_tema\"}\n\n" +
		"event: jawaban\ndata: {\"teks\":\"Kursi 288 [1].\",\"rujukan\":[{\"nomor\":1,\"provenance_id\":\"prv_1\"}]}\n\n"

	entri, err := bacaAliran(strings.NewReader(aliran))

	harap := tanya.EntriCache{Alat: []string{"indikator_tema"}, Teks: "Kursi 288 [1].", Rujukan: []tanya.Rujukan{{Nomor: 1, ProvenanceID: "prv_1"}}}
	if err != nil || !reflect.DeepEqual(entri, harap) {
		t.Errorf("entri = %+v, err = %v", entri, err)
	}
}

func TestBacaAliran_PeristiwaGalat_Gagal(t *testing.T) {
	_, err := bacaAliran(strings.NewReader("event: galat\ndata: {\"pesan\":\"belum bisa\"}\n\n"))

	if err == nil || !strings.Contains(err.Error(), "belum bisa") {
		t.Errorf("err = %v, harap pesan galat", err)
	}
}

func TestBacaAliran_TanpaJawaban_Gagal(t *testing.T) {
	_, err := bacaAliran(strings.NewReader("event: alat\ndata: {\"alat\":\"cari_wilayah\"}\n\n"))

	if !errors.Is(err, ErrTanpaJawaban) {
		t.Errorf("err = %v, harap ErrTanpaJawaban", err)
	}
}
