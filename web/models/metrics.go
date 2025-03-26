package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/cayo-rodrigues/safe"
	"github.com/gofiber/fiber/v2"
)

type Metrics struct {
	ID        int
	Entity    *Entity
	StartDate time.Time
	EndDate   time.Time
	CreatedBy int
	CreatedAt time.Time
	UpdatedAt time.Time
	Errors    safe.ErrorMessages
	*MetricsResult
}

func (m *Metrics) AsNotification() *Notification {
	return &Notification{
		ID:            m.ID,
		Status:        m.ReqStatus,
		OperationType: "Cálculo de Métricas",
		PageEndpoint:  fmt.Sprintf("/metrics/%d", m.ID),
		CreatedAt:     m.CreatedAt,
		UserID:        m.CreatedBy,
	}
}

func (m *Metrics) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *Metrics) GetStatus() string {
	return m.ReqStatus
}

type MetricsResult struct {
	ID              int              `json:"-"`
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
	InvoiceNumber   string           `json:"invoice_id"` // AJUSTAR NOME NA SS-API para invoice_number
	InvoicePDF      string           `json:"invoice_pdf"`
	IssueDate       time.Time        `json:"issue_date"`
	InvoiceSender   string           `json:"invoice_sender"`
	Months          []*MetricsResult `json:"months"`
	Records         []*MetricsResult `json:"records"`
	Total           *MetricsResult   `json:"total"`
	MetricsID       int              `json:"-"`
	EntityID        int              `json:"-"`
	CreatedBy       int              `json:"-"`
	CreatedAt       time.Time        `json:"-"`
}

func NewMetrics() *Metrics {
	return &Metrics{
		Entity: NewEntity(),
		MetricsResult: &MetricsResult{
			Months:  []*MetricsResult{},
			Records: []*MetricsResult{},
			Total:   &MetricsResult{},
		},
		EndDate: time.Now(),
	}

}

func NewMetricsResult() *MetricsResult {
	return &MetricsResult{}
}

func NewMetricsFromForm(c *fiber.Ctx) *Metrics {
	m := NewMetrics()

	startDate, err := utils.ParseDateAsBR(strings.TrimSpace(c.FormValue("start_date")))
	if err != nil {
		log.Println("Error converting input start date string to time.Time:", err)
	}
	endDate, err := utils.ParseDateAsBR(strings.TrimSpace(c.FormValue("end_date")))
	if err != nil {
		log.Println("Error converting input end date string to time.Time:", err)
	}
	entityID, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		entityID = 0
	}

	m.Entity.ID = entityID
	m.StartDate = startDate
	m.EndDate = endDate

	return m
}

func (m *Metrics) IsValid() bool {
	fields := safe.Fields{
		{
			Name:  "StartDate",
			Value: m.StartDate,
			Rules: safe.Rules{safe.Required(), safe.NotAfter(m.EndDate), safe.MaxDaysRange(m.EndDate, 365)},
		},
		{
			Name:  "EndDate",
			Value: m.EndDate,
			Rules: safe.Rules{safe.Required()},
		},
	}
	errors, ok := safe.Validate(fields)
	m.Errors = errors
	return ok
}

func (m *Metrics) Values() []any {
	return []any{
		&m.ID, &m.StartDate, &m.EndDate, &m.ReqStatus, &m.ReqMsg,
		&m.Entity.ID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt,
	}
}

func (r *MetricsResult) Values() []any {
	return []any{
		&r.ID, &r.Type, &r.MonthName,
		&r.TotalIncome, &r.TotalExpenses, &r.AvgIncome,
		&r.AvgExpenses, &r.Diff, &r.IsPositive,
		&r.TotalRecords, &r.PositiveRecords, &r.NegativeRecords,
		&r.MetricsID, &r.CreatedBy, &r.CreatedAt,
		&r.IssueDate, &r.InvoiceNumber, &r.EntityID,
		&r.InvoicePDF, &r.InvoiceSender,
	}
}
