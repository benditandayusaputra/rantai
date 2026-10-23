package main

import (
	"fmt"
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

func (pelaporLog) LaporkanDayaTampung(daftar []usecase.LaporanDayaTampung) {
	for _, l := range daftar {
		log.Printf("daya tampung %s: kursi di dokumen %d, terimpor %d, selisih %d", l.Jalur, l.KursiDokumen, l.KursiTerimpor, l.KursiDokumen-l.KursiTerimpor)
		tidakTerpetakan := make([]string, len(l.TidakTerpetakan))
		for i, b := range l.TidakTerpetakan {
			tidakTerpetakan[i] = fmt.Sprintf("NPSN %q %s (%d kursi)", b.NPSN, b.Nama, b.Kursi)
		}
		cetakDaftar("sekolah tidak terpetakan", tidakTerpetakan)
	}
}

func (pelaporLog) LaporkanKecamatanTanpaCalonLulusan(kode []string) {
	cetakDaftar("kecamatan tanpa data calon lulusan", kode)
}

func (pelaporLog) LaporkanKecamatanDemoTerpilih(kode string) {
	log.Printf("kecamatan demo terpilih: %s; kunci kode ini di kecamatan_demo skenario_demo.yaml", kode)
}
