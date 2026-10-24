package domain

import (
	"errors"
	"testing"
	"time"
)

var tanggalParipurnaUji = time.Date(2026, time.October, 19, 0, 0, 0, 0, time.UTC)

func sikapUji() SikapFraksi {
	return SikapFraksi{
		Fraksi: "Fraksi Partai Golkar", Agenda: AgendaAPBD, JudulAgenda: "Ranperda APBD 2027", Tanggal: tanggalParipurnaUji,
		Ringkasan: "Meminta subsidi transportasi berbasis kinerja.", URL: "https://dprd-dkijakartaprov.go.id/berita/",
	}
}

func aspirasiUji() AspirasiReses {
	return AspirasiReses{
		MasaReses: "Masa Reses ke-3", TanggalParipurna: tanggalParipurnaUji, Bidang: "kesehatan",
		Uraian: "Pembaruan alat kesehatan Posyandu.", URL: "https://dprd-dkijakartaprov.go.id/berita/",
	}
}

func TestSikapFraksiValidasi(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*SikapFraksi)
		harap error
	}{
		{nama: "lengkap, sah", ubah: func(*SikapFraksi) {}},
		{nama: "agenda tidak dikenal", ubah: func(s *SikapFraksi) { s.Agenda = "rapat" }, harap: ErrAgendaSikapTidakSah},
		{nama: "ringkasan kosong", ubah: func(s *SikapFraksi) { s.Ringkasan = " " }, harap: ErrEntriParipurnaKosong},
		{nama: "tanpa tanggal", ubah: func(s *SikapFraksi) { s.Tanggal = time.Time{} }, harap: ErrTanggalParipurnaKosong},
		{nama: "tautan bukan URL", ubah: func(s *SikapFraksi) { s.URL = "berita resmi" }, harap: ErrTautanSumberTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			s := sikapUji()
			k.ubah(&s)
			if err := s.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestAspirasiResesValidasi(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*AspirasiReses)
		harap error
	}{
		{nama: "lengkap tanpa wilayah, sah", ubah: func(*AspirasiReses) {}},
		{nama: "bidang huruf besar", ubah: func(a *AspirasiReses) { a.Bidang = "Kesehatan" }, harap: ErrBidangTidakSah},
		{nama: "uraian kosong", ubah: func(a *AspirasiReses) { a.Uraian = "" }, harap: ErrEntriParipurnaKosong},
		{nama: "tautan tanpa skema", ubah: func(a *AspirasiReses) { a.URL = "dprd-dkijakartaprov.go.id" }, harap: ErrTautanSumberTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			a := aspirasiUji()
			k.ubah(&a)
			if err := a.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
