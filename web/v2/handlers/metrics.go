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

	filters := c.Queries()

	metricsList, err := services.ListMetrics(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	a := services.GroupListByDate(metricsList)

	m := models.NewMetrics()

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter")
	return Render(c, layouts.Base(pages.MetricsPage(a, m, entities)))
}

func GenerateMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

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
		c.Append("HX-Trigger-After-Swap", "reload-metrics-list")
	}

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	return Render(c, forms.MetricsForm(metrics, entities))
}

func ListMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	filters := c.Queries()
	metrics, err := services.ListMetrics(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	a := services.GroupListByDate(metrics)

	return Render(c, components.MetricsList(a))
}

func GetMetricsForm(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	baseMetricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	baseMetrics, err := services.RetrieveMetrics(c.Context(), baseMetricsID, userID)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
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
