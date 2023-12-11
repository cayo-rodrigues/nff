package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type FiltersService struct{}

func NewFiltersService() *FiltersService {
	return &FiltersService{}
}

func (s *FiltersService) BuildQueryFilters(query *strings.Builder, filters map[string]string, userID int, table string) []interface{} {
	query.WriteString("WHERE ")
	query.WriteString(table)
	query.WriteString(".created_by = $1")

	params := []interface{}{userID}
	paramCounter := 2

	now := time.Now()
	fromDate, ok := filters["from_date"]
	if !ok || fromDate == "" {
		fromDate = utils.FormatedNDaysBefore(now, 30)
	}
	toDate, ok := filters["to_date"]
	if !ok || toDate == "" {
		toDate = utils.FormatDate(now)
	}

	query.WriteString(" AND CAST(")
	query.WriteString(table)
	query.WriteString(".created_at AS DATE) BETWEEN $")
	query.WriteString(strconv.Itoa(paramCounter))
	params = append(params, fromDate)
	paramCounter++

	query.WriteString(" AND $")
	query.WriteString(strconv.Itoa(paramCounter))
	params = append(params, toDate)
	paramCounter++

	if entityID, ok := filters["entity_id"]; ok && entityID != "" {
		counter := strconv.Itoa(paramCounter)

		if table == "invoices" {
			query.WriteString(" AND invoices.sender_id = $")
			query.WriteString(counter)
			query.WriteString(" OR invoices.recipient_id = $")
			query.WriteString(counter)
		} else {
			query.WriteString(" AND ")
			query.WriteString(table)
			query.WriteString(".entity_id = $")
			query.WriteString(counter)
		}

		params = append(params, entityID)
		paramCounter++
	}

	if senderID, ok := filters["sender_id"]; ok && senderID != "" {
		query.WriteString(" AND ")
		query.WriteString(table)
		query.WriteString(".sender_id = $")
		query.WriteString(strconv.Itoa(paramCounter))
		params = append(params, senderID)
		paramCounter++
	}

	if recipientID, ok := filters["recipient_id"]; ok && recipientID != "" {
		query.WriteString(" AND ")
		query.WriteString(table)
		query.WriteString(".recipient_id = $")
		query.WriteString(strconv.Itoa(paramCounter))
		params = append(params, recipientID)
		paramCounter++
	}

	query.WriteString(" ORDER BY ")
	query.WriteString(table)
	query.WriteString(".created_at DESC")

	return params
}
