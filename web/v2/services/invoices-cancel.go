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

func CreateCancelingFromInvoiceID(ctx context.Context, invoiceID, userID int) (*models.InvoiceCancel, error) {
	invoice, err := RetrieveInvoice(ctx, invoiceID, userID)
	if err != nil {
		return nil, err
	}

	c := models.NewInvoiceCancel()

	c.CreatedBy = userID
	c.Entity = invoice.Sender
	c.InvoiceNumber = invoice.Number
	c.Year = invoice.CreatedAt.Year()
	c.Justification = "A nota possui dados incorretos"

	err = storage.CreateInvoiceCanceling(ctx, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
