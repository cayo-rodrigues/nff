package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(ctx *fiber.Ctx) error {
	return ctx.Render("home", fiber.Map{
		"IsAuthenticated": true,
	}, "layouts/base")
}
