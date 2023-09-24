package models

import "github.com/gofiber/fiber/v2"

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
	}
}

func NewInvoiceCancelFromForm(c *fiber.Ctx) (*InvoiceCancel, error) {
	invoiceCancel := NewEmptyInvoiceCancel()
	return invoiceCancel, nil
}

func (i *InvoiceCancel) IsValid() bool {
	isValid := true
	return isValid
}
