package http

import (
	"crypto/subtle"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

const (
	kebijakanKontenAPI = "default-src 'none'; frame-ancestors 'none'"
	kebijakanIzinAPI   = "geolocation=(), camera=(), microphone=()"
	kebijakanPerujuk   = "strict-origin-when-cross-origin"
	umurHSTSDetik      = 63072000
	headerIPKlienAsli  = "X-Klien-Asli"
	headerRahasiaProxy = "X-Rantai-Proxy"
)

func pasangPengaman(app *fiber.App, cfg config.Server) {
	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: kebijakanKontenAPI,
		PermissionPolicy:      kebijakanIzinAPI,
		ReferrerPolicy:        kebijakanPerujuk,
		HSTSMaxAge:            umurHSTSDetik,
	}))
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.CORS.Origin, AllowCredentials: len(cfg.CORS.Origin) > 0}))
	app.Use(pembatasLaju(cfg.BatasLaju, cfg.ProxyRahasia))
}

func pembatasLaju(b config.BatasLaju, rahasiaProxy string) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:          b.MaksPermintaan,
		Expiration:   b.Jendela,
		KeyGenerator: kunciKlien(rahasiaProxy),
		LimitReached: func(fiber.Ctx) error { return fiber.ErrTooManyRequests },
	})
}

func kunciKlien(rahasiaProxy string) func(fiber.Ctx) string {
	return func(c fiber.Ctx) string {
		ipAsli := c.Get(headerIPKlienAsli)
		rahasia := c.Get(headerRahasiaProxy)
		if rahasiaProxy == "" || ipAsli == "" || subtle.ConstantTimeCompare([]byte(rahasia), []byte(rahasiaProxy)) != 1 {
			return c.IP()
		}
		return ipAsli
	}
}
