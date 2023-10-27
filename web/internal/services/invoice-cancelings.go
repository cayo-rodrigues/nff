package services

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

func ListInvoiceCancelings(ctx context.Context) ([]*models.InvoiceCancel, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices_cancelings ORDER BY id DESC")
	defer rows.Close()

	cancelings := []*models.InvoiceCancel{}

	for rows.Next() {
		canceling := models.NewEmptyInvoiceCancel()
		err := canceling.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice canceling rows: ", err)
			return nil, utils.InternalServerErr
		}

		entity, err := RetrieveEntity(ctx, canceling.Entity.Id)
		if err != nil {
			log.Println("Error linking invoice canceling to entity: ", err)
			return nil, utils.InternalServerErr
		}
		canceling.Entity = entity

		cancelings = append(cancelings, canceling)
	}

	return cancelings, nil

}

func CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id)
			VALUES ($1, $2, $3, $4)
		RETURNING id, req_status, req_msg`,
		canceling.Number, canceling.Year, canceling.Justification, canceling.Entity.Id,
	)
	err := row.Scan(&canceling.Id, &canceling.ReqStatus, &canceling.ReqMsg)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoiceCanceling(ctx context.Context, cancelingId int) (*models.InvoiceCancel, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices_cancelings WHERE id = $1",
		cancelingId,
	)

	canceling := models.NewEmptyInvoiceCancel()
	err := canceling.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice canceling with id %v not found: %v", cancelingId, err)
		return nil, utils.CancelingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning canceling row: ", err)
		return nil, utils.InternalServerErr
	}

	entity, err := RetrieveEntity(ctx, canceling.Entity.Id)
	if err != nil {
		log.Println("Error linking invoice canceling to entity: ", err)
		return nil, utils.InternalServerErr
	}
	canceling.Entity = entity

	return canceling, nil
}

func UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	result, err := db.PG.Exec(
		ctx,
		"UPDATE invoices_cancelings SET req_status = $1, req_msg = $2 WHERE id = $3",
		canceling.ReqStatus, canceling.ReqMsg, canceling.Id,
	)
	if err != nil {
		log.Println("Error when running update invoice canceling query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Canceling with id %v not found when running update query", canceling.Id)
		return utils.CancelingNotFoundErr
	}

	return nil
}

