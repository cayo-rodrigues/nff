package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListInvoices(ctx context.Context, userID int, filters *models.Filters) ([]*models.Invoice, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
	`)

	query.WriteString(filters.String())

	db := database.GetDB()

	rows, _ := db.SQLite.QueryContext(ctx, query.String(), filters.Values()...)
	defer rows.Close()

	invoices := []*models.Invoice{}

	for rows.Next() {
		invoice := models.NewInvoice()
		err := Scan(rows, invoice, invoice.Sender, invoice.Recipient)
		if err != nil {
			log.Println("Error scaning invoice rows: ", err)
			return nil, utils.InternalServerErr
		}

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO invoices (
				number, protocol, operation, cfop, is_final_customer, is_icms_contributor,
				shipping, add_shipping_to_total, gta, extra_notes, custom_file_name_prefix, file_name,
				sender_id, recipient_id, created_by, sender_ie, recipient_ie
			)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, req_status, req_msg, created_at, updated_at`,
		invoice.Number, invoice.Protocol, invoice.Operation, invoice.Cfop, invoice.IsFinalCustomer, invoice.IsIcmsContributor,
		invoice.Shipping, invoice.AddShippingToTotal, invoice.Gta, invoice.ExtraNotes, invoice.CustomFileNamePrefix, invoice.FileName,
		invoice.Sender.ID, invoice.Recipient.ID, invoice.CreatedBy, invoice.SenderIe, invoice.RecipientIe,
	)
	err := row.Scan(&invoice.ID, &invoice.ReqStatus, &invoice.ReqMsg, &invoice.CreatedAt, &invoice.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert invoice query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

// TODO
// DEIXAR BONITO
func RetrieveInvoice(ctx context.Context, invoiceID int, userID int, invoiceNumber string) (*models.Invoice, error) {
	db := database.GetDB()

	searchColumn := "id"
	queryValues := []any{}
	if invoiceID != 0 && invoiceNumber == "" {
		queryValues = append(queryValues, invoiceID)
	}
	if invoiceID == 0 && invoiceNumber != "" {
		searchColumn = "number"
		queryValues = append(queryValues, invoiceNumber)
	}
	queryValues = append(queryValues, userID)

	row := db.SQLite.QueryRowContext(
		ctx,
		fmt.Sprintf(
			`SELECT *
			FROM invoices
				JOIN entities AS senders ON invoices.sender_id = senders.id
				JOIN entities AS recipients ON invoices.recipient_id = recipients.id
		WHERE invoices.%s = ? AND invoices.created_by = ?`, searchColumn,
		),
		queryValues...,
	)
	invoice := models.NewInvoice()
	err := Scan(row, invoice, invoice.Sender, invoice.Recipient)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Invoice with id %v not found: %v", invoiceID, err)
		return nil, utils.InvoiceNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning invoice row: ", err)
		return nil, utils.InternalServerErr
	}

	items, err := ListInvoiceItems(ctx, invoice.ID, userID)
	if err != nil {
		log.Println("Error linking invoice to items: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Items = items

	return invoice, nil
}

func UpdateInvoice(ctx context.Context, invoice *models.Invoice) error {
	db := database.GetDB()

	invoice.UpdatedAt = time.Now()

	result, err := db.SQLite.ExecContext(
		ctx,
		`UPDATE invoices
			SET number = ?, protocol = ?, req_status = ?, req_msg = ?,
				invoice_pdf = ?, custom_file_name_prefix = ?, updated_at = ?, 
				sender_ie = ?, recipient_ie = ?
		WHERE id = ? AND created_by = ?`,
		invoice.Number, invoice.Protocol, invoice.ReqStatus, invoice.ReqMsg,
		invoice.PDF, invoice.CustomFileNamePrefix, invoice.UpdatedAt,
		invoice.SenderIe, invoice.RecipientIe,
		invoice.ID, invoice.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update invoice query. Invoice id: %d. Err: %v\n", invoice.ID, err)
		return utils.InternalServerErr
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Error when getting rows affected by update invoice query. Invoice id: %d. Err: %v\n", invoice.ID, err)
		return err
	}

	if rowsAffected == 0 {
		log.Printf("Invoice with id %v not found when running update query", invoice.ID)
		return utils.InvoiceNotFoundErr
	}

	return nil
}
