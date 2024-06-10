package models

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceCancel struct {
	ID            int           `json:"-"`
	InvoiceNumber string        `json:"invoice_number"`
	Year          int           `json:"year"`
	Justification string        `json:"justification"`
	Entity        *Entity       `json:"entity"`
	ReqStatus     string        `json:"-"`
	ReqMsg        string        `json:"-"`
	CreatedBy     int           `json:"-"`
	CreatedAt     time.Time     `json:"-"`
	UpdatedAt     time.Time     `json:"-"`
	Errors        ErrorMessages `json:"-"`
}

func (c *InvoiceCancel) GetCreatedAt() time.Time {
	return c.CreatedAt
}

func (c *InvoiceCancel) GetStatus() string {
	return c.ReqStatus
}

func NewInvoiceCancel() *InvoiceCancel {
	return &InvoiceCancel{
		Entity: NewEntity(),
		Year:   time.Now().Year(),
	}
}

func NewInvoiceCancelWithSamples(entities []*Entity) *InvoiceCancel {
	canceling := NewInvoiceCancel()
	if len(entities) > 0 {
		canceling.Entity = entities[0]
	}

	return canceling
}

func NewInvoiceCancelFromForm(c *fiber.Ctx) *InvoiceCancel {
	var err error

	invoiceCancel := NewInvoiceCancel()

	invoiceCancel.InvoiceNumber = strings.TrimSpace(c.FormValue("invoice_number"))
	invoiceCancel.Year, err = utils.TrimSpaceInt(c.FormValue("year"))
	if err != nil {
		log.Println("Error converting invoice canceling year from string to int: ", err)
	}
	invoiceCancel.Justification = strings.TrimSpace(c.FormValue("justification"))

	entityID, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		entityID = 0
	}
	invoiceCancel.Entity.ID = entityID

	return invoiceCancel
}

func (c *InvoiceCancel) IsValid() bool {
	fields := Fields{
		{
			Name:  "InvoiceNumber",
			Value: c.InvoiceNumber,
			Rules: Rules(Required, Match(SiareNFANumberRegex)),
		},
		{
			Name:  "Year",
			Value: c.Year,
			Rules: Rules(Required, Max(time.Now().Year())),
		},
		{
			Name:  "Justification",
			Value: c.Justification,
			Rules: Rules(Required, Max(128)),
		},
	}
	errors, ok := Validate(fields)
	c.Errors = errors
	return ok
}

func (c *InvoiceCancel) Values() []any {
	// TODO
	// TROCAR NUMBER PARA INVOICE_NUMBER
	// AJUSTAR SS-API DE ACORDO
	return []any{
		&c.ID, &c.InvoiceNumber, &c.Year, &c.Justification, &c.ReqStatus, &c.ReqMsg,
		&c.Entity.ID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	}
}
