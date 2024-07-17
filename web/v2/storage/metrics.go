package storage

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

func ListMetrics(ctx context.Context, userID int, filters *models.Filters) ([]*models.Metrics, error) {
	var query strings.Builder

	query.WriteString(`
		SELECT *
			FROM metrics_history
			JOIN entities ON entities.id = metrics_history.entity_id
	`)

	query.WriteString(filters.String())

	db := database.GetDB()

	rows, _ := db.PG.Query(ctx, query.String(), filters.Values()...)
	defer rows.Close()

	metricsList := []*models.Metrics{}

	for rows.Next() {
		metrics := models.NewMetrics()
		err := Scan(rows, metrics, metrics.Entity)
		if err != nil {
			log.Println("Error scaning metrics query rows: ", err)
			return nil, utils.InternalServerErr
		}

		metricsList = append(metricsList, metrics)
	}

	return metricsList, nil
}

func CreateMetrics(ctx context.Context, metrics *models.Metrics) error {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO metrics_history
			(start_date, end_date, entity_id, created_by)
			VALUES ($1, $2, $3, $4)
		RETURNING id, req_status, req_msg`,
		metrics.StartDate, metrics.EndDate, metrics.Entity.ID, metrics.CreatedBy,
	)
	err := row.Scan(&metrics.ID, &metrics.ReqStatus, &metrics.ReqMsg)
	if err != nil {
		log.Println("Error when running insert metrics query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveMetrics(ctx context.Context, queryId int, userID int) (*models.Metrics, error) {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`SELECT *
			FROM metrics_history
				JOIN entities ON entities.id = metrics_history.entity_id
		WHERE metrics_history.id = $1 AND metrics_history.created_by = $2`,
		queryId, userID,
	)

	metrics := models.NewMetrics()
	err := Scan(row, metrics, metrics.Entity)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Metrics query with id %v not found: %v", queryId, err)
		return nil, utils.MetricsNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning metrics row: ", err)
		return nil, utils.InternalServerErr
	}

	results, err := ListMetricsResults(ctx, metrics.ID, userID)
	if err != nil {
		log.Println("Error linking metrics to results: ", err)
		return nil, utils.InternalServerErr
	}
	for _, result := range results {
		switch result.Type {
		case "total":
			metrics.Total = result
		case "month":
			metrics.Months = append(metrics.Months, result)
		case "record":
			metrics.Records = append(metrics.Records, result)
		}
	}

	return metrics, nil
}

func UpdateMetrics(ctx context.Context, metrics *models.Metrics) error {
	db := database.GetDB()

	result, err := db.PG.Exec(
		ctx,
		`UPDATE metrics_history SET
			req_status = $1, req_msg = $2, updated_at = $3
		WHERE id = $4 AND created_by = $5`,
		metrics.ReqStatus, metrics.ReqMsg, time.Now(),
		metrics.ID, metrics.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update metrics query. Metrics id: %d. Err: %v\n", metrics.ID, err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Metrics with id %v not found when running update query", metrics.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}

