package storage

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListInvoiceCancelings(ctx context.Context, userID int, filters map[string]string) ([]*models.InvoiceCancel, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
	`)

	params := BuildQueryFilters(&query, filters, userID, "invoices_cancelings")

	db := database.GetDB()

	rows, _ := db.PG.Query(ctx, query.String(), params...)
	defer rows.Close()

	cancelings := []*models.InvoiceCancel{}

	for rows.Next() {
		canceling := models.NewInvoiceCancel()
		err := Scan(rows, canceling)
		if err != nil {
			log.Println("Error scaning invoice canceling rows: ", err)
			return nil, utils.InternalServerErr
		}

		cancelings = append(cancelings, canceling)
	}

	return cancelings, nil

}

func CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id, created_by)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id, req_status, req_msg`,
		canceling.InvoiceNumber, canceling.Year, canceling.Justification, canceling.Entity.ID, canceling.CreatedBy,
	)
	err := row.Scan(&canceling.ID, &canceling.ReqStatus, &canceling.ReqMsg)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoiceCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
		WHERE invoices_cancelings.id = $1 AND invoices_cancelings.created_by = $2`,
		cancelingID, userID,
	)

	canceling := models.NewInvoiceCancel()
	err := Scan(row, canceling)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice canceling with id %v not found: %v", cancelingID, err)
		return nil, utils.CancelingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning canceling row: ", err)
		return nil, utils.InternalServerErr
	}

	return canceling, nil
}

func UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	db := database.GetDB()

	result, err := db.PG.Exec(
		ctx,
		` UPDATE invoices_cancelings
			SET req_status = $1, req_msg = $2, updated_at = $3
		WHERE id = $4 AND created_by = $5`,
		canceling.ReqStatus, canceling.ReqMsg, time.Now(), canceling.ID, canceling.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update invoice canceling query. Canceling id: %d. Err: %v\n", canceling.ID, err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Canceling with id %v not found when running update query", canceling.ID)
		return utils.CancelingNotFoundErr
	}

	return nil
}

