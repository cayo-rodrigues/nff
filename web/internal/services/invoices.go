package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

type InvoiceService struct {
	entityService interfaces.EntityService
	itemsService  interfaces.ItemsService
}

func NewInvoiceService(entityService interfaces.EntityService, itemsService interfaces.ItemsService) *InvoiceService {
	return &InvoiceService{
		entityService: entityService,
		itemsService:  itemsService,
	}
}

// TODO accept filters
func (s *InvoiceService) ListInvoices(ctx context.Context, userID int) ([]*models.Invoice, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices WHERE created_by = $1 ORDER BY id DESC", userID)
	defer rows.Close()

	invoices := []*models.Invoice{}

	for rows.Next() {
		invoice := models.NewEmptyInvoice()
		err := invoice.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice rows: ", err)
			return nil, utils.InternalServerErr
		}

		// TODO async data aggregation with go routines

		sender, err := s.entityService.RetrieveEntity(ctx, invoice.Sender.ID, invoice.CreatedBy)
		if err != nil {
			log.Println("Error linking invoice to sender: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Sender = sender

		recipient, err := s.entityService.RetrieveEntity(ctx, invoice.Recipient.ID, invoice.CreatedBy)
		if err != nil {
			log.Println("Error linking invoice to recipient: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Recipient = recipient

		items, err := s.itemsService.ListInvoiceItems(ctx, invoice.ID, invoice.CreatedBy)
		if err != nil {
			log.Println("Error linking invoice to items: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Items = items

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices
			(number, protocol, operation, cfop, is_final_customer, is_icms_contributor, shipping, add_shipping_to_total, gta, extra_notes, custom_file_name, sender_id, recipient_id, created_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, req_status, req_msg`,
		invoice.Number, invoice.Protocol, invoice.Operation, invoice.Cfop, invoice.IsFinalCustomer, invoice.IsIcmsContributor,
		invoice.Shipping, invoice.AddShippingToTotal, invoice.Gta, invoice.ExtraNotes, invoice.CustomFileName, invoice.Sender.ID,
		invoice.Recipient.ID, invoice.CreatedBy,
	)
	err := row.Scan(&invoice.ID, &invoice.ReqStatus, &invoice.ReqMsg)
	if err != nil {
		log.Println("Error when running insert invoice query: ", err)
		return utils.InternalServerErr
	}

	err = s.itemsService.BulkCreateInvoiceItems(ctx, invoice.Items, invoice.ID, invoice.CreatedBy)
	if err != nil {
		log.Println("Error running create invoice items query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *InvoiceService) RetrieveInvoice(ctx context.Context, invoiceID int, userID int) (*models.Invoice, error) {
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices WHERE invoices.id = $1 AND created_by = $2",
		invoiceID, userID,
	)

	invoice := models.NewEmptyInvoice()
	err := invoice.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice with id %v not found: %v", invoiceID, err)
		return nil, utils.InvoiceNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning invoice row: ", err)
		return nil, utils.InternalServerErr
	}

	// TODO async data aggregation with go routines

	sender, err := s.entityService.RetrieveEntity(ctx, invoice.Sender.ID, invoice.CreatedBy)
	if err != nil {
		log.Println("Error linking invoice to sender: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Sender = sender

	recipient, err := s.entityService.RetrieveEntity(ctx, invoice.Recipient.ID, invoice.CreatedBy)
	if err != nil {
		log.Println("Error linking invoice to recipient: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Recipient = recipient

	items, err := s.itemsService.ListInvoiceItems(ctx, invoice.ID, userID)
	if err != nil {
		log.Println("Error linking invoice to items: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Items = items

	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	result, err := db.PG.Exec(
		ctx,
		"UPDATE invoices SET number = $1, protocol = $2, req_status = $3, req_msg = $4, invoice_pdf = $5, custom_file_name = $6, updated_at = $7 WHERE id = $8 AND created_by = $9",
		invoice.Number, invoice.Protocol, invoice.ReqStatus, invoice.ReqMsg, invoice.PDF, invoice.CustomFileName, time.Now(), invoice.ID, invoice.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update invoice query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Invoice with id %v not found when running update query", invoice.ID)
		return utils.InvoiceNotFoundErr
	}

	return nil
}
