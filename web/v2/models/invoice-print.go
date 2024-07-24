package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoicePrint struct {
	ID                   int
	InvoiceID            string        `json:"invoice_id"` // ATUALIZAR SS-API PARA USAR invoice_number_or_protocol
	InvoiceIDType        string        `json:"invoice_id_type"`
	InvoicePDF           string        `json:"invoice_pdf"`
	Entity               *Entity       `json:"entity"`
	ReqStatus            string        `json:"-"`
	ReqMsg               string        `json:"-"`
	CreatedBy            int           `json:"-"`
	CreatedAt            time.Time     `json:"-"`
	UpdatedAt            time.Time     `json:"-"`
	CustomFileNamePrefix string        `json:"custom_file_name"` // ATUALIZAR SS-API PARA USAR custom_file_name_prefix
	FileName             string        `json:"file_name"`
	Errors               ErrorMessages `json:"-"`
}

func (p *InvoicePrint) AsNotification() *Notification {
	return &Notification{
		ID:            p.ID,
		Status:        p.ReqStatus,
		OperationType: "Impressão/Download de NFA",
		PageEndpoint:  "/invoices/print",
		InvoicePDF:    p.InvoicePDF,
		CreatedAt:     p.CreatedAt,
		UserID:        p.CreatedBy,
	}
}

func (p *InvoicePrint) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *InvoicePrint) GetStatus() string {
	return p.ReqStatus
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
	invoicePrint.CustomFileNamePrefix = strings.TrimSpace(c.FormValue("custom_file_name_prefix"))

	entityID, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		entityID = 0
	}
	invoicePrint.Entity.ID = entityID

	return invoicePrint
}

func (p *InvoicePrint) IsValid() bool {
	fields := Fields{
		{
			Name:  "InvoiceID",
			Value: p.InvoiceID,
			Rules: Rules(Required),
		},
		{
			Name:  "CustomFileNamePrefix",
			Value: p.CustomFileNamePrefix,
			Rules: Rules(Max(64)),
		},
	}
	errors, ok := Validate(fields)
	p.Errors = errors

	if !p.InvoiceIDFormatIsValid() {
		return false
	}

	return ok
}

func (p *InvoicePrint) InvoiceIDFormatIsValid() bool {
	if p.Errors == nil {
		p.Errors = make(ErrorMessages)
	}
	_, hasVal := p.Errors["InvoiceID"]
	if hasVal {
		return true
	}

	isNumber := SiareNFANumberRegex.MatchString(p.InvoiceID)
	if isNumber {
		p.InvoiceIDType = InvoiceIDTypes.NFANumber()
		return true
	}

	isProtocol := SiareNFAProtocolRegex.MatchString(p.InvoiceID)
	if isProtocol {
		p.InvoiceIDType = InvoiceIDTypes.NFAProtocol()
		return true
	}

	p.Errors["InvoiceID"] = utils.InvalidFormatMsg
	return false
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
