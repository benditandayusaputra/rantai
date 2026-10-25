package domain

import (
	"fmt"
	"strings"
)

const (
	BatasKabarRingkasan   = 50
	bahasaTautanRingkasan = "/id"
	subjekRingkasan       = "Ringkasan mingguan Rantai: %d kabar dari kecamatan yang kamu pantau"
	pembukaRingkasan      = "Halo,\n\nIni kabar sepekan terakhir dari kecamatan yang kamu pantau di Rantai.\n"
	penutupRingkasan      = "\nSetiap angka di Rantai punya sumber; buka tautannya untuk melihat rincian dan tanggal ambil.\n" +
		"Berhenti menerima ringkasan ini di halaman Pantau: %s\n"
	penandaSimulasi = " (Data simulasi)"
)

var labelJenisKejadian = map[JenisKejadian]string{
	KejadianPaketBaru:            "Paket baru",
	KejadianMelewatiJadwal:       "Melewati jadwal",
	KejadianPemeliharaanBerakhir: "Masa pemeliharaan segera berakhir",
	KejadianLaporanTerverifikasi: "Laporan terverifikasi",
	KejadianHakJawab:             "Hak jawab baru",
	KejadianIndikatorDiperbarui:  "Indikator diperbarui",
}

type Situs struct {
	Alamat string `yaml:"alamat"`
}

type PenerimaRingkasan struct {
	PenggunaID int64
	Email      string
}

type EmailRingkasan struct {
	Subjek string
	Isi    string
}

func SusunEmailRingkasan(kabar []Notifikasi, alamatSitus string) EmailRingkasan {
	var isi strings.Builder
	isi.WriteString(pembukaRingkasan)
	kecamatan := ""
	for _, k := range kabar {
		if k.KodeKecamatan != kecamatan {
			kecamatan = k.KodeKecamatan
			fmt.Fprintf(&isi, "\n%s\n", k.NamaKecamatan)
		}
		fmt.Fprintf(&isi, "- [%s] %s%s\n  %s\n", labelJenisKejadian[k.Jenis], k.Judul, simulasi(k), alamatSitus+bahasaTautanRingkasan+k.Tautan)
	}
	fmt.Fprintf(&isi, penutupRingkasan, alamatSitus+bahasaTautanRingkasan+"/pantau")
	return EmailRingkasan{Subjek: fmt.Sprintf(subjekRingkasan, len(kabar)), Isi: isi.String()}
}

func simulasi(k Notifikasi) string {
	if k.IsSimulation {
		return penandaSimulasi
	}
	return ""
}
