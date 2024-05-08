package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
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

	metricsList, err := services.ListMetrics(c.Context(), userID)
	if err != nil {
		return err
	}

	m := models.NewMetrics()

	return Render(c, layouts.Base(pages.MetricsPage(metricsList, m, entities)))
}

func GenerateMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	metrics := models.NewMetricsFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), metrics.Entity.ID, userID)
	if err != nil {
		return err
	}
	metrics.Entity = entity

	if metrics.IsValid() {
		err := services.CreateMetrics(c.Context(), metrics, userID)
		if err != nil {
			return err
		}
		c.Set("HX-Trigger-After-Swap", "reload-metrics-list")
	}

	return Render(c, forms.MetricsForm(metrics, entities))
}

func ListMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	metrics, err := services.ListMetrics(c.Context(), userID)
	if err != nil {
		return err
	}

	return Render(c, components.MetricsList(metrics))
}

func GetMetricsForm(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	baseMetricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	baseMetrics, err := services.RetrieveMetrics(c.Context(), baseMetricsID, userID)
	if err != nil {
		return err
	}

	c.Set("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.MetricsForm(baseMetrics, entities))
}

func RetrieveMetricsResultsDetails(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	metricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID, userID)
	if err != nil {
		return err
	}

	return Render(c, components.MetricsResultsDetails(metrics))
}
