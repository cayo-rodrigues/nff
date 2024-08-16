package storage

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

func ListMetricsResults(ctx context.Context, metricsID int, userID int) ([]*models.MetricsResult, error) {
	db := database.GetDB()

	rows, _ := db.PG.Query(ctx, "SELECT * FROM metrics_results WHERE metrics_id = $1 AND created_by = $2 ORDER BY issue_date, id", metricsID, userID)
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

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO metrics_results (
				type, month_name, total_income, total_expenses,
				avg_income, avg_expenses, diff, is_positive,
				total_records, positive_records, negative_records,
				metrics_id, created_by,
				issue_date, invoice_id, entity_id
			)
			VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11,
				$12, $13,
				$14, $15, $16
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
		log.Println("Error when running insert metrics query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func BulkCreateMetricsResults(ctx context.Context, results []*models.MetricsResult, resultType string, metricsID, userID, entityID int) error {
	rows := [][]interface{}{}
	for _, result := range results {
		result.Type = resultType
		result.MetricsID = metricsID
		result.CreatedBy = userID
		result.EntityID = entityID

		rows = append(rows, []interface{}{
			result.Type, result.MonthName, result.TotalIncome, result.TotalExpenses,
			result.AvgIncome, result.AvgExpenses, result.Diff, result.IsPositive,
			result.TotalRecords, result.PositiveRecords, result.NegativeRecords,
			result.MetricsID, result.CreatedBy,
			result.IssueDate, result.InvoiceNumber, result.EntityID,
			result.InvoiceSender,
		})
	}

	db := database.GetDB()

	_, err := db.PG.CopyFrom(
		ctx,
		pgx.Identifier{"metrics_results"},
		[]string{
			"type", "month_name", "total_income", "total_expenses",
			"avg_income", "avg_expenses", "diff", "is_positive",
			"total_records", "positive_records", "negative_records",
			"metrics_id", "created_by",
			"issue_date", "invoice_id", "entity_id",
			"invoice_sender",
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("Error when running bulk insert metrics results query. Metrics id: %d. Result type: %s. Err: %v\n", metricsID, resultType, err)
		return utils.InternalServerErr
	}

	return nil
}

func UpdateMetricsResultRecord(ctx context.Context, result *models.MetricsResult) error {
	db := database.GetDB()

	cmd, err := db.PG.Exec(
		ctx,
		`UPDATE metrics_results
			SET invoice_pdf = $1
		WHERE id = $2 AND created_by = $3 AND type = 'record'`,
		result.InvoicePDF,
		result.ID, result.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update metrics query: ", err)
		return utils.InternalServerErr
	}
	if cmd.RowsAffected() == 0 {
		log.Printf("Metrics result with id %v not found when running update query", result.ID)
		return utils.MetricsNotFoundErr
	}

	return nil
}

func RetrieveMetricsResult(ctx context.Context, resultID int, userID int) (*models.MetricsResult, error) {
	db := database.GetDB()

	row := db.PG.QueryRow(ctx, "SELECT * FROM metrics_results WHERE id = $1 AND created_by = $2", resultID, userID)

	result := models.NewMetricsResult()
	err := Scan(row, result)
	if err != nil {
		log.Println("Error scaning metrics result row: ", err)
		return nil, utils.InternalServerErr
	}

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Metrics result with id %v not found: %v", resultID, err)
		return nil, utils.MetricsResultNotFoundErr
	}

	return result, nil
}
