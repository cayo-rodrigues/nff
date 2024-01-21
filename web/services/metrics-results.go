package services

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

type MetricsResultService struct{}

func NewMetricsResultService() *MetricsResultService {
	return &MetricsResultService{}
}

func (s *MetricsResultService) ListResults(ctx context.Context, metricsID int, userID int) ([]*models.MetricsResult, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM metrics_results WHERE metrics_id = $1 AND created_by = $2 ORDER BY issue_date", metricsID, userID)
	defer rows.Close()

	results := []*models.MetricsResult{}

	for rows.Next() {
		result := &models.MetricsResult{}
		err := result.Scan(rows)
		if err != nil {
			log.Println("Error scaning metrics result rows: ", err)
			return nil, utils.InternalServerErr
		}

		results = append(results, result)
	}

	return results, nil
}

func (s *MetricsResultService) BulkCreateResults(ctx context.Context, results []*models.MetricsResult, resultType string, metricsID int, userID int) error {
	rows := [][]interface{}{}
	for _, result := range results {
		result.MetricsID = metricsID
		result.CreatedBy = userID
		result.Type = resultType

		rows = append(rows, []interface{}{
			result.Type, result.MonthName, result.TotalIncome, result.TotalExpenses,
			result.AvgIncome, result.AvgExpenses, result.Diff, result.IsPositive,
			result.TotalRecords, result.PositiveRecords, result.NegativeRecords,
			result.MetricsID, result.CreatedBy,
			result.IssueDate,
		})
	}
	_, err := db.PG.CopyFrom(
		ctx,
		pgx.Identifier{"metrics_results"},
		[]string{
			"type", "month_name", "total_income", "total_expenses",
			"avg_income", "avg_expenses", "diff", "is_positive",
			"total_records", "positive_records", "negative_records",
			"metrics_id", "created_by",
			"issue_date",
		},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Println("Error when running bulk insert metrics results query: ", err)
		return utils.InternalServerErr
	}

	return nil
}
