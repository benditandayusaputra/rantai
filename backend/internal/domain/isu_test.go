package domain

import "testing"

func umpanIsuUji(id int64, jenis JenisSumberUmpan, tag ...string) UmpanTampil {
	return UmpanTampil{ID: id, JenisSumber: jenis, TagIsu: tag}
}

func TestKelompokkanIsu(t *testing.T) {
	kasus := []struct {
		nama  string
		umpan []UmpanTampil
		harap []IsuDisuarakan
	}{
		{nama: "tanpa umpan, daftar kosong", umpan: []UmpanTampil{}, harap: []IsuDisuarakan{}},
		{
			nama: "rilis resmi dan media dihitung terpisah",
			umpan: []UmpanTampil{
				umpanIsuUji(1, SumberUmpanMedia, "sekolah"),
				umpanIsuUji(2, SumberUmpanResmiDPRD, "sekolah", "spmb"),
				umpanIsuUji(3, SumberUmpanMedia, "sekolah"),
			},
			harap: []IsuDisuarakan{
				{Tag: "sekolah", JumlahResmi: 1, JumlahMedia: 2},
				{Tag: "spmb", JumlahResmi: 1},
			},
		},
		{
			nama: "jumlah sama, urut nama tag",
			umpan: []UmpanTampil{
				umpanIsuUji(1, SumberUmpanMedia, "zonasi"),
				umpanIsuUji(2, SumberUmpanMedia, "guru"),
			},
			harap: []IsuDisuarakan{{Tag: "guru", JumlahMedia: 1}, {Tag: "zonasi", JumlahMedia: 1}},
		},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil := KelompokkanIsu(k.umpan)
			if hasil == nil || len(hasil) != len(k.harap) {
				t.Fatalf("hasil = %+v", hasil)
			}
			for i, h := range k.harap {
				if hasil[i].Tag != h.Tag || hasil[i].JumlahResmi != h.JumlahResmi || hasil[i].JumlahMedia != h.JumlahMedia {
					t.Errorf("isu %d = %+v, harap %+v", i, hasil[i], h)
				}
			}
		})
	}
}

func TestKelompokkanIsu_ItemMelebihiBatas_DihitungSemuaTampilTerbaru(t *testing.T) {
	umpan := []UmpanTampil{}
	for id := range int64(BatasItemPerIsu + 2) {
		umpan = append(umpan, umpanIsuUji(id, SumberUmpanResmiDPRD, "guru"))
	}

	isu := KelompokkanIsu(umpan)[0]

	if isu.JumlahSumber() != BatasItemPerIsu+2 || len(isu.Item) != BatasItemPerIsu || isu.Item[0].ID != 0 {
		t.Errorf("jumlah %d, item %d, pertama %d", isu.JumlahSumber(), len(isu.Item), isu.Item[0].ID)
	}
}
