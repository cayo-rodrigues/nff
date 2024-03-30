package storage

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

func ListMetricsResults(ctx context.Context, metricsID int, userID int) ([]*models.MetricsResult, error) {
	db := database.GetDB()

	rows, _ := db.PG.Query(ctx, "SELECT * FROM metrics_results WHERE metrics_id = $1 AND created_by = $2 ORDER BY issue_date", metricsID, userID)
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

func BulkCreateMetricsResults(ctx context.Context, results []*models.MetricsResult, resultType string, metricsID, userID, entityID int) error {
	rows := [][]interface{}{}
	for _, result := range results {
		result.MetricsID = metricsID
		result.CreatedBy = userID
		result.Type = resultType
		result.EntityID = entityID

		rows = append(rows, []interface{}{
			result.Type, result.MonthName, result.TotalIncome, result.TotalExpenses,
			result.AvgIncome, result.AvgExpenses, result.Diff, result.IsPositive,
			result.TotalRecords, result.PositiveRecords, result.NegativeRecords,
			result.MetricsID, result.CreatedBy,
			result.IssueDate, result.InvoiceNumber, result.EntityID,
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
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Printf("Error when running bulk insert metrics results query. Metrics id: %d. Result type: %s. Err: %v\n", metricsID, resultType, err)
		return utils.InternalServerErr
	}

	return nil
}
