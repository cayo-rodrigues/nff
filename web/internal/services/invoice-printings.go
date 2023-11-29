package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

type PrintingService struct {
	entityService interfaces.EntityService
}

func NewPrintingService(entityService interfaces.EntityService) *PrintingService {
	return &PrintingService{
		entityService: entityService,
	}
}

func (s *PrintingService) ListInvoicePrintings(ctx context.Context, userID int) ([]*models.InvoicePrint, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices_printings WHERE created_by = $1 ORDER BY id DESC", userID)
	defer rows.Close()

	printings := []*models.InvoicePrint{}

	for rows.Next() {
		printing := models.NewEmptyInvoicePrint()
		err := printing.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice printing rows: ", err)
			return nil, utils.InternalServerErr
		}

		entity, err := s.entityService.RetrieveEntity(ctx, printing.Entity.ID, printing.CreatedBy)
		if err != nil {
			log.Println("Error linking invoice printing to entity: ", err)
			return nil, utils.InternalServerErr
		}
		printing.Entity = entity

		printings = append(printings, printing)
	}

	return printings, nil
}

func (s *PrintingService) CreateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO invoices_printings
			(invoice_id, invoice_id_type, custom_file_name, entity_id, created_by)
			VALUES ($1, $2, $3, $4, $5)
		RETURNING id, req_status, req_msg`,
		printing.InvoiceId, printing.InvoiceIdType, printing.CustomFileName, printing.Entity.ID, printing.CreatedBy,
	)
	err := row.Scan(&printing.ID, &printing.ReqStatus, &printing.ReqMsg)
	if err != nil {
		log.Println("Error when running insert printing query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *PrintingService) RetrieveInvoicePrinting(ctx context.Context, printingID int, userID int) (*models.InvoicePrint, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices_printings WHERE id = $1 AND created_by = $2",
		printingID, userID,
	)

	printing := models.NewEmptyInvoicePrint()
	err := printing.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice printing with id %v not found: %v", printingID, err)
		return nil, utils.PrintingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning printing row: ", err)
		return nil, utils.InternalServerErr
	}

	entity, err := s.entityService.RetrieveEntity(ctx, printing.Entity.ID, printing.CreatedBy)
	if err != nil {
		log.Println("Error linking invoice printing to entity: ", err)
		return nil, utils.InternalServerErr
	}
	printing.Entity = entity

	return printing, nil
}

func (s *PrintingService) UpdateInvoicePrinting(ctx context.Context, printing *models.InvoicePrint) error {
	result, err := db.PG.Exec(
		ctx,
		"UPDATE invoices_printings SET req_status = $1, req_msg = $2, invoice_pdf = $3, custom_file_name = $4, updated_at = $5 WHERE id = $6 AND created_by = $7",
		printing.ReqStatus, printing.ReqMsg, printing.InvoicePDF, printing.CustomFileName, time.Now(), printing.ID, printing.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update invoice printing query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Printing with id %v not found when running update query", printing.ID)
		return utils.PrintingNotFoundErr
	}

	return nil
}
