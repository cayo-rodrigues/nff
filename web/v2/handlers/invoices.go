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

	invoicesByDate := services.GroupListByDate(invoices)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter, highlight-current-page")
	return Render(c, layouts.Base(pages.InvoicesPage(invoicesByDate)))
}

func CreateInvoicePage(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	if baseInvoiceIDStr := c.Query("base-invoice-id"); baseInvoiceIDStr != "" {
		baseInvoiceID, err := strconv.Atoi(baseInvoiceIDStr)
		if err != nil {
			return err
		}
		baseInvoice, err := services.RetrieveInvoice(c.Context(), baseInvoiceID, userID)
		if err != nil {
			return err
		}
		return Render(c, layouts.Base(pages.InvoiceFormPage(baseInvoice, entities)))
	}

	invoice := models.NewInvoiceWithSamples(entities)

	c.Append("HX-Trigger-After-Settle", "highlight-current-page")
	return Render(c, layouts.Base(pages.InvoiceFormPage(invoice, entities)))
}

func ChooseInvoiceOperationPage(c *fiber.Ctx) error {
	return Render(c, layouts.Base(pages.ChooseInvoiceOperationPage()))
}

func GetSenderIeInput(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entityID, err := strconv.Atoi(c.Query("sender"))
	if err != nil {
		return err
	}
	entity, err := services.RetrieveEntity(c.Context(), entityID, userID)

	return Render(c, shared.SelectInput(&shared.InputData{
		ID:      "sender_ie",
		Label:   "IE do Remetente",
		Value:   entity.Ie,
		Options: &shared.InputOptions{StringOptions: entity.AllIes()},
	}))
}

func CreateInvoice(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

	invoice := models.NewInvoiceFromForm(c)

	sender, err := services.RetrieveEntity(c.Context(), invoice.Sender.ID, userID)
	if err != nil {
		return err
	}

	recipient, err := services.RetrieveEntity(c.Context(), invoice.Recipient.ID, userID)
	if err != nil {
		return err
	}

	invoice.Sender = sender
	invoice.Recipient = recipient

	if !invoice.IsValid() {
		entities, err := services.ListEntities(c.Context(), userID)
		if err != nil {
			return err
		}
		return RetargetToForm(c, "invoice", forms.InvoiceForm(invoice, entities))
	}

	err = services.CreateInvoice(c.Context(), invoice, userID)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient()
	go ssapi.IssueInvoice(invoice)

	c.Append("HX-Trigger-After-Swap", "reload-invoice-list")
	return RetargetToPageHandler(c, "/invoices", InvoicesPage)
}

func RetrieveInvoiceItemsDetails(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	invoice, err := services.RetrieveInvoice(c.Context(), invoiceID, userID)
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
