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

func ListInvoices(ctx context.Context, userID int) ([]*models.Invoice, error) {
	return storage.ListInvoices(ctx, userID, map[string]string{})
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
