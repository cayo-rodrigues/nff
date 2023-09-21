package models

import (
	"log"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type InvoiceItemFormError struct {
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           string
	ValuePerUnity      string
}

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

type InvoiceItemFormSelectFields struct {
	Groups               *[2]string
	Origins              *[3]string
	UnitiesOfMeasurement *[23]string
}

type InvoiceItem struct {
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           float64
	ValuePerUnity      float64
	FormSelectFields   *InvoiceItemFormSelectFields
	Errors             *InvoiceItemFormError
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
}

func NewEmptyInvoiceItem() *InvoiceItem {
	return &InvoiceItem{
		FormSelectFields: &InvoiceItemFormSelectFields{
			Groups:               &globals.InvoiceItemGroups,
			Origins:              &globals.InvoiceItemOrigins,
			UnitiesOfMeasurement: &globals.InvoiceItemUnitiesOfMeaasurement,
		},
		Errors: &InvoiceItemFormError{},
	}
}

func NewEmptyInvoice() *Invoice {
	return &Invoice{
		Errors: &InvoiceFormError{},
		Items:  &[]InvoiceItem{*NewEmptyInvoiceItem()},
	}
}

func NewInvoiceFromForm(r *http.Request) (*Invoice, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing invoice form: ", err)
		return nil, utils.InternalServerErr
	}

	invoice := NewEmptyInvoice()

	invoice.Operation = r.PostFormValue("operation")
	invoice.Cfop, err = strconv.Atoi(r.PostFormValue("cfop"))
	if err != nil {
		log.Println("Error converting invoice cfop from string to int: ", err)
		return nil, utils.InternalServerErr
	}
	invoice.IsIcmsContributor = r.PostFormValue("is_icms_contributor")
	invoice.IsFinalCustomer = r.PostFormValue("is_final_customer")
	if r.PostFormValue("shipping") != "" {
		invoice.Shipping, err = strconv.ParseFloat(r.PostFormValue("shipping"), 64)
		if err != nil {
			log.Println("Error converting invoice shipping from string to float64: ", err)
			return nil, utils.InternalServerErr
		}
	}
	invoice.AddShippingToTotal = r.PostFormValue("add_shipping_to_total")
	invoice.Gta = r.PostFormValue("gta")

	items := []InvoiceItem{}
	itemsQuantity := len(r.PostForm["description"]) // it could be any field
	for i := 0; i < itemsQuantity; i++ {
		item := NewEmptyInvoiceItem()

		item.Group = r.PostForm["group"][i]
		item.Description = r.PostForm["description"][i]
		item.Origin = r.PostForm["origin"][i]
		item.UnityOfMeasurement = r.PostForm["unity_of_measurement"][i]
		item.Quantity, err = strconv.ParseFloat(r.PostForm["quantity"][i], 64)
		if err != nil {
			log.Printf("Error converting invoice item %d quantity from string to float64: %v", i, err)
			return nil, utils.InternalServerErr
		}
		item.ValuePerUnity, err = strconv.ParseFloat(r.PostForm["value_per_unity"][i], 64)
		if err != nil {
			log.Printf("Error converting invoice item %d value_per_unity from string to float64: %v", i, err)
			return nil, utils.InternalServerErr
		}

		items = append(items, *item)
	}

	invoice.Items = &items

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

func (i *InvoiceItem) IsValid() bool {
	isValid := true
	if i.Group == "" {
		i.Errors.Group = "Campo obrigatório"
		isValid = false
	}

	if i.Description == "" {
		i.Errors.Description = "Campo obrigatório"
		isValid = false
	}

	if i.Origin == "" {
		i.Errors.Origin = "Campo obrigatório"
		isValid = false
	}

	if i.UnityOfMeasurement == "" {
		i.Errors.UnityOfMeasurement = "Campo obrigatório"
		isValid = false
	}

	if i.Quantity == 0.0 {
		i.Errors.Quantity = "Valor inaceitável"
		isValid = false
	}

	if i.ValuePerUnity == 0.0 {
		i.Errors.ValuePerUnity = "Valor inaceitável"
		isValid = false
	}

	if i.Group != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemGroups {
			if i.Group == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Group = "Valor inaceitável"
			isValid = false
		}
	}

	if i.Origin != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemOrigins {
			if i.Origin == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Origin = "Valor inaceitável"
			isValid = false
		}
	}

	if i.UnityOfMeasurement != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemUnitiesOfMeaasurement {
			if i.UnityOfMeasurement == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.UnityOfMeasurement = "Valor inaceitável"
			isValid = false
		}
	}

	return isValid
}
