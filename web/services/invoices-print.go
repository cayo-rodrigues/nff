package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListPrintings(ctx context.Context, filters ...map[string]string) ([]*models.InvoicePrint, error) {
	userID := utils.GetUserID(ctx)
	f := models.NewFilters().Where("invoices_printings.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for _, filter := range filters {
		handleDateFilters("invoices_printings.created_at", filter, f)
		handleEntityFilters("invoices_printings.entity_id", filter, f)
	}

	f.OrderBy("invoices_printings.created_at").Desc()

	return storage.ListInvoicePrintings(ctx, userID, f)
}

func CreatePrinting(ctx context.Context, p *models.InvoicePrint) error {
	userID := utils.GetUserID(ctx)
	p.CreatedBy = userID
	return storage.CreateInvoicePrinting(ctx, p)
}

func RetrievePrinting(ctx context.Context, printingID int) (*models.InvoicePrint, error) {
	userID := utils.GetUserID(ctx)
	return storage.RetrieveInvoicePrinting(ctx, printingID, userID)
}

func CreatePrintingFromMetricsRecord(ctx context.Context, invoiceNumber string, entityID int) (*models.InvoicePrint, error) {
	entity, err := RetrieveEntity(ctx, entityID)
	if err != nil {
		return nil, err
	}

	p := models.NewInvoicePrint()
	p.InvoiceID = invoiceNumber
	p.InvoiceIDType = models.InvoiceIDTypes.NFANumber()
	p.Entity = entity

	err = CreatePrinting(ctx, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
