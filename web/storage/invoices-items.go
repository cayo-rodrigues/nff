package storage

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListInvoiceItems(ctx context.Context, invoiceID int, userID int) ([]*models.InvoiceItem, error) {
	db := database.GetDB()

	rows, _ := db.SQLite.QueryContext(ctx, "SELECT * FROM invoices_items WHERE invoice_id = ? AND created_by = ? ORDER BY id", invoiceID, userID)
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

	return WithTransaction(ctx, db.SQLite, func(tx *sql.Tx) error {
		query := `
		INSERT INTO invoices_items (
			item_group, description, origin, unity_of_measurement, quantity, value_per_unity, 
			invoice_id, created_by, ncm
		)
		VALUES `

		values := []interface{}{}
		valuePlaceholders := []string{}

		for _, item := range items {
			item.InvoiceID = invoiceID
			item.CreatedBy = userID

			valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?)")
			values = append(values, item.Group, item.Description, item.Origin,
				item.UnityOfMeasurement, item.Quantity, item.ValuePerUnity,
				item.InvoiceID, item.CreatedBy, item.NCM,
			)
		}

		query += strings.Join(valuePlaceholders, ", ")

		_, err := tx.ExecContext(ctx, query, values...)
		if err != nil {
			log.Println("Error executing bulk insert: ", err)
			return utils.InternalServerErr
		}

		return nil
	})
}
