package interfaces

import "github.com/cayo-rodrigues/nff/web/internal/models"

type SiareBGWorker interface {
	RequestInvoice(invoice *models.Invoice)
	RequestInvoiceCanceling(invoiceCancel *models.InvoiceCancel)
}
