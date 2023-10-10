package workers

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

// TODO accept filters
func ListInvoices(ctx context.Context) ([]*models.Invoice, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices ORDER BY id DESC")
	defer rows.Close()

	invoices := []*models.Invoice{}

	for rows.Next() {
		invoice := models.NewEmptyInvoice()
		err := invoice.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice rows: ", err)
			return nil, utils.InternalServerErr
		}

		sender, err := RetrieveEntity(ctx, invoice.Sender.Id)
		if err != nil {
			log.Println("Error linking invoice to sender: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Sender = sender

		recipient, err := RetrieveEntity(ctx, invoice.Recipient.Id)
		if err != nil {
			log.Println("Error linking invoice to recipient: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Recipient = recipient

		items, err := ListInvoiceItems(ctx, invoice.Id)
		if err != nil {
			log.Println("Error linking invoice to items: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.Items = items

		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

func CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices
			(number, protocol, operation, cfop, is_final_customer, is_icms_contributor, shipping, add_shipping_to_total, gta, sender_id, recipient_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, req_status, req_msg`,
		invoice.Number, invoice.Protocol, invoice.Operation, invoice.Cfop, invoice.IsFinalCustomer, invoice.IsIcmsContributor,
		invoice.Shipping, invoice.AddShippingToTotal, invoice.Gta, invoice.Sender.Id, invoice.Recipient.Id,
	)
	err := row.Scan(&invoice.Id, &invoice.ReqStatus, &invoice.ReqMsg)
	if err != nil {
		log.Println("Error when running insert invoice query: ", err)
		return utils.InternalServerErr
	}

	err = BulkCreateInvoiceItems(ctx, invoice.Items, invoice.Id)
	if err != nil {
		log.Println("Error running create invoice items query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoice(ctx context.Context, invoiceId int) (*models.Invoice, error) {
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices WHERE invoices.id = $1",
		invoiceId,
	)

	invoice := models.NewEmptyInvoice()
	err := invoice.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice with id %v not found: %v", invoiceId, err)
		return nil, utils.InvoiceNotFoundErr
	}
	if err != nil {
		log.Println("AQUI Error scaning invoice row: ", err)
		return nil, utils.InternalServerErr
	}

	sender, err := RetrieveEntity(ctx, invoice.Sender.Id)
	if err != nil {
		log.Println("Error linking invoice to sender: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Sender = sender

	recipient, err := RetrieveEntity(ctx, invoice.Recipient.Id)
	if err != nil {
		log.Println("Error linking invoice to recipient: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Recipient = recipient

	items, err := ListInvoiceItems(ctx, invoice.Id)
	if err != nil {
		log.Println("Error linking invoice to items: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.Items = items

	return invoice, nil
}

func UpdateInvoice(ctx context.Context, invoice *models.Invoice)  error {
	result, err := db.PG.Exec(
		ctx,
		"UPDATE invoices SET number = $1, protocol = $2, req_status = $3, req_msg = $4 WHERE id = $5",
		invoice.Number, invoice.Protocol, invoice.ReqStatus, invoice.ReqMsg, invoice.Id,
	)
	if err != nil {
		log.Println("Error when running update invoice query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Invoice with id %v not found when running update query", invoice.Id)
		return utils.InvoiceNotFoundErr
	}

	return nil
}
