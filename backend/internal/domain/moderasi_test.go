package domain

import (
	"errors"
	"testing"
)

func TestPermintaanModerasiLaporanValidasi(t *testing.T) {
	kasus := []struct {
		nama  string
		p     PermintaanModerasiLaporan
		harap error
	}{
		{nama: "verifikasi", p: PermintaanModerasiLaporan{Aksi: AksiVerifikasi}},
		{nama: "tolak dengan alasan dari daftar", p: PermintaanModerasiLaporan{Aksi: AksiTolak, Alasan: "duplikat"}},
		{nama: "tolak tanpa alasan sah", p: PermintaanModerasiLaporan{Aksi: AksiTolak, Alasan: "tidak suka"}, harap: ErrAlasanTolakTidakSah},
		{nama: "sembunyikan foto", p: PermintaanModerasiLaporan{Aksi: AksiSembunyikanFoto, FotoID: 3, Alasan: "wajah orang terlihat"}},
		{nama: "sembunyikan tanpa foto", p: PermintaanModerasiLaporan{Aksi: AksiSembunyikanFoto, Alasan: "x"}, harap: ErrFotoModerasiTidakDipilih},
		{nama: "sembunyikan tanpa alasan", p: PermintaanModerasiLaporan{Aksi: AksiSembunyikanFoto, FotoID: 3}, harap: ErrAlasanModerasiKosong},
		{nama: "aksi tidak dikenal", p: PermintaanModerasiLaporan{Aksi: "hapus"}, harap: ErrAksiModerasiTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.p.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestPermintaanModerasiLaporan_StatusDanAlasanRiwayat(t *testing.T) {
	tolak := PermintaanModerasiLaporan{Aksi: AksiTolak, Alasan: "spam", Pesan: "iklan"}
	verifikasi := PermintaanModerasiLaporan{Aksi: AksiVerifikasi}
	if tanpaPesan := (PermintaanModerasiLaporan{Aksi: AksiTolak, Alasan: "spam"}); tanpaPesan.AlasanRiwayat() != "spam" {
		t.Errorf("tanpa pesan: %q", tanpaPesan.AlasanRiwayat())
	}

	if tolak.StatusTujuan() != StatusDitolak || tolak.AlasanRiwayat() != "spam: iklan" {
		t.Errorf("tolak: %s, %q", tolak.StatusTujuan(), tolak.AlasanRiwayat())
	}
	if verifikasi.StatusTujuan() != StatusTerverifikasi || verifikasi.AlasanRiwayat() != AlasanVerifikasiModerator {
		t.Errorf("verifikasi: %s, %q", verifikasi.StatusTujuan(), verifikasi.AlasanRiwayat())
	}
}
