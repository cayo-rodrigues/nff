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

type MetricsService struct {
	entityService interfaces.EntityService
}

func NewMetricsService(entityService interfaces.EntityService) *MetricsService {
	return &MetricsService{
		entityService: entityService,
	}
}

func (s *MetricsService) ListMetrics(ctx context.Context, userID int) ([]*models.MetricsQuery, error) {
	rows, _ := db.PG.Query(
		ctx,
		`SELECT *
			FROM metrics_history
			JOIN entities ON entities.id = metrics_history.entity_id
		WHERE metrics_history.created_by = $1 ORDER BY metrics_history.id DESC`,
		userID,
	)
	defer rows.Close()

	queriesHistory := []*models.MetricsQuery{}

	for rows.Next() {
		metricsQuery := models.NewEmptyMetricsQuery()
		err := metricsQuery.FullScan(rows)
		if err != nil {
			log.Println("Error scaning metrics query rows: ", err)
			return nil, utils.InternalServerErr
		}

		queriesHistory = append(queriesHistory, metricsQuery)
	}

	return queriesHistory, nil
}

func (s *MetricsService) CreateMetrics(ctx context.Context, query *models.MetricsQuery) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO metrics_history
			(start_date, end_date, entity_id, created_by)
			VALUES ($1, $2, $3, $4)
		RETURNING *`,
		query.StartDate, query.EndDate, query.Entity.ID, query.CreatedBy,
	)
	err := query.Scan(row)
	if err != nil {
		log.Println("Error when running insert metrics query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *MetricsService) RetrieveMetrics(ctx context.Context, queryId int, userID int) (*models.MetricsQuery, error) {
	row := db.PG.QueryRow(
		ctx,
		`SELECT *
			FROM metrics_history
				JOIN entities ON entities.id = metrics_history.entity_id
		WHERE metrics_history.id = $1 AND metrics_history.created_by = $2`,
		queryId, userID,
	)

	query := models.NewEmptyMetricsQuery()
	err := query.FullScan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Metrics query with id %v not found: %v", queryId, err)
		return nil, utils.MetricsNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning metrics row: ", err)
		return nil, utils.InternalServerErr
	}

	return query, nil
}

func (s *MetricsService) UpdateMetrics(ctx context.Context, query *models.MetricsQuery) error {
	result, err := db.PG.Exec(
		ctx,
		`UPDATE metrics_history SET
			req_status = $1, req_msg = $2, total_income = $3, total_expenses = $4,
			avg_income = $5, avg_expenses = $6, diff = $7, is_positive = $8,
			total_records = $9, positive_records = $10, negative_records = $11, updated_at = $12
		WHERE id = $13 AND created_by = $14`,
		query.Results.ReqStatus, query.Results.ReqMsg, query.Results.TotalIncome, query.Results.TotalExpenses,
		query.Results.AvgIncome, query.Results.AvgExpenses, query.Results.Diff, query.Results.IsPositive,
		query.Results.TotalRecords, query.Results.PositiveRecords, query.Results.NegativeRecords, time.Now(),
		query.ID, query.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update metrics query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Metrics with id %v not found when running update query", query.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}
