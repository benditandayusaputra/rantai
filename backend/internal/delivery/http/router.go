package http

import (
	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

const prefiksAPI = "/api/v1"

func NewApp(cfg config.Server) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: tanganiError,
		ReadTimeout:  cfg.BatasWaktuBaca,
		WriteTimeout: cfg.BatasWaktuTulis,
	})
	api := app.Group(prefiksAPI)
	api.Get("/kesehatan", periksaKesehatan)
	return app
}
