package services

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

type MetricsResultService struct{}

func NewMetricsResultService() *MetricsResultService {
	return &MetricsResultService{}
}

func (s *MetricsResultService) ListResults(ctx context.Context, metricsID int, userID int) ([]*models.MetricsResult, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM metrics_results WHERE metrics_id = $1 AND created_by = $2 ORDER BY type", metricsID, userID)
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
