package domain

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var waktuUji = time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)

func TestNormalisasiEmail(t *testing.T) {
	kasus := []struct {
		nama, masukan, harap string
		harapErr             error
	}{
		{"huruf besar dan spasi", "  Warga@Contoh.ID ", "warga@contoh.id", nil},
		{"tanpa at", "warga.contoh.id", "", ErrEmailTidakValid},
		{"dengan nama tampilan", "Warga <warga@contoh.id>", "", ErrEmailTidakValid},
		{"sisipan header", "warga@contoh.id\r\nBcc: lain@contoh.id", "", ErrEmailTidakValid},
		{"terlalu panjang", strings.Repeat("a", PanjangMaksEmail) + "@contoh.id", "", ErrEmailTidakValid},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil, err := NormalisasiEmail(k.masukan)
			if hasil != k.harap || !errors.Is(err, k.harapErr) {
				t.Errorf("hasil %q err %v, harap %q %v", hasil, err, k.harap, k.harapErr)
			}
		})
	}
}

func TestPeriksaFormatKodeOTP(t *testing.T) {
	kasus := map[string]error{"123456": nil, "12345": ErrKodeOTPTidakValid, "12345a": ErrKodeOTPTidakValid, "1234567": ErrKodeOTPTidakValid}
	for kode, harap := range kasus {
		if err := PeriksaFormatKodeOTP(kode); !errors.Is(err, harap) {
			t.Errorf("%s: err %v, harap %v", kode, err, harap)
		}
	}
}

func TestAturanOTP_BolehKirim_TigaPerJendela(t *testing.T) {
	aturan := AturanOTP{MaksPermintaan: 3}
	kasus := map[int]error{0: nil, 2: nil, 3: ErrBatasPermintaanOTP, 4: ErrBatasPermintaanOTP}
	for jumlah, harap := range kasus {
		if err := aturan.BolehKirim(jumlah); !errors.Is(err, harap) {
			t.Errorf("%d permintaan: err %v, harap %v", jumlah, err, harap)
		}
	}
}

func TestAturanOTP_Periksa(t *testing.T) {
	aturan := AturanOTP{MaksPercobaan: 5}
	hash := []byte("hash-benar")
	sah := OTPTersimpan{KodeHash: hash, Kedaluwarsa: waktuUji.Add(time.Minute), Percobaan: 1}
	kasus := []struct {
		nama  string
		ubah  func(o *OTPTersimpan)
		kode  []byte
		harap error
	}{
		{"kode benar", func(*OTPTersimpan) {}, hash, nil},
		{"percobaan kelima masih boleh", func(o *OTPTersimpan) { o.Percobaan = 5 }, hash, nil},
		{"kode salah", func(*OTPTersimpan) {}, []byte("hash-salah"), ErrKodeOTPSalah},
		{"percobaan keenam ditolak walau benar", func(o *OTPTersimpan) { o.Percobaan = 6 }, hash, ErrPercobaanOTPHabis},
		{"tepat kedaluwarsa", func(o *OTPTersimpan) { o.Kedaluwarsa = waktuUji }, hash, ErrOTPKedaluwarsa},
		{"sudah dipakai", func(o *OTPTersimpan) { o.IsDipakai = true }, hash, ErrOTPTidakBerlaku},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			otp := sah
			k.ubah(&otp)
			if err := aturan.Periksa(otp, k.kode, waktuUji); !errors.Is(err, k.harap) {
				t.Errorf("err %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestSesiRefresh_Periksa(t *testing.T) {
	kasus := []struct {
		nama  string
		sesi  SesiRefresh
		harap error
	}{
		{"berlaku", SesiRefresh{Kedaluwarsa: waktuUji.Add(time.Hour)}, nil},
		{"kedaluwarsa", SesiRefresh{Kedaluwarsa: waktuUji}, ErrSesiTidakBerlaku},
		{"sudah dirotasi", SesiRefresh{Kedaluwarsa: waktuUji.Add(time.Hour), IsDicabut: true}, ErrSesiDipakaiUlang},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.sesi.Periksa(waktuUji); !errors.Is(err, k.harap) {
				t.Errorf("err %v, harap %v", err, k.harap)
			}
		})
	}
}

func TestPeran_DibuatLewatPerintah_HanyaModeratorDanAdmin(t *testing.T) {
	for _, p := range []Peran{PeranWarga, PeranAnggota, PeranInstansi, PeranModerator, PeranAdmin} {
		harap := p == PeranModerator || p == PeranAdmin
		if p.IsDibuatLewatPerintah() != harap || !p.IsDikenal() {
			t.Errorf("%s: lewat perintah %v, dikenal %v", p, p.IsDibuatLewatPerintah(), p.IsDikenal())
		}
	}
}
