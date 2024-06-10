package services

import (
	"context"
	"time"

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

type CreatedAtGetter interface {
	GetCreatedAt() time.Time
}

func GroupListByDate[T CreatedAtGetter](list []T) []map[string][]T {
	groupedList := []map[string][]T{}

	if len(list) == 0 {
		return groupedList
	}

	lastSeenDate := utils.FormatDateAsBR(list[0].GetCreatedAt())
	dailyList := map[string][]T{}

	for _, item := range list {
		key := utils.FormatDateAsBR(item.GetCreatedAt())
		if key != lastSeenDate {
			lastSeenDate = key
			groupedList = append(groupedList, dailyList)
			dailyList = map[string][]T{}
		}
		dailyList[key] = append(dailyList[key], item)
	}

	groupedList = append(groupedList, dailyList)

	return groupedList
}
