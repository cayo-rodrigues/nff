package models

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
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

func (c *InvoiceCancel) Scan(rows db.Scanner) error {
	return rows.Scan(
		&c.ID, &c.Number, &c.Year, &c.Justification, &c.ReqStatus, &c.ReqMsg,
		&c.Entity.ID, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt,
	)
}
