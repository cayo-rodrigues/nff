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

func ListInvoiceCancelings(ctx context.Context, userID int, filters *models.Filters) ([]*models.InvoiceCancel, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
	`)

	query.WriteString(filters.String())

	db := database.GetDB()

	rows, _ := db.SQLite.QueryContext(ctx, query.String(), filters.Values()...)
	defer rows.Close()

	cancelings := []*models.InvoiceCancel{}

	for rows.Next() {
		canceling := models.NewInvoiceCancel()
		err := Scan(rows, canceling, canceling.Entity)
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

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id, created_by)
			VALUES (?, ?, ?, ?, ?)
		RETURNING id, req_status, req_msg, created_at, updated_at`,
		canceling.InvoiceNumber, canceling.Year, canceling.Justification, canceling.Entity.ID, canceling.CreatedBy,
	)
	err := row.Scan(&canceling.ID, &canceling.ReqStatus, &canceling.ReqMsg, &canceling.CreatedAt, &canceling.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoiceCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
		WHERE invoices_cancelings.id = ? AND invoices_cancelings.created_by = ?`,
		cancelingID, userID,
	)

	canceling := models.NewInvoiceCancel()
	err := Scan(row, canceling, canceling.Entity)
	if errors.Is(err, sql.ErrNoRows) {
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

	result, err := db.SQLite.ExecContext(
		ctx,
		` UPDATE invoices_cancelings
			SET req_status = ?, req_msg = ?, updated_at = ?
		WHERE id = ? AND created_by = ?`,
		canceling.ReqStatus, canceling.ReqMsg, time.Now(), canceling.ID, canceling.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update invoice canceling query. Canceling id: %d. Err: %v\n", canceling.ID, err)
		return utils.InternalServerErr
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when getting rows affected by update invoice canceling query. Canceling id: %d. Err: %v\n", canceling.ID, err)
		return utils.InternalServerErr
	}
	if rowsAffected == 0 {
		log.Printf("Canceling with id %v not found when running update query", canceling.ID)
		return utils.CancelingNotFoundErr
	}

	return nil
}

