package config

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	prefork := os.Getenv("FIBER_PREFORK") == "true"
	var app = fiber.New(fiber.Config{
		AppName:      os.Getenv("APP_NAME"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      prefork,
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
