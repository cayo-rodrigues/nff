package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func CreateInvoice(ctx context.Context, invoice *models.Invoice, userID int) error  {
	invoice.CreatedBy = userID	
	return storage.CreateInvoice(ctx, invoice)
}

func ListInvoices(ctx context.Context, userID int) ([]*models.Invoice, error) {
	return storage.ListInvoices(ctx, userID, map[string]string{})	
}
