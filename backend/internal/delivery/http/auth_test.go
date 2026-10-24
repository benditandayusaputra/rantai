package http

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/token"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const tokenRefreshUji = "refresh-uji"

var penggunaUji = domain.Pengguna{ID: 7, Email: "warga@contoh.id", Peran: domain.PeranWarga}

type masukPalsu struct {
	akses          token.Akses
	refreshDiminta string
}

func (m *masukPalsu) KirimOTP(_ context.Context, email string) error {
	if email == "sering@contoh.id" {
		return domain.ErrBatasPermintaanOTP
	}
	return nil
}

func (m *masukPalsu) Verifikasi(_ context.Context, _, kode string) (usecase.HasilMasuk, error) {
	if kode != "123456" {
		return usecase.HasilMasuk{}, domain.ErrKodeOTPSalah
	}
	return m.hasil()
}

func (m *masukPalsu) Segarkan(_ context.Context, tokenRefresh string) (usecase.HasilMasuk, error) {
	m.refreshDiminta = tokenRefresh
	if tokenRefresh != tokenRefreshUji {
		return usecase.HasilMasuk{}, domain.ErrSesiTidakBerlaku
	}
	return m.hasil()
}

func (m *masukPalsu) Keluar(context.Context, string) error { return nil }

func (m *masukPalsu) Saya(_ context.Context, id int64) (domain.Pengguna, error) {
	if id != penggunaUji.ID {
		return domain.Pengguna{}, domain.ErrPenggunaTidakDitemukan
	}
	return penggunaUji, nil
}

func (m *masukPalsu) hasil() (usecase.HasilMasuk, error) {
	akses, err := m.akses.Terbitkan(domain.KlaimAkses{PenggunaID: penggunaUji.ID, Peran: penggunaUji.Peran})
	return usecase.HasilMasuk{Pengguna: penggunaUji, TokenAkses: akses, TokenRefresh: tokenRefreshUji, RefreshKedaluwarsa: time.Now().Add(time.Hour)}, err
}

func appUjiAuth(t *testing.T) (*fiber.App, *masukPalsu) {
	t.Helper()
	akses, err := token.Baru(strings.Repeat("r", token.PanjangMinRahasia), 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	masuk := &masukPalsu{akses: akses}
	cfg := config.Server{Auth: config.ServerAuth{
		BatasLaju: config.BatasLaju{MaksPermintaan: 100, Jendela: time.Minute},
		OTP:       domain.AturanOTP{JendelaPermintaan: time.Hour},
	}}
	return NewApp(cfg, Layanan{Masuk: masuk, Akses: akses}), masuk
}

func kirimAuth(t *testing.T, app *fiber.App, metode, path, badan string, ubah func(h map[string]string)) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), metode, path, strings.NewReader(badan))
	header := map[string]string{fiber.HeaderContentType: fiber.MIMEApplicationJSON}
	ubah(header)
	for k, v := range header {
		permintaan.Header.Set(k, v)
	}
	return kirim(t, app, permintaan)
}

func tanpaUbah(map[string]string) {}

func TestAuth_VerifikasiBerhasil_TokenAksesDiBadanRefreshDiCookieAman(t *testing.T) {
	app, _ := appUjiAuth(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/auth/otp/verifikasi", `{"email":"warga@contoh.id","kode":"123456"}`, tanpaUbah)

	data, _ := hasil.respons.Data.(map[string]any)
	cookie := hasil.header(fiber.HeaderSetCookie)
	for _, harap := range []string{"rantai_refresh=" + tokenRefreshUji, "path=/api/v1/auth", "HttpOnly", "secure", "SameSite=Lax"} {
		if !strings.Contains(cookie, harap) {
			t.Errorf("cookie %q tanpa %q", cookie, harap)
		}
	}
	if hasil.kode != fiber.StatusOK || data["token_akses"] == "" || strings.Contains(hasil.header(fiber.HeaderSetCookie), data["token_akses"].(string)) {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestAuth_KodeSalah_401TanpaCookie(t *testing.T) {
	app, _ := appUjiAuth(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/auth/otp/verifikasi", `{"email":"warga@contoh.id","kode":"000000"}`, tanpaUbah)

	if hasil.kode != fiber.StatusUnauthorized || hasil.header(fiber.HeaderSetCookie) != "" || hasil.respons.Message != "kode salah" {
		t.Errorf("kode %d, cookie %q, pesan %q", hasil.kode, hasil.header(fiber.HeaderSetCookie), hasil.respons.Message)
	}
}

func TestAuth_KirimOTPMelebihiBatasEmail_429DenganRetryAfter(t *testing.T) {
	app, _ := appUjiAuth(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/auth/otp/kirim", `{"email":"sering@contoh.id"}`, tanpaUbah)

	if hasil.kode != fiber.StatusTooManyRequests || hasil.header(fiber.HeaderRetryAfter) != "3600" {
		t.Errorf("kode %d, retry-after %q", hasil.kode, hasil.header(fiber.HeaderRetryAfter))
	}
}

func TestAuth_Refresh_MembacaCookieDanMenggantinya(t *testing.T) {
	app, masuk := appUjiAuth(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/auth/refresh", "", func(h map[string]string) {
		h[fiber.HeaderCookie] = namaCookieRefresh + "=" + tokenRefreshUji
	})

	if hasil.kode != fiber.StatusOK || masuk.refreshDiminta != tokenRefreshUji || !strings.Contains(hasil.header(fiber.HeaderSetCookie), tokenRefreshUji) {
		t.Errorf("kode %d, diminta %q, cookie %q", hasil.kode, masuk.refreshDiminta, hasil.header(fiber.HeaderSetCookie))
	}
}

func TestAuth_RefreshTidakBerlaku_401DanCookieDihapus(t *testing.T) {
	app, _ := appUjiAuth(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/auth/refresh", "", tanpaUbah)

	if hasil.kode != fiber.StatusUnauthorized || !strings.HasPrefix(hasil.header(fiber.HeaderSetCookie), namaCookieRefresh+"=;") {
		t.Errorf("kode %d, cookie %q", hasil.kode, hasil.header(fiber.HeaderSetCookie))
	}
}

func TestAuth_Saya_WajibTokenAkses(t *testing.T) {
	app, masuk := appUjiAuth(t)
	sesi, err := masuk.hasil()
	if err != nil {
		t.Fatal(err)
	}
	kasus := []struct {
		nama, otorisasi string
		harap           int
	}{
		{"tanpa token", "", fiber.StatusUnauthorized},
		{"token rusak", "Bearer rusak", fiber.StatusUnauthorized},
		{"token sah", "Bearer " + sesi.TokenAkses, fiber.StatusOK},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil := kirimAuth(t, app, fiber.MethodGet, "/api/v1/auth/saya", "", func(h map[string]string) {
				h[fiber.HeaderAuthorization] = k.otorisasi
			})
			data, _ := hasil.respons.Data.(map[string]any)
			if hasil.kode != k.harap || (k.harap == fiber.StatusOK && data["email"] != penggunaUji.Email) {
				t.Errorf("kode %d, data %v", hasil.kode, data)
			}
		})
	}
}
