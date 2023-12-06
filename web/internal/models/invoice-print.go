package models

import (
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoicePrintFormSelectFields struct {
	Entities       []*Entity
	InvoiceIdTypes *[2]string
}

type InvoicePrintFormErrors struct {
	InvoiceId      string
	InvoiceIdType  string
	Entity         string
	CustomFileName string
}

type InvoicePrint struct {
	ID             int
	InvoiceId      string                  `json:"invoice_id"`
	InvoiceIdType  string                  `json:"invoice_id_type"`
	InvoicePDF     string                  `json:"invoice_pdf"`
	Entity         *Entity                 `json:"entity"`
	Errors         *InvoicePrintFormErrors `json:"-"`
	ReqStatus      string                  `json:"-"`
	ReqMsg         string                  `json:"-"`
	CreatedBy      int                     `json:"-"`
	CreatedAt      time.Time               `json:"-"`
	UpdatedAt      time.Time               `json:"-"`
	CustomFileName string                  `json:"custom_file_name"`
}

func NewEmptyInvoicePrint() *InvoicePrint {
	return &InvoicePrint{
		Entity: NewEmptyEntity(),
		Errors: &InvoicePrintFormErrors{},
	}
}

func NewInvoicePrintFromForm(c *fiber.Ctx) *InvoicePrint {
	invoicePrint := NewEmptyInvoicePrint()

	invoicePrint.InvoiceId = c.FormValue("invoice_id")
	invoicePrint.InvoiceIdType = c.FormValue("invoice_id_type")
	invoicePrint.CustomFileName = c.FormValue("custom_file_name")

	return invoicePrint
}

func (i *InvoicePrint) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	valueTooLongMsg := globals.ValueTooLongMsg
	invalidFormatMsg := globals.InvalidFormatMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg

	hasEntity := i.Entity != nil
	hasInvoiceId := i.InvoiceId != ""
	hasInvoiceIdType := i.InvoiceIdType != ""
	hasCustomFileName := i.CustomFileName != ""

	hasValidInvoiceNumberFormat := globals.ReSiareNFANumber.MatchString(i.InvoiceId)
	hasValidInvoiceProtocolFormat := globals.ReSiareNFAProtocol.MatchString(i.InvoiceId)

	customFileNameTooLong := utf8.RuneCount([]byte(i.CustomFileName)) > 64

	idTypeIsNumber := i.InvoiceIdType == globals.InvoiceIdTypes[0]
	idTypeIsProtocol := i.InvoiceIdType == globals.InvoiceIdTypes[1]

	fields := [6]*utils.Field{
		{ErrCondition: !hasEntity, ErrField: &i.Errors.Entity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceId, ErrField: &i.Errors.InvoiceId, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceIdType, ErrField: &i.Errors.InvoiceIdType, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: hasInvoiceId && idTypeIsNumber && !hasValidInvoiceNumberFormat, ErrField: &i.Errors.InvoiceId, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasInvoiceId && idTypeIsProtocol && !hasValidInvoiceProtocolFormat, ErrField: &i.Errors.InvoiceId, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasCustomFileName && customFileNameTooLong, ErrField: &i.Errors.CustomFileName, ErrMsg: &valueTooLongMsg},
	}

	var wg sync.WaitGroup
	for _, field := range fields {
		wg.Add(1)
		go utils.ValidateField(field, &isValid, &wg)
	}

	wg.Add(1)
	go utils.ValidateListField(i.InvoiceIdType, globals.InvoiceIdTypes[:], &i.Errors.InvoiceIdType, &unacceptableValueMsg, &isValid, &wg)

	wg.Wait()

	return isValid
}

func (p *InvoicePrint) Scan(rows db.Scanner) error {
	return rows.Scan(
		&p.ID, &p.InvoiceId, &p.InvoiceIdType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileName,
	)
}

func (p *InvoicePrint) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&p.ID, &p.InvoiceId, &p.InvoiceIdType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileName,

		&p.Entity.ID, &p.Entity.Name, &p.Entity.UserType, &p.Entity.CpfCnpj, &p.Entity.Ie, &p.Entity.Email, &p.Entity.Password,
		&p.Entity.Address.PostalCode, &p.Entity.Address.Neighborhood, &p.Entity.Address.StreetType, &p.Entity.Address.StreetName, &p.Entity.Address.Number,
		&p.Entity.CreatedBy, &p.Entity.CreatedAt, &p.Entity.UpdatedAt,
	)
}
