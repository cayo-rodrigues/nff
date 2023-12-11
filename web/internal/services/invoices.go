package services

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
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

func (s *InvoiceService) ListInvoices(ctx context.Context, userID int, filters map[string]string) ([]*models.Invoice, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
			WHERE invoices.created_by = $1
	`)

	params := []interface{}{userID}
	paramCounter := 2

	now := time.Now()
	fromDate, ok := filters["from_date"]
	if !ok || fromDate == "" {
		fromDate = utils.FormatDate(now.Add(-10 * 24 * time.Hour))
	}
	toDate, ok := filters["to_date"]
	if !ok || toDate == "" {
		toDate = utils.FormatDate(now)
	}

	query.WriteString(" AND CAST(invoices.created_at AS DATE) BETWEEN $")
	query.WriteString(strconv.Itoa(paramCounter))
	params = append(params, fromDate)
	paramCounter++

	query.WriteString(" AND $")
	query.WriteString(strconv.Itoa(paramCounter))
	params = append(params, toDate)
	paramCounter++

	if entityID, ok := filters["entity_id"]; ok && entityID != ""{
		counter := strconv.Itoa(paramCounter)

		query.WriteString(" AND sender_id = $")
		query.WriteString(counter)

		query.WriteString(" OR recipient_id = $")
		query.WriteString(counter)

		params = append(params, entityID)
		paramCounter++
	}

	if senderID, ok := filters["sender_id"]; ok && senderID != "" {
		query.WriteString(" AND sender_id = $")
		query.WriteString(strconv.Itoa(paramCounter))
		params = append(params, senderID)
		paramCounter++
	}

	if recipientID, ok := filters["recipient_id"]; ok && recipientID != "" {
		query.WriteString(" AND recipient_id = $")
		query.WriteString(strconv.Itoa(paramCounter))
		params = append(params, recipientID)
		paramCounter++
	}

	query.WriteString(" ORDER BY invoices.created_at DESC")

	rows, _ := db.PG.Query(ctx, query.String(), params...)
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
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices (
				number, protocol, operation, cfop, is_final_customer, is_icms_contributor,
				shipping, add_shipping_to_total, gta, extra_notes, custom_file_name,
				sender_id, recipient_id, created_by
			)
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
		`SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
		WHERE invoices.id = $1 AND invoices.created_by = $2`,
		invoiceID, userID,
	)
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
	result, err := db.PG.Exec(
		ctx,
		`UPDATE invoices
			SET number = $1, protocol = $2, req_status = $3, req_msg = $4,
				invoice_pdf = $5, custom_file_name = $6, updated_at = $7
		WHERE id = $8 AND created_by = $9`,
		invoice.Number, invoice.Protocol, invoice.ReqStatus, invoice.ReqMsg,
		invoice.PDF, invoice.CustomFileName, time.Now(), invoice.ID, invoice.CreatedBy,
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
