package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListCancelings(ctx context.Context, userID int, filters ...map[string]string) ([]*models.InvoiceCancel, error) {
	f := models.NewFilters().Where("invoices_cancelings.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for _, filter := range filters {
		handleDateFilters("invoices_cancelings.created_at", filter, f)
	}

	f.OrderBy("invoices_cancelings.created_at").Desc()

	return storage.ListInvoiceCancelings(ctx, userID, f)
}

func CreateCanceling(ctx context.Context, c *models.InvoiceCancel, userID int) error {
	c.CreatedBy = userID
	return storage.CreateInvoiceCanceling(ctx, c)
}

func RetrieveCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	return storage.RetrieveInvoiceCanceling(ctx, cancelingID, userID)
}
