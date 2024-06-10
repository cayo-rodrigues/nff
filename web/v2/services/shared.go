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
