package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/siare"
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

	metricsByDate := services.GroupListByDate(metricsList)

	m := models.NewMetrics()

	isAuthenticated := utils.IsAuthenticated(c)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter, highlight-current-page")
	return Render(c, layouts.Base(pages.MetricsPage(metricsByDate, m, entities), isAuthenticated))
}

func GenerateMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	metrics := models.NewMetricsFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), metrics.Entity.ID, userID)
	if err != nil {
		return err
	}
	metrics.Entity = entity

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	if !metrics.IsValid() {
		return Render(c, forms.MetricsForm(metrics, entities))
	}

	err = services.CreateMetrics(c.Context(), metrics, userID)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient()
	go ssapi.GetMetrics(metrics)

	c.Append("HX-Trigger-After-Swap", "reload-metrics-list")
	return Render(c, forms.MetricsForm(metrics, entities))
}

func ListMetrics(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	filters := c.Queries()
	metrics, err := services.ListMetrics(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	metricsByDate := services.GroupListByDate(metrics)

	return Render(c, components.MetricsList(metricsByDate))
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
