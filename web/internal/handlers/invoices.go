package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
)

type InvoicesPage struct{}

type InvoicesPageData struct {
	IsAuthenticated  bool
	Invoice          *models.Invoice
	GeneralError     string
	FormMsg          string
	FormSuccess      bool
	FormSelectFields *models.InvoiceFormSelectFields
}

func (page *InvoicesPage) NewEmptyData() *InvoicesPageData {
	return &InvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceFormSelectFields{
			Operations:   &globals.InvoiceOperations,
			Cfops:        &globals.InvoiceCfops,
			BooleanField: &globals.InvoiceBooleanField,
			IcmsOptions:  &globals.InvoiceIcmsOptions,
		},
	}
}

func (page *InvoicesPage) Render(c *fiber.Ctx) error {
	data := page.NewEmptyData()

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		data.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	data.FormSelectFields.Entities = entities
	data.Invoice = models.NewEmptyInvoice()

	return c.Render("invoices", data, "layouts/base")
}

func (page *InvoicesPage) RequireInvoice(c *fiber.Ctx) error {
	data := page.NewEmptyData()
	data.FormSuccess = true
	data.FormMsg = "Requerimento efetuado com sucesso! Acompanhe o progresso na sessão abaixo."

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	data.FormSelectFields.Entities = entities

	senderId, err := strconv.Atoi(c.FormValue("sender"))
	if err != nil {
		log.Println("Error converting sender id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}
	recipientId, err := strconv.Atoi(c.FormValue("recipient"))
	if err != nil {
		log.Println("Error converting recipient id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	sender, err := workers.RetrieveEntity(c.Context(), senderId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	recipient, err := workers.RetrieveEntity(c.Context(), recipientId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoice, err := models.NewInvoiceFromForm(c)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoice.Sender = sender
	invoice.Recipient = recipient

	if !invoice.IsValid() {
		data.FormMsg = "Corrija os campos abaixo."
		data.FormSuccess = false
	}

	data.Invoice = invoice

	// i would call ss-api here in case data.FormSuccess == true

	return c.Render("partials/invoice-form", data)
}

func (page *InvoicesPage) GetItemFormSection(c *fiber.Ctx) error {
	item := models.NewEmptyInvoiceItem()
	return c.Render("partials/invoice-form-item-section", item)
}
