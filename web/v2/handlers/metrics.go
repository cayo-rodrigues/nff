package handlers

import (
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func MetricsPage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	metricsList := []*models.Metrics{}
	m := models.NewMetrics()

	return Render(c, layouts.Base(pages.MetricsPage(metricsList, m, entities)))
}
