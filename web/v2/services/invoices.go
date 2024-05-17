package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func CreateInvoice(ctx context.Context, invoice *models.Invoice, userID int) error {
	invoice.CreatedBy = userID
	err := storage.CreateInvoice(ctx, invoice)
	if err != nil {
		return err
	}

	return storage.BulkCreateInvoiceItems(ctx, invoice.Items, invoice.ID, userID)
}

func ListInvoices(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Invoice, error) {
	f := models.NewFilters().Where("invoices.created_by = ").Placeholder(userID)

	for _, filter := range filters {
		fromDate, fromDateOk := filter["from_date"]
		toDate, toDateOk := filter["to_date"]

		if fromDateOk && toDateOk {
			f.And().AsDate("invoices.created_at").Between(fromDate, toDate)
		}
	}

	f.OrderBy("invoices.created_at").Desc()
	return storage.ListInvoices(ctx, userID, f)
}

func RetrieveInvoice(ctx context.Context, invoiceID int, userID int) (*models.Invoice, error) {
	invoice, err := storage.RetrieveInvoice(ctx, invoiceID, userID)
	if err != nil {
		return nil, err
	}

	items, err := storage.ListInvoiceItems(ctx, invoiceID, userID)
	if err != nil {
		return nil, err
	}

	invoice.Items = items
	return invoice, nil
}
