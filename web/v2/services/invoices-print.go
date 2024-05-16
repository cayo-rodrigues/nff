package services

import (
	"context"
	"time"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListPrintings(ctx context.Context, userID int, filters ...map[string]string) ([]*models.InvoicePrint, error) {
	f := models.NewFilters().Where("invoices_printings.created_by = ").Placeholder(userID)

	for _, filter := range filters {
		_, fromDateOk := filter["from_date"]
		_, toDateOk := filter["to_date"]

		if fromDateOk && toDateOk {
			now := time.Now()
			fromDate := utils.FormatedNDaysBefore(now, utils.DefaultFiltersDaysRange)
			toDate := utils.FormatDate(now)
			f.And().AsDate("invoices_printings.created_at").Between(fromDate, toDate)
		}
	}

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
