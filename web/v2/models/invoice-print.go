package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type InvoicePrint struct {
	ID                   int
	InvoiceID            string        `json:"invoice_number_or_protocol"`
	InvoiceIDType        string        `json:"invoice_id_type"`
	InvoicePDF           string        `json:"invoice_pdf"`
	Entity               *Entity       `json:"entity"`
	ReqStatus            string        `json:"-"`
	ReqMsg               string        `json:"-"`
	CreatedBy            int           `json:"-"`
	CreatedAt            time.Time     `json:"-"`
	UpdatedAt            time.Time     `json:"-"`
	CustomFileNamePrefix string        `json:"custom_file_name_prefix"`
	FileName             string        `json:"file_name"`
	Errors               ErrorMessages `json:"-"`
}

func (i *InvoicePrint) GetStatus() string {
	return i.ReqStatus
}

func NewInvoicePrint() *InvoicePrint {
	return &InvoicePrint{
		Entity: NewEntity(),
	}
}

func NewInvoicePrintWithSamples(entities []*Entity) *InvoicePrint {
	printing := NewInvoicePrint()
	if len(entities) > 0 {
		printing.Entity = entities[0]
	}

	return printing
}

func NewInvoicePrintFromForm(c *fiber.Ctx) *InvoicePrint {
	invoicePrint := NewInvoicePrint()

	invoicePrint.InvoiceID = strings.TrimSpace(c.FormValue("invoice_id"))
	invoicePrint.InvoiceIDType = strings.TrimSpace(c.FormValue("invoice_id_type"))
	invoicePrint.CustomFileNamePrefix = strings.TrimSpace(c.FormValue("custom_file_name_prefix"))

	entityID, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		entityID = 0
	}
	invoicePrint.Entity.ID = entityID

	return invoicePrint
}

func (p *InvoicePrint) IsValid() bool {
	var invoiceIDRegex *regexp.Regexp
	if p.InvoiceIDType == "number" {
		invoiceIDRegex = SiareNFANumberRegex
	} else {
		invoiceIDRegex = SiareNFAProtocolRegex
	}

	fields := Fields{
		{
			Name:  "InvoiceID",
			Value: p.InvoiceID,
			Rules: Rules(Required, Match(invoiceIDRegex)),
		},
		{
			Name:  "InvoiceIDType",
			Value: p.InvoiceIDType,
			Rules: Rules(Required, OneOf(InvoiceIDTypes[:])),
		},
		{
			Name:  "CustomFileNamePrefix",
			Value: p.CustomFileNamePrefix,
			Rules: Rules(Max(64)),
		},
	}
	errors, ok := Validate(fields)
	p.Errors = errors
	return ok
}

func (p *InvoicePrint) Values() []any {
	// TODO
	// TROCAR CUSTOM_FILE_NAME PARA CUSTOM_FILE_NAME_PREFIX
	// ADICIONAR COLUNA FILE_NAME
	// AJUSTAR SS-API DE ACORDO
	return []any{
		&p.ID, &p.InvoiceID, &p.InvoiceIDType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileNamePrefix, &p.FileName,
	}
}
