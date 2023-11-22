package models

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceFormError struct {
	Sender             string
	Recipient          string
	Operation          string
	Cfop               string
	IsFinalCustomer    string
	IsIcmsContributor  string
	Shipping           string
	AddShippingToTotal string
	Gta                string
	Items              string
}

type InvoiceFormSelectFields struct {
	Entities     []*Entity
	BooleanField *[2]string
	Cfops        *[14]int
	Operations   *[2]string
	IcmsOptions  *[3]string
}

type Invoice struct {
	ID                 int               `json:"-"`
	Number             string            `json:"invoice_id"`
	Protocol           string            `json:"-"`
	Operation          string            `json:"operation"`
	Cfop               int               `json:"cfop,string"`
	IsFinalCustomer    string            `json:"is_final_customer"`
	IsIcmsContributor  string            `json:"icms"`
	Shipping           float64           `json:"shipping"`
	AddShippingToTotal string            `json:"add_shipping_to_total_value"`
	Gta                string            `json:"gta"`
	Sender             *Entity           `json:"sender"`
	Recipient          *Entity           `json:"recipient"`
	Items              []*InvoiceItem    `json:"items"`
	Errors             *InvoiceFormError `json:"-"`
	ReqStatus          string            `json:"-"`
	ReqMsg             string            `json:"-"`
	PDF                string            `json:"invoice_pdf"`
	CreatedBy          int               `json:"-"`
	CreatedAt          time.Time         `json:"-"`
	UpdatedAt          time.Time         `json:"-"`
}

func NewEmptyInvoice() *Invoice {
	return &Invoice{
		Errors:    &InvoiceFormError{},
		Items:     []*InvoiceItem{NewEmptyInvoiceItem()},
		Sender:    NewEmptyEntity(),
		Recipient: NewEmptyEntity(),
	}
}

func NewInvoiceFromForm(c *fiber.Ctx) (*Invoice, error) {
	var err error

	invoice := NewEmptyInvoice()

	invoice.Operation = c.FormValue("operation")
	invoice.Cfop, err = strconv.Atoi(c.FormValue("cfop"))
	if err != nil {
		log.Println("Error converting invoice cfop from string to int: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.IsIcmsContributor = c.FormValue("is_icms_contributor")
	invoice.IsFinalCustomer = c.FormValue("is_final_customer")
	if c.FormValue("shipping") != "" {
		invoice.Shipping, err = strconv.ParseFloat(c.FormValue("shipping"), 64)
		if err != nil {
			log.Println("Error converting invoice shipping from string to float64: ", err)
			return nil, utils.InternalServerErr
		}
	}
	invoice.AddShippingToTotal = c.FormValue("add_shipping_to_total")
	invoice.Gta = c.FormValue("gta")

	items, err := NewInvoiceItemsFromForm(c)
	if err != nil {
		log.Println("Error getting invoice items from form: ", err)
		return nil, err
	}

	invoice.Items = items

	return invoice, nil
}

func (i *Invoice) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	unacceptableValueMsg := "Valor inaceitável"
	mustHaveItemsMsg := "A NF deve ter pelo menos 1 produto"
	invalidItemsMsg := "Dados dos produtos inválidos"
	itemsCount := len(i.Items)
	validationsCount := 14 + itemsCount

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(i.Sender == nil, &i.Errors.Sender, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Recipient == nil, &i.Errors.Recipient, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Operation == "", &i.Errors.Operation, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Cfop == 0, &i.Errors.Cfop, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.IsIcmsContributor == "", &i.Errors.IsIcmsContributor, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.IsFinalCustomer == "", &i.Errors.IsFinalCustomer, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Shipping == 0, &i.Errors.Shipping, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.AddShippingToTotal == "", &i.Errors.AddShippingToTotal, &mandatoryFieldMsg, ch, &wg)

	go utils.ValidateListField(i.Operation, globals.InvoiceOperations[:], &i.Errors.Operation, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.Cfop, globals.InvoiceCfops[:], &i.Errors.Cfop, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.IsIcmsContributor, globals.InvoiceIcmsOptions[:], &i.Errors.IsIcmsContributor, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.IsFinalCustomer, globals.InvoiceBooleanField[:], &i.Errors.IsFinalCustomer, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.AddShippingToTotal, globals.InvoiceBooleanField[:], &i.Errors.AddShippingToTotal, &unacceptableValueMsg, ch, &wg)

	go utils.ValidateField(itemsCount == 0, &i.Errors.Items, &mustHaveItemsMsg, ch, &wg)

	for _, item := range i.Items {
		go utils.ValidateField(!item.IsValid(), &i.Errors.Items, &invalidItemsMsg, ch, &wg)
	}

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

func (i *Invoice) Scan(rows db.Scanner) error {
	return rows.Scan(
		&i.ID, &i.Number, &i.Protocol, &i.Operation, &i.Cfop, &i.IsFinalCustomer, &i.IsIcmsContributor,
		&i.Shipping, &i.AddShippingToTotal, &i.Gta, &i.PDF, &i.ReqStatus, &i.ReqMsg,
		&i.Sender.ID, &i.Recipient.ID, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
	)
}
