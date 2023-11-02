package models

import (
	"log"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type MetricsFormSelectFields struct {
	Entities []*Entity
}

type MetricsFormErrors struct {
	Entity    string
	StartDate string
	EndDate   string
}

type MetricsQuery struct {
	Id        int
	Entity    *Entity
	StartDate time.Time
	EndDate   time.Time
	Results   *MetricsResult
	Errors    *MetricsFormErrors
}

type MetricsResult struct {
	TotalIncome     string `json:"total_income"`
	TotalExpenses   string `json:"total_expenses"`
	AvgIncome       string `json:"average_income"`
	AvgExpenses     string `json:"average_expenses"`
	Diff            string `json:"diff"`
	IsPositive      bool   `json:"is_positive"`
	TotalRecords    int    `json:"total_records"`
	PositiveRecords int    `json:"positive_records"`
	NegativeRecords int    `json:"negative_records"`
	ReqStatus       string `json:"status"`
	ReqMsg          string `json:"msg"`
}

func NewEmptyMetricsQuery() *MetricsQuery {
	return &MetricsQuery{
		Entity:  NewEmptyEntity(),
		Results: &MetricsResult{},
		Errors:  &MetricsFormErrors{},
	}

}

func NewMetricsQueryFromForm(c *fiber.Ctx) *MetricsQuery {
	startDate, err := utils.ParseDate(c.FormValue("start_date"))
	if err != nil {
		log.Println("Error converting input start date string to time.Time:", err)
	}
	endDate, err := utils.ParseDate(c.FormValue("end_date"))
	if err != nil {
		log.Println("Error converting input end date string to time.Time:", err)
	}
	return &MetricsQuery{
		StartDate: startDate,
		EndDate:   endDate,
		Results:   &MetricsResult{},
		Errors:    &MetricsFormErrors{},
	}
}

func (q *MetricsQuery) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	ilogicalDatesMsg := "Data inicial deve ser anterior à final"
	validationsCount := 4

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	startDateIsEmpty := q.StartDate.IsZero()
	endDateIsEmpty := q.EndDate.IsZero()

	go utils.ValidateField(startDateIsEmpty, &q.Errors.StartDate, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(endDateIsEmpty, &q.Errors.EndDate, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(!startDateIsEmpty && !endDateIsEmpty && q.StartDate.After(q.EndDate), &q.Errors.StartDate, &ilogicalDatesMsg, ch, &wg)
	go utils.ValidateField(q.Entity == nil, &q.Errors.Entity, &mandatoryFieldMsg, ch, &wg)

	wg.Wait()
	close(ch)

	for i := 0; i < validationsCount; i++ {
		if validationPassed := <-ch; !validationPassed {
			isValid = false
			break
		}
	}

	return isValid
}

func (q *MetricsQuery) Scan(rows db.Scanner) error {
	return rows.Scan(
		&q.Id, &q.StartDate, &q.EndDate,
		&q.Results.TotalIncome, &q.Results.TotalExpenses, &q.Results.AvgIncome,
		&q.Results.AvgExpenses, &q.Results.Diff, &q.Results.IsPositive,
		&q.Results.TotalRecords, &q.Results.PositiveRecords, &q.Results.NegativeRecords,
		&q.Results.ReqStatus, &q.Results.ReqMsg, &q.Entity.Id,
	)
}
