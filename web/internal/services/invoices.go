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
	query := `
		SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
			WHERE invoices.created_by = $1
		ORDER BY invoices.id DESC
	`
	rows, _ := db.PG.Query(ctx, query, userID)
	defer rows.Close()

	invoices := []*models.Invoice{}

	for rows.Next() {
		invoice := models.NewEmptyInvoice()
		err := invoice.FullScan(rows)
		if err != nil {
			log.Println("Error scaning invoice rows: ", err)
			return nil, utils.InternalServerErr
		}

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
	query := `
		INSERT INTO invoices (
				number, protocol, operation, cfop, is_final_customer, is_icms_contributor,
				shipping, add_shipping_to_total, gta, extra_notes, custom_file_name,
				sender_id, recipient_id, created_by
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, req_status, req_msg
	`
	query_args := [14]any{
		invoice.Number, invoice.Protocol, invoice.Operation, invoice.Cfop, invoice.IsFinalCustomer, invoice.IsIcmsContributor,
		invoice.Shipping, invoice.AddShippingToTotal, invoice.Gta, invoice.ExtraNotes, invoice.CustomFileName, invoice.Sender.ID,
		invoice.Recipient.ID, invoice.CreatedBy,
	}
	row := db.PG.QueryRow(ctx, query, query_args)
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
	query := `
		SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
		WHERE invoices.id = $1 AND invoices.created_by = $2
	`
	row := db.PG.QueryRow(ctx, query, invoiceID, userID)

	invoice := models.NewEmptyInvoice()
	err := invoice.FullScan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice with id %v not found: %v", invoiceID, err)
		return nil, utils.InvoiceNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning invoice row: ", err)
		return nil, utils.InternalServerErr
	}

	items, err := s.itemsService.ListInvoiceItems(ctx, invoice.ID, userID)
	if err != nil {
		log.Println("Error linking invoice to items: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Items = items

	return invoice, nil
}

func (s *InvoiceService) UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	query := `
		UPDATE invoices
			SET number = $1, protocol = $2, req_status = $3, req_msg = $4,
				invoice_pdf = $5, custom_file_name = $6, updated_at = $7
		WHERE id = $8 AND created_by = $9
	`
	query_args := [9]any{
		invoice.Number, invoice.Protocol, invoice.ReqStatus, invoice.ReqMsg,
		invoice.PDF, invoice.CustomFileName, time.Now(), invoice.ID, invoice.CreatedBy,
	}
	result, err := db.PG.Exec(ctx, query, query_args)
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
