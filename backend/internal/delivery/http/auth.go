package http

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

const (
	namaCookieRefresh = "rantai_refresh"
	pathCookieRefresh = prefiksAPI + "/auth"
	awalanBearer      = "Bearer "
	kunciKlaimAkses   = "klaim_akses"
)

type LayananMasuk interface {
	KirimOTP(ctx context.Context, email string) error
	Verifikasi(ctx context.Context, email, kode string) (usecase.HasilMasuk, error)
	Segarkan(ctx context.Context, tokenRefresh string) (usecase.HasilMasuk, error)
	Keluar(ctx context.Context, tokenRefresh string) error
	Saya(ctx context.Context, penggunaID int64) (domain.Pengguna, error)
}

type PemeriksaAkses interface {
	Periksa(token string) (domain.KlaimAkses, error)
}

type permintaanOTPJSON struct {
	Email string `json:"email"`
	Kode  string `json:"kode"`
}

type penggunaJSON struct {
	ID           int64        `json:"id"`
	Email        string       `json:"email"`
	NamaTampilan string       `json:"nama_tampilan"`
	Peran        domain.Peran `json:"peran"`
	InstansiID   *int64       `json:"instansi_id"`
	AnggotaID    *string      `json:"anggota_id"`
	IsSimulation bool         `json:"is_simulation"`
}

type sesiJSON struct {
	TokenAkses string       `json:"token_akses"`
	Pengguna   penggunaJSON `json:"pengguna"`
}

type handlerAuth struct {
	layanan    LayananMasuk
	jendelaOTP time.Duration
}

func daftarkanAuth(api fiber.Router, layanan Layanan, cfg config.ServerAuth) {
	if layanan.Masuk == nil || layanan.Akses == nil {
		return
	}
	h := handlerAuth{layanan: layanan.Masuk, jendelaOTP: cfg.OTP.JendelaPermintaan}
	batas := limiter.New(limiter.Config{
		Max:          cfg.BatasLaju.MaksPermintaan,
		Expiration:   cfg.BatasLaju.Jendela,
		LimitReached: func(fiber.Ctx) error { return fiber.ErrTooManyRequests },
	})
	auth := api.Group("/auth")
	auth.Post("/otp/kirim", batas, h.kirimOTP)
	auth.Post("/otp/verifikasi", batas, h.verifikasi)
	auth.Post("/refresh", batas, h.segarkan)
	auth.Post("/keluar", batas, h.keluar)
	auth.Get("/saya", wajibMasuk(layanan.Akses), h.saya)
}

func wajibMasuk(pemeriksa PemeriksaAkses) fiber.Handler {
	return func(c fiber.Ctx) error {
		token, ada := strings.CutPrefix(c.Get(fiber.HeaderAuthorization), awalanBearer)
		if !ada {
			return domain.ErrTokenAksesTidakBerlaku
		}
		klaim, err := pemeriksa.Periksa(token)
		if err != nil {
			return err
		}
		fiber.Locals(c, kunciKlaimAkses, klaim)
		return c.Next()
	}
}

func klaimAkses(c fiber.Ctx) domain.KlaimAkses {
	return fiber.Locals[domain.KlaimAkses](c, kunciKlaimAkses)
}

func (h handlerAuth) kirimOTP(c fiber.Ctx) error {
	masukan, err := bacaPermintaanOTP(c)
	if err != nil {
		return err
	}
	err = h.layanan.KirimOTP(c.Context(), masukan.Email)
	if errors.Is(err, domain.ErrBatasPermintaanOTP) {
		c.Set(fiber.HeaderRetryAfter, strconv.Itoa(int(h.jendelaOTP.Seconds())))
	}
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}

func (h handlerAuth) verifikasi(c fiber.Ctx) error {
	masukan, err := bacaPermintaanOTP(c)
	if err != nil {
		return err
	}
	hasil, err := h.layanan.Verifikasi(c.Context(), masukan.Email, masukan.Kode)
	if err != nil {
		return err
	}
	return kirimSesi(c, hasil)
}

func (h handlerAuth) segarkan(c fiber.Ctx) error {
	hasil, err := h.layanan.Segarkan(c.Context(), c.Cookies(namaCookieRefresh))
	if err != nil {
		hapusCookieRefresh(c)
		return err
	}
	return kirimSesi(c, hasil)
}

func (h handlerAuth) keluar(c fiber.Ctx) error {
	if err := h.layanan.Keluar(c.Context(), c.Cookies(namaCookieRefresh)); err != nil {
		return err
	}
	hapusCookieRefresh(c)
	return kirimData(c, nil)
}

func (h handlerAuth) saya(c fiber.Ctx) error {
	pengguna, err := h.layanan.Saya(c.Context(), klaimAkses(c).PenggunaID)
	if err != nil {
		return err
	}
	return kirimData(c, penggunaJSON(pengguna))
}

func bacaPermintaanOTP(c fiber.Ctx) (permintaanOTPJSON, error) {
	var masukan permintaanOTPJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return permintaanOTPJSON{}, domain.InputTidakValid("badan permintaan wajib JSON berisi email")
	}
	return masukan, nil
}

func kirimSesi(c fiber.Ctx, hasil usecase.HasilMasuk) error {
	pasangCookieRefresh(c, hasil.TokenRefresh, hasil.RefreshKedaluwarsa)
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, sesiJSON{TokenAkses: hasil.TokenAkses, Pengguna: penggunaJSON(hasil.Pengguna)})
}

func pasangCookieRefresh(c fiber.Ctx, token string, kedaluwarsa time.Time) {
	c.Cookie(&fiber.Cookie{
		Name: namaCookieRefresh, Value: token, Path: pathCookieRefresh, Expires: kedaluwarsa,
		HTTPOnly: true, Secure: true, SameSite: fiber.CookieSameSiteLaxMode,
	})
}

func hapusCookieRefresh(c fiber.Ctx) {
	pasangCookieRefresh(c, "", time.Unix(0, 0))
}
