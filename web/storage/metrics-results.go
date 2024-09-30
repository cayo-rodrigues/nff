package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListMetricsResults(ctx context.Context, metricsID int, userID int) ([]*models.MetricsResult, error) {
	db := database.GetDB()

	rows, _ := db.SQLite.QueryContext(ctx, "SELECT * FROM metrics_results WHERE metrics_id = ? AND created_by = ? ORDER BY issue_date, id", metricsID, userID)
	defer rows.Close()

	results := []*models.MetricsResult{}

	for rows.Next() {
		result := &models.MetricsResult{}
		err := Scan(rows, result)
		if err != nil {
			log.Println("Error scaning metrics result rows: ", err)
			return nil, utils.InternalServerErr
		}

		results = append(results, result)
	}

	return results, nil
}

func CreateMetricsResult(ctx context.Context, result *models.MetricsResult, resultType string, metricsID, userID, entityID int) error {
	result.Type = resultType
	result.MetricsID = metricsID
	result.CreatedBy = userID
	result.EntityID = entityID

	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO metrics_results (
				type, month_name, total_income, total_expenses,
				avg_income, avg_expenses, diff, is_positive,
				total_records, positive_records, negative_records,
				metrics_id, created_by,
				issue_date, invoice_id, entity_id
			)
			VALUES (
				?, ?, ?, ?,
				?, ?, ?, ?,
				?, ?, ?,
				?, ?,
				?, ?, ?
			)
		RETURNING *`,
		result.Type, result.MonthName, result.TotalIncome, result.TotalExpenses,
		result.AvgIncome, result.AvgExpenses, result.Diff, result.IsPositive,
		result.TotalRecords, result.PositiveRecords, result.NegativeRecords,
		result.MetricsID, result.CreatedBy,
		result.IssueDate, result.InvoiceNumber, result.EntityID,
	)
	err := Scan(row, result)
	if err != nil {
		log.Println("Error when running insert metrics results query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func BulkCreateMetricsResults(ctx context.Context, results []*models.MetricsResult, resultType string, metricsID, userID, entityID int) error {
	db := database.GetDB()

	return WithTransaction(ctx, db.SQLite, func(tx *sql.Tx) error {
		query := `
			INSERT INTO metrics_results (
				type, month_name, total_income, total_expenses,
				avg_income, avg_expenses, diff, is_positive,
				total_records, positive_records, negative_records,
				metrics_id, created_by, issue_date, invoice_id, entity_id, invoice_sender
			) VALUES `

		values := []interface{}{}
		valuePlaceholders := []string{}

		for _, result := range results {
			result.Type = resultType
			result.MetricsID = metricsID
			result.CreatedBy = userID
			result.EntityID = entityID

			valuePlaceholders = append(valuePlaceholders, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			values = append(values,
				result.Type, result.MonthName, result.TotalIncome, result.TotalExpenses,
				result.AvgIncome, result.AvgExpenses, result.Diff, result.IsPositive,
				result.TotalRecords, result.PositiveRecords, result.NegativeRecords,
				result.MetricsID, result.CreatedBy,
				result.IssueDate, result.InvoiceNumber, result.EntityID,
				result.InvoiceSender,
			)
		}

		query += strings.Join(valuePlaceholders, ", ")
		_, err := tx.ExecContext(ctx, query, values...)
		if err != nil {
			log.Printf("Error executing bulk insert metrics results query. Metrics ID: %d. Result type: %s. Err: %v\n", metricsID, resultType, err)
			return utils.InternalServerErr
		}

		return nil
	})
}

func UpdateMetricsResultRecord(ctx context.Context, result *models.MetricsResult) error {
	db := database.GetDB()

	cmd, err := db.SQLite.ExecContext(
		ctx,
		`UPDATE metrics_results
			SET invoice_pdf = ?
		WHERE id = ? AND created_by = ? AND type = 'record'`,
		result.InvoicePDF,
		result.ID, result.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update metrics query: ", err)
		return utils.InternalServerErr
	}
	rowsAffected, err := cmd.RowsAffected()
	if err != nil {
		log.Println("Error when getting rows affected by update metrics query: ", err)
		return utils.InternalServerErr
	}
	if rowsAffected == 0 {
		log.Printf("Metrics result with id %v not found when running update query", result.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}

func RetrieveMetricsResult(ctx context.Context, resultID int, userID int) (*models.MetricsResult, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(ctx, "SELECT * FROM metrics_results WHERE id = ? AND created_by = ?", resultID, userID)

	result := models.NewMetricsResult()
	err := Scan(row, result)
	if err != nil {
		log.Println("Error scaning metrics result row: ", err)
		return nil, utils.InternalServerErr
	}

	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Metrics result with id %v not found: %v", resultID, err)
		return nil, utils.MetricsResultNotFoundErr
	}

	return result, nil
}
