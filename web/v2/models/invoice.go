package models

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type Invoice struct {
	ID                   int            `json:"-"`
	Number               string         `json:"invoice_id"`
	Protocol             string         `json:"-"`
	Operation            string         `json:"operation"`
	Cfop                 int            `json:"cfop,string"`
	IsFinalCustomer      string         `json:"is_final_customer"`
	IsIcmsContributor    string         `json:"icms"`
	Shipping             float64        `json:"shipping"`
	AddShippingToTotal   string         `json:"add_shipping_to_total_value"`
	Gta                  string         `json:"gta"`
	Sender               *Entity        `json:"sender"`
	SenderIe             string         `json:"sender_ie"`
	Recipient            *Entity        `json:"recipient"`
	ReqStatus            string         `json:"-"`
	ReqMsg               string         `json:"-"`
	PDF                  string         `json:"invoice_pdf"`
	CreatedBy            int            `json:"-"`
	CreatedAt            time.Time      `json:"-"`
	UpdatedAt            time.Time      `json:"-"`
	ExtraNotes           string         `json:"extra_notes"`
	CustomFileNamePrefix string         `json:"custom_file_name_prefix"`
	FileName             string         `json:"file_name"`
	Items                []*InvoiceItem `json:"items"`
	Errors               ErrorMessages  `json:"-"`
}

func (i *Invoice) GetStatus() string {
	return i.ReqStatus
}

func NewInvoice() *Invoice {
	return &Invoice{
		Sender:    NewEntity(),
		Recipient: NewEntity(),
		Items:     []*InvoiceItem{},
	}
}

func NewInvoiceWithSamples(entities []*Entity) *Invoice {
	invoice := NewInvoice()
	if len(entities) > 0 {
		invoice.Sender = entities[0]
	}
	if len(invoice.Items) == 0 {
		invoice.Items = append(invoice.Items, NewInvoiceItem())
	}

	return invoice
}

func NewInvoiceFromForm(c *fiber.Ctx) *Invoice {
	var err error

	invoice := NewInvoice()

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
	invoice.CustomFileNamePrefix = strings.TrimSpace(c.FormValue("custom_file_name_prefix"))
	invoice.SenderIe = strings.TrimSpace(c.FormValue("sender_ie"))

	senderID, err := strconv.Atoi(c.FormValue("sender"))
	if err != nil {
		senderID = 0
	}

	recipientID, err := strconv.Atoi(c.FormValue("recipient"))
	if err != nil {
		recipientID = 0
	}

	invoice.Sender.ID = senderID
	invoice.Recipient.ID = recipientID

	invoice.Items = NewInvoiceItemsFromForm(c)

	return invoice
}

func (i *Invoice) IsValid() bool {
	fields := Fields{
		{
			Name:  "Operation",
			Value: i.Operation,
			Rules: Rules(OneOf(InvoiceOperations[:])),
		},
		{
			Name:  "Cfop",
			Value: i.Cfop,
			Rules: Rules(OneOf(InvoiceCfops[:])),
		},
		{
			Name:  "IsIcmsContributor",
			Value: i.IsIcmsContributor,
			Rules: Rules(OneOf(InvoiceIcmsOptions[:])),
		},
		{
			Name:  "IsFinalCustomer",
			Value: i.IsFinalCustomer,
			Rules: Rules(OneOf(InvoiceBooleanField[:])),
		},
		{
			Name:  "Shipping",
			Value: i.Shipping,
			Rules: Rules(Required),
		},
		{
			Name:  "AddShippingToTotal",
			Value: i.AddShippingToTotal,
			Rules: Rules(OneOf(InvoiceBooleanField[:])),
		},
		{
			Name:  "Gta",
			Value: i.Gta,
			Rules: Rules(Match(GTARegex), Max(16)),
		},
		{
			Name:  "ExtraNotes",
			Value: i.ExtraNotes,
			Rules: Rules(Max(512)),
		},
		{
			Name:  "CustomFileNamePrefix",
			Value: i.CustomFileNamePrefix,
			Rules: Rules(Max(64)),
		},
		{
			Name:  "SenderIe",
			Value: i.SenderIe,
			Rules: Rules(Required, Match(IEMGRegex)),
		},
		{
			Name:  "Sender",
			Value: i.Sender.ID,
			Rules: Rules(Required),
		},
		{
			Name:  "Recipient",
			Value: i.Recipient.ID,
			Rules: Rules(Required),
		},
	}
	errors, ok := Validate(fields)
	i.Errors = errors

	for _, item := range i.Items {
		ok = item.IsValid()
	}

	return ok
}

func (i *Invoice) Values() []any {
	return []any{
		&i.ID, &i.Number, &i.Protocol, &i.Operation, &i.Cfop, &i.IsFinalCustomer, &i.IsIcmsContributor,
		&i.Shipping, &i.AddShippingToTotal, &i.Gta, &i.PDF, &i.ReqStatus, &i.ReqMsg,
		&i.Sender.ID, &i.Recipient.ID, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.ExtraNotes, &i.CustomFileNamePrefix, &i.SenderIe, &i.FileName,
	}
}
