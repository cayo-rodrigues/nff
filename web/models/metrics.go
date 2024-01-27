package models

import (
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type MetricsSelectFields struct {
	Entities []*Entity
}

func NewMetricsSelectFields() *MetricsSelectFields {
	return &MetricsSelectFields{
		Entities: []*Entity{},
	}
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
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
	Errors    *MetricsFormErrors
	*MetricsResult
}

type MetricsResult struct {
	ID              int
	Type            string           `json:"type"`
	TotalIncome     string           `json:"total_income"`
	TotalExpenses   string           `json:"total_expenses"`
	AvgIncome       string           `json:"average_income"`
	AvgExpenses     string           `json:"average_expenses"`
	Diff            string           `json:"diff"`
	IsPositive      bool             `json:"is_positive"`
	TotalRecords    int              `json:"total_records"`
	PositiveRecords int              `json:"positive_records"`
	NegativeRecords int              `json:"negative_records"`
	ReqStatus       string           `json:"status"`
	ReqMsg          string           `json:"msg"`
	MonthName       string           `json:"month_name"`
	InvoiceNumber   string           `json:"invoice_id"`
	InvoicePDF      string           `json:"invoice_pdf"`
	IssueDate       time.Time        `json:"issue_date"`
	Months          []*MetricsResult `json:"months"`
	Records         []*MetricsResult `json:"records"`
	MetricsID       int
	EntityID        int
	CreatedBy       int
	CreatedAt       time.Time
}

func NewEmptyMetricsQuery() *MetricsQuery {
	return &MetricsQuery{
		Entity:        NewEmptyEntity(),
		MetricsResult: &MetricsResult{},
		Errors:        &MetricsFormErrors{},
	}

}

func NewMetricsQueryFromForm(c *fiber.Ctx) *MetricsQuery {
	startDate, err := utils.ParseDate(strings.TrimSpace(c.FormValue("start_date")))
	if err != nil {
		log.Println("Error converting input start date string to time.Time:", err)
	}
	endDate, err := utils.ParseDate(strings.TrimSpace(c.FormValue("end_date")))
	if err != nil {
		log.Println("Error converting input end date string to time.Time:", err)
	}

	return &MetricsQuery{
		StartDate:     startDate,
		EndDate:       endDate,
		MetricsResult: &MetricsResult{},
		Errors:        &MetricsFormErrors{},
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

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	return isValid
}

func (q *MetricsQuery) Scan(rows db.Scanner) error {
	return rows.Scan(
		&q.ID, &q.StartDate, &q.EndDate,
		&q.TotalIncome, &q.TotalExpenses, &q.AvgIncome,
		&q.AvgExpenses, &q.Diff, &q.IsPositive,
		&q.TotalRecords, &q.PositiveRecords, &q.NegativeRecords,
		&q.ReqStatus, &q.ReqMsg, &q.Entity.ID,
		&q.CreatedBy, &q.CreatedAt, &q.UpdatedAt,
	)
}

func (q *MetricsQuery) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&q.ID, &q.StartDate, &q.EndDate,
		&q.TotalIncome, &q.TotalExpenses, &q.AvgIncome,
		&q.AvgExpenses, &q.Diff, &q.IsPositive,
		&q.TotalRecords, &q.PositiveRecords, &q.NegativeRecords,
		&q.ReqStatus, &q.ReqMsg, &q.Entity.ID,
		&q.CreatedBy, &q.CreatedAt, &q.UpdatedAt,

		&q.Entity.ID, &q.Entity.Name, &q.Entity.UserType, &q.Entity.CpfCnpj, &q.Entity.Ie, &q.Entity.Email, &q.Entity.Password,
		&q.Entity.Address.PostalCode, &q.Entity.Address.Neighborhood, &q.Entity.Address.StreetType, &q.Entity.Address.StreetName, &q.Entity.Address.Number,
		&q.Entity.CreatedBy, &q.Entity.CreatedAt, &q.Entity.UpdatedAt,
	)
}

func (r *MetricsResult) Scan(rows db.Scanner) error {
	var issueDate any

	err := rows.Scan(
		&r.ID, &r.Type, &r.MonthName,
		&r.TotalIncome, &r.TotalExpenses, &r.AvgIncome,
		&r.AvgExpenses, &r.Diff, &r.IsPositive,
		&r.TotalRecords, &r.PositiveRecords, &r.NegativeRecords,
		&r.MetricsID, &r.CreatedBy, &r.CreatedAt,
		&issueDate, &r.InvoiceNumber, &r.EntityID,
	)

	if v, ok := issueDate.(time.Time); ok {
		r.IssueDate = v
	}

	return err
}
