package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type CancelingService struct {
	entityService interfaces.EntityService
}

func NewCancelingService(entityService interfaces.EntityService) *CancelingService {
	return &CancelingService{
		entityService: entityService,
	}
}

func (s *CancelingService) ListInvoiceCancelings(ctx context.Context, userID int) ([]*models.InvoiceCancel, error) {
	query := `
		SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
		WHERE invoices_cancelings.created_by = $1
		ORDER BY invoices_cancelings.id DESC
	`
	rows, _ := db.PG.Query(ctx, query, userID)
	defer rows.Close()

	cancelings := []*models.InvoiceCancel{}

	for rows.Next() {
		canceling := models.NewEmptyInvoiceCancel()
		err := canceling.FullScan(rows)
		if err != nil {
			log.Println("Error scaning invoice canceling rows: ", err)
			return nil, utils.InternalServerErr
		}

		cancelings = append(cancelings, canceling)
	}

	return cancelings, nil

}

func (s *CancelingService) CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	query := `
		INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id, created_by)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id, req_status, req_msg
	`
	queryArgs := [5]any{canceling.Number, canceling.Year, canceling.Justification, canceling.Entity.ID, canceling.CreatedBy}

	row := db.PG.QueryRow(ctx, query, queryArgs[:]...)
	err := row.Scan(&canceling.ID, &canceling.ReqStatus, &canceling.ReqMsg)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *CancelingService) RetrieveInvoiceCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	query := `
		SELECT *
			FROM invoices_cancelings
				JOIN entities ON entities.id = invoices_cancelings.entity_id
		WHERE invoices_cancelings.id = $1 AND invoices_cancelings.created_by = $2
	`
	row := db.PG.QueryRow(ctx, query, cancelingID, userID)

	canceling := models.NewEmptyInvoiceCancel()
	err := canceling.FullScan(row)
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

func (s *CancelingService) UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	query := `
		UPDATE invoices_cancelings
			SET req_status = $1, req_msg = $2, updated_at = $3
		WHERE id = $4 AND created_by = $5
	`
	queryArgs := [5]any{canceling.ReqStatus, canceling.ReqMsg, time.Now(), canceling.ID, canceling.CreatedBy}

	result, err := db.PG.Exec(ctx, query, queryArgs[:]...)
	if err != nil {
		log.Println("Error when running update invoice canceling query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Canceling with id %v not found when running update query", canceling.ID)
		return utils.CancelingNotFoundErr
	}

	return nil
}
