package handlers

import (
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	isAuthenticated := utils.IsAuthenticated(c)
	page := pages.HomePage()
	return Render(c, layouts.Base(page, isAuthenticated))
}
