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
	var entities []*models.Entity
	var metricsList []*models.Metrics

	err := utils.Concurrent(
		func() error {
			var err error
			entities, err = services.ListEntities(c.Context())
			return err
		},
		func() error {
			var err error
			filters := c.Queries()
			metricsList, err = services.ListMetrics(c.Context(), filters)
			return err
		},
	)

	if err != nil {
		return err
	}

	metricsByDate := services.GroupListByDate(metricsList)

	m := models.NewMetrics()

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(pages.MetricsPage(metricsByDate, m, entities)))
}

func MetricsDetailsPage(c *fiber.Ctx) error {
	metricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Settle", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(pages.MetricsDetailsPage(metrics)))
}

func GenerateMetrics(c *fiber.Ctx) error {
	metrics := models.NewMetricsFromForm(c)

	var entities []*models.Entity

	err := utils.Concurrent(
		func() error {
			var err error
			metrics.Entity, err = services.RetrieveEntity(c.Context(), metrics.Entity.ID)
			return err
		},
		func() error {
			var err error
			entities, err = services.ListEntities(c.Context())
			return err
		},
	)
	if err != nil {
		return err
	}

	if !metrics.IsValid() {
		return Render(c, forms.MetricsForm(metrics, entities))
	}

	err = services.CreateMetrics(c.Context(), metrics)
	if err != nil {
		return err
	}

	decryptionKey, err := services.GetEncryptionKeyFromSession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.GetMetrics(metrics)

	c.Append("HX-Trigger-After-Swap", "reload-metrics-list")
	return Render(c, forms.MetricsForm(metrics, entities))
}

func ListMetrics(c *fiber.Ctx) error {
	filters := c.Queries()
	metrics, err := services.ListMetrics(c.Context(), filters)
	if err != nil {
		return err
	}

	metricsByDate := services.GroupListByDate(metrics)

	return Render(c, components.MetricsList(metricsByDate))
}

func GetMetricsForm(c *fiber.Ctx) error {
	baseMetricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	var baseMetrics *models.Metrics
	var entities []*models.Entity

	err = utils.Concurrent(
		func() error {
			var err error
			baseMetrics, err = services.RetrieveMetrics(c.Context(), baseMetricsID)
			return err
		},
		func() error {
			var err error
			entities, err = services.ListEntities(c.Context())
			return err
		},
	)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.MetricsForm(baseMetrics, entities))
}

func RetrieveMetricsResultsDetails(c *fiber.Ctx) error {
	metricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID)
	if err != nil {
		return err
	}

	return Render(c, components.MetricsResultsDetails(metrics))
}

func RetrieveMetricsCard(c *fiber.Ctx) error {
	metricsID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID)
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

	// TODO
	// Garantir que result.ReqStatus possua um valor

	if result.InvoicePDF != "" {
		return Render(c, components.DownloadInvoiceFromRecordSuccessIcon(result))
	}

	if result.ReqStatus == "error" {
		return Render(c, components.DownloadInvoiceFromRecordErrorIcon())
	}

	return Render(c, components.DownloadInvoiceFromRecordLoadingIcon(result))
}
