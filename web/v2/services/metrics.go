package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListMetrics(ctx context.Context, userID int) ([]*models.Metrics, error) {
	return storage.ListMetrics(ctx, userID, map[string]string{})
}

func CreateMetrics(ctx context.Context, c *models.Metrics, userID int) error {
	c.CreatedBy = userID
	return storage.CreateMetrics(ctx, c)
}

func RetrieveMetrics(ctx context.Context, printingID int, userID int) (*models.Metrics, error) {
	return storage.RetrieveMetrics(ctx, printingID, userID)
}
