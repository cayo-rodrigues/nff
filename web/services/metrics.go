package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/interfaces"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

type MetricsService struct {
	resultsService interfaces.MetricsResultService
	filtersService interfaces.FiltersService
}

func NewMetricsService(resultsService interfaces.MetricsResultService, filtersService interfaces.FiltersService) *MetricsService {
	return &MetricsService{
		resultsService: resultsService,
		filtersService: filtersService,
	}
}

func (s *MetricsService) ListMetrics(ctx context.Context, userID int, filters map[string]string) ([]*models.MetricsQuery, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM metrics_history
			JOIN entities ON entities.id = metrics_history.entity_id
	`)

	params := s.filtersService.BuildQueryFilters(&query, filters, userID, "metrics_history")

	rows, _ := db.PG.Query(ctx, query.String(), params...)
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

	results, err := s.resultsService.ListResults(ctx, query.ID, userID)
	if err != nil {
		log.Println("Error linking metrics to results: ", err)
		return nil, utils.InternalServerErr
	}
	for _, result := range results {
		switch result.Type {
		case "month":
			query.Months = append(query.Months, result)
		}
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
		query.ReqStatus, query.ReqMsg, query.TotalIncome, query.TotalExpenses,
		query.AvgIncome, query.AvgExpenses, query.Diff, query.IsPositive,
		query.TotalRecords, query.PositiveRecords, query.NegativeRecords, time.Now(),
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
