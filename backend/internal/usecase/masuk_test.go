package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const emailUji = "warga@contoh.id"

type otpPalsu struct {
	OTPBaru
	id        int64
	percobaan int
	dipakai   bool
}

type sesiPalsu struct {
	SesiBaru
	dicabut bool
}

type akunPalsu struct {
	otp          []*otpPalsu
	sesi         []*sesiPalsu
	pengguna     map[string]domain.Pengguna
	kodeTerkirim string
}

func (a *akunPalsu) HitungOTPSejak(_ context.Context, email string, sejak time.Time) (int, error) {
	jumlah := 0
	for _, o := range a.otp {
		if o.Email == email && !o.DibuatPada.Before(sejak) {
			jumlah++
		}
	}
	return jumlah, nil
}

func (a *akunPalsu) SimpanOTP(_ context.Context, otp OTPBaru) error {
	a.otp = append(a.otp, &otpPalsu{OTPBaru: otp, id: int64(len(a.otp) + 1)})
	return nil
}

func (a *akunPalsu) CatatPercobaanOTP(_ context.Context, email string) (domain.OTPTersimpan, error) {
	for i := len(a.otp) - 1; i >= 0; i-- {
		if o := a.otp[i]; o.Email == email {
			o.percobaan++
			return domain.OTPTersimpan{ID: o.id, KodeHash: o.KodeHash, Kedaluwarsa: o.Kedaluwarsa, Percobaan: o.percobaan, IsDipakai: o.dipakai}, nil
		}
	}
	return domain.OTPTersimpan{}, domain.ErrOTPTidakBerlaku
}

func (a *akunPalsu) TandaiOTPDipakai(_ context.Context, id int64) error {
	a.otp[id-1].dipakai = true
	return nil
}

func (a *akunPalsu) PastikanPengguna(_ context.Context, email string) (domain.Pengguna, error) {
	if p, ada := a.pengguna[email]; ada {
		return p, nil
	}
	p := domain.Pengguna{ID: int64(len(a.pengguna) + 1), Email: email, Peran: domain.PeranWarga}
	a.pengguna[email] = p
	return p, nil
}

func (a *akunPalsu) AmbilPengguna(_ context.Context, id int64) (domain.Pengguna, error) {
	for _, p := range a.pengguna {
		if p.ID == id {
			return p, nil
		}
	}
	return domain.Pengguna{}, domain.ErrPenggunaTidakDitemukan
}

func (a *akunPalsu) SimpanSesi(_ context.Context, sesi SesiBaru) error {
	a.sesi = append(a.sesi, &sesiPalsu{SesiBaru: sesi})
	return nil
}

func (a *akunPalsu) CabutSesi(_ context.Context, tokenHash []byte) (domain.SesiRefresh, error) {
	for _, s := range a.sesi {
		if bytes.Equal(s.TokenHash, tokenHash) {
			lama := domain.SesiRefresh{PenggunaID: s.PenggunaID, Kedaluwarsa: s.Kedaluwarsa, IsDicabut: s.dicabut}
			s.dicabut = true
			return lama, nil
		}
	}
	return domain.SesiRefresh{}, domain.ErrSesiTidakBerlaku
}

func (a *akunPalsu) CabutSemuaSesi(_ context.Context, penggunaID int64) error {
	for _, s := range a.sesi {
		if s.PenggunaID == penggunaID {
			s.dicabut = true
		}
	}
	return nil
}

func (a *akunPalsu) KirimOTP(_ context.Context, _, kode string) error {
	a.kodeTerkirim = kode
	return nil
}

type penerbitPalsu struct{}

func (penerbitPalsu) Terbitkan(k domain.KlaimAkses) (string, error) {
	return fmt.Sprintf("akses-%d-%s", k.PenggunaID, k.Peran), nil
}

type jamUji struct {
	saat time.Time
}

func (j *jamUji) sekarang() time.Time { return j.saat }

func masukUji() (Masuk, *akunPalsu, *jamUji) {
	akun := &akunPalsu{pengguna: map[string]domain.Pengguna{}}
	jam := &jamUji{saat: time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)}
	return Masuk{
		OTP: akun, Pengguna: akun, Sesi: akun, Pengirim: akun, Akses: penerbitPalsu{},
		Aturan: AturanMasuk{
			OTP:         domain.AturanOTP{Umur: 10 * time.Minute, MaksPercobaan: 5, MaksPermintaan: 3, JendelaPermintaan: time.Hour},
			UmurRefresh: 7 * 24 * time.Hour,
		},
		Rahasia:  []byte("rahasia-uji"),
		Sekarang: jam.sekarang,
	}, akun, jam
}

func masukDenganKode(t *testing.T, u Masuk, akun *akunPalsu) HasilMasuk {
	t.Helper()
	if err := u.KirimOTP(t.Context(), emailUji); err != nil {
		t.Fatal(err)
	}
	hasil, err := u.Verifikasi(t.Context(), emailUji, akun.kodeTerkirim)
	if err != nil {
		t.Fatal(err)
	}
	return hasil
}

func kodeSalah(benar string) string {
	if benar == "000000" {
		return "111111"
	}
	return "000000"
}

func TestMasuk_KodeBenar_SesiDanTokenAksesTerbit(t *testing.T) {
	u, akun, _ := masukUji()

	hasil := masukDenganKode(t, u, akun)

	if len(akun.kodeTerkirim) != domain.PanjangKodeOTP || hasil.TokenAkses != "akses-1-warga" || hasil.TokenRefresh == "" ||
		hasil.Pengguna.Email != emailUji || len(akun.sesi) != 1 || bytes.Contains(akun.otp[0].KodeHash, []byte(akun.kodeTerkirim)) {
		t.Errorf("hasil %+v, kode %q", hasil, akun.kodeTerkirim)
	}
}

func TestMasuk_EmailHurufBesar_DinormalisasiSebelumDisimpan(t *testing.T) {
	u, akun, _ := masukUji()

	if err := u.KirimOTP(t.Context(), " Warga@Contoh.ID "); err != nil {
		t.Fatal(err)
	}
	_, err := u.Verifikasi(t.Context(), emailUji, akun.kodeTerkirim)

	if err != nil || akun.otp[0].Email != emailUji {
		t.Errorf("email tersimpan %q, err %v", akun.otp[0].Email, err)
	}
}

func TestKirimOTP_PermintaanKeempatDalamSejam_Ditolak(t *testing.T) {
	u, _, jam := masukUji()
	for range 3 {
		if err := u.KirimOTP(t.Context(), emailUji); err != nil {
			t.Fatal(err)
		}
	}

	errKeempat := u.KirimOTP(t.Context(), emailUji)
	jam.saat = jam.saat.Add(time.Hour + time.Second)
	errSetelahSejam := u.KirimOTP(t.Context(), emailUji)

	if !errors.Is(errKeempat, domain.ErrBatasPermintaanOTP) || errSetelahSejam != nil {
		t.Errorf("keempat %v, setelah sejam %v", errKeempat, errSetelahSejam)
	}
}

func TestVerifikasi_LimaKaliSalah_KodeBenarPunDitolak(t *testing.T) {
	u, akun, _ := masukUji()
	if err := u.KirimOTP(t.Context(), emailUji); err != nil {
		t.Fatal(err)
	}
	salah := kodeSalah(akun.kodeTerkirim)
	for i := range 5 {
		if _, err := u.Verifikasi(t.Context(), emailUji, salah); !errors.Is(err, domain.ErrKodeOTPSalah) {
			t.Fatalf("percobaan %d: err %v", i+1, err)
		}
	}

	_, err := u.Verifikasi(t.Context(), emailUji, akun.kodeTerkirim)

	if !errors.Is(err, domain.ErrPercobaanOTPHabis) {
		t.Errorf("err %v, harap %v", err, domain.ErrPercobaanOTPHabis)
	}
}

func TestVerifikasi_LewatSepuluhMenit_Kedaluwarsa(t *testing.T) {
	u, akun, jam := masukUji()
	if err := u.KirimOTP(t.Context(), emailUji); err != nil {
		t.Fatal(err)
	}
	jam.saat = jam.saat.Add(10 * time.Minute)

	_, err := u.Verifikasi(t.Context(), emailUji, akun.kodeTerkirim)

	if !errors.Is(err, domain.ErrOTPKedaluwarsa) {
		t.Errorf("err %v, harap %v", err, domain.ErrOTPKedaluwarsa)
	}
}

func TestVerifikasi_KodeDipakaiDuaKali_KeduaDitolak(t *testing.T) {
	u, akun, _ := masukUji()
	masukDenganKode(t, u, akun)

	_, err := u.Verifikasi(t.Context(), emailUji, akun.kodeTerkirim)

	if !errors.Is(err, domain.ErrOTPTidakBerlaku) {
		t.Errorf("err %v, harap %v", err, domain.ErrOTPTidakBerlaku)
	}
}

func TestVerifikasi_KodeLamaSetelahMintaBaru_Ditolak(t *testing.T) {
	u, akun, _ := masukUji()
	if err := u.KirimOTP(t.Context(), emailUji); err != nil {
		t.Fatal(err)
	}
	kodeLama := akun.kodeTerkirim
	if err := u.KirimOTP(t.Context(), emailUji); err != nil {
		t.Fatal(err)
	}
	if akun.kodeTerkirim == kodeLama {
		t.Skip("kode acak kebetulan sama")
	}

	_, err := u.Verifikasi(t.Context(), emailUji, kodeLama)

	if !errors.Is(err, domain.ErrKodeOTPSalah) {
		t.Errorf("err %v, harap %v", err, domain.ErrKodeOTPSalah)
	}
}

func TestSegarkan_TokenBerotasi_TokenLamaTidakBerlaku(t *testing.T) {
	u, akun, _ := masukUji()
	awal := masukDenganKode(t, u, akun)

	baru, err := u.Segarkan(t.Context(), awal.TokenRefresh)
	if err != nil {
		t.Fatal(err)
	}
	_, errUlang := u.Segarkan(t.Context(), awal.TokenRefresh)

	if baru.TokenRefresh == awal.TokenRefresh || baru.TokenAkses == "" || !errors.Is(errUlang, domain.ErrSesiDipakaiUlang) {
		t.Errorf("token baru %q, err ulang %v", baru.TokenRefresh, errUlang)
	}
}

func TestSegarkan_TokenLamaDipakaiUlang_SemuaSesiDicabut(t *testing.T) {
	u, akun, _ := masukUji()
	awal := masukDenganKode(t, u, akun)
	baru, err := u.Segarkan(t.Context(), awal.TokenRefresh)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := u.Segarkan(t.Context(), awal.TokenRefresh); err == nil {
		t.Fatal("token lama diterima")
	}

	_, err = u.Segarkan(t.Context(), baru.TokenRefresh)

	if !errors.Is(err, domain.ErrSesiDipakaiUlang) {
		t.Errorf("token hasil rotasi masih berlaku setelah pemakaian ulang: %v", err)
	}
}

func TestSegarkan_LewatTujuhHari_Kedaluwarsa(t *testing.T) {
	u, akun, jam := masukUji()
	awal := masukDenganKode(t, u, akun)
	jam.saat = jam.saat.Add(7 * 24 * time.Hour)

	_, err := u.Segarkan(t.Context(), awal.TokenRefresh)

	if !errors.Is(err, domain.ErrSesiTidakBerlaku) {
		t.Errorf("err %v, harap %v", err, domain.ErrSesiTidakBerlaku)
	}
}

func TestKeluar_SesiDicabutDanTokenTakDikenalDiabaikan(t *testing.T) {
	u, akun, _ := masukUji()
	awal := masukDenganKode(t, u, akun)

	errKeluar := u.Keluar(t.Context(), awal.TokenRefresh)
	errAsing := u.Keluar(t.Context(), "token-asing")
	_, errSegarkan := u.Segarkan(t.Context(), awal.TokenRefresh)

	if errKeluar != nil || errAsing != nil || errSegarkan == nil {
		t.Errorf("keluar %v, asing %v, segarkan setelah keluar %v", errKeluar, errAsing, errSegarkan)
	}
}
