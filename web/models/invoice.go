package models

import (
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
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
	ExtraNotes         string
	CustomFileName     string
	SenderIe           string
}

type InvoiceSelectFields struct {
	Entities     []*Entity
	BooleanField *globals.BooleanField
	Cfops        *globals.SiareInvoiceCfops
	Operations   *globals.SiareInvoiceOperations
	IcmsOptions  *globals.SiareInvoiceIcmsOptions
}

func NewInvoiceSelectFields() *InvoiceSelectFields {
	return &InvoiceSelectFields{
		Operations:   &globals.InvoiceOperations,
		Cfops:        &globals.InvoiceCfops,
		BooleanField: &globals.InvoiceBooleanField,
		IcmsOptions:  &globals.InvoiceIcmsOptions,
	}
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
	ExtraNotes         string            `json:"extra_notes"`
	CustomFileName     string            `json:"custom_file_name"`
	SenderIe           string            `json:"sender_ie"`
}

func NewEmptyInvoice() *Invoice {
	return &Invoice{
		Errors:    &InvoiceFormError{},
		Items:     []*InvoiceItem{NewEmptyInvoiceItem()},
		Sender:    NewEmptyEntity(),
		Recipient: NewEmptyEntity(),
	}
}

func NewInvoiceFromForm(c *fiber.Ctx) *Invoice {
	var err error

	invoice := NewEmptyInvoice()

	invoice.Operation = strings.TrimSpace(c.FormValue("operation"))
	invoice.Cfop, err = utils.TrimSpaceInt(c.FormValue("cfop"))
	if err != nil {
		log.Println("Error converting invoice cfop from string to int: ", err)
	}
	invoice.IsIcmsContributor = strings.TrimSpace(c.FormValue("is_icms_contributor"))
	invoice.IsFinalCustomer = strings.TrimSpace(c.FormValue("is_final_customer"))
	invoice.Shipping, err = utils.TrimSpaceFloat64(c.FormValue("shipping"))
	if err != nil {
		log.Println("Error converting invoice shipping from string to float64: ", err)
	}
	invoice.AddShippingToTotal = strings.TrimSpace(c.FormValue("add_shipping_to_total"))
	invoice.Gta = strings.TrimSpace(c.FormValue("gta"))
	invoice.ExtraNotes = strings.TrimSpace(c.FormValue("extra_notes"))
	invoice.CustomFileName = strings.TrimSpace(c.FormValue("custom_file_name"))
	invoice.SenderIe = strings.TrimSpace(c.FormValue("sender_ie"))

	invoice.Items = NewInvoiceItemsFromForm(c)

	return invoice
}

func (i *Invoice) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	invalidFormatMsg := globals.InvalidFormatMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg
	mustHaveItemsMsg := globals.MustHaveItemsMsg
	invalidItemsMsg := globals.InvalidItemsMsg
	valueTooLongMsg := globals.ValueTooLongMsg

	hasSender := i.Sender != nil
	hasRecipient := i.Recipient != nil
	hasShipping := i.Shipping != 0
	hasItems := len(i.Items) >= 0
	hasGta := i.Gta != ""
	hasCustomFileName := i.CustomFileName != ""
	hasExtraNotes := i.ExtraNotes != ""

	hasValidGtaFormat := globals.ReGta.MatchString(i.Gta)

	gtaTooLong := utf8.RuneCount([]byte(i.Gta)) > 16
	customFileNameTooLong := utf8.RuneCount([]byte(i.CustomFileName)) > 64
	extraNotesTooLong := utf8.RuneCount([]byte(i.ExtraNotes)) > 512

	fields := [8]*utils.Field{
		{ErrCondition: !hasSender, ErrField: &i.Errors.Sender, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasRecipient, ErrField: &i.Errors.Recipient, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasShipping, ErrField: &i.Errors.Shipping, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasItems, ErrField: &i.Errors.Items, ErrMsg: &mustHaveItemsMsg},
		{ErrCondition: hasGta && !hasValidGtaFormat, ErrField: &i.Errors.Gta, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasGta && gtaTooLong, ErrField: &i.Errors.Gta, ErrMsg: &valueTooLongMsg},
		{ErrCondition: hasCustomFileName && customFileNameTooLong, ErrField: &i.Errors.CustomFileName, ErrMsg: &valueTooLongMsg},
		{ErrCondition: hasExtraNotes && extraNotesTooLong, ErrField: &i.Errors.ExtraNotes, ErrMsg: &valueTooLongMsg},
	}

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	for _, item := range i.Items {
		field := &utils.Field{
			ErrCondition: !item.IsValid(),
			ErrField:     &i.Errors.Items,
			ErrMsg:       &invalidItemsMsg,
		}
		utils.ValidateField(field, &isValid)
	}

	utils.ValidateListField(i.Operation, globals.InvoiceOperations[:], &i.Errors.Operation, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.Cfop, globals.InvoiceCfops[:], &i.Errors.Cfop, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.IsIcmsContributor, globals.InvoiceIcmsOptions[:], &i.Errors.IsIcmsContributor, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.IsFinalCustomer, globals.InvoiceBooleanField[:], &i.Errors.IsFinalCustomer, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.AddShippingToTotal, globals.InvoiceBooleanField[:], &i.Errors.AddShippingToTotal, &unacceptableValueMsg, &isValid)

	return isValid
}

func (i *Invoice) Scan(rows db.Scanner) error {
	return rows.Scan(
		&i.ID, &i.Number, &i.Protocol, &i.Operation, &i.Cfop, &i.IsFinalCustomer, &i.IsIcmsContributor,
		&i.Shipping, &i.AddShippingToTotal, &i.Gta, &i.PDF, &i.ReqStatus, &i.ReqMsg,
		&i.Sender.ID, &i.Recipient.ID, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.ExtraNotes, &i.CustomFileName, &i.SenderIe,
	)
}

func (i *Invoice) FullScan(rows db.Scanner) error {
	return rows.Scan(
		&i.ID, &i.Number, &i.Protocol, &i.Operation, &i.Cfop, &i.IsFinalCustomer, &i.IsIcmsContributor,
		&i.Shipping, &i.AddShippingToTotal, &i.Gta, &i.PDF, &i.ReqStatus, &i.ReqMsg,
		&i.Sender.ID, &i.Recipient.ID, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.ExtraNotes, &i.CustomFileName, &i.SenderIe,

		&i.Sender.ID, &i.Sender.Name, &i.Sender.UserType, &i.Sender.CpfCnpj, &i.Sender.Ie, &i.Sender.Email, &i.Sender.Password,
		&i.Sender.Address.PostalCode, &i.Sender.Address.Neighborhood, &i.Sender.Address.StreetType, &i.Sender.Address.StreetName, &i.Sender.Address.Number,
		&i.Sender.CreatedBy, &i.Sender.CreatedAt, &i.Sender.UpdatedAt, &i.Sender.OtherIes,

		&i.Recipient.ID, &i.Recipient.Name, &i.Recipient.UserType, &i.Recipient.CpfCnpj, &i.Recipient.Ie, &i.Recipient.Email, &i.Recipient.Password,
		&i.Recipient.Address.PostalCode, &i.Recipient.Address.Neighborhood, &i.Recipient.Address.StreetType, &i.Recipient.Address.StreetName, &i.Recipient.Address.Number,
		&i.Recipient.CreatedBy, &i.Recipient.CreatedAt, &i.Recipient.UpdatedAt, &i.Sender.OtherIes,
	)
}

func (i *Invoice) GetAllSenderIes() []string {
	availableIes := make([]string, 0)

	if i.SenderIe != "" {
		availableIes = append(availableIes, i.SenderIe)
	}

	if i.Sender.Ie != i.SenderIe {
		availableIes = append(availableIes, i.Sender.Ie)
	}

	for _, ie := range i.Sender.OtherIes {
		if ie != i.SenderIe {
			availableIes = append(availableIes, ie)
		}
	}

	return availableIes
}
