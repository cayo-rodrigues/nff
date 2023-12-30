package models

import (
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceCancelSelectFields struct {
	Entities []*Entity
}

type InvoiceCancelFormErrors struct {
	Number        string
	Year          string
	Justification string
	Entity        string
}

func NewInvoiceCancelSelectFields() *InvoiceCancelSelectFields {
	return &InvoiceCancelSelectFields{
		Entities: []*Entity{},
	}
}

type InvoiceCancel struct {
	ID            int                      `json:"-"`
	Number        string                   `json:"invoice_id"`
	Year          int                      `json:"year"`
	Justification string                   `json:"justification"`
	Entity        *Entity                  `json:"entity"`
	Errors        *InvoiceCancelFormErrors `json:"-"`
	ReqStatus     string                   `json:"-"`
	ReqMsg        string                   `json:"-"`
	CreatedBy     int                      `json:"-"`
	CreatedAt     time.Time                `json:"-"`
	UpdatedAt     time.Time                `json:"-"`
}

func NewEmptyInvoiceCancel() *InvoiceCancel {
	return &InvoiceCancel{
		Entity: NewEmptyEntity(),
		Errors: &InvoiceCancelFormErrors{},
		Year:   time.Now().Year(),
	}
}

func NewInvoiceCancelFromForm(c *fiber.Ctx) *InvoiceCancel {
	var err error

	invoiceCancel := NewEmptyInvoiceCancel()

	invoiceCancel.Number = strings.TrimSpace(c.FormValue("invoice_id"))
	invoiceCancel.Year, err = utils.TrimSpaceInt(c.FormValue("year"))
	if err != nil {
		log.Println("Error converting invoice canceling year from string to int: ", err)
	}
	invoiceCancel.Justification = strings.TrimSpace(c.FormValue("justification"))

	return invoiceCancel
}

func (i *InvoiceCancel) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg
	invalidFormatMsg := globals.InvalidFormatMsg
	valueTooLongMsg := globals.ValueTooLongMsg

	hasEntity := i.Entity != nil
	hasInvoiceNumber := i.Number != ""
	hasJustification := i.Justification != ""
	hasYear := i.Year != 0
	hasValidYear := i.Year <= time.Now().Year()

	hasValidInvoiceNumberFormat := globals.ReSiareNFANumber.MatchString(i.Number)

	justificationTooLong := utf8.RuneCount([]byte(i.Justification)) > 128

	fields := [7]*utils.Field{
		{ErrCondition: !hasEntity, ErrField: &i.Errors.Entity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceNumber, ErrField: &i.Errors.Number, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasJustification, ErrField: &i.Errors.Justification, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasYear, ErrField: &i.Errors.Year, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasValidYear, ErrField: &i.Errors.Year, ErrMsg: &unacceptableValueMsg},
		{ErrCondition: hasInvoiceNumber && !hasValidInvoiceNumberFormat, ErrField: &i.Errors.Number, ErrMsg: &invalidFormatMsg},
		{ErrCondition: justificationTooLong, ErrField: &i.Errors.Justification, ErrMsg: &valueTooLongMsg},
	}

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	return isValid
}

func (c *InvoiceCancel) Scan(rows db.Scanner) error {
	return rows.Scan(
		&c.ID, &c.Number, &c.Year, &c.Justification, &c.ReqStatus, &c.ReqMsg,
		&c.Entity.ID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
}
func (c *InvoiceCancel) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&c.ID, &c.Number, &c.Year, &c.Justification, &c.ReqStatus, &c.ReqMsg,
		&c.Entity.ID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,

		&c.Entity.ID, &c.Entity.Name, &c.Entity.UserType, &c.Entity.CpfCnpj, &c.Entity.Ie, &c.Entity.Email, &c.Entity.Password,
		&c.Entity.Address.PostalCode, &c.Entity.Address.Neighborhood, &c.Entity.Address.StreetType, &c.Entity.Address.StreetName, &c.Entity.Address.Number,
		&c.Entity.CreatedBy, &c.Entity.CreatedAt, &c.Entity.UpdatedAt,
	)
}
