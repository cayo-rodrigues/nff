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
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices_cancelings WHERE created_by = $1 ORDER BY id DESC", userID)
	defer rows.Close()

	cancelings := []*models.InvoiceCancel{}

	for rows.Next() {
		canceling := models.NewEmptyInvoiceCancel()
		err := canceling.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice canceling rows: ", err)
			return nil, utils.InternalServerErr
		}

		entity, err := s.entityService.RetrieveEntity(ctx, canceling.Entity.ID, canceling.CreatedBy)
		if err != nil {
			log.Println("Error linking invoice canceling to entity: ", err)
			return nil, utils.InternalServerErr
		}
		canceling.Entity = entity

		cancelings = append(cancelings, canceling)
	}

	return cancelings, nil

}

func (s *CancelingService) CreateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices_cancelings
			(invoice_number, year, justification, entity_id, created_by)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id, req_status, req_msg`,
		canceling.Number, canceling.Year, canceling.Justification, canceling.Entity.ID, canceling.CreatedBy,
	)
	err := row.Scan(&canceling.ID, &canceling.ReqStatus, &canceling.ReqMsg)
	if err != nil {
		log.Println("Error when running insert canceling query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *CancelingService) RetrieveInvoiceCanceling(ctx context.Context, cancelingID int, userID int) (*models.InvoiceCancel, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices_cancelings WHERE id = $1 AND created_by = $2",
		cancelingID, userID,
	)

	canceling := models.NewEmptyInvoiceCancel()
	err := canceling.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice canceling with id %v not found: %v", cancelingID, err)
		return nil, utils.CancelingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning canceling row: ", err)
		return nil, utils.InternalServerErr
	}

	entity, err := s.entityService.RetrieveEntity(ctx, canceling.Entity.ID, canceling.CreatedBy)
	if err != nil {
		log.Println("Error linking invoice canceling to entity: ", err)
		return nil, utils.InternalServerErr
	}
	canceling.Entity = entity

	return canceling, nil
}

func (s *CancelingService) UpdateInvoiceCanceling(ctx context.Context, canceling *models.InvoiceCancel) error {
	result, err := db.PG.Exec(
		ctx,
		"UPDATE invoices_cancelings SET req_status = $1, req_msg = $2, updated_at = $3 WHERE id = $4 AND created_by = $5",
		canceling.ReqStatus, canceling.ReqMsg, time.Now(), canceling.ID, canceling.CreatedBy,
	)
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
