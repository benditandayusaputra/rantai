package lpse

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/net/html"
)

const (
	kolomNamaPemenang    = "Nama Pemenang"
	kolomHargaKontrak    = "Harga Kontrak"
	kolomHasilNegosiasi  = "Hasil Negosiasi"
	kolomHargaTerkoreksi = "Harga Terkoreksi"
	kolomHargaPenawaran  = "Harga Penawaran"
	tahapTandaTangan     = "Penandatanganan Kontrak"
	kolomTahapJadwal     = 1
	kolomSelesaiJadwal   = 3
	jumlahKolomJadwalMin = 4
)

var urutanKolomNilai = []string{kolomHargaKontrak, kolomHasilNegosiasi, kolomHargaTerkoreksi, kolomHargaPenawaran}

type Pemenang struct {
	Nama  string
	Nilai float64
}

func UraiPemenang(isi []byte) ([]Pemenang, error) {
	dokumen, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return nil, fmt.Errorf("urai html pemenang: %w", err)
	}
	daftar := []Pemenang{}
	for _, tabel := range cariSemua(dokumen, func(n *html.Node) bool { return isElemen(n, "table") }) {
		daftar = append(daftar, pemenangDariTabel(tabel)...)
	}
	return daftar, nil
}

func pemenangDariTabel(tabel *html.Node) []Pemenang {
	baris := cariSemua(tabel, func(n *html.Node) bool { return isElemen(n, "tr") })
	if len(baris) == 0 {
		return nil
	}
	kepala := indeksKepala(baris[0])
	indeksNama, ada := kepala[kolomNamaPemenang]
	if !ada {
		return nil
	}
	daftar := []Pemenang{}
	for _, b := range baris[1:] {
		sel := selSemua(b)
		if len(sel) > indeksNama && sel[indeksNama] != "" {
			daftar = append(daftar, Pemenang{Nama: sel[indeksNama], Nilai: nilaiPemenang(sel, kepala)})
		}
	}
	return daftar
}

func indeksKepala(baris *html.Node) map[string]int {
	indeks := map[string]int{}
	for i, judul := range selSemua(baris) {
		indeks[judul] = i
	}
	return indeks
}

func nilaiPemenang(sel []string, kepala map[string]int) float64 {
	for _, kolom := range urutanKolomNilai {
		if i, ada := kepala[kolom]; ada && i < len(sel) {
			if nilai, ok := uraiRupiah(sel[i]); ok {
				return nilai
			}
		}
	}
	return 0
}

func selSemua(baris *html.Node) []string {
	sel := []string{}
	for anak := baris.FirstChild; anak != nil; anak = anak.NextSibling {
		if isElemen(anak, "td") || isElemen(anak, "th") {
			sel = append(sel, teks(anak))
		}
	}
	return sel
}

func UraiTanggalKontrak(isi []byte) (time.Time, bool, error) {
	dokumen, err := html.Parse(bytes.NewReader(isi))
	if err != nil {
		return time.Time{}, false, fmt.Errorf("urai html jadwal: %w", err)
	}
	for _, baris := range cariSemua(dokumen, func(n *html.Node) bool { return isElemen(n, "tr") }) {
		sel := selBaris(baris)
		if len(sel) >= jumlahKolomJadwalMin && sel[kolomTahapJadwal] == tahapTandaTangan {
			tanggal, err := uraiTanggal(sel[kolomSelesaiJadwal])
			return tanggal, err == nil, err
		}
	}
	return time.Time{}, false, nil
}
