package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestHakJawabBaru_Validasi(t *testing.T) {
	sah := HakJawabBaru{TargetJenis: TargetHakJawabSinyal, TargetID: "12", Isi: "Paket sudah sesuai prosedur.", TautanDokumen: "https://dinas.jakarta.go.id/ba.pdf"}
	kasus := []struct {
		nama  string
		ubah  func(b *HakJawabBaru)
		harap error
	}{
		{"sah", func(*HakJawabBaru) {}, nil},
		{"paket berkode teks", func(b *HakJawabBaru) { b.TargetJenis, b.TargetID = TargetHakJawabPaket, "tender-123" }, nil},
		{"tanpa tautan", func(b *HakJawabBaru) { b.TautanDokumen = "" }, nil},
		{"jenis asing", func(b *HakJawabBaru) { b.TargetJenis = "penyedia" }, ErrTargetHakJawabTidakSah},
		{"id sinyal bukan angka", func(b *HakJawabBaru) { b.TargetID = "S-01" }, ErrTargetHakJawabTidakSah},
		{"isi kosong", func(b *HakJawabBaru) { b.Isi = "" }, ErrIsiHakJawabTidakSah},
		{"isi 2.001 karakter", func(b *HakJawabBaru) { b.Isi = strings.Repeat("é", PanjangMaksHakJawab+1) }, ErrIsiHakJawabTidakSah},
		{"isi tepat 2.000 karakter", func(b *HakJawabBaru) { b.Isi = strings.Repeat("é", PanjangMaksHakJawab) }, nil},
		{"tautan javascript", func(b *HakJawabBaru) { b.TautanDokumen = "javascript:alert(1)" }, ErrTautanDokumenTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			b := sah
			k.ubah(&b)
			if err := b.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestPengguna_PeriksaBolehMenjawab(t *testing.T) {
	dinas, dinasLain := int64(3), int64(4)
	anggota := "anggota-1"
	petugas := Pengguna{Peran: PeranInstansi, InstansiID: &dinas}
	kasus := []struct {
		nama     string
		pengguna Pengguna
		target   PemilikTarget
		harap    error
	}{
		{"petugas instansi pemilik", petugas, PemilikTarget{InstansiID: &dinas}, nil},
		{"petugas instansi lain", petugas, PemilikTarget{InstansiID: &dinasLain}, ErrBukanPenjawabTarget},
		{"anggota yang dipautkan", Pengguna{Peran: PeranAnggota, AnggotaID: &anggota}, PemilikTarget{AnggotaID: []string{"anggota-1"}}, nil},
		{"anggota tanpa tautan", Pengguna{Peran: PeranAnggota}, PemilikTarget{AnggotaID: []string{"anggota-1"}}, ErrBukanPenjawabTarget},
		{"warga", Pengguna{Peran: PeranWarga}, PemilikTarget{InstansiID: &dinas}, ErrBukanPenjawabTarget},
		{"akun simulasi pada data nyata", Pengguna{Peran: PeranInstansi, InstansiID: &dinas, IsSimulation: true}, PemilikTarget{InstansiID: &dinas}, ErrBukanPenjawabTarget},
		{"akun nyata pada data simulasi", petugas, PemilikTarget{InstansiID: &dinas, IsSimulation: true}, ErrBukanPenjawabTarget},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.pengguna.PeriksaBolehMenjawab(k.target); !errors.Is(err, k.harap) {
				t.Errorf("err %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestPengguna_PeriksaPetugasPemilik_AnggotaDitolak(t *testing.T) {
	anggota := "anggota-1"
	err := Pengguna{Peran: PeranAnggota, AnggotaID: &anggota}.PeriksaPetugasPemilik(PemilikTarget{AnggotaID: []string{anggota}})

	if !errors.Is(err, ErrBukanAkunInstansi) {
		t.Errorf("err %v", err)
	}
}

func TestSerahTerima_Validasi(t *testing.T) {
	hariIni := time.Date(2026, 10, 24, 0, 0, 0, 0, time.UTC)
	kasus := []struct {
		nama  string
		s     SerahTerima
		harap error
	}{
		{"hari ini", SerahTerima{Tanggal: hariIni}, nil},
		{"besok", SerahTerima{Tanggal: hariIni.AddDate(0, 0, 1)}, ErrTanggalSerahTerimaTidakSah},
		{"kosong", SerahTerima{}, ErrTanggalSerahTerimaTidakSah},
		{"nomor terlalu panjang", SerahTerima{Tanggal: hariIni, NomorBeritaAcara: strings.Repeat("1", PanjangMaksNomorBeritaAcara+1)}, ErrNomorBeritaAcaraTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.s.Validasi(hariIni); !errors.Is(err, k.harap) {
				t.Errorf("err %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestLaporan_IsBisaDiteruskan_HanyaTerverifikasiBerpaket(t *testing.T) {
	kasus := []struct {
		laporan LaporanUntukVerifikasi
		harap   bool
	}{
		{LaporanUntukVerifikasi{Status: StatusTerverifikasi, PaketKode: "simulasi-1"}, true},
		{LaporanUntukVerifikasi{Status: StatusTerverifikasi}, false},
		{LaporanUntukVerifikasi{Status: StatusMenungguSaksi, PaketKode: "simulasi-1"}, false},
	}
	for _, k := range kasus {
		if k.laporan.IsBisaDiteruskan() != k.harap {
			t.Errorf("%+v: harap %v", k.laporan, k.harap)
		}
	}
}
