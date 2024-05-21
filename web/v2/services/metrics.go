package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func ListMetrics(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Metrics, error) {
	f := models.NewFilters().Where("metrics_history.created_by = ").Placeholder(userID)

	if filters == nil {
		filters = make([]map[string]string, 1)
	}

	for i, filter := range filters {
		if i == 0 {
			handleDateFilters("metrics_history.created_at", filter, f)
		}
	}

	f.OrderBy("metrics_history.created_at").Desc()

	return storage.ListMetrics(ctx, userID, f)
}

func CreateMetrics(ctx context.Context, m *models.Metrics, userID int) error {
	m.CreatedBy = userID
	return storage.CreateMetrics(ctx, m)
}

func RetrieveMetrics(ctx context.Context, printingID int, userID int) (*models.Metrics, error) {
	return storage.RetrieveMetrics(ctx, printingID, userID)
}

func GroupListByDate(list []*models.Metrics) []map[string][]*models.Metrics {
	groupedList := []map[string][]*models.Metrics{}

	if len(list) == 0 {
		return groupedList
	}

	lastSeenDate := utils.FormatDateAsBR(list[0].CreatedAt)
	dailyList := map[string][]*models.Metrics{}

	for _, item := range list {
		key := utils.FormatDateAsBR(item.CreatedAt)
		if key != lastSeenDate {
			lastSeenDate = key
			groupedList = append(groupedList, dailyList)
			dailyList = map[string][]*models.Metrics{}
		}
		dailyList[key] = append(dailyList[key], item)
	}

	groupedList = append(groupedList, dailyList)

	return groupedList
}
