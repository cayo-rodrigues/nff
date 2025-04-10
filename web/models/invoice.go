package models

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/cayo-rodrigues/safe"
	"github.com/gofiber/fiber/v2"
)

type Invoice struct {
	ID                   int                `json:"-"`
	Number               string             `json:"invoice_id"` // TROCAR PARA invoice_number NA SS-API
	Protocol             string             `json:"-"`
	Operation            string             `json:"operation"`
	IsInterstate         string             `json:"is_interstate"`
	Cfop                 string             `json:"cfop"`
	IsFinalCustomer      string             `json:"is_final_customer"`
	IsIcmsContributor    string             `json:"icms"`
	Shipping             float64            `json:"shipping"`
	AddShippingToTotal   string             `json:"add_shipping_to_total_value"`
	Gta                  string             `json:"gta"`
	Sender               *Entity            `json:"sender"`
	SenderIe             string             `json:"sender_ie"`
	RecipientIe          string             `json:"recipient_ie"`
	Recipient            *Entity            `json:"recipient"`
	ReqStatus            string             `json:"-"`
	ReqMsg               string             `json:"-"`
	PDF                  string             `json:"invoice_pdf"`
	CreatedBy            int                `json:"-"`
	CreatedAt            time.Time          `json:"-"`
	UpdatedAt            time.Time          `json:"-"`
	ExtraNotes           string             `json:"extra_notes"`
	CustomFileNamePrefix string             `json:"custom_file_name"` // TROCAR PARA custom_file_name_prefix NA SS-API
	FileName             string             `json:"file_name"`
	Items                []*InvoiceItem     `json:"items"`
	Errors               safe.ErrorMessages `json:"-"`
}

func (i *Invoice) AsNotification() *Notification {
	return &Notification{
		ID:            i.ID,
		Status:        i.ReqStatus,
		OperationType: "Emissão de NFA",
		PageEndpoint:  "/invoices",
		InvoicePDF:    i.PDF,
		CreatedAt:     i.CreatedAt,
		UserID:        i.CreatedBy,
	}
}

func (i *Invoice) GetCreatedAt() time.Time {
	return i.CreatedAt
}

func (i *Invoice) GetStatus() string {
	return i.ReqStatus
}

func (i *Invoice) IsCancelable() bool {
	return i.Number != "" && i.ReqStatus != "canceled"
}

func NewInvoice() *Invoice {
	return &Invoice{
		Sender:    NewEntity(),
		Recipient: NewEntity(),
		Items:     []*InvoiceItem{},
	}
}

func NewInvoiceWithSamples(entitiesByType *EntitiesByType) *Invoice {
	invoice := NewInvoice()
	if len(entitiesByType.Senders) > 0 {
		invoice.Sender = entitiesByType.Senders[0]
	}
	if len(entitiesByType.All) > 0 {
		invoice.Recipient = entitiesByType.All[0]
	}
	if len(invoice.Items) == 0 {
		invoice.Items = append(invoice.Items, NewInvoiceItem())
	}

	invoice.Operation = InvoiceOperations.VENDA()

	return invoice
}

func NewInvoiceFromForm(c *fiber.Ctx) *Invoice {
	var err error

	invoice := NewInvoice()

	invoice.Operation = strings.TrimSpace(c.FormValue("operation"))
	invoice.Cfop = strings.TrimSpace(c.FormValue("cfop"))
	invoice.IsInterstate = strings.TrimSpace(c.FormValue("is_interstate"))
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
	invoice.RecipientIe = strings.TrimSpace(c.FormValue("recipient_ie"))

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
	cfops := InvoiceCfops
	if i.IsInterstate == "Sim" {
		cfops = InterstateInvoiceCfops
	}
	fields := safe.Fields{
		{
			Name:  "Operation",
			Value: i.Operation,
			Rules: safe.Rules{safe.OneOf(InvoiceOperations[:])},
		},
		{
			Name:  "IsInterstate",
			Value: i.IsInterstate,
			Rules: safe.Rules{safe.OneOf(InvoiceBooleanField[:])},
		},
		{
			Name:  "Cfop",
			Value: i.Cfop,
			Rules: safe.Rules{safe.OneOf(cfops.ByOperation(i.Operation)[:])},
		},
		{
			Name:  "IsIcmsContributor",
			Value: i.IsIcmsContributor,
			Rules: safe.Rules{safe.OneOf(InvoiceIcmsOptions[:])},
		},
		{
			Name:  "IsFinalCustomer",
			Value: i.IsFinalCustomer,
			Rules: safe.Rules{safe.OneOf(InvoiceBooleanField[:])},
		},
		{
			Name:  "Shipping",
			Value: i.Shipping,
			Rules: safe.Rules{safe.Required()},
		},
		{
			Name:  "AddShippingToTotal",
			Value: i.AddShippingToTotal,
			Rules: safe.Rules{safe.OneOf(InvoiceBooleanField[:])},
		},
		{
			Name:  "Gta",
			Value: i.Gta,
			Rules: safe.Rules{safe.Match(GTARegex), safe.Max(16)},
		},
		{
			Name:  "ExtraNotes",
			Value: i.ExtraNotes,
			Rules: safe.Rules{safe.Max(512)},
		},
		{
			Name:  "CustomFileNamePrefix",
			Value: i.CustomFileNamePrefix,
			Rules: safe.Rules{safe.Max(64)},
		},
		{
			Name:  "SenderIe",
			Value: i.SenderIe,
			Rules: safe.Rules{safe.Required(), safe.Match(IEMGRegex)},
		},
		{
			Name:  "RecipientIe",
			Value: i.RecipientIe,
			Rules: safe.Rules{
				safe.RequiredUnless(safe.All(i.Recipient.Address.Values()...)),
				safe.Match(IEMGRegex),
			},
		},
		{
			Name:  "Sender",
			Value: i.Sender.ID,
			Rules: safe.Rules{safe.Required(), safe.NotEqualTo(i.Recipient.ID).WithMessage(utils.ConflictingEntitiesMsg)},
		},
		{
			Name:  "Recipient",
			Value: i.Recipient.ID,
			Rules: safe.Rules{safe.Required(), safe.NotEqualTo(i.Sender.ID).WithMessage(utils.ConflictingEntitiesMsg)},
		},
	}
	errors, ok := safe.Validate(fields)
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
		&i.ExtraNotes, &i.CustomFileNamePrefix, &i.SenderIe, &i.FileName, &i.RecipientIe,
		&i.IsInterstate,
	}
}
