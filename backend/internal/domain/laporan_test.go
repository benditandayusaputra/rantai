package domain

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func laporanMasukUji() LaporanMasuk {
	return LaporanMasuk{
		PenggunaID: 1, Titik: Titik{Lat: -6.1214, Lng: 106.9012}, AkurasiMeter: 12,
		Kondisi: KondisiSelesaiRusak, Catatan: "Plafon retak.", Foto: [][]byte{{1, 2, 3}},
	}
}

func TestLaporanMasukValidasi(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*LaporanMasuk)
		harap error
	}{
		{nama: "lengkap, sah", ubah: func(*LaporanMasuk) {}},
		{nama: "akurasi tepat 100 m, sah", ubah: func(l *LaporanMasuk) { l.AkurasiMeter = 100 }},
		{nama: "akurasi lebih dari 100 m", ubah: func(l *LaporanMasuk) { l.AkurasiMeter = 100.5 }, harap: ErrAkurasiLokasiBuruk},
		{nama: "kondisi tidak dikenal", ubah: func(l *LaporanMasuk) { l.Kondisi = "rusak" }, harap: ErrKondisiLaporanTidakSah},
		{nama: "catatan 501 karakter", ubah: func(l *LaporanMasuk) { l.Catatan = strings.Repeat("é", 501) }, harap: ErrCatatanTerlaluPanjang},
		{nama: "tanpa foto", ubah: func(l *LaporanMasuk) { l.Foto = nil }, harap: ErrJumlahFotoTidakSah},
		{nama: "lima foto", ubah: func(l *LaporanMasuk) { l.Foto = make([][]byte, 5) }, harap: ErrJumlahFotoTidakSah},
		{nama: "foto lebih dari 8 MB", ubah: func(l *LaporanMasuk) { l.Foto = [][]byte{make([]byte, BatasUkuranFotoByte+1)} }, harap: ErrFotoTerlaluBesar},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			l := laporanMasukUji()
			k.ubah(&l)
			if err := l.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}

func jarakMeter(a, b Titik) float64 {
	dLat := (a.Lat - b.Lat) * meterPerDerajatLintang
	dLng := (a.Lng - b.Lng) * meterPerDerajatLintang * math.Cos(a.Lat*math.Pi/180)
	return math.Hypot(dLat, dLng)
}

func TestBulatkanLokasiPublik_TitikAsli_BergeserPalingJauhSetengahDiagonalGrid(t *testing.T) {
	asli := Titik{Lat: -6.121437, Lng: 106.901288}

	publik := BulatkanLokasiPublik(asli)

	if publik == asli || jarakMeter(asli, publik) > GridLokasiPublikMeter*math.Sqrt2/2 {
		t.Errorf("publik %+v berjarak %.1f m dari asli", publik, jarakMeter(asli, publik))
	}
}

func TestBulatkanLokasiPublik_DuaTitikSeSel_HasilSama(t *testing.T) {
	sel := BulatkanLokasiPublik(Titik{Lat: -6.1214, Lng: 106.9012})

	if geser := BulatkanLokasiPublik(Titik{Lat: sel.Lat + 0.00005, Lng: sel.Lng - 0.00005}); geser != sel {
		t.Errorf("titik dalam sel yang sama dibulatkan ke %+v, harap %+v", geser, sel)
	}
}
