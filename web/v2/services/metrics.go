package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListMetrics(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Metrics, error) {
	f := models.NewFilters().Where("metrics_history.created_by = ").Placeholder(userID)

	for _, filter := range filters {
		fromDate, fromDateOk := filter["from_date"]
		toDate, toDateOk := filter["to_date"]

		if fromDateOk && toDateOk {
			f.And().AsDate("metrics_history.created_at").Between(fromDate, toDate)
		}
	}

	f.OrderBy("metrics_history.created_at").Desc()

	return storage.ListMetrics(ctx, userID, f)
}

func CreateMetrics(ctx context.Context, m *models.Metrics, userID int) error {
	m.CreatedBy = userID
	return storage.CreateMetrics(ctx, m)
}

func RetrieveMetrics(ctx context.Context, printingID int, userID int) (*models.Metrics, error) {
	return storage.RetrieveMetrics(ctx, printingID, userID)
}
