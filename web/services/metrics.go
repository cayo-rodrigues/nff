package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListMetrics(ctx context.Context, filters ...map[string]string) ([]*models.Metrics, error) {
	userID := utils.GetUserID(ctx)

	f := models.NewFilters().Where("metrics_history.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for i, filter := range filters {
		if i == 0 {
			handleDateFilters("metrics_history.created_at", filter, f)
			handleEntityFilters("metrics_history.entity_id", filter, f)
		}
	}

	f.OrderBy("metrics_history.created_at").Desc()

	return storage.ListMetrics(ctx, userID, f)
}

func CreateMetrics(ctx context.Context, m *models.Metrics) error {
	userID := utils.GetUserID(ctx)
	m.CreatedBy = userID
	return storage.CreateMetrics(ctx, m)
}

func RetrieveMetrics(ctx context.Context, metricsID int) (*models.Metrics, error) {
	userID := utils.GetUserID(ctx)
	return storage.RetrieveMetrics(ctx, metricsID, userID)
}

func RetrieveMetricsResult(ctx context.Context, resultID int) (*models.MetricsResult, error) {
	userID := utils.GetUserID(ctx)
	return storage.RetrieveMetricsResult(ctx, resultID, userID)
}

