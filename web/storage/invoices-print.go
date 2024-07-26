package storage

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
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

	rows, _ := db.PG.Query(ctx, query.String(), filters.Values()...)
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

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices_printings
			(invoice_id, invoice_id_type, custom_file_name_prefix, entity_id, created_by)
			VALUES ($1, $2, $3, $4, $5)
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

	row := db.PG.QueryRow(
		ctx,
		`SELECT *
			FROM invoices_printings
				JOIN entities ON entities.id = invoices_printings.entity_id
		WHERE invoices_printings.id = $1 AND invoices_printings.created_by = $2`,
		printingID, userID,
	)

	printing := models.NewInvoicePrint()
	err := Scan(row, printing, printing.Entity)
	if errors.Is(err, pgx.ErrNoRows) {
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

	result, err := db.PG.Exec(
		ctx,
		`UPDATE invoices_printings
			SET req_status = $1, req_msg = $2, invoice_pdf = $3, custom_file_name_prefix = $4,
				file_name = $5, updated_at = $6
		WHERE id = $7 AND created_by = $8`,
		printing.ReqStatus, printing.ReqMsg, printing.InvoicePDF, printing.CustomFileNamePrefix,
		printing.FileName, printing.UpdatedAt,
		printing.ID, printing.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update invoice printing query. Printing id: %d. Err: %v\n", printing.ID, err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Printing with id %v not found when running update query", printing.ID)
		return utils.PrintingNotFoundErr
	}

	return nil
}
