package models

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoicePrintSelectFields struct {
	Entities       []*Entity
	InvoiceIDTypes *globals.SiareInvoiceIDTypes
}

func NewInvoicePrintSelectFields() *InvoicePrintSelectFields {
	return &InvoicePrintSelectFields{
		Entities:       []*Entity{},
		InvoiceIDTypes: &globals.InvoiceIDTypes,
	}
}

type InvoicePrintFormErrors struct {
	InvoiceID      string
	InvoiceIDType  string
	Entity         string
	CustomFileName string
}

type InvoicePrint struct {
	ID             int
	InvoiceID      string                  `json:"invoice_id"`
	InvoiceIDType  string                  `json:"invoice_id_type"`
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

	invoicePrint.InvoiceID = strings.TrimSpace(c.FormValue("invoice_id"))
	invoicePrint.InvoiceIDType = strings.TrimSpace(c.FormValue("invoice_id_type"))
	invoicePrint.CustomFileName = strings.TrimSpace(c.FormValue("custom_file_name"))

	return invoicePrint
}

func (i *InvoicePrint) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	valueTooLongMsg := globals.ValueTooLongMsg
	invalidFormatMsg := globals.InvalidFormatMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg

	hasEntity := i.Entity != nil
	hasInvoiceId := i.InvoiceID != ""
	hasInvoiceIdType := i.InvoiceIDType != ""
	hasCustomFileName := i.CustomFileName != ""

	hasValidInvoiceNumberFormat := globals.ReSiareNFANumber.MatchString(i.InvoiceID)
	hasValidInvoiceProtocolFormat := globals.ReSiareNFAProtocol.MatchString(i.InvoiceID)

	customFileNameTooLong := utf8.RuneCount([]byte(i.CustomFileName)) > 64

	idTypeIsNumber := i.InvoiceIDType == globals.InvoiceIDTypes[0]
	idTypeIsProtocol := i.InvoiceIDType == globals.InvoiceIDTypes[1]

	fields := [6]*utils.Field{
		{ErrCondition: !hasEntity, ErrField: &i.Errors.Entity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceId, ErrField: &i.Errors.InvoiceID, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceIdType, ErrField: &i.Errors.InvoiceIDType, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: hasInvoiceId && idTypeIsNumber && !hasValidInvoiceNumberFormat, ErrField: &i.Errors.InvoiceID, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasInvoiceId && idTypeIsProtocol && !hasValidInvoiceProtocolFormat, ErrField: &i.Errors.InvoiceID, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasCustomFileName && customFileNameTooLong, ErrField: &i.Errors.CustomFileName, ErrMsg: &valueTooLongMsg},
	}

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	utils.ValidateListField(i.InvoiceIDType, globals.InvoiceIDTypes[:], &i.Errors.InvoiceIDType, &unacceptableValueMsg, &isValid)

	return isValid
}

func (p *InvoicePrint) Scan(rows db.Scanner) error {
	return rows.Scan(
		&p.ID, &p.InvoiceID, &p.InvoiceIDType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileName,
	)
}

func (p *InvoicePrint) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&p.ID, &p.InvoiceID, &p.InvoiceIDType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.ID,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CustomFileName,

		&p.Entity.ID, &p.Entity.Name, &p.Entity.UserType, &p.Entity.CpfCnpj, &p.Entity.Ie, &p.Entity.Email, &p.Entity.Password,
		&p.Entity.Address.PostalCode, &p.Entity.Address.Neighborhood, &p.Entity.Address.StreetType, &p.Entity.Address.StreetName, &p.Entity.Address.Number,
		&p.Entity.CreatedBy, &p.Entity.CreatedAt, &p.Entity.UpdatedAt, &p.Entity.OtherIes,
	)
}
