package models

import (
	"log"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
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
	Id                 int
	Number             string
	Protocol           string
	Operation          string
	Cfop               int
	IsFinalCustomer    string
	IsIcmsContributor  string
	Shipping           float64
	AddShippingToTotal string
	Gta                string
	Sender             *Entity
	Recipient          *Entity
	Items              *[]InvoiceItem
	Errors             *InvoiceFormError
	OverviewType       string
}

func NewEmptyInvoice() *Invoice {
	return &Invoice{
		Errors:       &InvoiceFormError{},
		Items:        &[]InvoiceItem{*NewEmptyInvoiceItem()},
		Sender:       NewEmptyEntity(),
		Recipient:    NewEmptyEntity(),
		OverviewType: "invoice",
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
	// TODO parallel validation
	isValid := true

	if i.Sender == nil {
		i.Errors.Sender = "Campo obrigatório"
		isValid = false
	}

	if i.Recipient == nil {
		i.Errors.Recipient = "Campo obrigatório"
		isValid = false
	}

	if i.Operation == "" {
		i.Errors.Operation = "Campo obrigatório"
		isValid = false
	}

	if i.Cfop == 0 {
		i.Errors.Cfop = "Campo obrigatório"
		isValid = false
	}

	if i.IsIcmsContributor == "" {
		i.Errors.Cfop = "Campo obrigatório"
		isValid = false
	}

	if i.IsFinalCustomer == "" {
		i.Errors.IsFinalCustomer = "Campo obrigatório"
		isValid = false
	}

	if i.AddShippingToTotal == "" {
		i.Errors.AddShippingToTotal = "Campo obrigatório"
		isValid = false
	}

	if i.Operation != "" {
		hasValidOption := false
		for _, operation := range &globals.InvoiceOperations {
			if i.Operation == operation {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Operation = "Valor inaceitável"
			isValid = false
		}
	}

	if i.Cfop != 0 {
		hasValidOption := false
		for _, cfop := range &globals.InvoiceCfops {
			if i.Cfop == cfop {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Cfop = "Valor inaceitável"
			isValid = false
		}
	}

	if i.IsIcmsContributor != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceIcmsOptions {
			if i.IsIcmsContributor == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.IsIcmsContributor = "Valor inaceitável"
			isValid = false
		}
	}

	if i.IsFinalCustomer != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceBooleanField {
			if i.IsFinalCustomer == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.IsFinalCustomer = "Valor inaceitável"
			isValid = false
		}
	}

	if i.AddShippingToTotal != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceBooleanField {
			if i.AddShippingToTotal == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.AddShippingToTotal = "Valor inaceitável"
			isValid = false
		}
	}

	if len(*i.Items) == 0 {
		i.Errors.Items = "A NF deve ter pelo menos 1 produto"
		isValid = false
	}

	for _, item := range *i.Items {
		if !item.IsValid() {
			i.Errors.Items = "Dados dos produtos inválidos"
			isValid = false
		}
	}

	return isValid
}

func (i *Invoice) Scan(rows sql.Scanner) error {
	return rows.Scan(
		&i.Id, &i.Number, &i.Protocol, &i.Operation, &i.Cfop, &i.IsFinalCustomer, &i.IsIcmsContributor,
		&i.Shipping, &i.AddShippingToTotal, &i.Gta, &i.Sender.Id, &i.Recipient.Id,
	)
}
