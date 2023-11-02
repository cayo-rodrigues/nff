package interfaces

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/internal/models"
)

type EntityService interface {
	ListEntities(ctx context.Context) ([]*models.Entity, error)
	RetrieveEntity(ctx context.Context, entityId int) (*models.Entity, error)
	CreateEntity(ctx context.Context, entity *models.Entity) error
	UpdateEntity(ctx context.Context, entity *models.Entity) error
	DeleteEntity(ctx context.Context, entityId int) error
}

type CancelingService interface {
	ListInvoiceCancelings(ctx context.Context) ([]*models.InvoiceCancel, error)
	CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error
	RetrieveInvoiceCanceling(ctx context.Context, cancelingId int) (*models.InvoiceCancel, error)
	UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error
}

type InvoiceService interface {
	ListInvoices(ctx context.Context) ([]*models.Invoice, error)
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	RetrieveInvoice(ctx context.Context, invoiceId int) (*models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
}

type ItemsService interface {
	ListInvoiceItems(ctx context.Context, invoiceId int) ([]*models.InvoiceItem, error)
	BulkCreateInvoiceItems(ctx context.Context, items []*models.InvoiceItem, invoiceId int) error
}

type MetricsService interface {
	ListMetrics(ctx context.Context) ([]*models.MetricsQuery, error)
	CreateMetrics(ctx context.Context, query *models.MetricsQuery) error
	RetrieveMetrics(ctx context.Context, queryId int) (*models.MetricsQuery, error)
	UpdateMetrics(ctx context.Context, query *models.MetricsQuery) error
}
