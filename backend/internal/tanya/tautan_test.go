package tanya

import (
	"encoding/json"
	"reflect"
	"testing"
)

const hasilWakilUji = `{"kecamatan":{"kode":"31.75.07","nama":"Duren Sawit"},"dapil":{"nama":"Dapil DKI Jakarta 5",
"kecamatan":[{"kode":"31.75.07","nama":"Duren Sawit"}],
"anggota":[{"id":"justin-adrian-sh-mh","nama":"Justin Adrian, S.H., M.H.","fraksi":"Fraksi PSI"},
{"id":"ali","nama":"Ali","fraksi":"Fraksi X"}]}}`

const hasilUmpanUji = `[{"id":7,"judul":"Komisi E Pantau Sarana GOR","url":"https://dprd-dkijakartaprov.go.id/komisi-e"},
{"id":8,"judul":"Rilis Tautan Tidak Aman","url":"javascript:alert(1)"}]`

func TestSusunTautan_NamaDiJawaban_DitautkanDariHasilAlat(t *testing.T) {
	hasilAlat := []json.RawMessage{json.RawMessage(hasilWakilUji), json.RawMessage(hasilUmpanUji)}
	kasus := []struct {
		nama  string
		teks  string
		harap []Tautan
	}{
		{nama: "tanpa nama yang dikenal", teks: "Data belum tersedia.", harap: []Tautan{}},
		{
			nama: "anggota, dapil, dan rilis yang disebut", teks: "Justin Adrian, S.H., M.H. dari Dapil DKI Jakarta 5 hadir di Komisi E Pantau Sarana GOR.",
			harap: []Tautan{
				{Teks: "Komisi E Pantau Sarana GOR", Jenis: TautanRilis, Tujuan: "https://dprd-dkijakartaprov.go.id/komisi-e"},
				{Teks: "Justin Adrian, S.H., M.H.", Jenis: TautanAnggota, Tujuan: "justin-adrian-sh-mh"},
				{Teks: "Dapil DKI Jakarta 5", Jenis: TautanDapil, Tujuan: "31.75.07"},
			},
		},
		{nama: "nama terlalu pendek dan URL tidak aman dilewati", teks: "Ali membaca Rilis Tautan Tidak Aman.", harap: []Tautan{}},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if hasil := SusunTautan(k.teks, hasilAlat); !reflect.DeepEqual(hasil, k.harap) {
				t.Errorf("tautan = %+v, harap %+v", hasil, k.harap)
			}
		})
	}
}
