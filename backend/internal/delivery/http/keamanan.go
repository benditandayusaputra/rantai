package http

import (
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
)

func pasangPengaman(app *fiber.App, cfg config.Server) {
	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: kebijakanKontenAPI,
		PermissionPolicy:      kebijakanIzinAPI,
		ReferrerPolicy:        kebijakanPerujuk,
		HSTSMaxAge:            umurHSTSDetik,
	}))
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.CORS.Origin, AllowCredentials: len(cfg.CORS.Origin) > 0}))
	app.Use(limiter.New(limiter.Config{
		Max:          cfg.BatasLaju.MaksPermintaan,
		Expiration:   cfg.BatasLaju.Jendela,
		LimitReached: func(fiber.Ctx) error { return fiber.ErrTooManyRequests },
	}))
}
