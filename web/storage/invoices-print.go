package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListInvoicePrintings(ctx context.Context, userID int, filters *models.Filters) ([]*models.InvoicePrint, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM invoices_printings
				JOIN entities ON entities.id = invoices_printings.entity_id
	`)

	query.WriteString(filters.String())

	db := database.GetDB()

	rows, _ := db.SQLite.QueryContext(ctx, query.String(), filters.Values()...)
	defer rows.Close()

	printings := []*models.InvoicePrint{}

	for rows.Next() {
		printing := models.NewInvoicePrint()
		err := Scan(rows, printing, printing.Entity)
		if err != nil {
			log.Println("Error scaning invoice printing rows: ", err)
			return nil, utils.InternalServerErr
		}

		printings = append(printings, printing)
	}

	return printings, nil
}

func CreateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO invoices_printings
			(invoice_id, invoice_id_type, custom_file_name_prefix, entity_id, created_by)
			VALUES (?, ?, ?, ?, ?)
		RETURNING id, req_status, req_msg, created_at, updated_at`,
		printing.InvoiceID, printing.InvoiceIDType, printing.CustomFileNamePrefix, printing.Entity.ID, printing.CreatedBy,
	)
	err := row.Scan(&printing.ID, &printing.ReqStatus, &printing.ReqMsg, &printing.CreatedAt, &printing.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert printing query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoicePrinting(ctx context.Context, printingID int, userID int) (*models.InvoicePrint, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`SELECT *
			FROM invoices_printings
				JOIN entities ON entities.id = invoices_printings.entity_id
		WHERE invoices_printings.id = ? AND invoices_printings.created_by = ?`,
		printingID, userID,
	)

	printing := models.NewInvoicePrint()
	err := Scan(row, printing, printing.Entity)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Invoice printing with id %v not found: %v", printingID, err)
		return nil, utils.PrintingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning printing row: ", err)
		return nil, utils.InternalServerErr
	}

	return printing, nil
}

func UpdateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error {
	db := database.GetDB()

	printing.UpdatedAt = time.Now()

	result, err := db.SQLite.ExecContext(
		ctx,
		`UPDATE invoices_printings
			SET req_status = ?, req_msg = ?, invoice_pdf = ?, custom_file_name_prefix = ?,
				file_name = ?, updated_at = ?
		WHERE id = ? AND created_by = ?`,
		printing.ReqStatus, printing.ReqMsg, printing.InvoicePDF, printing.CustomFileNamePrefix,
		printing.FileName, printing.UpdatedAt,
		printing.ID, printing.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update invoice printing query. Printing id: %d. Err: %v\n", printing.ID, err)
		return utils.InternalServerErr
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when getting rows affected by update invoice printing query. Printing id: %d. Err: %v\n", printing.ID, err)
		return utils.InternalServerErr
	}
	if rowsAffected == 0 {
		log.Printf("Printing with id %v not found when running update query", printing.ID)
		return utils.PrintingNotFoundErr
	}

	return nil
}
