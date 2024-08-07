package storage

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

func ListInvoiceItems(ctx context.Context, invoiceID int, userID int) ([]*models.InvoiceItem, error) {
	db := database.GetDB()

	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices_items WHERE invoice_id = $1 AND created_by = $2 ORDER BY id", invoiceID, userID)
	defer rows.Close()

	items := []*models.InvoiceItem{}

	for rows.Next() {
		item := models.NewInvoiceItem()
		err := Scan(rows, item)
		if err != nil {
			log.Println("Error scaning invoice item rows: ", err)
			return nil, utils.InternalServerErr
		}

		items = append(items, item)
	}

	return items, nil

}

func BulkCreateInvoiceItems(ctx context.Context, items []*models.InvoiceItem, invoiceID int, userID int) error {
	db := database.GetDB()

	rows := [][]interface{}{}
	for _, item := range items {
		item.InvoiceID = invoiceID
		item.CreatedBy = userID
		rows = append(rows, []interface{}{
			item.Group, item.Description, item.Origin, item.UnityOfMeasurement,
			item.Quantity, item.ValuePerUnity, item.InvoiceID, item.CreatedBy,
			item.NCM,
		})
	}
	_, err := db.PG.CopyFrom(
		ctx,
		pgx.Identifier{"invoices_items"},
		[]string{"item_group", "description", "origin", "unity_of_measurement", "quantity", "value_per_unity", "invoice_id", "created_by", "ncm"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Println("Error when running bulk insert invoice items query: ", err)
		return utils.InternalServerErr
	}

	return nil
}
