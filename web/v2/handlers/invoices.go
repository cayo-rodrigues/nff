package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/ui/shared"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func InvoicesPage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	filters := c.Queries()

	invoices, err := services.ListInvoices(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	invoicesByDate := services.GroupListByDate(invoices)

	isAuthenticated := utils.IsAuthenticated(c)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter")
	return Render(c, layouts.Base(pages.InvoicesPage(invoicesByDate), isAuthenticated))
}

func CreateInvoicePage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
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
	isAuthenticated := utils.IsAuthenticated(c)
		return Render(c, layouts.Base(pages.InvoiceFormPage(baseInvoice, entities), isAuthenticated))
	}

	invoice := models.NewInvoiceWithSamples(entities)

	isAuthenticated := utils.IsAuthenticated(c)

	return Render(c, layouts.Base(pages.InvoiceFormPage(invoice, entities), isAuthenticated))
}

func ChooseInvoiceOperationPage(c *fiber.Ctx) error {
	isAuthenticated := utils.IsAuthenticated(c)
	return Render(c, layouts.Base(pages.ChooseInvoiceOperationPage(), isAuthenticated))
}

func GetSenderIeInput(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
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
	userID := utils.GetCurrentUserID(c)

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

	c.Append("HX-Trigger-After-Swap", "reload-invoice-list")
	return RetargetToPageHandler(c, "/invoices", InvoicesPage)
}

func RetrieveInvoiceItemsDetails(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
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
	userID := utils.GetCurrentUserID(c)
	filters := c.Queries()
	invoices, err := services.ListInvoices(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	invoicesByDate := services.GroupListByDate(invoices)

	return Render(c, components.InvoiceList(invoicesByDate))
}
