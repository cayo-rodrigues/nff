package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	userID := utils.GetUserID(ctx)
	invoice.CreatedBy = userID
	err := storage.CreateInvoice(ctx, invoice)
	if err != nil {
		return err
	}

	return storage.BulkCreateInvoiceItems(ctx, invoice.Items, invoice.ID, userID)
}

func ListInvoices(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Invoice, error) {
	f := models.NewFilters().Where("invoices.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for _, filter := range filters {
		handleDateFilters("invoices.created_at", filter, f)
		handleEntityFilters("invoices.", filter, f)
	}

	f.OrderBy("invoices.created_at").Desc()
	return storage.ListInvoices(ctx, userID, f)
}

func RetrieveInvoice(ctx context.Context, invoiceID int) (*models.Invoice, error) {
	userID := utils.GetUserData(ctx).ID

	// TODO
	// DEIXAR BONITO
	invoice, err := storage.RetrieveInvoice(ctx, invoiceID, userID, "")
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

func RetrieveInvoiceByNumber(ctx context.Context, invoiceNumber string, userID int) (*models.Invoice, error) {
	// TODO
	// DEIXAR BONITO
	invoice, err := storage.RetrieveInvoice(ctx, 0, userID, invoiceNumber)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
