package workers

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

func ListInvoiceCancelings(ctx context.Context) (*[]models.InvoiceCancel, error) {
	rows, _ := sql.DB.Query(ctx, "SELECT * FROM invoices_cancelings ORDER BY id DESC")
	defer rows.Close()

	cancelings := []models.InvoiceCancel{}

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

		cancelings = append(cancelings, *canceling)
	}

	return &cancelings, nil

}

func CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	row := sql.DB.QueryRow(
		ctx,
		`INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id)
			VALUES ($1, $2, $3, $4)
		RETURNING id`,
		canceling.Number, canceling.Year, canceling.Justification, canceling.Entity.Id,
	)
	err := row.Scan(&canceling.Id)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveInvoiceCanceling(ctx context.Context, cancelingId int) (*models.InvoiceCancel, error) {
	row := sql.DB.QueryRow(
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
