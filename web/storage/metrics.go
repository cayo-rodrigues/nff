package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
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

	rows, _ := db.SQLite.QueryContext(ctx, query.String(), filters.Values()...)
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

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO metrics_history
			(start_date, end_date, entity_id, created_by)
			VALUES (?, ?, ?, ?)
		RETURNING id, req_status, req_msg, created_at, updated_at`,
		metrics.StartDate, metrics.EndDate, metrics.Entity.ID, metrics.CreatedBy,
	)
	err := row.Scan(&metrics.ID, &metrics.ReqStatus, &metrics.ReqMsg, &metrics.CreatedAt, &metrics.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert metrics query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func RetrieveMetrics(ctx context.Context, queryId int, userID int) (*models.Metrics, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`SELECT *
			FROM metrics_history
				JOIN entities ON entities.id = metrics_history.entity_id
		WHERE metrics_history.id = ? AND metrics_history.created_by = ?`,
		queryId, userID,
	)

	metrics := models.NewMetrics()
	err := Scan(row, metrics, metrics.Entity)
	if errors.Is(err, sql.ErrNoRows) {
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

	result, err := db.SQLite.ExecContext(
		ctx,
		`UPDATE metrics_history SET
			req_status = ?, req_msg = ?, updated_at = ?
		WHERE id = ? AND created_by = ?`,
		metrics.ReqStatus, metrics.ReqMsg, time.Now(),
		metrics.ID, metrics.CreatedBy,
	)
	if err != nil {
		log.Printf("Error when running update metrics query. Metrics id: %d. Err: %v\n", metrics.ID, err)
		return utils.InternalServerErr
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error when getting rows affected by update metrics query. Metrics id: %d. Err: %v\n", metrics.ID, err)
		return utils.InternalServerErr
	}
	if rowsAffected == 0 {
		log.Printf("Metrics with id %v not found when running update query", metrics.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}

