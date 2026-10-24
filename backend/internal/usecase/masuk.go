package usecase

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const panjangTokenRefreshByte = 32

var batasKodeOTP = big.NewInt(1_000_000)

type OTPBaru struct {
	Email       string
	KodeHash    []byte
	Kedaluwarsa time.Time
	DibuatPada  time.Time
}

type SesiBaru struct {
	PenggunaID  int64
	TokenHash   []byte
	Kedaluwarsa time.Time
}

type PenyimpanOTP interface {
	HitungOTPSejak(ctx context.Context, email string, sejak time.Time) (int, error)
	SimpanOTP(ctx context.Context, otp OTPBaru) error
	CatatPercobaanOTP(ctx context.Context, email string) (domain.OTPTersimpan, error)
	TandaiOTPDipakai(ctx context.Context, id int64) error
}

type PenyimpanPengguna interface {
	PastikanPengguna(ctx context.Context, email string) (domain.Pengguna, error)
	AmbilPengguna(ctx context.Context, id int64) (domain.Pengguna, error)
}

type PenyimpanSesi interface {
	SimpanSesi(ctx context.Context, sesi SesiBaru) error
	CabutSesi(ctx context.Context, tokenHash []byte) (domain.SesiRefresh, error)
	CabutSemuaSesi(ctx context.Context, penggunaID int64) error
}

type PengirimOTP interface {
	KirimOTP(ctx context.Context, email, kode string) error
}

type PenerbitAkses interface {
	Terbitkan(klaim domain.KlaimAkses) (string, error)
}

type AturanMasuk struct {
	OTP         domain.AturanOTP
	UmurRefresh time.Duration
}

type HasilMasuk struct {
	Pengguna           domain.Pengguna
	TokenAkses         string
	TokenRefresh       string
	RefreshKedaluwarsa time.Time
}

type Masuk struct {
	OTP      PenyimpanOTP
	Pengguna PenyimpanPengguna
	Sesi     PenyimpanSesi
	Pengirim PengirimOTP
	Akses    PenerbitAkses
	Aturan   AturanMasuk
	Rahasia  []byte
	Sekarang func() time.Time
}

func (u Masuk) KirimOTP(ctx context.Context, email string) error {
	bersih, err := domain.NormalisasiEmail(email)
	if err != nil {
		return err
	}
	sekarang := u.Sekarang()
	jumlah, err := u.OTP.HitungOTPSejak(ctx, bersih, sekarang.Add(-u.Aturan.OTP.JendelaPermintaan))
	if err != nil {
		return err
	}
	if err := u.Aturan.OTP.BolehKirim(jumlah); err != nil {
		return err
	}
	kode, err := buatKodeOTP()
	if err != nil {
		return err
	}
	otp := OTPBaru{Email: bersih, KodeHash: u.hashKode(bersih, kode), Kedaluwarsa: sekarang.Add(u.Aturan.OTP.Umur), DibuatPada: sekarang}
	if err := u.OTP.SimpanOTP(ctx, otp); err != nil {
		return err
	}
	return u.Pengirim.KirimOTP(ctx, bersih, kode)
}

func (u Masuk) Verifikasi(ctx context.Context, email, kode string) (HasilMasuk, error) {
	bersih, err := domain.NormalisasiEmail(email)
	if err != nil {
		return HasilMasuk{}, err
	}
	if err := domain.PeriksaFormatKodeOTP(kode); err != nil {
		return HasilMasuk{}, err
	}
	otp, err := u.OTP.CatatPercobaanOTP(ctx, bersih)
	if err != nil {
		return HasilMasuk{}, err
	}
	if err := u.Aturan.OTP.Periksa(otp, u.hashKode(bersih, kode), u.Sekarang()); err != nil {
		return HasilMasuk{}, err
	}
	if err := u.OTP.TandaiOTPDipakai(ctx, otp.ID); err != nil {
		return HasilMasuk{}, err
	}
	pengguna, err := u.Pengguna.PastikanPengguna(ctx, bersih)
	if err != nil {
		return HasilMasuk{}, err
	}
	return u.bukaSesi(ctx, pengguna)
}

func (u Masuk) Segarkan(ctx context.Context, tokenRefresh string) (HasilMasuk, error) {
	if tokenRefresh == "" {
		return HasilMasuk{}, domain.ErrSesiTidakBerlaku
	}
	sesi, err := u.Sesi.CabutSesi(ctx, hashToken(tokenRefresh))
	if err != nil {
		return HasilMasuk{}, err
	}
	if err := sesi.Periksa(u.Sekarang()); err != nil {
		return HasilMasuk{}, u.tanganiSesiGagal(ctx, sesi, err)
	}
	pengguna, err := u.Pengguna.AmbilPengguna(ctx, sesi.PenggunaID)
	if err != nil {
		return HasilMasuk{}, err
	}
	return u.bukaSesi(ctx, pengguna)
}

func (u Masuk) Keluar(ctx context.Context, tokenRefresh string) error {
	if tokenRefresh == "" {
		return nil
	}
	_, err := u.Sesi.CabutSesi(ctx, hashToken(tokenRefresh))
	if errors.Is(err, domain.ErrSesiTidakBerlaku) {
		return nil
	}
	return err
}

func (u Masuk) Saya(ctx context.Context, penggunaID int64) (domain.Pengguna, error) {
	return u.Pengguna.AmbilPengguna(ctx, penggunaID)
}

func (u Masuk) tanganiSesiGagal(ctx context.Context, sesi domain.SesiRefresh, errSesi error) error {
	if !errors.Is(errSesi, domain.ErrSesiDipakaiUlang) {
		return errSesi
	}
	if err := u.Sesi.CabutSemuaSesi(ctx, sesi.PenggunaID); err != nil {
		return errors.Join(errSesi, err)
	}
	return errSesi
}

func (u Masuk) bukaSesi(ctx context.Context, pengguna domain.Pengguna) (HasilMasuk, error) {
	token, err := buatTokenRefresh()
	if err != nil {
		return HasilMasuk{}, err
	}
	kedaluwarsa := u.Sekarang().Add(u.Aturan.UmurRefresh)
	if err := u.Sesi.SimpanSesi(ctx, SesiBaru{PenggunaID: pengguna.ID, TokenHash: hashToken(token), Kedaluwarsa: kedaluwarsa}); err != nil {
		return HasilMasuk{}, err
	}
	akses, err := u.Akses.Terbitkan(domain.KlaimAkses{PenggunaID: pengguna.ID, Peran: pengguna.Peran})
	if err != nil {
		return HasilMasuk{}, err
	}
	return HasilMasuk{Pengguna: pengguna, TokenAkses: akses, TokenRefresh: token, RefreshKedaluwarsa: kedaluwarsa}, nil
}

func (u Masuk) hashKode(email, kode string) []byte {
	mac := hmac.New(sha256.New, u.Rahasia)
	mac.Write([]byte(email + "\x00" + kode))
	return mac.Sum(nil)
}

func buatKodeOTP() (string, error) {
	n, err := rand.Int(rand.Reader, batasKodeOTP)
	if err != nil {
		return "", fmt.Errorf("buat kode masuk: %w", err)
	}
	return fmt.Sprintf("%0*d", domain.PanjangKodeOTP, n.Int64()), nil
}

func buatTokenRefresh() (string, error) {
	acak := make([]byte, panjangTokenRefreshByte)
	if _, err := rand.Read(acak); err != nil {
		return "", fmt.Errorf("buat token refresh: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(acak), nil
}

func hashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
