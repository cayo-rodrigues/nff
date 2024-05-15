package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListMetrics(ctx context.Context, userID int) ([]*models.Metrics, error) {
	f := models.NewFilters().Where("metrics_history.created_by = ").Placeholder(userID)
	f.OrderBy("metrics_history.created_at")
	return storage.ListMetrics(ctx, userID, f)
}

func CreateMetrics(ctx context.Context, c *models.Metrics, userID int) error {
	c.CreatedBy = userID
	return storage.CreateMetrics(ctx, c)
}

func RetrieveMetrics(ctx context.Context, printingID int, userID int) (*models.Metrics, error) {
	return storage.RetrieveMetrics(ctx, printingID, userID)
}
