package main

import (
	"testing"
	"time"
)

func TestWaktuJalanBerikutnya_SelaluPukulDuaWIBBerikutnya(t *testing.T) {
	kasus := []struct {
		nama     string
		sekarang time.Time
		harapan  time.Time
	}{
		{"sebelum pukul dua di hari yang sama", time.Date(2026, 10, 27, 1, 30, 0, 0, zonaWIB), time.Date(2026, 10, 27, 2, 0, 0, 0, zonaWIB)},
		{"tepat pukul dua pindah ke besok", time.Date(2026, 10, 27, 2, 0, 0, 0, zonaWIB), time.Date(2026, 10, 28, 2, 0, 0, 0, zonaWIB)},
		{"siang hari pindah ke besok", time.Date(2026, 10, 27, 14, 0, 0, 0, zonaWIB), time.Date(2026, 10, 28, 2, 0, 0, 0, zonaWIB)},
		{"jam UTC yang sudah lewat tengah malam WIB", time.Date(2026, 10, 26, 18, 30, 0, 0, time.UTC), time.Date(2026, 10, 27, 2, 0, 0, 0, zonaWIB)},
		{"akhir bulan pindah ke bulan berikutnya", time.Date(2026, 10, 31, 23, 0, 0, 0, zonaWIB), time.Date(2026, 11, 1, 2, 0, 0, 0, zonaWIB)},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if hasil := waktuJalanBerikutnya(k.sekarang); !hasil.Equal(k.harapan) {
				t.Fatalf("harapan %s, dapat %s", k.harapan, hasil)
			}
		})
	}
}
