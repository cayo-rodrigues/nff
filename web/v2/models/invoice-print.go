package models

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type InvoicePrint struct {
	ID                   int
	InvoiceNumber        string         `json:"invoice_number"`
	InvoiceProtocol      string         `json:"invoice_protocol"`
	InvoiceIDType        string         `json:"invoice_id_type"`
	InvoicePDF           string         `json:"invoice_pdf"`
	Entity               *Entity        `json:"entity"`
	ReqStatus            string         `json:"-"`
	ReqMsg               string         `json:"-"`
	CreatedBy            int            `json:"-"`
	CreatedAt            time.Time      `json:"-"`
	UpdatedAt            time.Time      `json:"-"`
	CustomFileNamePrefix string         `json:"custom_file_name_prefix"`
	FileName             string         `json:"file_name"`
	Errors               *ErrorMessages `json:"-"`
}

func NewInvoicePrint() *InvoicePrint {
	return &InvoicePrint{
		Entity: NewEntity(),
	}
}

func NewInvoicePrintFromForm(c *fiber.Ctx) *InvoicePrint {
	invoicePrint := NewInvoicePrint()

	invoicePrint.InvoiceNumber = strings.TrimSpace(c.FormValue("invoice_number"))
	invoicePrint.InvoiceProtocol = strings.TrimSpace(c.FormValue("invoice_protocol"))
	invoicePrint.CustomFileNamePrefix = strings.TrimSpace(c.FormValue("custom_file_name_prefix"))

	return invoicePrint
}

func (p *InvoicePrint) IsValid() bool {
	fields := Fields{
		{
			Name:  "InvoiceNumber",
			Value: p.InvoiceNumber,
			Rules: Rules(RequiredUnlessAtLeastOneIsPresent(p.InvoiceProtocol), Match(SiareNFANumberRegex)),
		},
		{
			Name:  "InvoiceProtocol",
			Value: p.InvoiceProtocol,
			Rules: Rules(RequiredUnlessAtLeastOneIsPresent(p.InvoiceNumber), Match(SiareNFAProtocolRegex)),
		},
		{
			Name:  "CustomFileNamePrefix",
			Value: p.CustomFileNamePrefix,
			Rules: Rules(Max(64)),
		},
	}
	errors, isValid := Validate(fields)
	p.Errors = &errors
	return isValid
}

func (p *InvoicePrint) Values() []any {
	// TODO
	// TROCAR INVOICE_ID PARA INVOICE_NUMBER
	// TROCAR CUSTOM_FILE_NAME PARA CUSTOM_FILE_NAME_PREFIX
	// ADICIONAR COLUNAS FILE_NAME E INVOICE_PROTOCOL
	// AJUSTAR SS-API DE ACORDO
	return []any{
		&p.ID, &p.InvoiceNumber, &p.InvoiceIDType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileNamePrefix, &p.FileName, &p.InvoiceProtocol,
	}
}
