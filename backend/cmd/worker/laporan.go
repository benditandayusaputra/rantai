package main

import (
	"log"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

type pelaporLog struct{}

func (pelaporLog) LaporkanSekolah(l usecase.LaporanImporSekolah) {
	log.Printf("sekolah dari sumber: %d, tersimpan: %d", l.JumlahSumber, l.JumlahTersimpan)
	cetakDaftar("sekolah tanpa koordinat (kecamatan dari nama, lokasi kosong)", l.TanpaKoordinat)
	cetakDaftar("koordinat di luar batas kecamatan DKI (kecamatan dari nama)", l.KoordinatDiLuarBatas)
	cetakDaftar("sekolah tidak terpetakan (tidak disimpan)", l.TidakTerpetakan)
	cetakDaftar("kecamatan tanpa data sekolah", l.KecamatanTanpaSekolah)
}

func cetakDaftar(judul string, daftar []string) {
	log.Printf("%s: %d", judul, len(daftar))
	if len(daftar) > 0 {
		log.Print("  " + strings.Join(daftar, "\n  "))
	}
}
