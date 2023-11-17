package models

import (
	"sync"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoicePrintFormSelectFields struct {
	Entities       []*Entity
	InvoiceIdTypes *[2]string
}

type InvoicePrintFormErrors struct {
	InvoiceId     string
	InvoiceIdType string
	Entity        string
}

type InvoicePrint struct {
	Id            int
	InvoiceId     string                  `json:"invoice_id"`
	InvoiceIdType string                  `json:"invoice_id_type"`
	InvoicePDF    string                  `json:"invoice_pdf"`
	Entity        *Entity                 `json:"entity"`
	Errors        *InvoicePrintFormErrors `json:"-"`
	ReqStatus     string                  `json:"-"`
	ReqMsg        string                  `json:"-"`
}

func NewEmptyInvoicePrint() *InvoicePrint {
	return &InvoicePrint{
		Entity: NewEmptyEntity(),
		Errors: &InvoicePrintFormErrors{},
	}
}

func (p *InvoicePrint) Scan(rows db.Scanner) error {
	return rows.Scan(
		&p.Id, &p.InvoiceId, &p.InvoiceIdType, &p.InvoicePDF,
		&p.ReqStatus, &p.ReqMsg, &p.Entity.Id,
	)
}

func NewInvoicePrintFromForm(c *fiber.Ctx) *InvoicePrint {
	invoicePrint := NewEmptyInvoicePrint()

	invoicePrint.InvoiceId = c.FormValue("invoice_id")
	invoicePrint.InvoiceIdType = c.FormValue("invoice_id_type")

	return invoicePrint
}

func (i *InvoicePrint) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	// unacceptableValueMsg := "Valor inaceitável"
	validationsCount := 3

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(i.Entity == nil, &i.Errors.Entity, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.InvoiceId == "", &i.Errors.InvoiceId, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.InvoiceIdType == "", &i.Errors.InvoiceIdType, &mandatoryFieldMsg, ch, &wg)
	// TODO
	// validate invoice id format for protocol and number

	wg.Wait()
	close(ch)

	for i := 0; i < validationsCount; i++ {
		if validationPassed := <-ch; !validationPassed {
			isValid = false
			break
		}
	}

	return isValid
}
