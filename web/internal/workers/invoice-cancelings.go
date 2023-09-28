package workers

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)


func ListInvoiceCancelings(ctx context.Context) (*[]models.InvoiceCancel, error) {
	dbpool := sql.GetDatabasePool()
	rows, _ := dbpool.Query(ctx, "SELECT * FROM invoices_cancelings ORDER BY id")
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
	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
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
