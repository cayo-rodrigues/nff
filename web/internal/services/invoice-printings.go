package services

import (
	"context"
	"errors"
	"log"

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

func (s *PrintingService) ListInvoicePrintings(ctx context.Context) ([]*models.InvoicePrint, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM invoices_printings ORDER BY id DESC")
	defer rows.Close()

	printings := []*models.InvoicePrint{}

	for rows.Next() {
		printing := models.NewEmptyInvoicePrint()
		err := printing.Scan(rows)
		if err != nil {
			log.Println("Error scaning invoice printing rows: ", err)
			return nil, utils.InternalServerErr
		}

		entity, err := s.entityService.RetrieveEntity(ctx, printing.Entity.ID)
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
			(invoice_id, invoice_id_type, entity_id)
			VALUES ($1, $2, $3)
		RETURNING id, req_status, req_msg`,
		printing.InvoiceId, printing.InvoiceIdType, printing.Entity.ID,
	)
	err := row.Scan(&printing.ID, &printing.ReqStatus, &printing.ReqMsg)
	if err != nil {
		log.Println("Error when running insert printing query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *PrintingService) RetrieveInvoicePrinting(ctx context.Context, printingId int) (*models.InvoicePrint, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM invoices_printings WHERE id = $1",
		printingId,
	)

	printing := models.NewEmptyInvoicePrint()
	err := printing.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Invoice printing with id %v not found: %v", printingId, err)
		return nil, utils.PrintingNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning printing row: ", err)
		return nil, utils.InternalServerErr
	}

	entity, err := s.entityService.RetrieveEntity(ctx, printing.Entity.ID)
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
		"UPDATE invoices_printings SET req_status = $1, req_msg = $2, invoice_pdf = $3 WHERE id = $4",
		printing.ReqStatus, printing.ReqMsg, printing.InvoicePDF, printing.ID,
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
