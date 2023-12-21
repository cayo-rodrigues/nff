package interfaces

import (
	"context"
	"strings"

	"github.com/cayo-rodrigues/nff/web/internal/models"
)

type EntityService interface {
	ListEntities(ctx context.Context, userID int) ([]*models.Entity, error)
	CreateEntity(ctx context.Context, entity *models.Entity) error
	RetrieveEntity(ctx context.Context, entityID int, userID int) (*models.Entity, error)
	UpdateEntity(ctx context.Context, entity *models.Entity) error
	DeleteEntity(ctx context.Context, entityID int, userID int) error
}

type InvoiceService interface {
	ListInvoices(ctx context.Context, userID int, filters map[string]string) ([]*models.Invoice, error)
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	RetrieveInvoice(ctx context.Context, invoiceID int, userID int) (*models.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *models.Invoice) error
}

type ItemsService interface {
	ListInvoiceItems(ctx context.Context, invoiceID int, userID int) ([]*models.InvoiceItem, error)
	BulkCreateInvoiceItems(ctx context.Context, items []*models.InvoiceItem, invoiceID int, userID int) error
}

type CancelingService interface {
	ListInvoiceCancelings(ctx context.Context, userID int, filters map[string]string) ([]*models.InvoiceCancel, error)
	CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error
	RetrieveInvoiceCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error)
	UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error
}

type PrintingService interface {
	ListInvoicePrintings(ctx context.Context, userID int, filters map[string]string) ([]*models.InvoicePrint, error)
	CreateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error
	RetrieveInvoicePrinting(ctx context.Context, printingID int, userID int) (*models.InvoicePrint, error)
	UpdateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error
}

type MetricsService interface {
	ListMetrics(ctx context.Context, userID int) ([]*models.MetricsQuery, error)
	CreateMetrics(ctx context.Context, query *models.MetricsQuery) error
	RetrieveMetrics(ctx context.Context, queryID int, userID int) (*models.MetricsQuery, error)
	UpdateMetrics(ctx context.Context, query *models.MetricsQuery) error
}

type UserService interface {
	RetrieveUser(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type FiltersService interface {
	BuildQueryFilters(query *strings.Builder, filters map[string]string, userID int, table string) []interface{}
}
