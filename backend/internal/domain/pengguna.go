package domain

import (
	"crypto/subtle"
	"net/mail"
	"strings"
	"time"
)

type Peran string

const (
	PeranWarga     Peran = "warga"
	PeranAnggota   Peran = "anggota"
	PeranInstansi  Peran = "instansi"
	PeranModerator Peran = "moderator"
	PeranAdmin     Peran = "admin"

	PanjangMaksEmail = 254
	PanjangKodeOTP   = 6
)

var (
	ErrEmailTidakValid         = InputTidakValid("alamat email tidak valid")
	ErrKodeOTPTidakValid       = InputTidakValid("kode masuk terdiri dari 6 angka")
	ErrBatasPermintaanOTP      = TerlaluBanyak("terlalu banyak permintaan kode untuk email ini, coba lagi nanti")
	ErrOTPTidakBerlaku         = TidakBerwenang("kode tidak berlaku, minta kode baru")
	ErrOTPKedaluwarsa          = TidakBerwenang("kode sudah kedaluwarsa, minta kode baru")
	ErrPercobaanOTPHabis       = TidakBerwenang("percobaan habis, minta kode baru")
	ErrKodeOTPSalah            = TidakBerwenang("kode salah")
	ErrSesiTidakBerlaku        = TidakBerwenang("sesi berakhir, silakan masuk lagi")
	ErrSesiDipakaiUlang        = TidakBerwenang("sesi berakhir, silakan masuk lagi")
	ErrTokenAksesTidakBerlaku  = TidakBerwenang("silakan masuk lebih dulu")
	ErrPeranTidakDikenal       = InputTidakValid("peran tidak dikenal")
	ErrPeranBukanUntukPerintah = InputTidakValid("hanya peran moderator dan admin yang dibuat lewat perintah")
	ErrPenggunaTidakDitemukan  = TidakDitemukan("pengguna tidak ditemukan")
	daftarPeran                = map[Peran]bool{PeranWarga: true, PeranAnggota: true, PeranInstansi: true, PeranModerator: true, PeranAdmin: true}
	peranDibuatLewatPerintah   = map[Peran]bool{PeranModerator: true, PeranAdmin: true}
)

type Pengguna struct {
	ID           int64
	Email        string
	NamaTampilan string
	Peran        Peran
	InstansiID   *int64
	AnggotaID    *string
}

type KlaimAkses struct {
	PenggunaID int64
	Peran      Peran
}

type AturanOTP struct {
	Umur              time.Duration `yaml:"umur"`
	MaksPercobaan     int           `yaml:"maks_percobaan"`
	MaksPermintaan    int           `yaml:"maks_permintaan"`
	JendelaPermintaan time.Duration `yaml:"jendela_permintaan"`
}

type OTPTersimpan struct {
	ID          int64
	KodeHash    []byte
	Kedaluwarsa time.Time
	Percobaan   int
	IsDipakai   bool
}

type SesiRefresh struct {
	PenggunaID  int64
	Kedaluwarsa time.Time
	IsDicabut   bool
}

func (p Peran) IsDikenal() bool {
	return daftarPeran[p]
}

func (p Peran) IsDibuatLewatPerintah() bool {
	return peranDibuatLewatPerintah[p]
}

func NormalisasiEmail(email string) (string, error) {
	bersih := strings.ToLower(strings.TrimSpace(email))
	if len(bersih) > PanjangMaksEmail {
		return "", ErrEmailTidakValid
	}
	alamat, err := mail.ParseAddress(bersih)
	if err != nil || alamat.Address != bersih {
		return "", ErrEmailTidakValid
	}
	return bersih, nil
}

func PeriksaFormatKodeOTP(kode string) error {
	if len(kode) != PanjangKodeOTP || strings.Trim(kode, "0123456789") != "" {
		return ErrKodeOTPTidakValid
	}
	return nil
}

func (a AturanOTP) BolehKirim(jumlahDalamJendela int) error {
	if jumlahDalamJendela >= a.MaksPermintaan {
		return ErrBatasPermintaanOTP
	}
	return nil
}

func (a AturanOTP) Periksa(otp OTPTersimpan, kodeHash []byte, sekarang time.Time) error {
	switch {
	case otp.IsDipakai:
		return ErrOTPTidakBerlaku
	case !sekarang.Before(otp.Kedaluwarsa):
		return ErrOTPKedaluwarsa
	case otp.Percobaan > a.MaksPercobaan:
		return ErrPercobaanOTPHabis
	case subtle.ConstantTimeCompare(otp.KodeHash, kodeHash) != 1:
		return ErrKodeOTPSalah
	}
	return nil
}

func (s SesiRefresh) Periksa(sekarang time.Time) error {
	if s.IsDicabut {
		return ErrSesiDipakaiUlang
	}
	if !sekarang.Before(s.Kedaluwarsa) {
		return ErrSesiTidakBerlaku
	}
	return nil
}
