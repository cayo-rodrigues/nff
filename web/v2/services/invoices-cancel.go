package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListCancelings(ctx context.Context, userID int, filters ...map[string]string) ([]*models.InvoiceCancel, error) {
	f := models.NewFilters().Where("invoices_cancelings.created_by = ").Placeholder(userID)

	for _, filter := range filters {
		fromDate, fromDateOk := filter["from_date"]
		toDate, toDateOk := filter["to_date"]

		if fromDateOk && toDateOk {
			f.And().AsDate("invoices_cancelings.created_at").Between(fromDate, toDate)
		}
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
