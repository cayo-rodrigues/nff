package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListCancelings(ctx context.Context, filters ...map[string]string) ([]*models.InvoiceCancel, error) {
	userID := utils.GetUserID(ctx)

	f := models.NewFilters().Where("invoices_cancelings.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for _, filter := range filters {
		handleDateFilters("invoices_cancelings.created_at", filter, f)
		handleEntityFilters("invoices_cancelings.entity_id", filter, f)
	}

	f.OrderBy("invoices_cancelings.created_at").Desc()

	return storage.ListInvoiceCancelings(ctx, userID, f)
}

func CreateCanceling(ctx context.Context, c *models.InvoiceCancel) error {
	c.CreatedBy = utils.GetUserID(ctx)
	return storage.CreateInvoiceCanceling(ctx, c)
}

func RetrieveCanceling(ctx context.Context, cancelingID int) (*models.InvoiceCancel, error) {
	userID := utils.GetUserID(ctx)
	return storage.RetrieveInvoiceCanceling(ctx, cancelingID, userID)
}

func CreateCancelingFromInvoiceID(ctx context.Context, invoiceID int) (*models.InvoiceCancel, error) {
	userID := utils.GetUserID(ctx)

	invoice, err := RetrieveInvoice(ctx, invoiceID)
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
