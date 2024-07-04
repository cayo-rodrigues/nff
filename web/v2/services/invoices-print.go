package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListPrintings(ctx context.Context, userID int, filters ...map[string]string) ([]*models.InvoicePrint, error) {
	f := models.NewFilters().Where("invoices_printings.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for _, filter := range filters {
		handleDateFilters("invoices_printings.created_at", filter, f)
	}

	f.OrderBy("invoices_printings.created_at").Desc()

	return storage.ListInvoicePrintings(ctx, userID, f)
}

func CreatePrinting(ctx context.Context, p *models.InvoicePrint, userID int) error {
	p.CreatedBy = userID
	return storage.CreateInvoicePrinting(ctx, p)
}

func RetrievePrinting(ctx context.Context, printingID int, userID int) (*models.InvoicePrint, error) {
	return storage.RetrieveInvoicePrinting(ctx, printingID, userID)
}

func CreatePrintingFromMetricsRecord(ctx context.Context, invoiceNumber string, entityID, userID int) (*models.InvoicePrint, error) {
	entity, err := RetrieveEntity(ctx, entityID, userID)
	if err != nil {
		return nil, err
	}

	p := models.NewInvoicePrint()
	p.InvoiceID = invoiceNumber
	p.InvoiceIDType = models.InvoiceIDTypes.NFANumber()
	p.Entity = entity

	err = CreatePrinting(ctx, p, userID)
	if err != nil {
		return nil, err
	}

	return p, nil
}
