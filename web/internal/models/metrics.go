package models

import (
	"log"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
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
	ID        int
	Entity    *Entity
	StartDate time.Time
	EndDate   time.Time
	Results   *MetricsResult
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
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

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	ilogicalDatesMsg := globals.IlogicalDatesMsg
	timeRangeTooLongMsg := globals.TimeRangeTooLongMsg

	startDateIsEmpty := q.StartDate.IsZero()
	endDateIsEmpty := q.EndDate.IsZero()
	startDateGreaterThanEndDate := q.StartDate.After(q.EndDate)
	timeRangeTooLong := int(q.EndDate.Sub(q.StartDate).Hours()/24) > 365
	hasEntity := q.Entity != nil

	fields := [5]*utils.Field{
		{ErrCondition: startDateIsEmpty, ErrField: &q.Errors.StartDate, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: endDateIsEmpty, ErrField: &q.Errors.EndDate, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: startDateGreaterThanEndDate, ErrField: &q.Errors.StartDate, ErrMsg: &ilogicalDatesMsg},
		{ErrCondition: timeRangeTooLong, ErrField: &q.Errors.StartDate, ErrMsg: &timeRangeTooLongMsg},
		{ErrCondition: !hasEntity, ErrField: &q.Errors.Entity, ErrMsg: &mandatoryFieldMsg},
	}

	var wg sync.WaitGroup
	for _, field := range fields {
		wg.Add(1)
		go utils.ValidateField(field, &isValid, &wg)
	}

	wg.Wait()

	return isValid
}

func (q *MetricsQuery) Scan(rows db.Scanner) error {
	return rows.Scan(
		&q.ID, &q.StartDate, &q.EndDate,
		&q.Results.TotalIncome, &q.Results.TotalExpenses, &q.Results.AvgIncome,
		&q.Results.AvgExpenses, &q.Results.Diff, &q.Results.IsPositive,
		&q.Results.TotalRecords, &q.Results.PositiveRecords, &q.Results.NegativeRecords,
		&q.Results.ReqStatus, &q.Results.ReqMsg, &q.Entity.ID,
		&q.CreatedBy, &q.CreatedAt, &q.UpdatedAt,
	)
}

func (q *MetricsQuery) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&q.ID, &q.StartDate, &q.EndDate,
		&q.Results.TotalIncome, &q.Results.TotalExpenses, &q.Results.AvgIncome,
		&q.Results.AvgExpenses, &q.Results.Diff, &q.Results.IsPositive,
		&q.Results.TotalRecords, &q.Results.PositiveRecords, &q.Results.NegativeRecords,
		&q.Results.ReqStatus, &q.Results.ReqMsg, &q.Entity.ID,
		&q.CreatedBy, &q.CreatedAt, &q.UpdatedAt,

		&q.Entity.ID, &q.Entity.Name, &q.Entity.UserType, &q.Entity.CpfCnpj, &q.Entity.Ie, &q.Entity.Email, &q.Entity.Password,
		&q.Entity.Address.PostalCode, &q.Entity.Address.Neighborhood, &q.Entity.Address.StreetType, &q.Entity.Address.StreetName, &q.Entity.Address.Number,
		&q.Entity.CreatedBy, &q.Entity.CreatedAt, &q.Entity.UpdatedAt,
	)
}
