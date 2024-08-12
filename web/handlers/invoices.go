package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/siare"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/ui/shared"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func InvoicesPage(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

	filters := c.Queries()

	invoices, err := services.ListInvoices(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	invoicesByDate := services.GroupListByDate(invoices)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(pages.InvoicesPage(invoicesByDate, entities)))
}

func GetInvoiceForm(c *fiber.Ctx) error {
	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	baseInvoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	if baseInvoiceID == 0 {
		invoice := models.NewInvoiceWithSamples(entities)
		return Render(c, forms.InvoiceForm(invoice, entities))
	}

	baseInvoice, err := services.RetrieveInvoice(c.Context(), baseInvoiceID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Settle", "open-invoice-form-dialog")
	return Render(c, forms.InvoiceForm(baseInvoice, entities))
}

func GetSenderIeInput(c *fiber.Ctx) error {
	entityID, err := strconv.Atoi(c.Query("sender"))
	if err != nil {
		return err
	}

	entity, err := services.RetrieveEntity(c.Context(), entityID)

	return Render(c, shared.SelectInput(&shared.InputData{
		ID:      "sender_ie",
		Label:   "IE do Remetente",
		Value:   entity.Ie,
		Options: &shared.InputOptions{StringOptions: entity.AllIes()},
	}))
}

func GetRecipientIeInput(c *fiber.Ctx) error {
	entityID, err := strconv.Atoi(c.Query("recipient"))
	if err != nil {
		return err
	}

	entity, err := services.RetrieveEntity(c.Context(), entityID)

	return Render(c, shared.SelectInput(&shared.InputData{
		ID:      "recipient_ie",
		Label:   "IE do Destinatário",
		Value:   entity.Ie,
		Options: &shared.InputOptions{StringOptions: entity.AllIes()},
	}))
}

func CreateInvoice(c *fiber.Ctx) error {

	invoice := models.NewInvoiceFromForm(c)

	sender, err := services.RetrieveEntity(c.Context(), invoice.Sender.ID)
	if err != nil {
		return err
	}

	recipient, err := services.RetrieveEntity(c.Context(), invoice.Recipient.ID)
	if err != nil {
		return err
	}

	invoice.Sender = sender
	invoice.Recipient = recipient

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	if !invoice.IsValid() {
		return Render(c, forms.InvoiceForm(invoice, entities))
	}

	err = services.CreateInvoice(c.Context(), invoice)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient()
	go ssapi.IssueInvoice(invoice)

	c.Append("HX-Trigger-After-Settle", "reload-invoice-list", "close-invoice-form-dialog")
	return Render(c, forms.InvoiceForm(invoice, entities))
}

func RetrieveInvoiceItemsDetails(c *fiber.Ctx) error {
	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	invoice, err := services.RetrieveInvoice(c.Context(), invoiceID)
	if err != nil {
		return err
	}

	return Render(c, components.InvoiceItemsDetails(invoice))
}

func ListInvoices(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	filters := c.Queries()
	invoices, err := services.ListInvoices(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	invoicesByDate := services.GroupListByDate(invoices)

	return Render(c, components.InvoiceList(invoicesByDate))
}

func RetrieveInvoiceCard(c *fiber.Ctx) error {
	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	invoice, err := services.RetrieveInvoice(c.Context(), invoiceID)
	if err != nil {
		return err
	}

	return Render(c, components.InvoiceCard(invoice))
}
