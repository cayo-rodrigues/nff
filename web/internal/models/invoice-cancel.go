package models

import (
	"log"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceCancelFormSelectFields struct {
	Entities []*Entity
}

type InvoiceCancelFormErrors struct {
	Number        string
	Year          string
	Justification string
	Entity        string
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

func NewInvoiceCancelFromForm(c *fiber.Ctx) (*InvoiceCancel, error) {
	var err error

	invoiceCancel := NewEmptyInvoiceCancel()

	invoiceCancel.Number = c.FormValue("invoice_id")
	if c.FormValue("year") != "" {
		invoiceCancel.Year, err = strconv.Atoi(c.FormValue("year"))
		if err != nil {
			log.Println("Error converting invoice canceling year from string to int: ", err)
			return nil, utils.InternalServerErr
		}
	}
	invoiceCancel.Justification = c.FormValue("justification")

	return invoiceCancel, nil
}

func (i *InvoiceCancel) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg
	valueTooLongMsg := globals.ValueTooLongMsg

	hasEntity := i.Entity != nil
	hasInvoiceNumber := i.Number != ""
	hasJustification := i.Justification != ""
	hasYear := i.Year != 0
	hasValidYear := i.Year <= time.Now().Year()

	invoiceNumberTooLong := utf8.RuneCount([]byte(i.Number)) > 9
	justificationTooLong := utf8.RuneCount([]byte(i.Justification)) > 128

	fields := [7]*utils.Field{
		{ErrCondition: !hasEntity, ErrField: &i.Errors.Entity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasInvoiceNumber, ErrField: &i.Errors.Number, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasJustification, ErrField: &i.Errors.Justification, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasYear, ErrField: &i.Errors.Year, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasValidYear, ErrField: &i.Errors.Year, ErrMsg: &unacceptableValueMsg},
		{ErrCondition: invoiceNumberTooLong, ErrField: &i.Errors.Number, ErrMsg: &valueTooLongMsg},
		{ErrCondition: justificationTooLong, ErrField: &i.Errors.Justification, ErrMsg: &valueTooLongMsg},
	}

	var wg sync.WaitGroup
	for _, field := range fields {
		wg.Add(1)
		go utils.ValidateField(field, &isValid, &wg)
	}

	wg.Wait()

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
