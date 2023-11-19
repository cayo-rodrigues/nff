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

type MetricsService struct {
	entityService interfaces.EntityService
}

func NewMetricsService(entityService interfaces.EntityService) *MetricsService {
	return &MetricsService{
		entityService: entityService,
	}
}

func (s *MetricsService) ListMetrics(ctx context.Context) ([]*models.MetricsQuery, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM metrics_history ORDER BY id DESC")
	defer rows.Close()

	queriesHistory := []*models.MetricsQuery{}

	for rows.Next() {
		metricsQuery := models.NewEmptyMetricsQuery()
		err := metricsQuery.Scan(rows)
		if err != nil {
			log.Println("Error scaning metrics query rows: ", err)
			return nil, utils.InternalServerErr
		}

		entity, err := s.entityService.RetrieveEntity(ctx, metricsQuery.Entity.ID)
		if err != nil {
			log.Println("Error linking metrics query to entity: ", err)
			return nil, utils.InternalServerErr
		}
		metricsQuery.Entity = entity

		queriesHistory = append(queriesHistory, metricsQuery)
	}

	return queriesHistory, nil
}

func (s *MetricsService) CreateMetrics(ctx context.Context, query *models.MetricsQuery) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO metrics_history
			(start_date, end_date, entity_id)
			VALUES ($1, $2, $3)
		RETURNING *`,
		query.StartDate, query.EndDate, query.Entity.ID,
	)
	err := query.Scan(row)
	if err != nil {
		log.Println("Error when running insert metrics query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *MetricsService) RetrieveMetrics(ctx context.Context, queryId int) (*models.MetricsQuery, error) {
	// TODO maybe JOIN would be more efficient than two separated queries
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM metrics_history WHERE id = $1",
		queryId,
	)

	query := models.NewEmptyMetricsQuery()
	err := query.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Metrics query with id %v not found: %v", queryId, err)
		return nil, utils.MetricsNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning metrics row: ", err)
		return nil, utils.InternalServerErr
	}

	entity, err := s.entityService.RetrieveEntity(ctx, query.Entity.ID)
	if err != nil {
		log.Println("Error linking metrics query to entity: ", err)
		return nil, utils.InternalServerErr
	}
	query.Entity = entity

	return query, nil
}

func (s *MetricsService) UpdateMetrics(ctx context.Context, query *models.MetricsQuery) error {
	result, err := db.PG.Exec(
		ctx,
		`UPDATE metrics_history SET
			req_status = $1, req_msg = $2, total_income = $3, total_expenses = $4,
			avg_income = $5, avg_expenses = $6, diff = $7, is_positive = $8,
			total_records = $9, positive_records = $10, negative_records = $11
		WHERE id = $12`,
		query.Results.ReqStatus, query.Results.ReqMsg, query.Results.TotalIncome, query.Results.TotalExpenses,
		query.Results.AvgIncome, query.Results.AvgExpenses, query.Results.Diff, query.Results.IsPositive,
		query.Results.TotalRecords, query.Results.PositiveRecords, query.Results.NegativeRecords,
		query.ID,
	)
	if err != nil {
		log.Println("REQ STATUS =", query.Results.ReqStatus, len(query.Results.ReqStatus))
		log.Println("Error when running update metrics query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Metrics with id %v not found when running update query", query.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}
