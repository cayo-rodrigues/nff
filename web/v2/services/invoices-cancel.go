package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListCancelings(ctx context.Context, userID int) ([]*models.InvoiceCancel, error) {
	return storage.ListInvoiceCancelings(ctx, userID, map[string]string{})
}
