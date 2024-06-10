package services

import (
	"time"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func handleDateFilters(colName string, query map[string]string, f *models.Filters) {
	fromDate, fromDateOk := query["from_date"]
	toDate, toDateOk := query["to_date"]

	if !fromDateOk && !toDateOk {
		now := time.Now()
		fromDate = utils.FormatedNDaysBefore(now, utils.DefaultFiltersDaysRange)
		toDate = utils.FormatDate(now)
	}

	// verificar se o período selecionado é válido

	f.And().AsDate(colName).Between(fromDate, toDate)
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
