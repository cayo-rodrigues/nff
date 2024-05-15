package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func ListPrintings(ctx context.Context, userID int) ([]*models.InvoicePrint, error) {
	f := models.NewFilters().Where("invoices_printings.created_by = ").Placeholder(userID)
	f.OrderBy("invoices_printings.created_at").Desc()
	return storage.ListInvoicePrintings(ctx, userID, f)
}

func CreatePrinting(ctx context.Context, c *models.InvoicePrint, userID int) error {
	c.CreatedBy = userID
	return storage.CreateInvoicePrinting(ctx, c)
}

func RetrievePrinting(ctx context.Context, printingID int, userID int) (*models.InvoicePrint, error) {
	return storage.RetrieveInvoicePrinting(ctx, printingID, userID)
}
