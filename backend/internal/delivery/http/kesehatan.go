package http

import "github.com/gofiber/fiber/v3"

const statusSehat = "sehat"

type statusKesehatan struct {
	Status string `json:"status"`
}

func periksaKesehatan(c fiber.Ctx) error {
	return kirimData(c, statusKesehatan{Status: statusSehat})
}
