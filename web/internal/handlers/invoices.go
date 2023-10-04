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
	Invoices         []*models.Invoice
	Invoice          *models.Invoice
	GeneralError     string
	FormMsg          string
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
		return c.Render("invoices", data, "layouts/base")
	}

	data.FormSelectFields.Entities = entities
	data.Invoice = models.NewEmptyInvoice()

	// get the latest 10 invoices
	invoices, err := workers.ListInvoices(c.Context())
	if err != nil {
		data.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices", data, "layouts/base")
	}

	data.Invoices = invoices

	return c.Render("invoices", data, "layouts/base")
}

func (page *InvoicesPage) RequireInvoice(c *fiber.Ctx) error {
	data := page.NewEmptyData()
	// formSuccess := true

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

	data.Invoice = invoice
	data.Invoice.Sender = sender
	data.Invoice.Recipient = recipient

	if !data.Invoice.IsValid() {
		data.FormMsg = "Corrija os campos abaixo."
		// formSuccess = false
		c.Set("HX-Retarget", "#invoice-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/invoice-form", data)
	}

	err = workers.CreateInvoice(c.Context(), data.Invoice)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	// i would call ss-api here in case formSuccess == true

	c.Set("HX-Trigger-After-Settle", "invoice-required")
	return c.Render("partials/request-card", data.Invoice)
}

func (page *InvoicesPage) GetItemFormSection(c *fiber.Ctx) error {
	item := models.NewEmptyInvoiceItem()
	return c.Render("partials/invoice-form-item-section", item)
}

func (page *InvoicesPage) GetRequestCardDetails(c *fiber.Ctx) error {
	invoiceId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.InvoiceNotFoundErr)
	}
	invoice, err := workers.RetrieveInvoice(c.Context(), invoiceId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", invoice)
}

func (page *InvoicesPage) GetInvoiceForm(c *fiber.Ctx) error {
	data := page.NewEmptyData()

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	data.FormSelectFields.Entities = entities

	invoiceId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.InvoiceNotFoundErr)
	}
	invoice, err := workers.RetrieveInvoice(c.Context(), invoiceId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	data.Invoice = invoice

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/invoice-form", data)
}
