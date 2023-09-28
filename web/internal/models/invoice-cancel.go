package models

import (
	"log"
	"strconv"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceCancelFormSelectFields struct {
	Entities *[]Entity
}

type InvoiceCancelFormErrors struct {
	Number        string
	Year          string
	Justification string
	Entity        string
}

type InvoiceCancel struct {
	Id            int
	Number        string
	Year          int
	Justification string
	Entity        *Entity
	Errors        *InvoiceCancelFormErrors
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
	invoiceCancel.Year, err = strconv.Atoi(c.FormValue("year"))
	if err != nil {
		log.Println("Error converting invoice canceling year from string to int: ", err)
		return nil, utils.InternalServerErr
	}
	invoiceCancel.Justification = c.FormValue("justification")

	return invoiceCancel, nil
}

func (i *InvoiceCancel) IsValid() bool {
	isValid := true

	if i.Entity == nil {
		i.Errors.Entity = "Campo obrigatório"
		isValid = false
	}

	if i.Number == "" {
		i.Errors.Number = "Campo obrigatório"
		isValid = false
	}

	if i.Justification == "" {
		i.Errors.Justification = "Campo obrigatório"
		isValid = false
	}

	if i.Year == 0 || i.Year > time.Now().Year() {
		i.Errors.Year = "Valor inaceitável"
		isValid = false
	}

	return isValid
}

func (c *InvoiceCancel) Scan(rows sql.Scanner) error {
	return rows.Scan(&c.Id, &c.Number, &c.Year, &c.Justification, &c.Entity.Id)
}
