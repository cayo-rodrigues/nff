package workers

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

func ListInvoiceItems(ctx context.Context, invoiceId int) (*[]models.InvoiceItem, error) {
	dbpool := sql.GetDatabasePool()
	rows, _ := dbpool.Query(ctx, "SELECT * FROM invoices_items WHERE invoice_id = $1 ORDER BY id", invoiceId)
	defer rows.Close()

	items := []models.InvoiceItem{}

	for rows.Next() {
		item := models.NewEmptyInvoiceItem()
		err := item.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice item rows: ", err)
			return nil, utils.InternalServerErr
		}

		items = append(items, *item)
	}

	return &items, nil

}

func BulkCreateInvoiceItems(ctx context.Context, items *[]models.InvoiceItem, invoiceId int) error {
	dbpool := sql.GetDatabasePool()
	rows := [][]interface{}{}
	for _, item := range *items {
		item.InvoiceId = invoiceId
		rows = append(rows, []interface{}{
			item.Group, item.Description, item.Origin, item.UnityOfMeasurement,
			item.Quantity, item.ValuePerUnity, item.InvoiceId,
		})
	}
	_, err := dbpool.CopyFrom(
		ctx,
		pgx.Identifier{"invoices_items"},
		[]string{"item_group", "description", "origin", "unity_of_measurement", "quantity", "value_per_unity", "invoice_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Println("Error when running bulk insert invoice items query: ", err)
		return utils.InternalServerErr
	}

	return nil
}
