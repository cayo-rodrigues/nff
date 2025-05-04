package handlers

import (
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/gofiber/fiber/v2"
)

func PricesPage(c *fiber.Ctx) error {
	return Render(c, layouts.Base(pages.PricesPage()))
}
