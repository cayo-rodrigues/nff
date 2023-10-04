package models

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/sql"
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
	Id            int
	Number        string
	Year          int
	Justification string
	Entity        *Entity
	Errors        *InvoiceCancelFormErrors
	OverviewType  string
}

func NewEmptyInvoiceCancel() *InvoiceCancel {
	return &InvoiceCancel{
		Entity:       NewEmptyEntity(),
		Errors:       &InvoiceCancelFormErrors{},
		Year:         time.Now().Year(),
		OverviewType: "canceling",
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

	mandatoryFieldMsg := "Campo obrigatório"
	unacceptableValueMsg := "Valor inaceitável"
	validationsCount := 4

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(i.Entity == nil, &i.Errors.Entity, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Number == "", &i.Errors.Number, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Justification == "", &i.Errors.Justification, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Year == 0 || i.Year > time.Now().Year(), &i.Errors.Year, &unacceptableValueMsg, ch, &wg)

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

func (c *InvoiceCancel) Scan(rows sql.Scanner) error {
	return rows.Scan(&c.Id, &c.Number, &c.Year, &c.Justification, &c.Entity.Id)
}
