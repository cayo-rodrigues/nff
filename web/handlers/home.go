package handlers

import (
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	page := pages.HomePage()
	return Render(c, layouts.Base(page))
}
