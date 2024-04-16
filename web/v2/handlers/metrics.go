package handlers

import (
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/gofiber/fiber/v2"
)

func MetricsPage(c *fiber.Ctx) error  {
	metrics := models.NewMetrics()
	return Render(c, layouts.Base(pages.MetricsPage(metrics)))
}
