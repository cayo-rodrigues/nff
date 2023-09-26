package workers

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

func ListInvoices(ctx context.Context) (*[]models.Invoice, error) {
	dbpool := sql.GetDatabasePool()
	rows, _ := dbpool.Query(ctx, "SELECT * FROM invoices ORDER BY id DESC")
	defer rows.Close()

	invoices := []models.Invoice{}

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
		}
		invoice.Sender = sender

		recipient, err := RetrieveEntity(ctx, invoice.Recipient.Id)
		if err != nil {
			log.Println("Error linking invoice to recipient: ", err)
		}
		invoice.Recipient = recipient

		invoices = append(invoices, *invoice)
	}

	return &invoices, nil
}

func CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
		ctx,
		`INSERT INTO invoices
			(number, protocol, operation, cfop, is_final_customer, is_icms_contributor, shipping, add_shipping_to_total, gta, sender_id, recipient_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		invoice.Number, invoice.Protocol, invoice.Operation, invoice.Cfop, invoice.IsFinalCustomer, invoice.IsIcmsContributor,
		invoice.Shipping, invoice.AddShippingToTotal, invoice.Gta, invoice.Sender.Id, invoice.Recipient.Id,
	)
	err := row.Scan(&invoice.Id)
	if err != nil {
		log.Println("Error when running insert invoice query: ", err)
		return utils.InternalServerErr
	}

	return nil
}
