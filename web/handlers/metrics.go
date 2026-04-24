package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

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

	entitiesByType := models.NewEntitiesByType(entities)

	metricsByDate := services.GroupListByDate(metricsList)

	m := models.NewMetrics()

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(pages.MetricsPage(metricsByDate, m, entitiesByType)))
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
	decryptionKey, err := services.GetEncryptionKeySession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

	metrics := models.NewMetricsFromForm(c)

	var entities []*models.Entity

	err = utils.Concurrent(
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
	entitiesByType := models.NewEntitiesByType(entities)

	if !metrics.IsValid() {
		return Render(c, forms.MetricsForm(metrics, entitiesByType.Senders))
	}

	err = services.CreateMetrics(c.Context(), metrics)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.GetMetrics(metrics)

	c.Append("HX-Trigger-After-Swap", "reload-metrics-list")
	return Render(c, forms.MetricsForm(metrics, entitiesByType.Senders))
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
	entitiesByType := models.NewEntitiesByType(entities)

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.MetricsForm(baseMetrics, entitiesByType.Senders))
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

func DownloadMetricsRecordsZip(c *fiber.Ctx) error {
	metricsID, err := c.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	var body struct {
		RecordIDs []int `json:"record_ids"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.RecordIDs) == 0 {
		return fiber.ErrBadRequest
	}

	metrics, err := services.RetrieveMetrics(c.Context(), metricsID)
	if err != nil {
		return err
	}

	selectedIDs := make(map[int]bool, len(body.RecordIDs))
	for _, id := range body.RecordIDs {
		selectedIDs[id] = true
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	httpClient := &http.Client{Timeout: 30 * time.Second}

	for _, record := range metrics.Records {
		if !selectedIDs[record.ID] || record.InvoicePDF == "" {
			continue
		}

		resp, err := httpClient.Get(record.InvoicePDF)
		if err != nil {
			log.Printf("Failed to fetch PDF for record %d: %v", record.ID, err)
			continue
		}

		f, err := zw.Create(fmt.Sprintf("NFA-%s.pdf", record.InvoiceNumber))
		if err != nil {
			resp.Body.Close()
			continue
		}

		io.Copy(f, resp.Body)
		resp.Body.Close()
	}

	if err := zw.Close(); err != nil {
		return err
	}

	c.Set("Content-Type", "application/zip")
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="metricas-%d.zip"`, metricsID))
	return c.Send(buf.Bytes())
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

	printing, err := services.RetrievePrinting(c.Context(), result.PrintingID)
	if err != nil {
		return err
	}

	return Render(c, components.DownloadInvoiceFromMetricsRecordButton(result, printing))
}
