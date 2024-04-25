package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListCancelings(ctx context.Context, userID int) ([]*models.InvoiceCancel, error) {
	return storage.ListInvoiceCancelings(ctx, userID, map[string]string{})
}

func CreateCanceling(ctx context.Context, c *models.InvoiceCancel, userID int) error {
	c.CreatedBy = userID
	return storage.CreateInvoiceCanceling(ctx, c)
}

func RetrieveCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	return storage.RetrieveInvoiceCanceling(ctx, cancelingID, userID)
}
