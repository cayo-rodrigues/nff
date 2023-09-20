package models

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
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
	Entities     *[]Entity
	BooleanField *[2]string
	Cfops        *[14]int
	Operations   *[2]string
	IcmsOptions  *[3]string
}

type InvoiceItem struct {
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           int
	ValuePerUnity      float64
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
	FormSelectFields   *InvoiceFormSelectFields
	Errors             *InvoiceFormError
}

func NewEmptyInvoice() *Invoice {
	return &Invoice{
		FormSelectFields: &InvoiceFormSelectFields{
			Operations:   &globals.InvoiceOperations,
			Cfops:        &globals.InvoiceCfops,
			BooleanField: &globals.InvoiceBooleanField,
			IcmsOptions:  &globals.InvoiceIcmsOptions,
		},
		Errors: &InvoiceFormError{},
	}
}

func NewInvoiceListFromForm(r *http.Request, entities *[]Entity) (*[]Invoice, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing invoice form: ", err)
		return nil, utils.InternalServerErr
	}

	invoices := []Invoice{}
	invoicesQuantity := len(r.PostForm["gta"]) // it could be actually any field
	for i := 0; i < invoicesQuantity; i++ {
		invoice := NewEmptyInvoice()

		invoice.FormSelectFields.Entities = entities

		senderId, err := strconv.Atoi(r.PostForm["sender"][i])
		if err != nil {
			log.Println("Error converting sender id from string to int: ", err)
			return nil, utils.InternalServerErr
		}
		recipientId, err := strconv.Atoi(r.PostForm["recipient"][i])
		if err != nil {
			log.Println("Error converting recipient id from string to int: ", err)
			return nil, utils.InternalServerErr
		}
		for _, entity := range *entities {
			if entity.Id == senderId {
				invoice.Sender = &entity
			}
			if entity.Id == recipientId {
				invoice.Recipient = &entity
			}
			if invoice.Sender != nil && invoice.Recipient != nil {
				break
			}
		}

		invoice.Operation = r.PostForm["operation"][i]
		invoice.Cfop, err = strconv.Atoi(r.PostForm["cfop"][i])
		if err != nil {
			log.Println("Error converting invoice cfop from string to int: ", err)
			return nil, utils.InternalServerErr
		}
		invoice.IsIcmsContributor = r.PostForm["is_icms_contributor"][i]
		invoice.IsFinalCustomer = r.PostForm["is_final_customer"][i]
		if r.PostForm["shipping"][i] != "" {
			invoice.Shipping, err = strconv.ParseFloat(r.PostForm["shipping"][i], 64)
			if err != nil {
				log.Println("Error converting invoice shipping from string to float64: ", err)
				return nil, utils.InternalServerErr
			}
		}
		invoice.AddShippingToTotal = r.PostForm["add_shipping_to_total"][i]
		invoice.Gta = r.PostForm["gta"][i]
		invoices = append(invoices, *invoice)
	}

	return &invoices, nil
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

	return isValid
}
