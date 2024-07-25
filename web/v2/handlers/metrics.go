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
	userID := utils.GetUserData(c.Context()).ID
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

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(pages.MetricsPage(metricsByDate, m, entities)))
}

func GenerateMetrics(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

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
	userID := utils.GetUserData(c.Context()).ID
	filters := c.Queries()
	metrics, err := services.ListMetrics(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	metricsByDate := services.GroupListByDate(metrics)

	return Render(c, components.MetricsList(metricsByDate))
}

func GetMetricsForm(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

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
	userID := utils.GetUserData(c.Context()).ID
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

func RetrieveMetricsCard(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	metricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID, userID)
	if err != nil {
		return err
	}

	return Render(c, components.MetricsCard(metrics))
}

func GetDownloadFromRecordStatusIcon(c *fiber.Ctx) error {
	recordID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	result, err := services.RetrieveMetricsResult(c.Context(), recordID)
	if err != nil {
		return err
	}

	if result.InvoicePDF != "" {
		return Render(c, components.DownloadInvoiceFromRecordSuccessIcon(result))
	}
	return Render(c, components.DownloadInvoiceFromRecordLoadingIcon(result))
}
