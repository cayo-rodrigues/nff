package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListMetrics(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Metrics, error) {
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

func CreateMetrics(ctx context.Context, m *models.Metrics, userID int) error {
	m.CreatedBy = userID
	return storage.CreateMetrics(ctx, m)
}

func RetrieveMetrics(ctx context.Context, metricsID int, userID int) (*models.Metrics, error) {
	return storage.RetrieveMetrics(ctx, metricsID, userID)
}

