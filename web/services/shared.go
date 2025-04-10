package services

import (
	"strings"
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

	// TODO
	// verificar se o período selecionado é válido

	f.And().Date(colName).Between(fromDate, toDate)
}

func handleEntityFilters(colName string, query map[string]string, f *models.Filters) {
	entityID, ok := query["entity_filter"]

	if !ok || entityID == "" {
		return
	}

	if strings.HasPrefix(colName, "invoices.") {
		f.And(colName + "sender_id = ").Placeholder(entityID).Or(colName + "recipient_id = ").Placeholder(entityID)
		return
	}

	f.And(colName + " = ").Placeholder(entityID)
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

